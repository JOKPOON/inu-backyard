package repository

import (
	"fmt"

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
	err := r.gorm.First(&course).Where("id = ?", id).Error
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
		clo.expected_passing_assignment_percentage,
		a.id AS assignment_id,
		a.max_score AS max_score,
		a.expected_score_percentage,
		sc.score AS score
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
		c.id = ?
	GROUP BY
		a.id,
		s.id,
		clo.id,
		clo.code,
		clo.expected_passing_assignment_percentage,
		sc.score
	ORDER BY
		s.id,
		clo.id;
	`, courseId).Scan(&cloResults).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get students passing CLOs: %w", err)
	}

	resp := &entity.StudentPassCLOResp{
		Clos:   []string{},
		Result: []entity.StudentCLO{},
	}

	cloCode := make(map[string]string)
	students := make(map[int]map[string][]bool)
	for _, cloResult := range cloResults {
		if _, ok := students[cloResult.StudentId]; !ok {
			students[cloResult.StudentId] = make(map[string][]bool)
		}

		if _, ok := students[cloResult.StudentId][cloResult.CLOId]; !ok {
			students[cloResult.StudentId][cloResult.CLOId] = make([]bool, 0)
		}

		if _, ok := cloCode[cloResult.CLOId]; !ok {
			cloCode[cloResult.CLOId] = cloResult.CLOCode
			resp.Clos = append(resp.Clos, cloResult.CLOCode)
		}

		score := cloResult.Score / float64(cloResult.MaxScore) * 100
		pass := score > cloResult.ExpectedScorePercent

		students[cloResult.StudentId][cloResult.CLOId] = append(students[cloResult.StudentId][cloResult.CLOId], pass)
	}

	for studentId, cloResults := range students {
		for cloId, results := range cloResults {
			for _, result := range results {
				println(studentId, cloId, cloCode[cloId], result)
			}
		}
	}

	for studentId, cloResults := range students {
		clos := []entity.CLO{}
		for cloId, results := range cloResults {
			pass := true
			for _, result := range results {
				if !result {
					pass = false
					break
				}
			}

			clos = append(clos, entity.CLO{
				Id:   cloId,
				Code: cloCode[cloId],
				Pass: pass,
			})
		}

		resp.Result = append(resp.Result, entity.StudentCLO{
			StudentID: studentId,
			CLOs:      clos,
		})
	}

	return resp, nil
}
