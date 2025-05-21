package usecase

import (
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils/slice"
)

type scoreUseCase struct {
	scoreRepo         entity.ScoreRepository
	enrollmentUseCase entity.EnrollmentUseCase
	assignmentUseCase entity.AssignmentUseCase
	courseUseCase     entity.CourseUseCase
	userUseCase       entity.UserUseCase
	studentUseCase    entity.StudentUseCase
}

func NewScoreUseCase(
	scoreRepo entity.ScoreRepository,
	enrollmentUseCase entity.EnrollmentUseCase,
	assignmentUseCase entity.AssignmentUseCase,
	courseUseCase entity.CourseUseCase,
	userUseCase entity.UserUseCase,
	studentUsecase entity.StudentUseCase,
) entity.ScoreUseCase {
	return &scoreUseCase{
		scoreRepo:         scoreRepo,
		enrollmentUseCase: enrollmentUseCase,
		assignmentUseCase: assignmentUseCase,
		courseUseCase:     courseUseCase,
		userUseCase:       userUseCase,
		studentUseCase:    studentUsecase,
	}
}

func (u scoreUseCase) GetAll() ([]entity.Score, error) {
	scores, err := u.scoreRepo.GetAll()
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get all scores", err)
	}

	return scores, nil
}

func (u scoreUseCase) GetById(id string) (*entity.Score, error) {
	score, err := u.scoreRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get score by id", err)
	}

	return score, nil
}

func (u scoreUseCase) GetByAssignmentId(assignmentId string, courseId string) (*entity.AssignmentScore, error) {
	assignment, err := u.assignmentUseCase.GetById(assignmentId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get assignment when finding score", err)
	} else if assignment == nil {
		return nil, errs.New(errs.ErrQueryScore, "assignment id %s not found while finding score", assignmentId)
	}

	assignmentGroup, err := u.assignmentUseCase.GetGroupByGroupId(assignment.AssignmentGroupId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get assignment group when finding score", err)
	}

	scores, err := u.scoreRepo.GetByAssignmentId(assignmentId, courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get all scores", err)
	}

	enrollments, err := u.enrollmentUseCase.GetByCourseId(assignmentGroup.CourseId, "")
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get enrollments when finding score", err)
	}

	clos, err := u.assignmentUseCase.GetLinkedCLOs(assignmentId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get linked CLOs when finding score", err)
	}

	assignmentClos := []entity.AssessmentClos{}
	for _, clo := range clos {
		assignmentClos = append(assignmentClos, entity.AssessmentClos{
			Id:          clo.Id,
			Code:        clo.Code,
			Description: clo.DescriptionTH,
		})
	}

	assignmentScore := &entity.AssignmentScore{
		CLOs:                   assignmentClos,
		IsIncludedInClo:        *assignment.IsIncludedInClo,
		MaxScore:               assignment.MaxScore,
		Description:            assignment.Description,
		ExpectedPassingStudent: assignment.ExpectedPassingStudentPercentage,
		ExpectedPassingScore:   assignment.ExpectedScorePercentage,
		Scores:                 scores,
		EnrolledAmount:         len(enrollments),
		SubmittedAmount:        len(scores),
	}

	return assignmentScore, nil
}

func (u scoreUseCase) GetByUserId(userId string) ([]entity.Score, error) {
	user, err := u.userUseCase.GetById(userId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get user id %s while get scores", user, err)
	} else if user == nil {
		return nil, errs.New(errs.ErrQueryScore, "user id %s not found while getting scores", userId, err)
	}

	scores, err := u.scoreRepo.GetByUserId(userId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryScore, "cannot get score by user id %s", userId, err)
	}

	return scores, nil
}

func (u scoreUseCase) GetByStudentId(studentId string) ([]entity.Score, error) {
	student, err := u.studentUseCase.GetById(studentId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get student id %s while get scores", student, err)
	} else if student == nil {
		return nil, errs.New(errs.ErrQueryScore, "student id %s not found while getting scores", studentId, err)
	}

	scores, err := u.scoreRepo.GetByStudentId(studentId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryScore, "cannot get score by student id %s", studentId, err)
	}

	return scores, nil
}

func (u scoreUseCase) CreateMany(userId string, assignmentId string, studentScores []entity.StudentScore) error {
	if len(studentScores) == 0 {
		return errs.New(errs.ErrCreateScore, "studentScores must not be empty")
	}

	user, err := u.userUseCase.GetById(userId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get user id %s to create score", userId, err)
	} else if user == nil {
		return errs.New(errs.ErrUserNotFound, "cannot get user id %s to create score", userId)
	}

	assignment, err := u.assignmentUseCase.GetById(assignmentId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get assignment id %s to create score", assignmentId, err)
	} else if assignment == nil {
		return errs.New(errs.ErrAssignmentNotFound, "cannot get assignment id %s to create score", assignmentId)
	}

	assignmentGroup, err := u.assignmentUseCase.GetGroupByGroupId(assignment.AssignmentGroupId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get assignment group id %s to create score", assignment.AssignmentGroupId, err)
	}

	course, err := u.courseUseCase.GetById(assignmentGroup.CourseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get course id %s to create score", assignment.CourseId, err)
	} else if course == nil {
		return errs.New(errs.ErrCourseNotFound, "cannot get course id %s to create score", assignment.CourseId)
	}

	for _, lecturer := range course.Lecturers {
		if user.IsRoles([]entity.UserRole{entity.UserRoleLecturer}) && user.Id != lecturer.Id {
			return errs.New(errs.ErrDeleteScore, "no permission to create score")
		}
	}

	for _, studentScore := range studentScores {
		if *studentScore.Score > float64(assignment.MaxScore) {
			return errs.New(errs.ErrCreateScore, "score %f of student id %s is more than max score of assignment (score: %d)", studentScore.Score, studentScore.StudentId, assignment.MaxScore)
		}
	}

	studentIds := []string{}
	for _, studentScore := range studentScores {
		studentIds = append(studentIds, studentScore.StudentId)
	}

	withStatus := entity.EnrollmentStatusEnroll
	joinedStudentIds, err := u.enrollmentUseCase.FilterJoinedStudent(studentIds, assignmentGroup.CourseId, &withStatus)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get existed student ids while creating score")
	}

	nonJoinedStudentIds := slice.Subtraction(studentIds, joinedStudentIds)
	if len(nonJoinedStudentIds) > 0 {
		return errs.New(errs.ErrCreateAssignment, "there are non joined student ids %v", nonJoinedStudentIds)
	}

	submittedScoreStudentIds, err := u.FilterSubmittedScoreStudents(assignmentId, studentIds)
	if err != nil {
		return errs.New(errs.SameCode, "cannot filter submitted score student while creating score")
	} else if len(submittedScoreStudentIds) != 0 {
		return errs.New(errs.ErrCreateAssignment, "there are already submitted score students %v", submittedScoreStudentIds)
	}

	scores := []entity.Score{}
	for _, studentScore := range studentScores {
		scores = append(scores, entity.Score{
			Id:           ulid.Make().String(),
			Score:        *studentScore.Score,
			StudentId:    studentScore.StudentId,
			UserId:       userId,
			AssignmentId: assignmentId,
		})
	}

	err = u.scoreRepo.CreateMany(scores)
	if err != nil {
		return errs.New(errs.ErrCreateScore, "cannot create score", err)
	}

	return nil
}

func (u scoreUseCase) Update(user entity.User, scoreId string, score float64) error {
	existScore, err := u.GetById(scoreId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get score by id %s ", scoreId, err)
	} else if existScore == nil {
		return errs.New(errs.ErrScoreNotFound, "score not found", err)
	}

	if user.IsRoles([]entity.UserRole{entity.UserRoleLecturer}) && user.Id != existScore.UserId {
		return errs.New(errs.ErrUpdateScore, "no permission to update score")
	}

	err = u.scoreRepo.Update(scoreId, &entity.Score{
		Score:        score,
		StudentId:    existScore.StudentId,
		UserId:       existScore.UserId,
		AssignmentId: existScore.AssignmentId,
	})
	if err != nil {
		return errs.New(errs.ErrUpdateScore, "cannot update score", err)
	}

	return nil
}

func (u scoreUseCase) Delete(user entity.User, id string) error {
	existScore, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get score by id %s ", id, err)
	} else if existScore == nil {
		return errs.New(errs.ErrScoreNotFound, "score not found to delete")
	}

	if user.IsRoles([]entity.UserRole{entity.UserRoleLecturer}) && user.Id != existScore.UserId {
		return errs.New(errs.ErrDeleteScore, "no permission to delete score")
	}

	err = u.scoreRepo.Delete(id)
	if err != nil {
		return errs.New(errs.ErrDeleteScore, "cannot delete score by id %s", id, err)
	}
	return nil
}

func (u scoreUseCase) FilterSubmittedScoreStudents(assignmentId string, studentIds []string) ([]string, error) {
	submittedScoreStudentIds, err := u.scoreRepo.FilterSubmittedScoreStudents(assignmentId, studentIds)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot query students", err)
	}

	return submittedScoreStudentIds, nil
}
