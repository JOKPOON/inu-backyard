package entity

import "time"

type Feedback struct {
	Id        string    `json:"id" gorm:"type:char(255);primaryKey"`
	CourseId  string    `json:"course_id" gorm:"type:char(255)"`
	StudentId string    `json:"student_id" gorm:"type:char(255)"`
	Comments  string    `json:"comments" gorm:"type:text"`
	Rating    int       `json:"rating" gorm:"check:rating >= 1 AND rating <= 5"`
	Date      time.Time `json:"date"`

	Course  Course  `json:"course" gorm:"foreignKey:CourseId"`
	Student Student `json:"student" gorm:"foreignKey:StudentId"`
}
