package repository

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

type enrollmentRepositoryGorm struct {
	gorm *gorm.DB
}

func NewEnrollmentRepositoryGorm(gorm *gorm.DB) entity.EnrollmentRepository {
	return &enrollmentRepositoryGorm{gorm: gorm}
}

func (r enrollmentRepositoryGorm) GetAll() ([]entity.Enrollment, error) {
	var enrollments []entity.Enrollment

	err := r.gorm.
		Model(&enrollments).
		Select("enrollment.*, student.first_name, student.last_name, student.email").
		Joins("LEFT JOIN student on student.id = enrollment.student_id").
		Scan(&enrollments).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get enrollments: %w", err)
	}

	return enrollments, nil
}

func (r enrollmentRepositoryGorm) GetById(id string) (*entity.Enrollment, error) {
	var enrollments *entity.Enrollment

	err := r.gorm.
		First(&enrollments, "id = ?", id).
		Select("enrollment.*, student.first_name, student.last_name, student.email").
		Joins("LEFT JOIN student on student.id = enrollment.student_id").
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get enrollment by id: %w", err)
	}

	return enrollments, nil
}

func (r enrollmentRepositoryGorm) GetByCourseId(courseId string) ([]entity.Enrollment, error) {
	var enrollments []entity.Enrollment
	err := r.gorm.
		Model(&enrollments).
		Select("enrollment.*, student.first_name, student.last_name, student.email").
		Joins("LEFT JOIN student on student.id = enrollment.student_id").
		Where("enrollment.course_id = ?", courseId).
		Scan(&enrollments).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get enrollment by id: %w", err)
	}

	return enrollments, nil
}

func (r enrollmentRepositoryGorm) GetByStudentId(studentId string) ([]entity.Enrollment, error) {
	var enrollments []entity.Enrollment
	err := r.gorm.Where("student_id = ?", studentId).Find(&enrollments).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get enrollments by student id: %w", err)
	}

	return enrollments, nil
}

func (r enrollmentRepositoryGorm) CreateMany(enrollments []entity.Enrollment) error {
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)
	return r.gorm.Create(&enrollments).Error
}

func (r enrollmentRepositoryGorm) Create(enrollment *entity.Enrollment) error {
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)
	return r.gorm.Create(&enrollment).Error
}

func (r enrollmentRepositoryGorm) Update(id string, enrollment *entity.Enrollment) error {
	//find old enrollment by name
	var oldEnrollment *entity.Enrollment
	err := r.gorm.Where("id = ?", id).First(&oldEnrollment).Error
	if err != nil {
		return fmt.Errorf("cannot get enrollment while updating enrollment: %w", err)
	}

	//update old enrollment with new name
	err = r.gorm.Model(&oldEnrollment).Updates(enrollment).Error
	if err != nil {
		return fmt.Errorf("cannot update enrollment by id: %w", err)
	}
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r enrollmentRepositoryGorm) Delete(id string) error {
	err := r.gorm.Where("id = ?", id).Delete(&entity.Enrollment{}).Error
	if err != nil {
		return fmt.Errorf("cannot delete enrollment by id: %w", err)
	}
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r enrollmentRepositoryGorm) FilterExisted(ids []string) ([]string, error) {
	var existedIds []string

	err := r.gorm.Raw("SELECT id FROM `enrollment` WHERE id in ?", ids).Scan(&existedIds).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query enrollments: %w", err)
	}

	return existedIds, nil
}

func (r enrollmentRepositoryGorm) FilterJoinedStudent(studentIds []string, courseId string, status *entity.EnrollmentStatus) ([]string, error) {
	// fmt.Println(*status)
	var existedIds []string

	query := "SELECT student_id FROM `enrollment` WHERE course_id = ? AND student_id in ?"
	args := []interface{}{courseId, studentIds}

	if status != nil {
		query += " AND status = ?"
		args = append(args, *status)
	}

	err := r.gorm.Raw(query, args...).Scan(&existedIds).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query student: %w", err)
	}

	return existedIds, nil
}
