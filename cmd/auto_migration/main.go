package main

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/database"
)

func main() {
	// for development purpose only
	gormDB, err := database.NewGorm(&database.GormConfig{
		User:         "root",
		Password:     "root",
		Host:         "localhost",
		DatabaseName: "inu_2",
		Port:         "3306",
	})
	if err != nil {
		panic(err)
	}

	err = gormDB.AutoMigrate(
		&entity.AssignmentGroup{},
		&entity.Assignment{},
		&entity.CourseLearningOutcome{},
		&entity.CourseStream{},
		&entity.Course{},
		&entity.Department{},
		&entity.Enrollment{},
		&entity.Faculty{},
		&entity.Feedback{},
		&entity.Grade{},
		&entity.GraduatedStudent{},
		&entity.User{},
		&entity.Prediction{},
		&entity.ProgramEducationalObjective{},
		&entity.ProgramImprovement{},
		&entity.ProgramLearningOutcome{},
		&entity.ProgramOutcome{},
		&entity.Programme{},
		&entity.Score{},
		&entity.Semester{},
		&entity.Session{},
		&entity.StudentOutcome{},
		&entity.Student{},
		&entity.SubProgramLearningOutcome{},
	)

	fmt.Println(err)
}
