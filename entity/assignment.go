package entity

type AssignmentRepository interface {
	GetAll() ([]Assignment, error)
	GetAllGroup() ([]AssignmentGroup, error)
	GetById(id string) (*Assignment, error)
	GetByCourseId(courseId string) ([]Assignment, error)
	GetByGroupId(groupId string) ([]Assignment, error)
	GetPassingStudentPercentage(assignmentId string) (float64, error)
	Create(assignment *Assignment) error
	CreateMany(assignment []Assignment) error
	Update(id string, assignment *Assignment) error
	Delete(id string) error

	CreateLinkCourseLearningOutcome(assignmentId string, courseLearningOutcomeId []string) error
	DeleteLinkCourseLearningOutcome(assignmentId string, courseLearningOutcomeId string) error

	GetGroupByQuery(query AssignmentGroup, withAssignment bool) ([]AssignmentGroup, error)
	GetGroupByGroupId(assignmentGroupId string) (*AssignmentGroup, error)
	CreateGroup(assignmentGroup *AssignmentGroup) error
	UpdateGroup(assignmentGroupId string, assignmentGroup *AssignmentGroup) error
	DeleteGroup(assignmentGroupId string) error
}

type AssignmentUseCase interface {
	GetAll() ([]Assignment, error)
	GetAllGroup() ([]AssignmentGroup, error)
	GetById(id string) (*Assignment, error)
	GetByCourseId(courseId string) ([]Assignment, error)
	GetByGroupId(assignmentGroupId string) ([]Assignment, error)
	GetPassingStudentPercentage(assignmentId string) (float64, error)
	Create(CreateAssignmentPayload) error
	Update(id string, payload UpdateAssignmentPayload) error
	Delete(id string) error

	CreateLinkCourseLearningOutcome(assignmentId string, courseLearningOutcomeId []string) error
	DeleteLinkCourseLearningOutcome(assignmentId string, courseLearningOutcomeId string) error

	GetGroupByGroupId(assignmentGroupId string) (*AssignmentGroup, error)
	GetGroupByCourseId(courseId string, withAssignment bool) ([]AssignmentGroup, error)
	CreateGroup(CreateAssignmentGroupPayload) error
	UpdateGroup(assignmentGroupId string, payload UpdateAssignmentGroupPayload) error
	DeleteGroup(assignmentGroupId string) error
}

type Assignment struct {
	Id                               string                   `json:"id" gorm:"primaryKey;type:char(255)"`
	Name                             string                   `json:"name"`
	Description                      string                   `json:"description"`
	MaxScore                         int                      `json:"max_score"`
	ExpectedScorePercentage          float64                  `json:"expected_score_percentage"`
	ExpectedPassingStudentPercentage float64                  `json:"expected_passing_student_percentage"`
	IsIncludedInClo                  *bool                    `json:"is_included_in_clo"`
	AssignmentGroupId                string                   `json:"assignment_group_id" gorm:"not null"`
	CourseId                         string                   `json:"course_id" gorm:"->;-:migration"`
	CourseLearningOutcomes           []*CourseLearningOutcome `json:"course_learning_outcomes" gorm:"many2many:clo_assignment"`
}

type AssignmentGroup struct {
	Id       string `json:"id" gorm:"primaryKey;type:char(255)"`
	Name     string `json:"name"`
	CourseId string `json:"course_id"`
	Weight   int    `json:"weight"`

	Assignments []Assignment `gorm:"foreignKey:AssignmentGroupId" json:"assignments,omitempty"`

	Course *Course `json:",omitempty"`
}

type CreateAssignmentGroupPayload struct {
	Name     string `json:"name" validate:"required"`
	Weight   int    `json:"weight" validate:"required"`
	CourseId string `json:"course_id" validate:"required"`
}

type UpdateAssignmentGroupPayload struct {
	Name   string `json:"name" validate:"required"`
	Weight int    `json:"weight" validate:"required"`
}

type CreateAssignmentPayload struct {
	Name                             string   `json:"name" validate:"required"`
	Description                      string   `json:"description"`
	MaxScore                         *int     `json:"max_score" validate:"required"`
	ExpectedScorePercentage          *float64 `json:"expected_score_percentage" validate:"required"`
	ExpectedPassingStudentPercentage *float64 `json:"expected_passing_student_percentage" validate:"required"`
	CourseLearningOutcomeIds         []string `json:"course_learning_outcome_ids" validate:"required"`
	IsIncludedInClo                  *bool    `json:"is_included_in_clo" validate:"required"`
	AssignmentGroupId                string   `json:"assignment_group_id" validate:"required"`
}

type GetAssignmentsByParamsPayload struct {
	CourseLearningOutcomeId string `json:"course_learning_outcome_id"`
}

type GetAssignmentsByCourseIdPayload struct {
	CourseId string `json:"course_id"`
}

type CreateBulkAssignmentsPayload struct {
	Assignments []CreateAssignmentPayload
}

type UpdateAssignmentPayload struct {
	Name                             string   `json:"name"`
	Description                      string   `json:"description"`
	MaxScore                         *int     `json:"maxScore"`
	ExpectedPassingStudentPercentage *float64 `json:"expected_passing_student_percentage"`
	ExpectedScorePercentage          *float64 `json:"expected_score_percentage"`
	IsIncludedInClo                  *bool    `json:"is_included_in_clo"`
}

type CreateLinkCourseLearningOutcomePayload struct {
	CourseLearningOutcomeIds []string `json:"course_learning_outcome_ids" validate:"required"`
}

// func GenerateGroupByAssignmentId(assignmentGroups []AssignmentGroup, assignments []Assignment) map[string]*AssignmentGroup {
// 	weightByAssignmentGroupId := make(map[string]*AssignmentGroup, len(assignmentGroups))
// 	for _, assignmentGroup := range assignmentGroups {
// 		weightByAssignmentGroupId[assignmentGroup.Id] = &assignmentGroup
// 	}

// 	weightByAssignmentId := make(map[string]*AssignmentGroup, len(assignments))
// 	for _, assignment := range assignments {
// 		assignmentGroup, ok := weightByAssignmentGroupId[assignment.AssignmentGroupId]
// 		if !ok {
// 			continue
// 		}

// 		weightByAssignmentId[assignment.Id] = assignmentGroup
// 	}

// 	return weightByAssignmentId
// }
