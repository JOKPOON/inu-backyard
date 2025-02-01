package repository

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

type programmeRepositoryGorm struct {
	gorm *gorm.DB
}

func NewProgrammeRepositoryGorm(gorm *gorm.DB) entity.ProgrammeRepository {
	return &programmeRepositoryGorm{gorm}
}

func (r programmeRepositoryGorm) GetAll() ([]entity.Programme, error) {
	var programs []entity.Programme

	err := r.gorm.Preload("ProgramOutcomes").Preload("ProgramLearningOutcomes").Preload("StudentOutcomes").Find(&programs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get programs: %w", err)
	}

	return programs, nil
}

func (r programmeRepositoryGorm) GetById(id string) (*entity.Programme, error) {
	var programme *entity.Programme

	err := r.gorm.Where("id = ?", id).First(&programme).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get programme by id: %w", err)
	}

	return programme, nil
}

func (r programmeRepositoryGorm) GetByName(nameTH string, nameEN string) ([]entity.Programme, error) {
	var programme []entity.Programme

	err := r.gorm.Find(&programme, "name_th = ? OR name_en = ?", nameTH, nameEN).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get programme by id: %w", err)
	}

	return programme, nil
}

func (r programmeRepositoryGorm) GetByNameAndYear(nameTH string, nameEN string, year string) (*entity.Programme, error) {
	var programme *entity.Programme

	err := r.gorm.Where("(name_th = ? OR name_en = ?) AND year = ?", nameTH, nameEN, year).First(&programme).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get programme by id: %w", err)
	}

	return programme, nil
}

func (r programmeRepositoryGorm) Create(programme *entity.Programme) error {
	err := r.gorm.Create(&programme).Error
	if err != nil {
		return fmt.Errorf("cannot create programme: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) Update(id string, programme *entity.Programme) error {
	err := r.gorm.Model(&entity.Programme{}).Where("id = ?", id).Updates(programme).Error
	if err != nil {
		return fmt.Errorf("cannot update programme: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) Delete(id string) error {
	err := r.gorm.Where("id = ?", id).Delete(&entity.Programme{}).Error
	if err != nil {
		return fmt.Errorf("cannot delete programme: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) FilterExisted(namesTH []string, namesEN []string) ([]string, error) {
	var existedNames []string

	err := r.gorm.Model(&entity.Programme{}).Where("name_th IN (?) OR name_en IN (?)", namesTH, namesEN).Pluck("name_th", &existedNames).Error
	if err != nil {
		return nil, fmt.Errorf("cannot filter existed programme names: %w", err)
	}

	return existedNames, nil
}

func (r programmeRepositoryGorm) CreateLinkWithPO(programmeId string, poId string) error {
	err := r.gorm.Exec(`
	INSERT INTO programme_po(programme_id, program_outcome_id)
	VALUES (?, ?);
	`, programmeId, poId).Error
	if err != nil {
		return fmt.Errorf("cannot create link between programme and program outcome: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) CreateLinkWithPLO(programmeId string, ploId string) error {
	err := r.gorm.Exec(`
	INSERT IF NOT EXISTS INTO programme_plo(programme_id, program_learning_outcome_id)
	VALUES (?, ?);
	`, programmeId, ploId).Error
	if err != nil {
		return fmt.Errorf("cannot create link between programme and program learning outcome: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) CreateLinkWithSO(programmeId string, soId string) error {
	err := r.gorm.Exec(`
	INSERT INTO programme_so (programme_id, student_outcome_id)
	VALUES (?, ?);
	`, programmeId, soId).Error
	if err != nil {
		return fmt.Errorf("cannot create link between programme and student outcome: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) DeleteLinkWithPO(programmeId string, poId string) error {
	err := r.gorm.Exec(`
	DELETE FROM programme_program_outcome
	WHERE programme_id = ? AND program_outcome_id = ?;
	`, programmeId, poId).Error
	if err != nil {
		return fmt.Errorf("cannot delete link between programme and program outcome: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) DeleteLinkWithPLO(programmeId string, ploId string) error {
	err := r.gorm.Exec(`
	DELETE FROM programme_program_learning_outcome
	WHERE programme_id = ? AND program_learning_outcome_id = ?;
	`, programmeId, ploId).Error
	if err != nil {
		return fmt.Errorf("cannot delete link between programme and program learning outcome: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) DeleteLinkWithSO(programmeId string, soId string) error {
	err := r.gorm.Exec(`
	DELETE FROM programme_student_outcome
	WHERE programme_id = ? AND student_outcome_id = ?;
	`, programmeId, soId).Error
	if err != nil {
		return fmt.Errorf("cannot delete link between programme and student outcome: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) GetAllCourseOutcomeLinked(programmeId string) ([]entity.CourseOutcomes, error) {
	type gormResult struct {
		CourseCode string `json:"course_code"`
		CourseName string `json:"course_name"`
		POCode     string `json:"po_code"`
		PLOCode    string `json:"plo_code"`
		SPLOCode   string `json:"splo_code"`
		SOCode     string `json:"so_code"`
		SSOCode    string `json:"sso_code"`
		CLOCode    string `json:"clo_code"`
	}

	var result []gormResult
	err := r.gorm.Raw(`
	SELECT
		c.code AS course_code,
		c.name AS course_name,
		po.code AS po_code,
		plo.code AS plo_code,
		splo.code AS splo_code,
		so.code AS so_code,
		sso.code AS sso_code,
		clo.code AS clo_code
	FROM
		course c
	LEFT JOIN course_learning_outcome clo ON
		clo.course_id = c.id
	LEFT JOIN clo_po cpo ON
		cpo.course_learning_outcome_id = clo.id
	LEFT JOIN program_outcome po ON
		po.id = cpo.program_outcome_id
	LEFT JOIN clo_subplo cplo ON
		cplo.course_learning_outcome_id = clo.id
	LEFT JOIN sub_program_learning_outcome splo ON
		splo.id = cplo.sub_program_learning_outcome_id
	LEFT JOIN program_learning_outcome plo ON
		plo.id = splo.program_learning_outcome_id
	LEFT JOIN clo_subso cso ON
		cso.course_learning_outcome_id = clo.id
	LEFT JOIN sub_student_outcome sso ON
		sso.id = cso.sub_student_outcome_id
	LEFT JOIN student_outcome so ON
		so.id = sso.student_outcome_id
	WHERE
		c.programme_id = ?;
	`, programmeId).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	// outcomes := &entity.ProgrammeOutcomes{
	// 	POs:      make(map[string][]string),
	// 	PLO_SPLO: make(map[string][]string),
	// 	SO_SSO:   make(map[string][]string),
	// }
	courseMap := make(map[string]*entity.CourseOutcomes)
	for _, res := range result {
		if _, exists := courseMap[res.CourseCode]; !exists {
			courseMap[res.CourseCode] = &entity.CourseOutcomes{
				CourseCode: res.CourseCode,
				CourseName: res.CourseName,
				CLOs:       []string{},
				POs:        []string{},
				PLOs:       make(map[string][]string),
				SOs:        make(map[string][]string),
			}
		}
		course := courseMap[res.CourseCode]

		// Add CLO if not already added
		if res.CLOCode != "" && !contains(course.CLOs, res.CLOCode) {
			course.CLOs = append(course.CLOs, res.CLOCode)
		}

		// Add PO if not already added
		if res.POCode != "" && !contains(course.POs, res.POCode) {
			course.POs = append(course.POs, res.POCode)
		}

		// Add PLO and SPLO
		if res.PLOCode != "" && res.SPLOCode != "" {
			if _, exists := course.PLOs[res.PLOCode]; !exists {
				course.PLOs[res.PLOCode] = []string{}
			}
			if !contains(course.PLOs[res.PLOCode], res.SPLOCode) {
				course.PLOs[res.PLOCode] = append(course.PLOs[res.PLOCode], res.SPLOCode)
			}
		}

		// Add SO and SSO
		if res.SOCode != "" && res.SSOCode != "" {
			if _, exists := course.SOs[res.SOCode]; !exists {
				course.SOs[res.SOCode] = []string{}
			}
			if !contains(course.SOs[res.SOCode], res.SSOCode) {
				course.SOs[res.SOCode] = append(course.SOs[res.SOCode], res.SSOCode)
			}
		}
	}

	var courses []entity.CourseOutcomes
	for _, course := range courseMap {
		courses = append(courses, *course)
	}

	return courses, nil
}

func (r programmeRepositoryGorm) GetAllCourseLinkedPO(programmeId string) (*entity.ProgrammeLinkedPO, error) {
	type gormResult struct {
		CourseCode string `json:"course_code"`
		CourseName string `json:"course_name"`
		POCode     string `json:"po_code"`
	}

	var result []gormResult
	err := r.gorm.Raw(`
	SELECT
		c.code AS course_code,
		c.name AS course_name,
		po.code AS po_code
	FROM
		course c
	LEFT JOIN course_learning_outcome clo ON
		clo.course_id = c.id
	LEFT JOIN clo_po cpo ON
		cpo.course_learning_outcome_id = clo.id
	LEFT JOIN program_outcome po ON
		po.id = cpo.program_outcome_id
	WHERE
		c.programme_id = ?;
	`, programmeId).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	POs, err := r.GetAllPO(programmeId)
	if err != nil {
		return nil, err
	}

	pos := make([]string, len(POs))
	for i, po := range POs {
		pos[i] = po.Code
	}

	program := &entity.ProgrammeLinkedPO{
		AllCourse: []string{},
		AllPOs:    pos,
	}

	courseMap := make(map[string]*entity.CourseLinkedPO)
	for _, res := range result {
		if _, exists := courseMap[res.CourseCode]; !exists {
			courseMap[res.CourseCode] = &entity.CourseLinkedPO{
				CourseCode: res.CourseCode,
				CourseName: res.CourseName,
				POs:        []string{},
			}
			program.AllCourse = append(program.AllCourse, res.CourseCode)
		}
		course := courseMap[res.CourseCode]

		if res.POCode != "" && !contains(course.POs, res.POCode) {
			course.POs = append(course.POs, res.POCode)
		}
	}

	var courses []entity.CourseLinkedPO
	for _, course := range courseMap {
		courses = append(courses, *course)
	}

	program.CourseLinkedPOs = courses

	return program, nil
}

func (r programmeRepositoryGorm) GetAllCourseLinkedPLO(programmeId string) (*entity.ProgrammeLinkedPLO, error) {
	type gormResult struct {
		CourseCode string `json:"course_code"`
		CourseName string `json:"course_name"`
		PLOCode    string `json:"plo_code"`
		SPLOCode   string `json:"splo_code"`
	}

	var result []gormResult
	err := r.gorm.Raw(`
	SELECT
		c.code AS course_code,
		c.name AS course_name,
		plo.code AS plo_code,
		splo.code AS splo_code
	FROM
		course c
	LEFT JOIN course_learning_outcome clo ON
		clo.course_id = c.id
	LEFT JOIN clo_subplo cplo ON
		cplo.course_learning_outcome_id = clo.id
	LEFT JOIN sub_program_learning_outcome splo ON
		splo.id = cplo.sub_program_learning_outcome_id
	LEFT JOIN program_learning_outcome plo ON
		plo.id = splo.program_learning_outcome_id
	WHERE
		c.programme_id = ?;
	`, programmeId).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	PLOs, err := r.GetAllPLO(programmeId)
	if err != nil {
		return nil, err
	}

	var plos = make(map[string][]string)
	for _, plo := range PLOs {
		plos[plo.Code] = []string{}
		for _, splo := range plo.SubProgramLearningOutcomes {
			plos[plo.Code] = append(plos[plo.Code], splo.Code)
		}
	}

	program := &entity.ProgrammeLinkedPLO{
		AllCourse: []string{},
		AllPLOs:   plos,
	}

	courseMap := make(map[string]*entity.CourseLinkedPLO)
	for _, res := range result {
		if _, exists := courseMap[res.CourseCode]; !exists {
			courseMap[res.CourseCode] = &entity.CourseLinkedPLO{
				CourseCode: res.CourseCode,
				CourseName: res.CourseName,
				PLOs:       make(map[string][]string),
			}
			program.AllCourse = append(program.AllCourse, res.CourseCode)
		}
		course := courseMap[res.CourseCode]

		// Add PLO and SPLO
		if res.PLOCode != "" && res.SPLOCode != "" {
			if _, exists := course.PLOs[res.PLOCode]; !exists {
				course.PLOs[res.PLOCode] = []string{}
			}
			if !contains(course.PLOs[res.PLOCode], res.SPLOCode) {
				course.PLOs[res.PLOCode] = append(course.PLOs[res.PLOCode], res.SPLOCode)
			}
		}
	}

	var courses []entity.CourseLinkedPLO
	for _, course := range courseMap {
		courses = append(courses, *course)
	}

	program.CourseLinkedPLOs = courses

	return program, nil
}

func (r programmeRepositoryGorm) GetAllCourseLinkedSO(programmeId string) (*entity.ProgrammeLinkedSO, error) {
	type gormResult struct {
		CourseCode string `json:"course_code"`
		CourseName string `json:"course_name"`
		SOCode     string `json:"so_code"`
		SSOCode    string `json:"sso_code"`
	}

	var result []gormResult
	err := r.gorm.Raw(`
	SELECT
		c.code AS course_code,
		c.name AS course_name,
		so.code AS so_code,
		sso.code AS sso_code
	FROM
		course c
	LEFT JOIN course_learning_outcome clo ON
		clo.course_id = c.id
	LEFT JOIN clo_subso cso ON
		cso.course_learning_outcome_id = clo.id
	LEFT JOIN sub_student_outcome sso ON
		sso.id = cso.sub_student_outcome_id
	LEFT JOIN student_outcome so ON
		so.id = sso.student_outcome_id
	WHERE
		c.programme_id = ?;
	`, programmeId).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	SOs, err := r.GetAllSO(programmeId)
	if err != nil {
		return nil, err
	}

	var sos = make(map[string][]string)
	for _, so := range SOs {
		sos[so.Code] = []string{}
		for _, sso := range so.SubStudentOutcomes {
			sos[so.Code] = append(sos[so.Code], sso.Code)
		}
	}

	program := &entity.ProgrammeLinkedSO{
		AllCourse: []string{},
		AllSOs:    sos,
	}

	courseMap := make(map[string]*entity.CourseLinkedSO)
	for _, res := range result {
		if _, exists := courseMap[res.CourseCode]; !exists {
			courseMap[res.CourseCode] = &entity.CourseLinkedSO{
				CourseCode: res.CourseCode,
				CourseName: res.CourseName,
				SOs:        make(map[string][]string),
			}
			program.AllCourse = append(program.AllCourse, res.CourseCode)
		}
		course := courseMap[res.CourseCode]

		// Add SO and SSO
		if res.SOCode != "" && res.SSOCode != "" {
			if _, exists := course.SOs[res.SOCode]; !exists {
				course.SOs[res.SOCode] = []string{}
			}
			if !contains(course.SOs[res.SOCode], res.SSOCode) {
				course.SOs[res.SOCode] = append(course.SOs[res.SOCode], res.SSOCode)
			}
		}
	}

	var courses []entity.CourseLinkedSO
	for _, course := range courseMap {
		courses = append(courses, *course)
	}

	program.CourseLinkedSOs = courses

	return program, nil
}

func (r programmeRepositoryGorm) GetAllPLO(programmeId string) ([]entity.ProgramLearningOutcome, error) {
	var plos []entity.ProgramLearningOutcome

	err := r.gorm.Raw(`
	SELECT
		plo.*
	FROM
		programme_plo pplo
	LEFT JOIN program_learning_outcome plo ON
		plo.id = pplo.program_learning_outcome_id
	WHERE
		pplo.programme_id = ?;
	`, programmeId).Scan(&plos).Error
	if err != nil {
		return nil, err
	}

	var ploIds []string
	for _, plo := range plos {
		ploIds = append(ploIds, plo.Id)
	}

	var splos []entity.SubProgramLearningOutcome
	err = r.gorm.Where("program_learning_outcome_id IN ?", ploIds).Find(&splos).Error
	if err != nil {
		return nil, err
	}

	for i, plo := range plos {
		for _, splo := range splos {
			if plo.Id == splo.ProgramLearningOutcomeId {
				plos[i].SubProgramLearningOutcomes = append(plos[i].SubProgramLearningOutcomes, splo)
			}
		}
	}

	return plos, nil
}

func (r programmeRepositoryGorm) GetAllSO(programmeId string) ([]entity.StudentOutcome, error) {
	var sos []entity.StudentOutcome

	err := r.gorm.Raw(`
	SELECT
		so.*
	FROM
		programme_so pso
	LEFT JOIN student_outcome so ON
		so.id = pso.student_outcome_id
	WHERE
		pso.programme_id = ?;
	`, programmeId).Scan(&sos).Error
	if err != nil {
		return nil, err
	}

	var soIds []string
	for _, so := range sos {
		soIds = append(soIds, so.Id)
	}

	var ssos []entity.SubStudentOutcome
	err = r.gorm.Where("student_outcome_id IN ?", soIds).Find(&ssos).Error
	if err != nil {
		return nil, err
	}

	for i, plo := range sos {
		for _, splo := range ssos {
			if plo.Id == splo.StudentOutcomeId {
				sos[i].SubStudentOutcomes = append(sos[i].SubStudentOutcomes, splo)
			}
		}
	}

	return sos, nil
}

func (r programmeRepositoryGorm) GetAllPO(programmeId string) ([]entity.ProgramOutcome, error) {
	var pos []entity.ProgramOutcome

	err := r.gorm.Raw(`
	SELECT
		po.*
	FROM
		programme_po ppo
	LEFT JOIN program_outcome po ON
		po.id = ppo.program_outcome_id
	WHERE
		ppo.programme_id = ?;
	`, programmeId).Scan(&pos).Error
	if err != nil {
		return nil, err
	}

	return pos, nil
}

func (r programmeRepositoryGorm) FilterExistedPO(programmeId string, poIds []string) ([]string, error) {
	var existedIds []string

	err := r.gorm.Raw(`
	SELECT
		ppo.program_outcome_id AS id
	FROM
		programme_po ppo
	WHERE
		ppo.programme_id = ? AND
		ppo.program_outcome_id IN ?;
	`, programmeId, poIds).Pluck("id", &existedIds).Error
	if err != nil {
		return nil, err
	}

	return existedIds, nil
}

func (r programmeRepositoryGorm) FilterExistedPLO(programmeId string, ploIds []string) ([]string, error) {
	var existedIds []string

	err := r.gorm.Raw(`
	SELECT
		pplo.program_learning_outcome_id AS id
	FROM
		programme_plo pplo
	WHERE
		pplo.programme_id = ? AND
		pplo.program_learning_outcome_id IN ?;
	`, programmeId, ploIds).Pluck("id", &existedIds).Error
	if err != nil {
		return nil, err
	}

	return existedIds, nil
}

func (r programmeRepositoryGorm) FilterExistedSO(programmeId string, soIds []string) ([]string, error) {
	var existedIds []string

	err := r.gorm.Raw(`
	SELECT
		pso.student_outcome_id AS id
	FROM
		programme_so pso
	WHERE
		pso.programme_id = ? AND
		pso.student_outcome_id IN ?;
	`, programmeId, soIds).Pluck("id", &existedIds).Error
	if err != nil {
		return nil, err
	}

	return existedIds, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
