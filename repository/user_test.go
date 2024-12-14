package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(mysql:3306)/inu_backyard_test?parseTime=true"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
	err = db.AutoMigrate(
		&entity.AssignmentGroup{},
		&entity.Assignment{},
		&entity.CourseLearningOutcome{},
		&entity.CourseStream{},
		&entity.Course{},
		&entity.Department{},
		&entity.Enrollment{},
		&entity.Faculty{},
		&entity.Grade{},
		&entity.GraduatedStudent{},
		&entity.User{},
		&entity.Prediction{},
		&entity.ProgramLearningOutcome{},
		&entity.ProgramOutcome{},
		&entity.Programme{},
		&entity.Score{},
		&entity.Semester{},
		&entity.Session{},
		&entity.Student{},
		&entity.SubProgramLearningOutcome{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	tx := db.Begin() // Start a new transaction
	t.Cleanup(func() {
		tx.Rollback() // Roll back the transaction after the test ends
	})
	return tx
}

func teardownTestDB(t *testing.T, db *gorm.DB) {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		t.Fatalf("Failed to get tables: %v", err)
	}

	for _, table := range tables {
		err := db.Migrator().DropTable(table)
		if err != nil {
			t.Fatalf("Failed to drop table %s: %v", table, err)
		}
	}
}

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	repo := NewUserRepositoryGorm(db)

	t.Run("TestCreateUser", func(t *testing.T) {
		user := &entity.User{
			Id:          "1",
			Email:       "test@example.com",
			FirstNameTH: "Test",
			LastNameTH:  "User",
		}

		err := repo.Create(user)
		assert.Nil(t, err, "Expected no error while creating user, got %v", err)

		var result entity.User
		db.First(&result, "id = ?", user.Id)

		assert.Equal(t, user.Email, result.Email, "Expected user email to match")
		assert.Equal(t, user.FirstNameTH, result.FirstNameTH, "Expected user first name to match")
		assert.Equal(t, user.LastNameTH, result.LastNameTH, "Expected user last name to match")
	})

	t.Run("TestGetById", func(t *testing.T) {
		user := &entity.User{
			Id:          "2",
			Email:       "test2@example.com",
			FirstNameTH: "Test2",
			LastNameTH:  "User",
		}
		db.Create(user)

		foundUser, err := repo.GetById("2")
		assert.Nil(t, err, "Expected no error while getting user by id, got %v", err)
		assert.NotNil(t, foundUser, "Expected to find a user")
		assert.Equal(t, user.Email, foundUser.Email, "Expected email to match")
	})

	t.Run("TestGetById_NotFound", func(t *testing.T) {
		foundUser, err := repo.GetById("3")
		assert.Nil(t, err, "Expected no error even when user is not found, got %v", err)
		assert.Nil(t, foundUser, "Expected foundUser to be nil when user does not exist")
	})

	t.Run("TestUpdateUser", func(t *testing.T) {
		user := &entity.User{
			Id:          "4",
			Email:       "update@example.com",
			FirstNameTH: "Update",
			LastNameTH:  "User",
		}
		db.Create(user)

		updatedUser := &entity.User{
			Email:       "updated@example.com",
			FirstNameTH: "Updated",
			LastNameTH:  "User",
		}
		err := repo.Update("4", updatedUser)
		assert.Nil(t, err, "Expected no error while updating user, got %v", err)

		var result entity.User
		db.First(&result, "id = ?", "4")

		assert.Equal(t, "updated@example.com", result.Email, "Expected updated email to match")
		assert.Equal(t, "Updated", result.FirstNameTH, "Expected updated first name to match")
		assert.Equal(t, "User", result.LastNameTH, "Expected updated last name to match")
	})

	t.Run("TestDeleteUser", func(t *testing.T) {
		user := &entity.User{
			Id:          "5",
			Email:       "delete@example.com",
			FirstNameTH: "Delete",
			LastNameTH:  "User",
		}
		db.Create(user)

		err := repo.Delete("5")
		assert.Nil(t, err, "Expected no error while deleting user, got %v", err)

		var result entity.User
		err = db.First(&result, "id = ?", "5").Error

		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Expected user to be deleted from the database")
	})
}
