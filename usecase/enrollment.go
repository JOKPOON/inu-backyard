package usecase

import (
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils/slice"
)

type enrollmentUseCase struct {
	enrollmentRepo entity.EnrollmentRepository
	studentUseCase entity.StudentUseCase
	courseUseCase  entity.CourseUseCase
}

func NewEnrollmentUseCase(enrollmentRepo entity.EnrollmentRepository, studentUseCase entity.StudentUseCase, courseUseCase entity.CourseUseCase) entity.EnrollmentUseCase {
	return &enrollmentUseCase{
		enrollmentRepo: enrollmentRepo,
		studentUseCase: studentUseCase,
		courseUseCase:  courseUseCase,
	}
}

func (u enrollmentUseCase) GetAll() ([]entity.Enrollment, error) {
	enrollments, err := u.enrollmentRepo.GetAll()
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get all enrollments", err)
	}

	return enrollments, nil
}

func (u enrollmentUseCase) GetById(id string) (*entity.Enrollment, error) {
	enrollment, err := u.enrollmentRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get enrollment by id %s", id, err)
	}

	return enrollment, nil
}

func (u enrollmentUseCase) GetByCourseId(courseId string) ([]entity.Enrollment, error) {
	course, err := u.courseUseCase.GetById(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course id %s while get enrollments", course, err)
	} else if course == nil {
		return nil, errs.New(errs.ErrEnrollmentNotFound, "course id %s not found while getting enrollments", courseId, err)
	}

	enrollment, err := u.enrollmentRepo.GetByCourseId(courseId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryEnrollment, "cannot get enrollment by course id %s", courseId, err)
	}

	return enrollment, nil
}

func (u enrollmentUseCase) GetByStudentId(studentId string) ([]entity.Enrollment, error) {
	student, err := u.studentUseCase.GetById(studentId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get student id %s while get enrollments", student, err)
	} else if student == nil {
		return nil, errs.New(errs.ErrQueryEnrollment, "student id %s not found while getting enrollments", studentId, err)
	}

	enrollment, err := u.enrollmentRepo.GetByStudentId(studentId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryEnrollment, "cannot get enrollment by student id %s", studentId, err)
	}

	return enrollment, nil
}

func (u enrollmentUseCase) CreateMany(payload entity.CreateEnrollmentsPayload) error {
	course, err := u.courseUseCase.GetById(payload.CourseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get course id %s while creating enrollment", payload.CourseId, err)
	} else if course == nil {
		return errs.New(errs.ErrCourseNotFound, "course id %s not found while creating enrollment", payload.CourseId)
	}

	duplicateStudentIds := slice.GetDuplicateValue(payload.StudentIds)
	if len(duplicateStudentIds) != 0 {
		return errs.New(errs.ErrCreateEnrollment, "duplicate student ids %v", duplicateStudentIds)
	}

	nonExistedStudentIds, err := u.studentUseCase.FilterNonExisted(payload.StudentIds)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get non existed student ids while creating enrollment")
	} else if len(nonExistedStudentIds) != 0 {
		return errs.New(errs.ErrCreateEnrollment, "there are non exist student ids %v", nonExistedStudentIds)
	}

	joinedStudentIds, err := u.FilterJoinedStudent(payload.StudentIds, payload.CourseId, nil)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get existed student ids while creating score")
	} else if len(joinedStudentIds) > 0 {
		return errs.New(errs.ErrCreateAssignment, "there are already joined student ids, %v", joinedStudentIds)
	}

	enrollments := []entity.Enrollment{}

	for _, studentId := range payload.StudentIds {
		enrollment := entity.Enrollment{
			Id:        ulid.Make().String(),
			CourseId:  payload.CourseId,
			Status:    payload.Status,
			StudentId: studentId,
		}

		enrollments = append(enrollments, enrollment)
	}

	err = u.enrollmentRepo.CreateMany(enrollments)
	if err != nil {
		return errs.New(errs.ErrCreateEnrollment, "cannot create enrollment", err)
	}

	return err
}

func (u enrollmentUseCase) Update(id string, status entity.EnrollmentStatus) error {
	existEnrollment, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get enrollment id %s to update", id, err)
	} else if existEnrollment == nil {
		return errs.New(errs.ErrEnrollmentNotFound, "enrollment id %s not found while update enrollment", id)
	}

	err = u.enrollmentRepo.Update(id, &entity.Enrollment{
		Status: status,
	})

	if err != nil {
		return errs.New(errs.ErrUpdateEnrollment, "cannot update enrollment by id %s", id, err)
	}

	return nil
}

func (u enrollmentUseCase) Delete(id string) error {
	enrollment, err := u.enrollmentRepo.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get enrollment id %s to delete", id, err)
	} else if enrollment == nil {
		return errs.New(errs.ErrEnrollmentNotFound, "cannot get enrollment id %s to delete", id)
	}

	err = u.enrollmentRepo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (u enrollmentUseCase) FilterJoinedStudent(studentIds []string, courseId string, status *entity.EnrollmentStatus) ([]string, error) {
	joinedIds, err := u.enrollmentRepo.FilterJoinedStudent(studentIds, courseId, status)
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot query enrollment", err)
	}

	return joinedIds, nil
}
