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
