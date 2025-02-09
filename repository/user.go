package repository

import (
	"fmt"
	"math"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

type userRepositoryGorm struct {
	gorm *gorm.DB
}

func NewUserRepositoryGorm(gorm *gorm.DB) entity.UserRepository {
	return &userRepositoryGorm{gorm: gorm}
}

func (r userRepositoryGorm) GetAll(query string, offset int, limit int) (*entity.Pagination, error) {
	var users []entity.User
	var pagination entity.Pagination
	var total int64

	queryBuilder := r.gorm.Model(&entity.User{})

	if query != "" {
		queryBuilder = queryBuilder.Where("first_name_th LIKE ? OR last_name_th LIKE ? OR first_name_en LIKE ? OR last_name_en LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	// Get total count for pagination
	if err := queryBuilder.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("cannot count users: %w", err)
	}

	// Fetch paginated data
	if err := queryBuilder.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("cannot query users: %w", err)
	}

	// Set pagination data
	pagination.Total = total
	pagination.Size = limit
	pagination.Page = offset/limit + 1
	pagination.TotalPage = int(math.Ceil(float64(total) / float64(limit)))
	pagination.Data = users

	return &pagination, nil
}

func (r userRepositoryGorm) GetBySessionId(sessionId string) (*entity.User, error) {
	var user *entity.User

	err := r.gorm.Joins("JOIN session ON session.user_id = user.id").Where("session.id = ?", sessionId).Find(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get user by session id: %w", err)
	}

	return user, nil
}

func (r userRepositoryGorm) GetById(id string) (*entity.User, error) {
	var user *entity.User

	err := r.gorm.Where("id = ?", id).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get user by id: %w", err)
	}

	return user, nil
}

func (r userRepositoryGorm) GetByEmail(email string) (*entity.User, error) {
	var user *entity.User

	err := r.gorm.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get user by email: %w", err)
	}

	return user, nil
}

func (r userRepositoryGorm) GetByParams(params *entity.User, limit int, offset int) ([]entity.User, error) {
	var users []entity.User

	err := r.gorm.Where(params).Limit(limit).Offset(offset).Find(&users).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get users by params: %w", err)
	}

	return users, nil
}

func (r userRepositoryGorm) Create(user *entity.User) error {
	err := r.gorm.Create(&user).Error
	if err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}

	return nil
}

func (r userRepositoryGorm) CreateMany(users []entity.User) error {
	err := r.gorm.Create(&users).Error
	if err != nil {
		return fmt.Errorf("cannot create users: %w", err)
	}

	return nil
}

func (r userRepositoryGorm) Update(id string, user *entity.User) error {
	err := r.gorm.Model(&entity.User{}).Where("id = ?", id).Updates(user).Error
	if err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}

	return nil
}

func (r userRepositoryGorm) Delete(id string) error {
	err := r.gorm.Delete(&entity.User{Id: id}).Error

	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	return nil
}
