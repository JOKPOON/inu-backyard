package repository

import (
	"fmt"
	"math"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

type courseRepositoryGorm struct {
	gorm *gorm.DB
}

func NewCourseRepositoryGorm(gorm *gorm.DB) entity.CourseRepository {
	return &courseRepositoryGorm{gorm: gorm}
}

func (r courseRepositoryGorm) GetAll() ([]entity.Course, error) {
	var courses []entity.Course
	err := r.gorm.Preload("Lecturers").Preload("Semester").Preload("Programme").Find(&courses).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get courses: %w", err)
	}

	return courses, nil
}

func (r courseRepositoryGorm) GetById(id string) (*entity.Course, error) {
	var course entity.Course
	err := r.gorm.Where("id = ?", id).First(&course).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get course by id: %w", err)
	}

	return &course, nil
}

func (r courseRepositoryGorm) GetByUserId(userId string) ([]entity.Course, error) {
	var courses []entity.Course
	err := r.gorm.Where("user_id = ?", userId).Find(&courses).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get course by user id: %w", err)
	}

	return courses, nil
}

func (r courseRepositoryGorm) Create(course *entity.Course) error {
	err := r.gorm.Create(&course).Error
	if err != nil {
		return fmt.Errorf("cannot create course: %w", err)
	}
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r courseRepositoryGorm) Update(id string, course *entity.Course) error {
	err := r.gorm.Model(&entity.Course{}).Where("id = ?", id).Updates(course).Error
	if err != nil {
		return fmt.Errorf("cannot update course: %w", err)
	}
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r courseRepositoryGorm) Delete(id string) error {
	err := r.gorm.Delete(&entity.Course{Id: id}).Error

	if err != nil {
		return fmt.Errorf("cannot delete course: %w", err)
	}
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r courseRepositoryGorm) CreateLinkWithLecturer(courseId string, lecturerIds []string) error {
	query := ""
	for _, lecturerId := range lecturerIds {
		query += fmt.Sprintf("('%s', '%s'),", lecturerId, courseId)
	}

	query = query[:len(query)-1]

	err := r.gorm.Exec(fmt.Sprintf("INSERT INTO `course_lecturer` (user_id, course_id) VALUES %s", query)).Error
	if err != nil {
		return fmt.Errorf("cannot create link between lecturer and course: %w", err)
	}

	return nil
}

func (r courseRepositoryGorm) DeleteLinkWithLecturer(courseId string, lecturerIds []string) error {
	query := ""
	for _, lecturerId := range lecturerIds {
		query += fmt.Sprintf("('%s', '%s'),", lecturerId, courseId)
	}

	query = query[:len(query)-1]

	err := r.gorm.Exec(fmt.Sprintf("DELETE FROM `course_lecturer` WHERE (user_id, course_id) IN (%s)", query)).Error
	if err != nil {
		return fmt.Errorf("cannot delete link between lecturer and course: %w", err)
	}

	return nil
}

func (r courseRepositoryGorm) GetStudentsPassingCLOs(courseId string) (*entity.StudentPassCLOResp, error) {
	var cloResults []entity.CLOResult
	err := r.gorm.Raw(`
	SELECT
		s.id AS student_id,
		clo.id AS clo_id,
		clo.code AS clo_code,
		COUNT(DISTINCT a.id) AS passed_assignments,
		(
		SELECT
			COUNT(*)
		FROM
			clo_assignment ca2
		WHERE
			ca2.course_learning_outcome_id = clo.id
	) AS total_assignments,
	clo.expected_passing_assignment_percentage
	FROM
		student s
	JOIN enrollment e ON
		s.id = e.student_id
	JOIN course c ON
		e.course_id = c.id
	JOIN course_learning_outcome clo ON
		c.id = clo.course_id
	JOIN clo_assignment ca ON
		clo.id = ca.course_learning_outcome_id
	JOIN assignment a ON
		ca.assignment_id = a.id
	JOIN score sc ON
		a.id = sc.assignment_id AND s.id = sc.student_id
	WHERE
		sc.score >=(
			a.max_score * a.expected_score_percentage / 100
		) AND c.id = ?
	GROUP BY
		s.id,
		clo.id,
		clo.code,
		clo.expected_passing_assignment_percentage
	ORDER BY
		s.id,
		clo.id;
	`, courseId).Scan(&cloResults).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get students passing CLOs: %w", err)
	}

	studentCLOs := make(map[int]*entity.StudentCLO)

	for _, res := range cloResults {
		passingThreshold := int(math.Ceil(float64(res.TotalAssignments) * (res.ExpectedPassingAssignmentPerc / 100)))
		if res.PassedAssignments >= passingThreshold {
			if _, exists := studentCLOs[res.StudentID]; !exists {
				studentCLOs[res.StudentID] = &entity.StudentCLO{
					StudentID: res.StudentID,
					PassCLO:   []string{},
				}
			}

			studentCLOs[res.StudentID].PassCLO = append(studentCLOs[res.StudentID].PassCLO, res.CLOCode)
		}
	}

	cloCodes := []string{}
	err = r.gorm.Raw(`
		SELECT clo.code AS clo_code
		FROM course c
		JOIN course_learning_outcome clo ON c.id = clo.course_id
		WHERE c.id = ?;
	`, courseId).Scan(&cloCodes).Error

	if err != nil {
		return nil, err
	}

	var students []entity.StudentCLO
	for _, student := range studentCLOs {
		students = append(students, *student)
	}

	output := &entity.StudentPassCLOResp{
		Clos:   cloCodes,
		Result: students,
	}

	return output, nil
}
