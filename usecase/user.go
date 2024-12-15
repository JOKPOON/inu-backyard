package usecase

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/middleware"
	"github.com/team-inu/inu-backyard/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo entity.UserRepository
	// courseUseCase entity.CourseUseCase
	// scoreUseCase  entity.ScoreUseCase
}

func NewUserUseCase(userRepo entity.UserRepository) entity.UserUseCase {
	// func NewUserUseCase(userRepo entity.UserRepository, courseUseCase entity.CourseUseCase, scoreUseCase entity.ScoreUseCase) entity.UserUseCase {
	// return &userUseCase{userRepo: userRepo, courseUseCase: courseUseCase, scoreUseCase: scoreUseCase}
	return &userUseCase{userRepo: userRepo}
}

func (u userUseCase) GetAll(pageIndex string, pageSize string) (*entity.Pagination, error) {
	offset, limit, err := utils.ValidatePagination(pageIndex, pageSize)
	if err != nil {
		return nil, errs.New(errs.ErrQueryUser, "cannot get all users", err)
	}

	users, err := u.userRepo.GetAll(offset, limit)
	if err != nil {
		return nil, errs.New(errs.ErrQueryUser, "cannot get all users", err)
	}

	return users, nil
}

func (u userUseCase) GetByEmail(email string) (*entity.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errs.New(errs.ErrQueryUser, "cannot get user by email %s", email, err)
	}

	return user, nil
}

func (u userUseCase) GetById(id string) (*entity.User, error) {
	user, err := u.userRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.ErrQueryUser, "cannot get user by id %s", id, err)
	}

	return user, nil
}

func (u userUseCase) GetBySessionId(sessionId string) (*entity.User, error) {
	user, err := u.userRepo.GetBySessionId(sessionId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryUser, "cannot get user by session id %s", sessionId, err)
	}

	return user, nil
}

func (u userUseCase) GetByParams(params *entity.User, limit int, offset int) ([]entity.User, error) {
	users, err := u.userRepo.GetByParams(params, limit, offset)

	if err != nil {
		return nil, errs.New(errs.ErrQueryUser, "cannot get users by params", err)
	}

	return users, nil
}

func (u userUseCase) Create(payload entity.CreateUserPayload) error {
	if payload.Password == "" {
		payload.Password = uuid.New().String()
	}

	hashPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return errs.New(errs.ErrCreateUser, "cannot create user", err)
	}

	user := &entity.User{
		Id:                 ulid.Make().String(),
		Email:              payload.Email,
		FirstNameTH:        payload.FirstNameTH,
		LastNameTH:         payload.LastNameTH,
		FirstNameEN:        payload.FirstNameEN,
		LastNameEN:         payload.LastNameEN,
		AcademicPositionTH: payload.AcademicPositionTH,
		AcademicPositionEN: payload.AcademicPositionEN,
		Password:           hashPassword,
		Role:               payload.Role,
		Degree:             payload.Degree,
	}
	fmt.Println(user.Role)
	if !user.IsRoles(entity.Roles) {
		return errs.New(errs.ErrCreateUser, "cannot create user", fmt.Errorf("role %s is not valid", user.Role))
	}

	err = u.userRepo.Create(user)
	if err != nil {
		return errs.New(errs.ErrCreateUser, "cannot create user", err)
	}

	return nil
}

func (u userUseCase) CreateMany(users []entity.User) error {
	//encrypt password
	for i := range users {
		bcryptPassword, err := bcrypt.GenerateFromPassword([]byte((users)[i].Password), bcrypt.DefaultCost)
		if err != nil {
			return errs.New(errs.ErrCreateUser, "cannot create user", err)
		}
		users[i].Id = ulid.Make().String()
		(users)[i].Password = string(bcryptPassword)
	}

	err := u.userRepo.CreateMany(users)
	if err != nil {
		return errs.New(errs.ErrCreateUser, "cannot create user", err)
	}

	return nil
}

func (u userUseCase) Update(id string, user *entity.User) error {
	existedUser, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get user id %s to update", id, err)
	} else if existedUser == nil {
		return errs.New(errs.ErrUserNotFound, "cannot get user id %s to update", id)
	}

	err = u.userRepo.Update(id, user)
	if err != nil {
		return errs.New(errs.ErrUpdateUser, "cannot update user by id %s", user.Id, err)
	}

	return nil
}

func (u userUseCase) Delete(id string) error {
	user, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get user id %s to delete", id, err)
	} else if user == nil {
		return errs.New(errs.ErrUserNotFound, "cannot get user id %s to delete", id)
	}

	// courses, err := u.courseUseCase.GetByUserId(id)
	// if err != nil {
	// 	return errs.New(errs.SameCode, "cannot get courses related to this user", err)
	// } else if len(courses) > 0 {
	// 	return errs.New(errs.ErrUserNotFound, "courses related to this user still exist", courses[0].Id)
	// }

	// scores, err := u.scoreUseCase.GetByUserId(id)
	// if err != nil {
	// 	return errs.New(errs.SameCode, "cannot get scores related to this user", err)
	// } else if len(scores) > 0 {
	// 	return errs.New(errs.ErrUserNotFound, "scores related to this user still exist", scores[0].Id)
	// }

	err = u.userRepo.Delete(id)

	if err != nil {
		return errs.New(errs.ErrDeleteUser, "cannot delete user by id %s", id, err)
	}

	return nil
}

func (r userUseCase) CheckUserRole(ctx *fiber.Ctx, userId string, role entity.UserRole) error {
	user := middleware.GetUserFromCtx(ctx)

	if user == nil {
		return errs.New(errs.ErrUserNotFound, "cannot get user from context")
	}

	fmt.Println(user.Role)

	if !user.IsRoles([]entity.UserRole{role}) {
		return errs.New(errs.ErrUserNotFound, "user id %s is not the same as the user id in the context", userId)
	}

	return nil
}

func (r userUseCase) GetUserFromCtx(ctx *fiber.Ctx) (*entity.User, error) {
	user := middleware.GetUserFromCtx(ctx)
	if user == nil {
		return &entity.User{}, errs.New(errs.ErrUserNotFound, "cannot get user from context")
	}

	return user, nil
}
