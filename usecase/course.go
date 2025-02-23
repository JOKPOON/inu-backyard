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

func NewCourseUseCase(
	courseRepo entity.CourseRepository,
	semesterUseCase entity.SemesterUseCase,
	userUseCase entity.UserUseCase,
) entity.CourseUseCase {
	return &courseUseCase{
		courseRepo:      courseRepo,
		semesterUseCase: semesterUseCase,
		userUseCase:     userUseCase,
	}
}

func (u courseUseCase) GetAll(query string, year string, program string) (*entity.GetAllCourseResponse, error) {
	courses, err := u.courseRepo.GetAll(query, year, program)
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot get all courses", err)
	} else if len(courses) == 0 {
		return nil, errs.New(errs.ErrCourseNotFound, "no course found")
	}

	res := entity.GetAllCourseResponse{}
	for _, c := range courses {
		var lec []entity.Lecturer
		for _, l := range c.Lecturers {
			lec = append(lec, entity.Lecturer{
				Id:     l.Id,
				NameTH: l.TitleTHShort + l.FirstNameTH + " " + l.LastNameTH,
				NameEN: l.TitleENShort + l.FirstNameEN + " " + l.LastNameEN,
			})
		}

		res.Courses = append(res.Courses, entity.CourseSimpleData{
			Id:           c.Id,
			Name:         c.Name,
			Code:         c.Code,
			Credit:       c.Credit,
			AcademicYear: c.AcademicYear,
			GraduateYear: c.GraduateYear,
			Description:  c.Description,
			Program: entity.Program{
				Id:     c.Programme.Id,
				NameTH: c.Programme.NameTH,
				NameEN: c.Programme.NameEN,
			},
			Lecturers: lec,
		})
	}

	return &res, nil
}

func (u courseUseCase) GetById(id string) (*entity.Course, error) {
	course, err := u.courseRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot get course by id %s", id, err)
	}

	return course, nil
}

func (u courseUseCase) GetByUserId(userId string, query string, year string, program string) (*entity.GetAllCourseResponse, error) {
	user, err := u.userUseCase.GetById(userId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get user id %s while get scores", user, err)
	} else if user == nil {
		return nil, errs.New(errs.ErrQueryCourse, "user id %s not found while getting scores", userId, err)
	}

	courses, err := u.courseRepo.GetByUserId(userId, query, year, program)
	if err != nil {
		return nil, errs.New(errs.ErrQueryCourse, "cannot get score by user id %s", userId, err)
	}

	res := entity.GetAllCourseResponse{}
	for _, c := range courses {
		var lec []entity.Lecturer
		for _, l := range c.Lecturers {
			lec = append(lec, entity.Lecturer{
				Id:     l.Id,
				NameTH: l.TitleTHShort + l.FirstNameTH + " " + l.LastNameTH,
				NameEN: l.TitleENShort + l.FirstNameEN + " " + l.LastNameEN,
			})
		}

		res.Courses = append(res.Courses, entity.CourseSimpleData{
			Id:           c.Id,
			Name:         c.Name,
			Code:         c.Code,
			Credit:       c.Credit,
			AcademicYear: c.AcademicYear,
			GraduateYear: c.GraduateYear,
			Description:  c.Description,
			Program: entity.Program{
				Id:     c.Programme.Id,
				NameTH: c.Programme.NameTH,
				NameEN: c.Programme.NameEN,
			},
			Lecturers: lec,
		})
	}

	return &res, nil
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

	for _, lecturerId := range payload.LecturerIds {
		lecturer, err := u.userUseCase.GetById(lecturerId)
		if err != nil {
			return errs.New(errs.SameCode, "cannot get user id %s while creating course", lecturerId, err)
		} else if lecturer == nil {
			return errs.New(errs.ErrUserNotFound, "user id %s not found while creating course", lecturerId)
		}
	}

	if !payload.CriteriaGrade.IsValid() {
		return errs.New(errs.ErrCreateCourse, "invalid criteria grade")
	}

	emptyJson, _ := json.Marshal(map[string]string{})
	course := entity.Course{
		Id:                           ulid.Make().String(),
		Name:                         payload.Name,
		Code:                         payload.Code,
		ProgrammeId:                  payload.ProgrammeId,
		Description:                  payload.Description,
		Credit:                       payload.Credit,
		AcademicYear:                 payload.AcademicYear,
		GraduateYear:                 payload.GraduateYear,
		ExpectedPassingCloPercentage: payload.ExpectedPassingCloPercentage,
		SemesterId:                   payload.SemesterId,
		CriteriaGrade:                payload.CriteriaGrade,
		PortfolioData:                emptyJson,
	}

	err = u.courseRepo.Create(&course)
	if err != nil {
		return errs.New(errs.ErrCreateCourse, "cannot create course", err)
	}

	err = u.courseRepo.CreateLinkWithLecturer(course.Id, payload.LecturerIds)
	if err != nil {
		return errs.New(errs.ErrCreateCourse, "cannot create link with lecturer", err)
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

	for _, lecturerId := range payload.LecturerIds {
		if !user.IsRoles([]entity.UserRole{entity.UserRoleHeadOfCurriculum}) && user.Id != lecturerId {
			return errs.New(errs.ErrCreateCourse, "No permission to edit this course")
		}
	}

	err = u.courseRepo.ReplaceLecturersForCourse(id, payload.LecturerIds)
	if err != nil {
		return errs.New(errs.ErrCreateCourse, "cannot create link with lecturer", err)
	}

	if !payload.CriteriaGrade.IsValid() {
		return errs.New(errs.ErrCreateCourse, "invalid criteria grade")
	}

	err = u.courseRepo.Update(id, &entity.Course{
		Name:                         payload.Name,
		Code:                         payload.Code,
		ProgrammeId:                  payload.ProgrammeId,
		Description:                  payload.Description,
		Credit:                       payload.Credit,
		AcademicYear:                 payload.AcademicYear,
		GraduateYear:                 payload.GraduateYear,
		CriteriaGrade:                payload.CriteriaGrade,
		ExpectedPassingCloPercentage: payload.ExpectedPassingCloPercentage,
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

func (u courseUseCase) GetStudentsPassingCLOs(courseId string) (*entity.StudentPassCLOResp, error) {
	resp, err := u.courseRepo.GetStudentsPassingCLOs(courseId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryCourse, "cannot get students passing CLOs by course id %s", courseId, err)
	}

	return resp, nil
}
