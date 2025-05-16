package entity

import "gorm.io/datatypes"

type CoursePortfolioRepository interface {
	EvaluatePassingAssignmentPercentage(courseId string) ([]AssignmentPercentage, error)
	EvaluatePassingPoPercentage(courseId string) ([]PoPercentage, error)
	EvaluatePassingCloPercentage(courseId string) ([]CloPercentage, error)
	EvaluatePassingCloStudents(courseId string) ([]CloPassingStudentGorm, error)
	EvaluatePassingPloStudents(courseId string) ([]PloPassingStudentGorm, error)
	EvaluatePassingPoStudents(courseId string) ([]PoPassingStudentGorm, error)
	EvaluateAllPloCourses() ([]PloCoursesGorm, error)
	EvaluateAllPoCourses() ([]PoCoursesGorm, error)
	EvaluateProgramLearningOutcomesByStudentId(studentId string) ([]StudentPlosGorm, error)
	EvaluateProgramOutcomesByStudentId(studentId string) ([]StudentPosGorm, error)
	GetCourseCloAssessment(programmeId string, fromSerm, toSerm int) ([]FlatRow, error)
	GetCourseLinkedOutcomes(programmeId string, fromSerm, toSerm int) ([]FlatRow, error)
	GetCourseOutcomesSuccessRate(programmeId string, fromSerm, toSerm int) ([]CourseOutcomeSuccessRate, error)
	GetCourseOutcomes(courseId string) (*CoursePortfolioOutcome, error)

	UpdateCoursePortfolio(courseId string, data datatypes.JSON) error
}

type CoursePortfolioUseCase interface {
	Generate(courseId string) (*CoursePortfolio, error)
	CalculateGradeDistribution(courseId string) (*GradeDistribution, error)
	EvaluateTabeeOutcomes(courseId string) ([]TabeeOutcome, error)
	GetCloPassingStudentsByCourseId(courseId string) ([]CloPassingStudent, error)
	GetStudentOutcomesStatusByCourseId(courseId string) ([]StudentOutcomeStatus, error)
	GetAllProgramLearningOutcomeCourses() ([]PloCourses, error)
	GetAllProgramOutcomeCourses() ([]PoCourses, error)
	GetOutcomesByStudentId(studentId string) ([]StudentOutcomes, error)
	GetCourseCloAssessment(programmeId string, fromSerm, toSerm int) (*FileResponse, error)
	GetCourseLinkedOutcomes(programmeId string, fromSerm, toSerm int) (*FileResponse, error)
	GetCourseOutcomesSuccessRate(programmeId string, fromSerm, toSerm int) (*FileResponse, error)
	GetCourseOutcomes(courseId string) (*CoursePortfolioOutcome, error)

	UpdateCoursePortfolio(courseId string, implement Implementation, educationOutcomes EducationOutcome, continuous ContinuousDevelopment) error
}

// [1] Info
type CourseInfo struct {
	Name        string   `json:"course_name"`
	Code        string   `json:"course_code"`
	LecturersTH []string `json:"lecturers_th"`
	LecturersEN []string `json:"lecturers_en"`
	Programme   string   `json:"programme"`
}

// [2] Summary
type CourseSummary struct {
	TeachingMethods []string `json:"teaching_methods"`
	OnlineTools     string   `json:"online_tools"`
	Objectives      []string `json:"objectives"`
}

// [3.1] Tabee Outcome
type Assessment struct {
	AssessmentTask        string  `json:"assessment_task"`
	PassingCriteria       float64 `json:"passing_criteria"`
	StudentPassPercentage float64 `json:"student_pass_percentage"`
}

type CourseOutcome struct {
	Name                                string       `json:"name"`
	Code                                string       `json:"code"`
	ExpectedPassingAssignmentPercentage float64      `json:"expected_passing_assignment_percentage"`
	PassingCloPercentage                float64      `json:"passing_clo_percentage"`
	Assessments                         []Assessment `json:"assessments"`
}

type TabeeOutcome struct {
	Name                  string          `json:"name"`
	Code                  string          `json:"code"`
	CourseOutcomes        []CourseOutcome `json:"course_outcomes"`
	MinimumPercentage     float64         `json:"minimum_percentage"`
	ExpectedCloPercentage float64         `json:"expected_clo_percentage"`
	Plos                  []NestedOutcome `json:"plos"`
}

// [3.2] Grade Distribution
type GradeFrequency struct {
	Name       string  `json:"name"`
	GradeScore float64 `json:"grade_score"`
	Frequency  int     `json:"frequency"`
}

type ScoreFrequency struct {
	Score     string `json:"score"`
	Frequency int    `json:"frequency"`
}

type Statistics struct {
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	Mode   float64 `json:"mode"`
	SD     float64 `json:"sd"`
}

type GradeDistribution struct {
	StudentAmount    int              `json:"student_amount"`
	GPA              float64          `json:"gpa"`
	GradeFrequencies []GradeFrequency `json:"grade_frequencies"`
	ScoreFrequencies []ScoreFrequency `json:"score_frequencies"`
	Statistics       Statistics       `json:"statistics"`
}

type Outcome struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type NestedOutcome struct {
	Code   string    `json:"code"`
	Name   string    `json:"name"`
	Nested []Outcome `json:"nested,omitempty"`
}

// [3] Result
type CourseResult struct {
	Plos []NestedOutcome `json:"plos"`
	Clos []Outcome       `json:"clos"`
	Pos  []Outcome       `json:"pos"`

	TabeeOutcomes     []TabeeOutcome    `json:"tabee_outcomes"`
	GradeDistribution GradeDistribution `json:"grade_distribution"`
}

// [4.1] SubjectComments
type Subject struct {
	CourseName string `json:"course_name"`
	Comment    string `json:"comments"`
}

type SubjectComments struct {
	UpstreamSubjects   []Subject `json:"upstream_subjects"`
	DownstreamSubjects []Subject `json:"downstream_subjects"`
	Other              string    `json:"other"`
}

// [4] Development
type CourseDevelopment struct {
	Plans           []string        `json:"plans"`
	DoAndChecks     []string        `json:"do_and_checks"`
	Acts            []string        `json:"acts"`
	SubjectComments SubjectComments `json:"subject_comments"`
	OtherComment    string          `json:"other_comment"`
}

// Course Portfolio
type CoursePortfolio struct {
	CourseInfo        CourseInfo        `json:"info"`
	CourseSummary     CourseSummary     `json:"summary"`
	CourseResult      CourseResult      `json:"result"`
	CourseDevelopment CourseDevelopment `json:"development"`
	Raw               datatypes.JSON    `json:"raw"`
}

type AssignmentPercentage struct {
	AssignmentId            string `gorm:"column:a_id"`
	Name                    string
	ExpectedScorePercentage float64
	PassingPercentage       float64
	CourseLearningOutcomeId string `gorm:"column:c_id"`
}

type PoPercentage struct {
	PassingPercentage float64
	ProgramOutcomeId  string `gorm:"column:p_id"`
}

type CloPercentage struct {
	PassingPercentage       float64
	CourseLearningOutcomeId string `gorm:"column:clo_id"`
}

type StudentData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	StudentId string `json:"student_id"`
	Pass      bool   `json:"pass"`
}

type CloPassingStudent struct {
	CourseLearningOutcomeId string        `json:"course_learning_outcome_id"`
	Students                []StudentData `json:"students"`
}

type CloPassingStudentGorm struct {
	FirstName               string
	LastName                string
	StudentId               string
	Pass                    bool
	CourseLearningOutcomeId string `gorm:"column:clo_id"`
	Code                    string
	Description             string
}

type CloData struct {
	Pass                    bool   `json:"pass"`
	CourseLearningOutcomeId string `json:"course_learning_outcome_id"`
	Code                    string `json:"code"`
	Description             string `json:"description"`
}

type PloData struct {
	Id              string `json:"id"`
	Code            string `json:"code"`
	DescriptionThai string `json:"description_thai"`
	ProgramYear     int    `json:"program_year"`
	Pass            bool   `json:"pass"`
}

type PloPassingStudentGorm struct {
	Code                     string
	DescriptionThai          string
	ProgramYear              int
	StudentId                string
	Pass                     bool
	ProgramLearningOutcomeId string `gorm:"column:plo_id"`
}

type PoData struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Pass bool   `json:"pass"`
}

type PoPassingStudentGorm struct {
	Code             string
	Name             string
	StudentId        string
	Pass             bool
	ProgramOutcomeId string `gorm:"column:p_id"`
}

type StudentOutcomeStatus struct {
	StudentId               string    `json:"student_id"`
	ProgramLearningOutcomes []PloData `json:"program_learning_outcomes"`
	ProgramOutcomes         []PoData  `json:"program_outcomes"`
	CourseLearningOutcomes  []CloData `json:"course_learning_outcomes"`
}

type CourseData struct {
	Id                string  `json:"id"`
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	PassingPercentage float64 `json:"passing_percentage"`
	Year              int     `json:"year"`
	SemesterSequence  string  `json:"semester_sequence"`
}

type PloCourses struct {
	ProgramLearningOutcomeId string       `json:"program_learning_outcome_id"`
	Courses                  []CourseData `json:"courses"`
}

type PloCoursesGorm struct {
	PassingPercentage        float64
	ProgramLearningOutcomeId string `gorm:"column:plo_id"`
	CourseId                 string
	Name                     string
	Code                     string
	Year                     int
	SemesterSequence         string
}

type PoCourses struct {
	ProgramOutcomeId string       `json:"program_outcome_id"`
	Courses          []CourseData `json:"courses"`
}

type PoCoursesGorm struct {
	PassingPercentage float64
	ProgramOutcomeId  string `gorm:"column:p_id"`
	CourseId          string
	Name              string
	Code              string
	Year              int
	SemesterSequence  string
}

type StudentCourseData struct {
	Id               string `json:"id"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	Pass             bool   `json:"pass"`
	Year             int    `json:"year"`
	SemesterSequence string `json:"semester_sequence"`
}

type StudentPloData struct {
	ProgramLearningOutcomeId string              `json:"program_learning_outcome_id"`
	Code                     string              `json:"code"`
	DescriptionThai          string              `json:"description_thai"`
	ProgramYear              int                 `json:"program_year"`
	Courses                  []StudentCourseData `json:"courses"`
}

type StudentPoData struct {
	ProgramOutcomeId string              `json:"program_outcome_id"`
	Code             string              `json:"code"`
	Name             string              `json:"name"`
	Courses          []StudentCourseData `json:"courses"`
}

type StudentOutcomes struct {
	StudentId               string           `json:"student_id"`
	ProgramLearningOutcomes []StudentPloData `json:"program_learning_outcomes"`
	ProgramOutcomes         []StudentPoData  `json:"program_outcomes"`
}

type StudentPlosGorm struct {
	StudentId                  string
	ProgramLearningOutcomeId   string `gorm:"column:plo_id"`
	ProgramLearningOutcomeCode string `gorm:"column:plo_code"`
	DescriptionThai            string
	ProgramYear                int
	CourseId                   string
	CourseCode                 string
	CourseName                 string
	Pass                       bool
	Year                       int
	SemesterSequence           string
}

type StudentPosGorm struct {
	StudentId          string
	ProgramOutcomeId   string `gorm:"column:p_id"`
	ProgramOutcomeCode string `gorm:"column:po_code"`
	ProgramOutcomeName string `gorm:"column:po_name"`
	CourseId           string
	CourseCode         string
	CourseName         string
	Pass               bool
	Year               int
	SemesterSequence   string
}

// ///

// type NameObject struct {
// 	Name string `json:"name"`
// }

// type CourseSummmaryForm struct {
// 	TeachingMethods []NameObject `json:"teaching_method"`
// 	OnlineTool      string       `json:"online_tool"`
// 	Objectives      []NameObject `json:"objectives"`
// }

// type CourseDevelopmentForm struct {
// 	Plans           []NameObject    `json:"plans"`
// 	DoAndChecks     []NameObject    `json:"do_and_checks"`
// 	Acts            []NameObject    `json:"acts"`
// 	SubjectComments SubjectComments `json:"subject_comments"`
// 	OtherComment    string          `json:"other_comment"`
// }

// type PortfolioData struct {
// 	Summary     CourseSummary     `json:"summary"`
// 	Development CourseDevelopment `json:"development"`
// }

type SaveCoursePortfolioPayload struct {
	Implementation        Implementation        `json:"implementation" validate:"required"`
	ContinuousDevelopment ContinuousDevelopment `json:"continuous_development" validate:"required"`
	EducationOutcomes     EducationOutcome      `json:"education_outcomes" validate:"required"`
}

type FlatRow struct {
	CourseCode     string `json:"course_code"`
	Semester       string `json:"semester"` // Assuming this is a string, adjust if it's an int
	CourseName     string `json:"course_name"`
	CloID          string `json:"clo_id"`
	CloDescription string `json:"clo_description"`
	AssessmentID   string `json:"assessment_id"`
	AssessmentName string `json:"assessment_name"`
	PLOCode        string `json:"plo_code"`
	SPLOCode       string `json:"splo_code"`
	SOCode         string `json:"so_code"`
	SSOCode        string `json:"sso_code"`
	POCode         string `json:"po_code"`
}

type FileResponse struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
	FileUrl  string `json:"file_url"`
}

type PortfolioData struct {
	Implementation        Implementation        `json:"implementation"`
	EducationOutcomes     EducationOutcome      `json:"education_outcomes"`
	ContinuousDevelopment ContinuousDevelopment `json:"continuous_development"`
}

type Implementation struct {
	Methods            []string `json:"methods"`
	OnlineMedia        []string `json:"online_media"`
	TeachingObjectives []string `json:"teaching_objectives"`
}

type ContinuousDevelopment struct {
	Plans       []string  `json:"plans"`
	DoAndChecks []string  `json:"do_and_checks"`
	Acts        []string  `json:"acts"`
	Feedbacks   Feedbacks `json:"feedbacks"`
}

type Feedbacks struct {
	UpstreamFeedbacks   []FeedbackJson `json:"upstream_subjects"`
	DownstreamFeedbacks []FeedbackJson `json:"downstream_subjects"`
}

type FeedbackJson struct {
	CourseCode string `json:"course_code"`
	CourseName string `json:"course_name"`
	Comment    string `json:"comment"`
	Sender     string `json:"sender"`
}

type EducationOutcome struct {
	GradeDistribution GradeDistribution `json:"grade_distribution"`
}

type AssignmentPassingRate struct {
	AssignmentID                        string
	AssignmentName                      string
	PassedPercentage                    float64
	ExpectedPassingAssignmentPercentage float64
}

type CloPassingRate struct {
	CLOID                               string
	CLOCode                             string
	PassedPercentage                    float64
	ExpectedPassingAssignmentPercentage float64
	Assignments                         map[string]AssignmentPassingRate
}

type PoPassingRate struct {
	POID                         string
	POCode                       string
	PassedPercentage             float64
	ExpectedPassingCloPercentage float64
	CLOPassingRate               map[string]CloPassingRate
}

type SploPassingRate struct {
	SPLOID                       string
	SPLOCode                     string
	PLOID                        string
	PLOCode                      string
	PassedPercentage             float64
	ExpectedPassingCloPercentage float64
	CLOPassingRate               map[string]CloPassingRate
}

type SsoPassingRate struct {
	SOID                         string
	SOCode                       string
	SSOID                        string
	SSOCode                      string
	PassedPercentage             float64
	ExpectedPassingCloPercentage float64
	CLOPassingRate               map[string]CloPassingRate
}

type CourseOutcomeSuccessRate struct {
	CourseId       string
	CourseCode     string
	CourseName     string
	CourseSemester string
	CLOs           map[string]float64
	PLOs           map[string]map[string]float64
	SOs            map[string]map[string]float64
	POs            map[string]float64
}

type CoursePortfolioOutcome struct {
	CLOs map[string]CloPassingRate  `json:"clos"`
	PLOs map[string]SploPassingRate `json:"plos"`
	POs  map[string]PoPassingRate   `json:"pos"`
	SOs  map[string]SsoPassingRate  `json:"sos"`
}
