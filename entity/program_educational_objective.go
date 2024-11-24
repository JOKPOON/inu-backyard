package entity

type ProgramEducationalObjective struct {
	Id            string `json:"id" gorm:"type:char(255);primaryKey"`
	Description   string `json:"description" gorm:"type:text;not null"`
	ProgrammeName string `json:"programme_name"  gorm:"type:char(255)"`

	Programme Programme `gorm:"foreignKey:ProgrammeName"`
}
