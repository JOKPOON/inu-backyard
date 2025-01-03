package usecase

import (
	"encoding/json"

	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
)

type courseUseCase struct {
	courseRepo      entity.CourseRepository
	semesterUseCase entity.SemesterUseCase
	userUseCase     entity.UserUseCase
}

func NewCourseUseCase(courseRepo entity.CourseRepository, semesterUseCase entity.SemesterUseCase, userUseCase entity.UserUseCase) entity.CourseUseCase {
	return &courseUseCase{courseRepo: courseRepo, semesterUseCase: semesterUseCase, userUseCase: userUseCase}
}

func (u courseUseCase) GetAll() ([]entity.Course, error) {
	courses, err := u.courseRepo.GetAll()
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot get all courses", err)
	}

	return courses, nil
}

func (u courseUseCase) GetById(id string) (*entity.Course, error) {
	course, err := u.courseRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot get course by id %s", id, err)
	}

	return course, nil
}

func (u courseUseCase) GetByUserId(userId string) ([]entity.Course, error) {
	user, err := u.userUseCase.GetById(userId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get user id %s while get scores", user, err)
	} else if user == nil {
		return nil, errs.New(errs.ErrQueryCourse, "user id %s not found while getting scores", userId, err)
	}

	course, err := u.courseRepo.GetByUserId(userId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryCourse, "cannot get score by user id %s", userId, err)
	}

	return course, nil
}

func (u courseUseCase) Create(user entity.User, payload entity.CreateCoursePayload) error {
	if !user.IsRoles([]entity.UserRole{entity.UserRoleHeadOfCurriculum}) {
		return errs.New(errs.ErrCreateCourse, "no permission to create course")
	}

	semester, err := u.semesterUseCase.GetById(payload.SemesterId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get semester id %s while creating course", payload.SemesterId, err)
	} else if semester == nil {
		return errs.New(errs.ErrSemesterNotFound, "semester id %s not found while creating course", payload.SemesterId)
	}

	lecturer, err := u.userUseCase.GetById(payload.UserId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get user id %s while creating course", payload.UserId, err)
	} else if lecturer == nil {
		return errs.New(errs.ErrUserNotFound, "user id %s not found while creating course", payload.UserId)
	}

	if !payload.CriteriaGrade.IsValid() {
		return errs.New(errs.ErrCreateCourse, "invalid criteria grade")
	}

	emptyJson, _ := json.Marshal(map[string]string{})
	course := entity.Course{
		Id:                           ulid.Make().String(),
		Name:                         payload.Name,
		Code:                         payload.Code,
		Curriculum:                   payload.Curriculum,
		Description:                  payload.Description,
		ExpectedPassingCloPercentage: payload.ExpectedPassingCloPercentage,
		AcademicYear:                 payload.AcademicYear,
		GraduateYear:                 payload.GraduateYear,
		ProgramYear:                  payload.ProgramYear,
		UserId:                       payload.UserId,
		SemesterId:                   payload.SemesterId,
		CriteriaGrade:                payload.CriteriaGrade,
		PortfolioData:                emptyJson,
	}

	err = u.courseRepo.Create(&course)
	if err != nil {
		return errs.New(errs.ErrCreateCourse, "cannot create course", err)
	}

	return nil
}

func (u courseUseCase) Update(user entity.User, id string, payload entity.UpdateCoursePayload) error {
	existCourse, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get course id %s to update", id, err)
	} else if existCourse == nil {
		return errs.New(errs.ErrCourseNotFound, "cannot get course id %s to update", id)
	}

	if !user.IsRoles([]entity.UserRole{entity.UserRoleHeadOfCurriculum}) && user.Id != existCourse.UserId {
		return errs.New(errs.ErrCreateCourse, "No permission to edit this course")
	}

	if !payload.CriteriaGrade.IsValid() {
		return errs.New(errs.ErrCreateCourse, "invalid criteria grade")
	}

	err = u.courseRepo.Update(id, &entity.Course{
		Name:                         payload.Name,
		Code:                         payload.Code,
		Curriculum:                   payload.Curriculum,
		Description:                  payload.Description,
		CriteriaGrade:                payload.CriteriaGrade,
		ExpectedPassingCloPercentage: payload.ExpectedPassingCloPercentage,
		AcademicYear:                 payload.AcademicYear,
		GraduateYear:                 payload.GraduateYear,
		ProgramYear:                  payload.ProgramYear,
		IsPortfolioCompleted:         payload.IsPortfolioCompleted,
	})
	if err != nil {
		return errs.New(errs.ErrUpdateCourse, "cannot update course by id %s", id, err)
	}

	return nil
}

func (u courseUseCase) Delete(user entity.User, id string) error {
	if !user.IsRoles([]entity.UserRole{entity.UserRoleHeadOfCurriculum}) {
		return errs.New(errs.ErrCreateCourse, "no permission to create course")
	}

	err := u.courseRepo.Delete(id)
	if err != nil {
		return errs.New(errs.ErrDeleteCourse, "cannot delete course", err)
	}

	return nil
}
