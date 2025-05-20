package repository

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	"gorm.io/gorm"
)

type courseRepositoryGorm struct {
	gorm *gorm.DB
}

func NewCourseRepositoryGorm(gorm *gorm.DB) entity.CourseRepository {
	return &courseRepositoryGorm{gorm: gorm}
}

func (r courseRepositoryGorm) GetAll(query, serm, program string) ([]entity.Course, error) {
	var courses []entity.Course
	tx := r.gorm.Preload("Lecturers").Preload("Semester").Preload("Programme")

	if query != "" {
		tx = tx.Where("name LIKE ? OR code LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if serm != "" {
		tx = tx.Where("semester_id = ?", serm)
	}
	if program != "" {
		tx = tx.Where("programme_id = ?", program)
	}

	err := tx.Find(&courses).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errs.New(errs.ErrCourseNotFound, "no course found")
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get all courses: %w", err)
	}

	return courses, nil
}

func (r courseRepositoryGorm) GetById(id string) (*entity.Course, error) {
	var course entity.Course
	err := r.gorm.Preload("Lecturers").Preload("Semester").Preload("Programme").Preload("Programme.Department").Preload("Programme.Department.Faculty").Where("id = ?", id).First(&course).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get course by id: %w", err)
	}

	return &course, nil
}

func (r courseRepositoryGorm) GetByUserId(userId string, query string, year string, program string) ([]entity.Course, error) {
	var courses []entity.Course
	db := r.gorm.
		Preload("Lecturers").Preload("Semester").Preload("Programme").
		Joins("JOIN course_lecturer cl ON cl.course_id = course.id").
		Where("user_id = ?", userId)

	if query != "" {
		db = db.Where("name LIKE ?", "%"+query+"%")
	}
	if year != "" {
		db = db.Where("academic_year = ?", year)
	}
	if program != "" {
		db = db.Where("programme_id = ?", program)
	}

	err := db.Find(&courses).Error
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

	return nil
}

func (r courseRepositoryGorm) Update(id string, course *entity.Course) error {
	err := r.gorm.Model(&entity.Course{}).Where("id = ?", id).Updates(course).Error
	if err != nil {
		return fmt.Errorf("cannot update course: %w", err)
	}

	return nil
}

func (r courseRepositoryGorm) Delete(id string) error {
	err := r.gorm.Delete(&entity.Course{Id: id}).Error

	if err != nil {
		return fmt.Errorf("cannot delete course: %w", err)
	}

	return nil
}

func (r courseRepositoryGorm) CreateLinkWithLecturer(courseId string, lecturerIds []string) error {
	query := ""
	for _, lecturerId := range lecturerIds {
		query += fmt.Sprintf("('%s', '%s'),", lecturerId, courseId)
	}

	query = query[:len(query)-1]

	err := r.gorm.Exec(fmt.Sprintf("INSERT IGNORE INTO `course_lecturer` (user_id, course_id) VALUES %s", query)).Error
	if err != nil {
		return fmt.Errorf("cannot create link between lecturer and course: %w", err)
	}

	return nil
}

func (r courseRepositoryGorm) ReplaceLecturersForCourse(courseId string, lecturerIds []string) error {
	tx := r.gorm.Begin() // Start transaction

	// Step 1: Delete lecturers who are no longer in the list
	if err := tx.Exec(`
		DELETE FROM course_lecturer
		WHERE course_id = ? AND user_id NOT IN (?)`, courseId, lecturerIds).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete old lecturers: %w", err)
	}

	// Step 2: Insert new lecturers (Ignore duplicates to avoid errors)
	if len(lecturerIds) > 0 {
		query := ""
		for _, lecturerId := range lecturerIds {
			query += fmt.Sprintf("('%s', '%s'),", lecturerId, courseId)
		}

		query = query[:len(query)-1] // Remove trailing comma

		if err := tx.Exec(fmt.Sprintf(`
			INSERT IGNORE INTO course_lecturer (user_id, course_id) VALUES %s`, query)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert new lecturers: %w", err)
		}
	}

	// Commit transaction
	return tx.Commit().Error
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
		s.first_name_th AS student_first_name_th,
		s.last_name_th AS student_last_name_th,
		s.first_name_en AS student_first_name_en,
		s.last_name_en AS student_last_name_en,
		clo.id AS clo_id,
		clo.code AS clo_code,
		clo.expected_passing_assignment_percentage,
		a.id AS assignment_id,
		a.max_score AS max_score,
		a.expected_score_percentage,
		sc.score AS score
	FROM
		student s
	LEFT JOIN enrollment e ON
		s.id = e.student_id
	LEFT JOIN course c ON
		e.course_id = c.id
	LEFT JOIN course_learning_outcome clo ON
		c.id = clo.course_id
	LEFT JOIN clo_assignment ca ON
		clo.id = ca.course_learning_outcome_id
	LEFT JOIN assignment a ON
		ca.assignment_id = a.id
	LEFT JOIN score sc ON
		a.id = sc.assignment_id AND s.id = sc.student_id
	WHERE
		c.id = ?
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
		clos := []entity.PassOutcome{}
		for cloId, results := range cloResults {
			pass := true
			for _, result := range results {
				if !result {
					pass = false
					break
				}
			}

			clos = append(clos, entity.PassOutcome{
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

	for _, cloResult := range cloResults {
		for id, result := range resp.Result {
			if cloResult.StudentId == result.StudentID {
				resp.Result[id].StudentNameEN = cloResult.StudentFirstNameEN + " " + cloResult.StudentLastNameEN
				resp.Result[id].StudentNameTH = cloResult.StudentFirstNameTH + " " + cloResult.StudentLastNameTH
				break
			}
		}
	}

	pass := r.GetStudentsPassingOutcomes(courseId, resp)

	return pass, nil
}

func (r courseRepositoryGorm) GetStudentsPassingOutcomes(
	courseId string,
	clos *entity.StudentPassCLOResp,
) *entity.StudentPassCLOResp {
	type Outcome struct {
		Id   string `json:"id"`
		Code string `json:"code"`
	}
	type StudentCLO struct {
		POs  []Outcome `json:"po"`
		PLOs []Outcome `json:"plo"`
		SOs  []Outcome `json:"so"`
	}
	type StudentData struct {
		StudentID                    int                   `json:"student_id"`
		CLOMapping                   map[string]StudentCLO `json:"clo_mapping"`
		ExpectedPassingCLOPercentage float64               `json:"expected_passing_clo_percentage"`
	}
	type RawQueryResult struct {
		StudentID                    int     `json:"student_id"`
		CloId                        string  `json:"clo_id"`
		PoId                         string  `json:"po_id"`
		PoCode                       string  `json:"po_code"`
		PloId                        string  `json:"plo_id"`
		PloCode                      string  `json:"plo_code"`
		SoId                         string  `json:"so_id"`
		SoCode                       string  `json:"so_code"`
		ExpectedPassingCLOPercentage float64 `json:"expected_passing_clo_percentage"`
	}

	var rawQueryResults []RawQueryResult

	err := r.gorm.Raw(`
		SELECT s.id AS student_id, clo.id AS clo_id, po.id AS po_id, po.code AS po_code,
			plo.id AS plo_id, plo.code AS plo_code, so.id AS so_id, so.code AS so_code,
			c.expected_passing_clo_percentage
		FROM student s
		JOIN enrollment e ON s.id = e.student_id
		JOIN course c ON e.course_id = c.id
		JOIN course_learning_outcome clo ON c.id = clo.course_id
		LEFT JOIN clo_po cpo ON clo.id = cpo.course_learning_outcome_id
		LEFT JOIN program_outcome po ON cpo.program_outcome_id = po.id
		LEFT JOIN clo_subplo csplo ON clo.id = csplo.course_learning_outcome_id
		LEFT JOIN sub_program_learning_outcome splo ON csplo.sub_program_learning_outcome_id = splo.id
		LEFT JOIN program_learning_outcome plo ON splo.program_learning_outcome_id = plo.id
		LEFT JOIN clo_subso csso ON clo.id = csso.course_learning_outcome_id
		LEFT JOIN sub_student_outcome sso ON csso.sub_student_outcome_id = sso.id
		LEFT JOIN student_outcome so ON sso.student_outcome_id = so.id
		WHERE c.id = ?
		ORDER BY s.id, clo.id, po.id, plo.id, so.id`, courseId).Scan(&rawQueryResults).Error
	if err != nil {
		return nil
	}

	studentMap := make(map[int]StudentData)
	for _, row := range rawQueryResults {
		student := studentMap[row.StudentID]
		if student.CLOMapping == nil {
			student = StudentData{
				StudentID:                    row.StudentID,
				CLOMapping:                   make(map[string]StudentCLO),
				ExpectedPassingCLOPercentage: row.ExpectedPassingCLOPercentage,
			}
		}

		clo := student.CLOMapping[row.CloId]

		addUniqueOutcome := func(outcomes *[]Outcome, id, code string) {
			for _, outcome := range *outcomes {
				if outcome.Id == id {
					return
				}
			}
			*outcomes = append(*outcomes, Outcome{Id: id, Code: code})
		}

		addUniqueOutcome(&clo.POs, row.PoId, row.PoCode)
		addUniqueOutcome(&clo.PLOs, row.PloId, row.PloCode)
		addUniqueOutcome(&clo.SOs, row.SoId, row.SoCode)

		student.CLOMapping[row.CloId] = clo
		studentMap[row.StudentID] = student
	}

	for id, student := range clos.Result {
		poCount, ploCount, soCount := make(map[string]int), make(map[string]int), make(map[string]int)
		poTotal, ploTotal, soTotal := make(map[string]int), make(map[string]int), make(map[string]int)
		poCode, ploCode, soCode := make(map[string]string), make(map[string]string), make(map[string]string)

		for _, clo := range student.CLOs {
			mapping, exists := studentMap[student.StudentID].CLOMapping[clo.Id]
			if !exists {
				continue
			}
			for _, po := range mapping.POs {
				poCode[po.Id] = po.Code
				poTotal[po.Id]++
				if clo.Pass {
					poCount[po.Id]++
				}
			}
			for _, plo := range mapping.PLOs {
				ploCode[plo.Id] = plo.Code
				ploTotal[plo.Id]++
				if clo.Pass {
					ploCount[plo.Id]++
				}
			}
			for _, so := range mapping.SOs {
				soCode[so.Id] = so.Code
				soTotal[so.Id]++
				if clo.Pass {
					soCount[so.Id]++
				}
			}
		}

		calculatePassOutcomes := func(count, total map[string]int, codes map[string]string) []entity.PassOutcome {
			var passOutcomes []entity.PassOutcome
			for id := range total {
				passOutcomes = append(passOutcomes, entity.PassOutcome{
					Id:   id,
					Code: codes[id],
					Pass: (float64(count[id])/float64(total[id]))*100 >= studentMap[student.StudentID].ExpectedPassingCLOPercentage,
				})
			}
			return passOutcomes
		}

		clos.Result[id].POs = calculatePassOutcomes(poCount, poTotal, poCode)
		clos.Result[id].PLOs = calculatePassOutcomes(ploCount, ploTotal, ploCode)
		clos.Result[id].SOs = calculatePassOutcomes(soCount, soTotal, soCode)
	}

	return clos
}
