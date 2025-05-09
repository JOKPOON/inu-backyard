package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

type studentRepositoryGorm struct {
	gorm *gorm.DB
}

func NewStudentRepositoryGorm(gorm *gorm.DB) entity.StudentRepository {
	return &studentRepositoryGorm{gorm: gorm}
}

func (r studentRepositoryGorm) GetById(id string) (*entity.Student, error) {
	var student *entity.Student

	err := r.gorm.Where("id = ?", id).First(&student).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get student by id: %w", err)
	}

	return student, nil
}

func (r studentRepositoryGorm) GetAll() ([]entity.Student, error) {
	var students []entity.Student
	return students, nil
}

func (r studentRepositoryGorm) GetByParams(query string, params *entity.Student, limit, offset int) (*entity.GetStudentsResponse, error) {
	var students []entity.Student
	var total int64
	var TotalPage int64

	db := r.gorm.Model(&entity.Student{})

	// Apply filters if params is provided
	if params != nil {
		if params.ProgrammeId != "" {
			db = db.Where("programme_id = ?", params.ProgrammeId)
		}
		if params.Year != "" {
			db = db.Where("year = ?", params.Year)
		}
		if params.DepartmentId != "" {
			db = db.Where("department_id = ?", params.DepartmentId) // fixed here
		}
	}

	// Apply search query
	if query != "" {
		search := "%" + query + "%"
		db = db.Where(`
			first_name_th LIKE ? OR
			last_name_th LIKE ? OR
			first_name_en LIKE ? OR
			last_name_en LIKE ? OR
			id LIKE ?`,
			search, search, search, search, search,
		)
	}
	// Count total records
	err := db.Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("cannot count students: %w", err)
	}

	// Calculate total pages
	TotalPage = total / int64(limit)
	if total%int64(limit) != 0 {
		TotalPage++
	}

	// Preload related Programme and Department data
	err = db.Preload("Programme").Preload("Department").
		Limit(limit).Offset(offset).
		Find(&students).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot query students: %w", err)
	}

	return &entity.GetStudentsResponse{
		Students:  students,
		Total:     int(total),
		Page:      offset + 1,
		TotalPage: int(TotalPage),
	}, nil
}

func (r studentRepositoryGorm) Create(student *entity.Student) error {
	err := r.gorm.Create(&student).Error
	if err != nil {
		return fmt.Errorf("cannot create student: %w", err)
	}

	return nil
}

func (r studentRepositoryGorm) CreateMany(students []entity.Student) error {
	err := r.gorm.Create(&students).Error
	if err != nil {
		return fmt.Errorf("cannot create student: %w", err)
	}

	return nil
}

func (r studentRepositoryGorm) Update(id string, student *entity.Student) error {
	err := r.gorm.Model(&entity.Student{}).Where("id = ?", id).Updates(map[string]interface{}{
		"id":              student.Id,
		"first_name":      student.FirstNameTH,
		"last_name":       student.LastNameTH,
		"email":           student.Email,
		"year":            student.Year,
		"programme_id":    student.ProgrammeId,
		"department_name": student.DepartmentId,
	}).Error
	if err != nil {
		return fmt.Errorf("cannot update student: %w", err)
	}

	return nil
}

func (r studentRepositoryGorm) Delete(id string) error {
	err := r.gorm.Delete(&entity.Student{Id: id}).Error

	if err != nil {
		return fmt.Errorf("cannot delete student: %w", err)
	}

	return nil
}

func (r studentRepositoryGorm) FilterExisted(studentIds []string) ([]string, error) {
	var existedIds []string

	err := r.gorm.Raw("SELECT id FROM `student` WHERE id in ?", studentIds).Scan(&existedIds).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query student: %w", err)
	}

	return existedIds, nil
}

func (r studentRepositoryGorm) GetAllSchools() ([]string, error) {
	var schools []sql.NullString

	err := r.gorm.Raw("SELECT DISTINCT school FROM student").Scan(&schools).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query student: %w", err)
	}

	nonNullSchool := make([]string, 0)
	for _, school := range schools {
		if school.Valid {
			nonNullSchool = append(nonNullSchool, school.String)
		}
	}

	return nonNullSchool, nil
}
func (r studentRepositoryGorm) GetAllAdmissions() ([]string, error) {
	var admissions []sql.NullString

	err := r.gorm.Raw("SELECT DISTINCT admission FROM student").Scan(&admissions).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query student: %w", err)
	}

	nonNullAdmission := make([]string, 0)
	for _, admission := range admissions {
		if admission.Valid {
			nonNullAdmission = append(nonNullAdmission, admission.String)
		}
	}

	return nonNullAdmission, nil
}
