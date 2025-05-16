package usecase

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	"github.com/team-inu/inu-backyard/utils"
	"github.com/xuri/excelize/v2"
)

// TODO: refactor (real)
type coursePortfolioUseCase struct {
	CoursePortfolioRepository    entity.CoursePortfolioRepository
	CourseUseCase                entity.CourseUseCase
	UserUseCase                  entity.UserUseCase
	EnrollmentUseCase            entity.EnrollmentUseCase
	AssignmentUseCase            entity.AssignmentUseCase
	ScoreUseCase                 entity.ScoreUseCase
	StudentUseCase               entity.StudentUseCase
	CourseLearningOutcomeUseCase entity.CourseLearningOutcomeUseCase
	CourseStreamUseCase          entity.CourseStreamsUseCase
}

func NewCoursePortfolioUseCase(
	coursePortfolioRepository entity.CoursePortfolioRepository,
	courseUseCase entity.CourseUseCase,
	userUseCase entity.UserUseCase,
	enrollmentUseCase entity.EnrollmentUseCase,
	assignmentUseCase entity.AssignmentUseCase,
	scoreUseCase entity.ScoreUseCase,
	studentUsecase entity.StudentUseCase,
	courseLearningOutcomeUseCase entity.CourseLearningOutcomeUseCase,
	courseStreamUseCase entity.CourseStreamsUseCase,
) entity.CoursePortfolioUseCase {
	return &coursePortfolioUseCase{
		CoursePortfolioRepository:    coursePortfolioRepository,
		CourseUseCase:                courseUseCase,
		UserUseCase:                  userUseCase,
		EnrollmentUseCase:            enrollmentUseCase,
		AssignmentUseCase:            assignmentUseCase,
		ScoreUseCase:                 scoreUseCase,
		StudentUseCase:               studentUsecase,
		CourseLearningOutcomeUseCase: courseLearningOutcomeUseCase,
		CourseStreamUseCase:          courseStreamUseCase,
	}
}

func generateOutcome(cloWithPos []entity.CourseLearningOutcomeWithPO) ([]entity.NestedOutcome, []entity.Outcome, []entity.Outcome) {
	addedSubPlo := make(map[string]bool, 0)

	plosByPloId := make(map[string]entity.NestedOutcome, 0)
	closByCloId := make(map[string]entity.Outcome, 0)
	posByPoId := make(map[string]entity.Outcome, 0)

	plos := make([]entity.NestedOutcome, 0)
	clos := make([]entity.Outcome, 0)
	pos := make([]entity.Outcome, 0)

	for _, c := range cloWithPos {
		plosFromMap, found := plosByPloId[c.ProgramLearningOutcomeCode]

		if !found {
			addedSubPlo[c.SubProgramLearningOutcomeCode] = true

			plosByPloId[c.ProgramLearningOutcomeCode] = entity.NestedOutcome{
				Code: c.ProgramLearningOutcomeCode,
				Name: c.ProgramLearningOutcomeName,
				Nested: []entity.Outcome{
					{
						Code: c.SubProgramLearningOutcomeCode,
						Name: c.SubProgramLearningOutcomeName,
					},
				},
			}
		} else {
			if _, isSubPloAdded := addedSubPlo[c.SubProgramLearningOutcomeCode]; !isSubPloAdded {
				addedSubPlo[c.SubProgramLearningOutcomeCode] = true

				plosFromMap.Nested = append(
					plosByPloId[c.ProgramLearningOutcomeCode].Nested,
					entity.Outcome{
						Code: c.SubProgramLearningOutcomeCode,
						Name: c.SubProgramLearningOutcomeName,
					},
				)

				plosByPloId[c.ProgramLearningOutcomeCode] = plosFromMap
			}
		}

		closByCloId[c.Code] = entity.Outcome{
			Code: c.Code,
			Name: c.DescriptionTH,
		}

		posByPoId[c.ProgramOutcomeName] = entity.Outcome{
			Code: c.ProgramOutcomeCode,
			Name: c.ProgramOutcomeName,
		}

	}
	for _, plo := range plosByPloId {
		plos = append(plos, plo)
	}
	for _, clo := range closByCloId {
		clos = append(clos, clo)
	}
	for _, po := range posByPoId {
		pos = append(pos, po)
	}

	return plos, clos, pos
}

func (u coursePortfolioUseCase) Generate(courseId string) (*entity.CoursePortfolio, error) {
	course, err := u.CourseUseCase.GetById(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course id %s while generate course portfolio", courseId, err)
	} else if course == nil {
		return nil, errs.New(errs.ErrCourseNotFound, "course id %s not found while generate course portfolio", courseId)
	}

	lecturersNameTH := make([]string, 0)
	lecturersNameEN := make([]string, 0)
	for _, lecturer := range course.Lecturers {
		if lecturer == nil {
			return nil, errs.New(errs.ErrCourseNotFound, "lecturer not found while generate course portfolio", courseId)
		}

		lecturersNameTH = append(lecturersNameTH, fmt.Sprintf("%s %s", lecturer.FirstNameTH, lecturer.LastNameTH))
		lecturersNameEN = append(lecturersNameEN, fmt.Sprintf("%s %s", lecturer.FirstNameEN, lecturer.LastNameEN))
	}

	courseInfo := entity.CourseInfo{
		Name:        course.Name,
		Code:        course.Code,
		LecturersTH: lecturersNameTH,
		LecturersEN: lecturersNameEN,
		Programme:   course.Programme.NameTH + " " + course.Programme.NameEN,
	}

	// gradeDistribution, err := u.CalculateGradeDistribution(courseId)
	// if err != nil {
	// 	return nil, errs.New(errs.SameCode, "cannot calculate grade distribution while generate course portfolio", err)
	// }

	// tabeeOutcomes, err := u.EvaluateTabeeOutcomes(courseId)
	// if err != nil {
	// 	return nil, errs.New(errs.SameCode, "cannot evaluate tabee outcomes while generate course portfolio", err)
	// }

	// closWithPos, err := u.CourseLearningOutcomeUseCase.GetByCourseId(courseId)
	// if err != nil {
	// 	return nil, errs.New(errs.SameCode, "cannot get clo while evaluate tabee outcome", err)
	// }

	// plos, clos, pos := generateOutcome(closWithPos)

	// courseResult := entity.CourseResult{
	// 	Plos:              plos,
	// 	Clos:              clos,
	// 	Pos:               pos,
	// 	GradeDistribution: *gradeDistribution,
	// 	TabeeOutcomes:     tabeeOutcomes,
	// }

	courseStreams, err := u.CourseStreamUseCase.GetByTargetCourseId(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course streams while generate course portfolio", err)
	}

	upstreamSubject := make([]entity.Subject, 0)
	downStreamSubject := make([]entity.Subject, 0)

	for _, stream := range courseStreams {
		switch stream.StreamType {
		case entity.DownCourseStreamType:
			fmt.Println(stream.FromCourse)
			upstreamSubject = append(upstreamSubject, entity.Subject{
				CourseName: fmt.Sprintf("%s %s", stream.FromCourse.Code, stream.FromCourse.Name),
				Comment:    stream.Comment,
			})

		case entity.UpCourseStreamType:
			downStreamSubject = append(downStreamSubject, entity.Subject{
				CourseName: fmt.Sprintf("%s %s", stream.TargetCourse.Code, stream.TargetCourse.Name),
				Comment:    stream.Comment,
			})

		}
	}

	portfolioData := entity.PortfolioData{}

	err = json.Unmarshal(course.PortfolioData, &portfolioData)
	if err != nil {
		return nil, errs.New(0, "cannot unmarshal data from db")
	}

	courseDevelopment := entity.CourseDevelopment{
		// Plans:       portfolioData.Development.Plans,
		// DoAndChecks: portfolioData.Development.DoAndChecks,
		// Acts:        portfolioData.Development.Acts,
		SubjectComments: entity.SubjectComments{
			UpstreamSubjects:   upstreamSubject,
			DownstreamSubjects: downStreamSubject,
			// Other:              portfolioData.Development.SubjectComments.Other,
		},
		// OtherComment: portfolioData.Development.OtherComment,
	}

	courseSummary := entity.CourseSummary{
		// TeachingMethods: portfolioData.Summary.TeachingMethods,
		// Objectives:      portfolioData.Summary.Objectives,
		// OnlineTools:     portfolioData.Summary.OnlineTools,
	}

	coursePortfolio := &entity.CoursePortfolio{
		CourseInfo: courseInfo,
		//CourseResult:      courseResult,
		CourseSummary:     courseSummary,
		CourseDevelopment: courseDevelopment,
		Raw:               course.PortfolioData,
	}

	return coursePortfolio, nil
}

func (u coursePortfolioUseCase) CalculateGradeDistribution(courseId string) (*entity.GradeDistribution, error) {
	// Retrieve course details
	course, err := u.CourseUseCase.GetById(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course by id %s while calculating grade distribution", courseId, err)
	} else if course == nil {
		return nil, errs.New(errs.ErrCourseNotFound, "course id %s not found while calculating grade distribution", courseId)
	}

	// Retrieve assignment groups for the course
	assignmentGroups, err := u.AssignmentUseCase.GetGroupByCourseId(courseId, "", true)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get assignment group while calculating grade distribution")
	}

	// Maps to store weight and scores
	sumGroupScoreByGroupId := make(map[string]int)
	sumScoreByStudentId := make(map[string]float64)

	for _, group := range assignmentGroups {
		for _, assignment := range group.Assignments {
			sumGroupScoreByGroupId[group.Id] += assignment.MaxScore

			// Get scores for the assignment
			assignmentScores, _ := u.ScoreUseCase.GetByAssignmentId(assignment.Id)
			for _, score := range assignmentScores.Scores {
				sumScoreByStudentId[score.StudentId] += (score.Score * float64(group.Weight)) / float64(assignment.MaxScore)
			}
		}
	}

	// Collect scores for statistics
	studentScores := make([]float64, 0, len(sumScoreByStudentId))
	for _, score := range sumScoreByStudentId {
		studentScores = append(studentScores, score)
	}

	// Initialize statistics
	stat := entity.Statistics{
		Min: 100, Max: 0, Mean: 0, Median: 0, Mode: 0, SD: 0,
	}

	if len(studentScores) > 0 {
		sort.Float64s(studentScores)

		// Min & Max
		stat.Min = studentScores[0]
		stat.Max = studentScores[len(studentScores)-1]

		// Mean
		totalScore := 0.0
		for _, score := range studentScores {
			totalScore += score
		}
		stat.Mean = totalScore / float64(len(studentScores))

		// Median
		mid := len(studentScores) / 2
		if len(studentScores)%2 == 0 {
			stat.Median = (studentScores[mid-1] + studentScores[mid]) / 2
		} else {
			stat.Median = studentScores[mid]
		}

		// Mode
		frequencyMap := make(map[float64]int)
		maxFreq := 0
		for _, score := range studentScores {
			frequencyMap[score]++
			if frequencyMap[score] > maxFreq {
				maxFreq = frequencyMap[score]
				stat.Mode = score
			}
		}

		// Standard Deviation
		var variance float64
		for _, score := range studentScores {
			variance += math.Pow(score-stat.Mean, 2)
		}
		stat.SD = math.Sqrt(variance / float64(len(studentScores)))
	}

	// Score Frequency Distribution
	scoreRanges := []string{"0-50", "51-55", "56-60", "61-65", "66-70", "71-75", "76-80", "81-85", "86-90", "91-95", "96-100"}
	frequencyByScore := make(map[string]int)

	for _, score := range sumScoreByStudentId {
		switch {
		case score <= 50:
			frequencyByScore["0-50"]++
		case score <= 55:
			frequencyByScore["51-55"]++
		case score <= 60:
			frequencyByScore["56-60"]++
		case score <= 65:
			frequencyByScore["61-65"]++
		case score <= 70:
			frequencyByScore["66-70"]++
		case score <= 75:
			frequencyByScore["71-75"]++
		case score <= 80:
			frequencyByScore["76-80"]++
		case score <= 85:
			frequencyByScore["81-85"]++
		case score <= 90:
			frequencyByScore["86-90"]++
		case score <= 95:
			frequencyByScore["91-95"]++
		default:
			frequencyByScore["96-100"]++
		}
	}

	scoreFrequencies := make([]entity.ScoreFrequency, 0, len(frequencyByScore))
	for _, rangeLabel := range scoreRanges {
		scoreFrequencies = append(scoreFrequencies, entity.ScoreFrequency{
			Score:     rangeLabel,
			Frequency: frequencyByScore[rangeLabel],
		})
	}

	// Grade Frequency Distribution
	gradeFrequencies := map[string]int{}
	for _, score := range sumScoreByStudentId {
		switch {
		case score >= course.A:
			gradeFrequencies["A"]++
		case score >= course.BP:
			gradeFrequencies["BP"]++
		case score >= course.B:
			gradeFrequencies["B"]++
		case score >= course.CP:
			gradeFrequencies["CP"]++
		case score >= course.C:
			gradeFrequencies["C"]++
		case score >= course.DP:
			gradeFrequencies["DP"]++
		case score >= course.D:
			gradeFrequencies["D"]++
		default:
			gradeFrequencies["F"]++
		}
	}

	// Convert to structured format
	gradeFrequenciesList := []entity.GradeFrequency{
		{Name: "A", GradeScore: course.A, Frequency: gradeFrequencies["A"]},
		{Name: "BP", GradeScore: course.BP, Frequency: gradeFrequencies["BP"]},
		{Name: "B", GradeScore: course.B, Frequency: gradeFrequencies["B"]},
		{Name: "CP", GradeScore: course.CP, Frequency: gradeFrequencies["CP"]},
		{Name: "C", GradeScore: course.C, Frequency: gradeFrequencies["C"]},
		{Name: "DP", GradeScore: course.DP, Frequency: gradeFrequencies["DP"]},
		{Name: "D", GradeScore: course.D, Frequency: gradeFrequencies["D"]},
		{Name: "F", GradeScore: 0, Frequency: gradeFrequencies["F"]},
	}

	// GPA Calculation
	totalStudentGPA := 0.0
	studentAmount := len(sumScoreByStudentId)

	for grade, frequency := range gradeFrequencies {
		totalStudentGPA += float64(frequency) * course.CriteriaGrade.GradeToGPA(grade)
	}

	gpa := 0.0
	if studentAmount > 0 {
		gpa = totalStudentGPA / float64(studentAmount)
	}

	// Final Grade Distribution Struct
	return &entity.GradeDistribution{
		StudentAmount:    studentAmount,
		ScoreFrequencies: scoreFrequencies,
		GradeFrequencies: gradeFrequenciesList,
		GPA:              gpa,
		Statistics:       stat,
	}, nil
}

func (u coursePortfolioUseCase) EvaluateTabeeOutcomes(courseId string) ([]entity.TabeeOutcome, error) {
	// assignmentPercentages, err := u.CoursePortfolioRepository.EvaluatePassingAssignmentPercentage(courseId)
	// if err != nil {
	// 	return nil, errs.New(errs.SameCode, "cannot evaluate passing assignment percentage by course id %s while evaluate tabee outcome", courseId, err)
	// }

	// assessmentsByCloId := make(map[string][]entity.Assessment, len(assignmentPercentages))
	// for _, assignmentPercentage := range assignmentPercentages {

	// 	cloId := assignmentPercentage.CourseLearningOutcomeId

	// 	assessmentsByCloId[cloId] = append(assessmentsByCloId[cloId], entity.Assessment{
	// 		AssessmentTask:        assignmentPercentage.Name,
	// 		PassingCriteria:       assignmentPercentage.ExpectedScorePercentage,
	// 		StudentPassPercentage: assignmentPercentage.PassingPercentage,
	// 	})
	// }

	// clos, err := u.CourseLearningOutcomeUseCase.GetByCourseId(courseId)
	// if err != nil {
	// 	return nil, errs.New(errs.SameCode, "cannot get clo while evaluate tabee outcome", err)
	// }

	// cloPassingPercentage, err := u.CoursePortfolioRepository.EvaluatePassingCloPercentage(courseId)
	// if err != nil {
	// 	return nil, errs.New(errs.SameCode, "cannot evaluate passing clo percentage", err)
	// }

	// passingCloPercentage := make(map[string]float64, 0)
	// for _, clo := range cloPassingPercentage {
	// 	passingCloPercentage[clo.CourseLearningOutcomeId] = clo.PassingPercentage
	// }

	// courseOutcomeByPoId := make(map[string][]entity.CourseOutcome, 0)
	// expectedPassingCloByPoId := make(map[string]float64, 0)

	// sploDuplicates := make(map[[3]string]entity.Outcome, 0)
	// splosByBothId := make(map[[2]string][]entity.Outcome, 0)
	// plosByBothId := make(map[[2]string]entity.NestedOutcome, 0)
	// for _, clo := range clos {
	// 	courseOutcomeByPoId[clo.ProgramOutcomeId] = append(courseOutcomeByPoId[clo.ProgramOutcomeId], entity.CourseOutcome{
	// 		Name:                                clo.Description,
	// 		Code:                                clo.Code,
	// 		ExpectedPassingAssignmentPercentage: clo.ExpectedPassingAssignmentPercentage,
	// 		PassingCloPercentage:                passingCloPercentage[clo.Id],
	// 		Assessments:                         assessmentsByCloId[clo.Id],
	// 	})
	// 	expectedPassingCloByPoId[clo.ProgramOutcomeId] = clo.ExpectedPassingCloPercentage

	// 	key := [2]string{clo.ProgramOutcomeId, clo.ProgramLearningOutcomeCode}

	// 	_, found := sploDuplicates[[3]string{clo.ProgramOutcomeId, clo.ProgramLearningOutcomeCode, clo.SubProgramLearningOutcomeCode}]

	// 	if !found {
	// 		splo := entity.Outcome{
	// 			Name: clo.SubProgramLearningOutcomeName,
	// 			Code: clo.SubProgramLearningOutcomeCode,
	// 		}
	// 		sploDuplicates[[3]string{clo.ProgramOutcomeId, clo.ProgramLearningOutcomeCode, clo.SubProgramLearningOutcomeCode}] = splo
	// 		splosByBothId[key] = append(splosByBothId[key], splo)
	// 		plosByBothId[key] = entity.NestedOutcome{
	// 			Name: clo.ProgramLearningOutcomeName,
	// 			Code: clo.ProgramLearningOutcomeCode,
	// 		}
	// 	}
	// }

	// plosByPoId := make(map[string][]entity.NestedOutcome, 0)
	// for key := range plosByBothId {
	// 	plo := plosByBothId[key]
	// 	plo.Nested = splosByBothId[key]
	// 	plosByBothId[key] = plo
	// 	plosByPoId[key[0]] = append(plosByPoId[key[0]], plosByBothId[key])
	// }

	// passingPoPercentages, err := u.CoursePortfolioRepository.EvaluatePassingPoPercentage(courseId)
	// if err != nil {
	// 	return nil, errs.New(errs.SameCode, "cannot evaluate passing po percentage by course id %s while evaluate tabee outcome", courseId, err)
	// }

	// passingPoPercentageByPoId := make(map[string]float64, len(passingPoPercentages))
	// for _, passingPoPercentage := range passingPoPercentages {
	// 	passingPoPercentageByPoId[passingPoPercentage.ProgramOutcomeId] = passingPoPercentage.PassingPercentage
	// }

	// tabeeOutcomesByPoId := make(map[string][]entity.TabeeOutcome, 0)
	// for _, clo := range clos {
	// 	checkIsSameOutcomeName := func(foundOutcome []entity.TabeeOutcome, clo entity.CourseLearningOutcomeWithPO) bool {
	// 		isNameSame := false

	// 		for _, tabeeOutcome := range foundOutcome {
	// 			if tabeeOutcome.Name == clo.ProgramOutcomeName {
	// 				isNameSame = true
	// 				break
	// 			}
	// 		}

	// 		return isNameSame
	// 	}

	// 	foundOutcome, found := tabeeOutcomesByPoId[clo.ProgramOutcomeId]
	// 	if !found {
	// 		tabeeOutcomesByPoId[clo.ProgramOutcomeId] = append(tabeeOutcomesByPoId[clo.ProgramOutcomeId], entity.TabeeOutcome{
	// 			Name:                  clo.ProgramOutcomeName,
	// 			Code:                  clo.ProgramLearningOutcomeCode,
	// 			CourseOutcomes:        courseOutcomeByPoId[clo.ProgramOutcomeId],
	// 			MinimumPercentage:     passingPoPercentageByPoId[clo.ProgramOutcomeId],
	// 			ExpectedCloPercentage: expectedPassingCloByPoId[clo.ProgramOutcomeId],
	// 			Plos:                  plosByPoId[clo.ProgramOutcomeId],
	// 		})
	// 		continue
	// 	}

	// 	isNameSame := checkIsSameOutcomeName(foundOutcome, clo)
	// 	if isNameSame {
	// 		continue
	// 	}

	// 	tabeeOutcomesByPoId[clo.ProgramOutcomeId] = append(tabeeOutcomesByPoId[clo.ProgramOutcomeId], entity.TabeeOutcome{
	// 		Name:                  clo.ProgramOutcomeName,
	// 		Code:                  clo.ProgramLearningOutcomeCode,
	// 		CourseOutcomes:        courseOutcomeByPoId[clo.ProgramOutcomeId],
	// 		MinimumPercentage:     passingPoPercentageByPoId[clo.ProgramOutcomeId],
	// 		ExpectedCloPercentage: expectedPassingCloByPoId[clo.ProgramOutcomeId],
	// 		Plos:                  plosByPoId[clo.ProgramOutcomeId],
	// 	})
	// }

	// tabeeOutcomes := make([]entity.TabeeOutcome, 0, len(tabeeOutcomesByPoId))
	// for _, tabeeOutcome := range tabeeOutcomesByPoId {
	// 	tabeeOutcomes = append(tabeeOutcomes, tabeeOutcome...)
	// }

	// return tabeeOutcomes, nil
	return nil, nil
}

func (u coursePortfolioUseCase) GetCloPassingStudentsByCourseId(courseId string) ([]entity.CloPassingStudent, error) {
	course, err := u.CourseUseCase.GetById(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course id %s while getting clo passing students", course, err)
	} else if course == nil {
		return nil, errs.New(errs.ErrCourseNotFound, "course id %s not found while getting clo passing students", courseId, err)
	}

	records, err := u.CoursePortfolioRepository.EvaluatePassingCloStudents(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot evaluate passing clo student by course id %s", courseId, err)
	}

	closMap := make(map[string][]entity.StudentData)

	for _, record := range records {
		closMap[record.CourseLearningOutcomeId] = append(closMap[record.CourseLearningOutcomeId], entity.StudentData{
			FirstName: record.FirstName,
			LastName:  record.LastName,
			StudentId: record.StudentId,
			Pass:      record.Pass,
		})
	}

	clos := make([]entity.CloPassingStudent, 0)

	for cloId := range closMap {
		clos = append(clos, entity.CloPassingStudent{
			CourseLearningOutcomeId: cloId,
			Students:                closMap[cloId],
		})
	}

	return clos, nil
}

func (u coursePortfolioUseCase) GetStudentOutcomesStatusByCourseId(courseId string) ([]entity.StudentOutcomeStatus, error) {
	course, err := u.CourseUseCase.GetById(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course id %s while getting clo passing students", course, err)
	} else if course == nil {
		return nil, errs.New(errs.ErrCourseNotFound, "course id %s not found while getting clo passing students", courseId, err)
	}

	ploRecords, err := u.CoursePortfolioRepository.EvaluatePassingPloStudents(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot evaluate passing plo student by course id %s", courseId, err)
	}

	poRecords, err := u.CoursePortfolioRepository.EvaluatePassingPoStudents(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot evaluate passing po student by course id %s", courseId, err)
	}

	cloRecords, err := u.CoursePortfolioRepository.EvaluatePassingCloStudents(courseId)

	studentPloMap := make(map[string][]entity.PloData)
	studentPoMap := make(map[string][]entity.PoData)
	studentCloMap := make(map[string][]entity.CloData)

	for _, record := range ploRecords {
		studentPloMap[record.StudentId] = append(studentPloMap[record.StudentId], entity.PloData{
			Id:              record.ProgramLearningOutcomeId,
			Code:            record.Code,
			DescriptionThai: record.DescriptionThai,
			ProgramYear:     record.ProgramYear,
			Pass:            record.Pass,
		})
	}

	for _, record := range poRecords {
		studentPoMap[record.StudentId] = append(studentPoMap[record.StudentId], entity.PoData{
			Id:   record.ProgramOutcomeId,
			Code: record.Code,
			Name: record.Name,
			Pass: record.Pass,
		})
	}

	for _, record := range cloRecords {
		studentCloMap[record.StudentId] = append(studentCloMap[record.StudentId], entity.CloData{
			Pass:                    record.Pass,
			CourseLearningOutcomeId: record.CourseLearningOutcomeId,
			Code:                    record.Code,
			Description:             record.Description,
		})
	}

	if len(studentPloMap) != len(studentPoMap) {
		return nil, errs.New(errs.SameCode, "number of students with plo is different from po by course id %s", courseId, err)
	}

	students := make([]entity.StudentOutcomeStatus, 0)

	for studentId := range studentPloMap {
		students = append(students, entity.StudentOutcomeStatus{
			StudentId:               studentId,
			ProgramLearningOutcomes: studentPloMap[studentId],
		})
	}

	for i := range students {
		students[i].ProgramOutcomes = studentPoMap[students[i].StudentId]
		students[i].CourseLearningOutcomes = studentCloMap[students[i].StudentId]
	}

	return students, nil
}

func (u coursePortfolioUseCase) GetAllProgramLearningOutcomeCourses() ([]entity.PloCourses, error) {
	records, err := u.CoursePortfolioRepository.EvaluateAllPloCourses()
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot evaluate plo courses %s", err)
	}

	plosMap := make(map[string][]entity.CourseData)

	for _, record := range records {
		if record.CourseId == "" {
			plosMap[record.ProgramLearningOutcomeId] = append(plosMap[record.ProgramLearningOutcomeId], entity.CourseData{})
		} else {
			plosMap[record.ProgramLearningOutcomeId] = append(plosMap[record.ProgramLearningOutcomeId], entity.CourseData{
				Id:                record.CourseId,
				Code:              record.Code,
				Name:              record.Name,
				PassingPercentage: record.PassingPercentage,
				Year:              record.Year,
				SemesterSequence:  record.SemesterSequence,
			})
		}
	}

	plos := make([]entity.PloCourses, 0)

	for ploId := range plosMap {
		plos = append(plos, entity.PloCourses{
			ProgramLearningOutcomeId: ploId,
			Courses:                  plosMap[ploId],
		})
	}

	return plos, nil
}

func (u coursePortfolioUseCase) GetAllProgramOutcomeCourses() ([]entity.PoCourses, error) {
	records, err := u.CoursePortfolioRepository.EvaluateAllPoCourses()
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot evaluate po courses %s", err)
	}

	posMap := make(map[string][]entity.CourseData)

	for _, record := range records {
		if record.CourseId == "" {
			posMap[record.ProgramOutcomeId] = append(posMap[record.ProgramOutcomeId], entity.CourseData{})
		} else {
			posMap[record.ProgramOutcomeId] = append(posMap[record.ProgramOutcomeId], entity.CourseData{
				Id:                record.CourseId,
				Code:              record.Code,
				Name:              record.Name,
				PassingPercentage: record.PassingPercentage,
				Year:              record.Year,
				SemesterSequence:  record.SemesterSequence,
			})
		}
	}

	pos := make([]entity.PoCourses, 0)

	for poId := range posMap {
		pos = append(pos, entity.PoCourses{
			ProgramOutcomeId: poId,
			Courses:          posMap[poId],
		})
	}

	return pos, nil
}

func (u coursePortfolioUseCase) UpdateCoursePortfolio(courseId string, implement entity.Implementation, educationOutcomes entity.EducationOutcome, continuous entity.ContinuousDevelopment) error {
	portfolioData := &entity.PortfolioData{
		Implementation:        implement,
		EducationOutcomes:     educationOutcomes,
		ContinuousDevelopment: continuous,
	}

	JsonByte, err := json.Marshal(*portfolioData)
	if err != nil {
		return errs.New(errs.SameCode, "cannot marshal course summary %s", err)
	}

	err = u.CoursePortfolioRepository.UpdateCoursePortfolio(courseId, JsonByte)
	if err != nil {
		return errs.New(errs.SameCode, "cannot update course portfolio %s", err)
	}

	return nil
}

func (u coursePortfolioUseCase) GetOutcomesByStudentId(studentId string) ([]entity.StudentOutcomes, error) {
	student, err := u.StudentUseCase.GetById(studentId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get student id %s while getting student outcomes", student, err)
	} else if student == nil {
		return nil, errs.New(errs.ErrStudentNotFound, "student id %s not found while getting student outcomes", studentId, err)
	}

	ploRecords, err := u.CoursePortfolioRepository.EvaluateProgramLearningOutcomesByStudentId(studentId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot evaluate student plos by student id %s", studentId, err)
	}

	poRecords, err := u.CoursePortfolioRepository.EvaluateProgramOutcomesByStudentId(studentId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot evaluate student pos by student id %s", studentId, err)
	}

	studentPloMap := make(map[string][]entity.StudentPloData)
	studentPoMap := make(map[string][]entity.StudentPoData)

	PloCourseMap := make(map[string][]entity.StudentCourseData)
	PoCourseMap := make(map[string][]entity.StudentCourseData)

	for _, record := range ploRecords {
		studentPloData, ok := studentPloMap[record.StudentId]
		if !ok {

			studentPloMap[record.StudentId] = append(studentPloMap[record.StudentId], entity.StudentPloData{
				ProgramLearningOutcomeId: record.ProgramLearningOutcomeId,
				Code:                     record.ProgramLearningOutcomeCode,
				DescriptionThai:          record.DescriptionThai,
				ProgramYear:              record.ProgramYear,
			})
		} else {
			isExist := false
			for i := range studentPloData {
				if studentPloData[i].ProgramLearningOutcomeId == record.ProgramLearningOutcomeId {
					isExist = true
					break
				}
			}
			if !isExist {
				studentPloMap[record.StudentId] = append(studentPloMap[record.StudentId], entity.StudentPloData{
					ProgramLearningOutcomeId: record.ProgramLearningOutcomeId,
					Code:                     record.ProgramLearningOutcomeCode,
					DescriptionThai:          record.DescriptionThai,
					ProgramYear:              record.ProgramYear,
				})
			}
		}

		ploData, ok := PloCourseMap[record.ProgramLearningOutcomeId]

		if !ok {

			PloCourseMap[record.ProgramLearningOutcomeId] = append(PloCourseMap[record.ProgramLearningOutcomeId], entity.StudentCourseData{
				Id:               record.CourseId,
				Code:             record.CourseCode,
				Name:             record.CourseName,
				Pass:             record.Pass,
				Year:             record.Year,
				SemesterSequence: record.SemesterSequence,
			})
		} else {
			isExist := false
			for i := range ploData {
				if ploData[i].Id == record.CourseId {
					isExist = true
					break
				}
			}
			if !isExist {
				PloCourseMap[record.ProgramLearningOutcomeId] = append(PloCourseMap[record.ProgramLearningOutcomeId], entity.StudentCourseData{
					Id:               record.CourseId,
					Code:             record.CourseCode,
					Name:             record.CourseName,
					Pass:             record.Pass,
					Year:             record.Year,
					SemesterSequence: record.SemesterSequence,
				})
			}
		}
	}

	for _, record := range poRecords {
		studentData, found := studentPoMap[record.StudentId]
		if !found {
			studentPoMap[record.StudentId] = append(studentPoMap[record.StudentId], entity.StudentPoData{
				ProgramOutcomeId: record.ProgramOutcomeId,
				Code:             record.ProgramOutcomeCode,
				Name:             record.ProgramOutcomeName,
			})

		} else {
			isExist := false
			for i := range studentData {
				if studentData[i].ProgramOutcomeId == record.ProgramOutcomeId {
					isExist = true
					break
				}
			}
			if !isExist {
				studentPoMap[record.StudentId] = append(studentPoMap[record.StudentId], entity.StudentPoData{
					ProgramOutcomeId: record.ProgramOutcomeId,
					Code:             record.ProgramOutcomeCode,
					Name:             record.ProgramOutcomeName,
				})
			}
		}

		poData, found := PoCourseMap[record.ProgramOutcomeId]

		if !found {
			PoCourseMap[record.ProgramOutcomeId] = append(PoCourseMap[record.ProgramOutcomeId], entity.StudentCourseData{
				Id:               record.CourseId,
				Code:             record.CourseCode,
				Name:             record.CourseName,
				Pass:             record.Pass,
				Year:             record.Year,
				SemesterSequence: record.SemesterSequence,
			})
		} else {
			isExist := false
			for i := range poData {
				if poData[i].Id == record.CourseId {
					isExist = true
					break
				}
			}
			if !isExist {
				PoCourseMap[record.ProgramOutcomeId] = append(PoCourseMap[record.ProgramOutcomeId], entity.StudentCourseData{
					Id:               record.CourseId,
					Code:             record.CourseCode,
					Name:             record.CourseName,
					Pass:             record.Pass,
					Year:             record.Year,
					SemesterSequence: record.SemesterSequence,
				})
			}
		}
	}

	students := make([]entity.StudentOutcomes, 0)

	for studentId := range studentPloMap {
		for ploIndex := range studentPloMap[studentId] {
			studentPloMap[studentId][ploIndex].Courses = PloCourseMap[studentPloMap[studentId][ploIndex].ProgramLearningOutcomeId]
		}
		for poIndex := range studentPoMap[studentId] {
			studentPoMap[studentId][poIndex].Courses = PoCourseMap[studentPoMap[studentId][poIndex].ProgramOutcomeId]
		}

		students = append(students, entity.StudentOutcomes{
			StudentId:               studentId,
			ProgramLearningOutcomes: studentPloMap[studentId],
			ProgramOutcomes:         studentPoMap[studentId],
		})
	}

	return students, nil
}

func (u coursePortfolioUseCase) GetCourseCloAssessment(programmeId string, fromSerm, toSerm int) (*entity.FileResponse, error) {
	rows, err := u.CoursePortfolioRepository.GetCourseCloAssessment(programmeId, fromSerm, toSerm)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course clo assessment %s", err)
	}

	// Process to nested structure
	output := ProcessCourseCloAssessment(rows)

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot marshal course clo assessment %s", err)
	}

	fmt.Println(string(jsonData))

	fileDir := filepath.Join("output", "course_clo_assessment")
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return nil, errs.New(errs.SameCode, "cannot create directory %s", err)
	}
	fileName := fmt.Sprintf("course_clo_assessment_%s.xlsx", time.Now().Format("20060102150405"))
	filepath := filepath.Join(fileDir, fileName)

	err = WriteCourseCloAssessmentToExcel(output, filepath)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot write to excel %s", err)
	}

	return &entity.FileResponse{
		FileName: fileName,
		FilePath: filepath,
		FileType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}, nil
}

func WriteCourseCloAssessmentToExcel(outputs []Output, filename string) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName(f.GetSheetName(0), sheet)

	// Create central style
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create style: %v", err)
	}

	// Collect unique outcomes
	uniquePLOs := map[string]map[string]bool{}
	uniqueSOs := map[string]map[string]bool{}
	uniquePOs := map[string]bool{}

	for _, output := range outputs {
		for _, clo := range output.CLOs {
			for _, assessment := range clo.Assessments {
				for plo, splos := range assessment.PLOs {
					if uniquePLOs[plo] == nil {
						uniquePLOs[plo] = map[string]bool{}
					}
					for splo := range splos {
						uniquePLOs[plo][splo] = true
					}
				}
				for so, ssos := range assessment.SOs {
					if uniqueSOs[so] == nil {
						uniqueSOs[so] = map[string]bool{}
					}
					for sso := range ssos {
						uniqueSOs[so][sso] = true
					}
				}
				for po := range assessment.POs {
					uniquePOs[po] = true
				}
			}
		}
	}

	// Write static headers
	staticHeaders := []string{"Course Code", "Semester", "Course Name", "CLO", "Assessment"}
	for i, h := range staticHeaders {
		col := i + 1
		cell1, _ := excelize.CoordinatesToCellName(col, 1)
		cell2, _ := excelize.CoordinatesToCellName(col, 2)
		if err := f.SetCellValue(sheet, cell1, h); err != nil {
			return err
		}
		if err := f.MergeCell(sheet, cell1, cell2); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheet, cell1, cell2, style); err != nil {
			return err
		}
	}

	// Write dynamic outcome headers
	colIndex := len(staticHeaders) + 1
	outcomeMap := map[string]int{}

	writeOutcomeHeaders := func(groups map[string]map[string]bool) error {
		for mainKey, subMap := range groups {
			startCol := colIndex
			for _, subKey := range getSortedKeys(subMap) {
				cellTop, _ := excelize.CoordinatesToCellName(colIndex, 1)
				cellBottom, _ := excelize.CoordinatesToCellName(colIndex, 2)
				if err := f.SetCellValue(sheet, cellTop, mainKey); err != nil {
					return err
				}
				if err := f.SetCellValue(sheet, cellBottom, subKey); err != nil {
					return err
				}
				if err := f.SetCellStyle(sheet, cellTop, cellBottom, style); err != nil {
					return err
				}
				outcomeMap[subKey] = colIndex
				colIndex++
			}
			startCell, _ := excelize.CoordinatesToCellName(startCol, 1)
			endCell, _ := excelize.CoordinatesToCellName(colIndex-1, 1)
			if err := f.MergeCell(sheet, startCell, endCell); err != nil {
				return err
			}
		}
		return nil
	}

	if err := writeOutcomeHeaders(uniquePLOs); err != nil {
		return err
	}
	if err := writeOutcomeHeaders(uniqueSOs); err != nil {
		return err
	}

	// Write POs
	startCol := colIndex
	for _, po := range getSortedKeysBoolMap(uniquePOs) {
		cellTop, _ := excelize.CoordinatesToCellName(colIndex, 1)
		cellBottom, _ := excelize.CoordinatesToCellName(colIndex, 2)
		if err := f.SetCellValue(sheet, cellTop, "PO"); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, cellBottom, po); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheet, cellTop, cellBottom, style); err != nil {
			return err
		}
		outcomeMap[po] = colIndex
		colIndex++
	}
	startPOCell, _ := excelize.CoordinatesToCellName(startCol, 1)
	endPOCell, _ := excelize.CoordinatesToCellName(colIndex-1, 1)
	if err := f.MergeCell(sheet, startPOCell, endPOCell); err != nil {
		return err
	}

	// Write data rows
	row := 3
	for _, output := range outputs {
		startRow := row
		for _, clo := range output.CLOs {
			cloStart := row
			for _, assessment := range clo.Assessments {
				f.SetCellValue(sheet, getCell(1, row), output.CourseCode)
				f.SetCellValue(sheet, getCell(2, row), output.Semester)
				f.SetCellValue(sheet, getCell(3, row), output.CourseName)
				f.SetCellValue(sheet, getCell(4, row), clo.CLO)
				f.SetCellValue(sheet, getCell(5, row), assessment.Name)

				for plo := range assessment.PLOs {
					for splo := range assessment.PLOs[plo] {
						if col, ok := outcomeMap[splo]; ok {
							f.SetCellValue(sheet, getCell(col, row), "X")
						}
					}
				}
				for so := range assessment.SOs {
					for sso := range assessment.SOs[so] {
						if col, ok := outcomeMap[sso]; ok {
							f.SetCellValue(sheet, getCell(col, row), "X")
						}
					}
				}
				for po := range assessment.POs {
					if col, ok := outcomeMap[po]; ok {
						f.SetCellValue(sheet, getCell(col, row), "X")
					}
				}
				row++
			}
			if row > cloStart {
				f.MergeCell(sheet, getCell(4, cloStart), getCell(4, row-1))
			}
		}
		if row > startRow {
			f.MergeCell(sheet, getCell(1, startRow), getCell(1, row-1))
			f.MergeCell(sheet, getCell(2, startRow), getCell(2, row-1))
			f.MergeCell(sheet, getCell(3, startRow), getCell(3, row-1))
		}
	}

	// Adjust column widths
	colWidths := make(map[int]int)
	for r := 1; r < row; r++ {
		for c := 1; c < colIndex; c++ {
			cell, _ := excelize.CoordinatesToCellName(c, r)
			val, err := f.GetCellValue(sheet, cell)
			if err != nil {
				continue
			}
			if len(val) > colWidths[c] {
				colWidths[c] = len(val)
			}
		}
	}

	for colNum, maxLen := range colWidths {
		colName, _ := excelize.ColumnNumberToName(colNum)
		width := float64(maxLen + 10) // Add padding
		if err := f.SetColWidth(sheet, colName, colName, width); err != nil {
			return fmt.Errorf("failed to set column width for %s: %v", colName, err)
		}
	}

	// Save file
	if err := f.SaveAs(filename); err != nil {
		return err
	}

	// Cleanup old files
	fileFolder := filepath.Dir(filename)
	if err := utils.DeleteOldFiles(fileFolder, 1); err != nil {
		return fmt.Errorf("cannot delete old files: %w", err)
	}

	return nil
}

// Helper Functions
func getCell(col, row int) string {
	cell, _ := excelize.CoordinatesToCellName(col, row)
	return cell
}

func getSortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getSortedKeysBoolMap(m map[string]bool) []string {
	return getSortedKeys(m)
}

func (u coursePortfolioUseCase) GetCourseLinkedOutcomes(programmeId string, fromSerm, toSerm int) (*entity.FileResponse, error) {
	rows, err := u.CoursePortfolioRepository.GetCourseLinkedOutcomes(programmeId, fromSerm, toSerm)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course linked outcomes %s", err)
	}

	// Process to nested structure
	output := ProcessCourseLinkedOutcome(rows)

	// jsonData, err := json.MarshalIndent(output, "", "  ")
	// if err != nil {
	// 	return nil, errs.New(errs.SameCode, "cannot marshal course linked outcomes %s", err)
	// }

	// fmt.Println(string(jsonData))

	fileDir := filepath.Join("output", "course_linked_outcomes")
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return nil, errs.New(errs.SameCode, "cannot create directory %s", err)
	}
	fileName := fmt.Sprintf("course_linked_outcomes_%s.xlsx", time.Now().Format("20060102150405"))
	filepath := filepath.Join(fileDir, fileName)

	err = WriteCourseLinkedOutcomes(output, filepath)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot write to excel %s", err)
	}

	return &entity.FileResponse{
		FileName: fileName,
		FilePath: filepath,
		FileType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}, nil

}

func ProcessCourseCloAssessment(rows []entity.FlatRow) []Output {
	courseMap := make(map[string]*Output)

	for _, row := range rows {
		courseKey := row.CourseCode + "_" + row.Semester

		// Initialize course
		if _, exists := courseMap[courseKey]; !exists {
			courseMap[courseKey] = &Output{
				CourseCode: row.CourseCode,
				Semester:   row.Semester,
				CourseName: row.CourseName,
			}
		}

		course := courseMap[courseKey]

		// Find or add CLO
		var clo *CLOGroup
		for i := range course.CLOs {
			if course.CLOs[i].CLO == row.CloDescription {
				clo = &course.CLOs[i]
				break
			}
		}
		if clo == nil {
			course.CLOs = append(course.CLOs, CLOGroup{CLO: row.CloDescription})
			clo = &course.CLOs[len(course.CLOs)-1]
		}

		if row.AssessmentID == "" {
			continue
		}

		// Find or add Assessment
		var assessment *Assessment
		for i := range clo.Assessments {
			if clo.Assessments[i].Name == row.AssessmentName {
				assessment = &clo.Assessments[i]
				break
			}
		}
		if assessment == nil {
			newAssessment := Assessment{
				Name: row.AssessmentName,
				PLOs: map[string]map[string]string{},
				SOs:  map[string]map[string]string{},
				POs:  map[string]string{},
			}
			clo.Assessments = append(clo.Assessments, newAssessment)
			assessment = &clo.Assessments[len(clo.Assessments)-1]
		}

		// Populate maps
		// Populate nested PLO map
		if row.PLOCode != "" && row.SPLOCode != "" {
			if _, exists := assessment.PLOs[row.PLOCode]; !exists {
				assessment.PLOs[row.PLOCode] = map[string]string{}
			}
			assessment.PLOs[row.PLOCode][row.SPLOCode] = "X"
		}

		// Populate nested SO map
		if row.SOCode != "" && row.SSOCode != "" {
			if _, exists := assessment.SOs[row.SOCode]; !exists {
				assessment.SOs[row.SOCode] = map[string]string{}
			}
			assessment.SOs[row.SOCode][row.SSOCode] = "X"
		}

		// Populate PO map
		if row.POCode != "" {
			assessment.POs[row.POCode] = "O"
		}

	}

	// Convert map to slice
	output := make([]Output, 0, len(courseMap))
	for _, v := range courseMap {
		output = append(output, *v)
	}
	return output
}

type Output struct {
	CourseCode string     `json:"course_code"`
	Semester   string     `json:"semester"`
	CourseName string     `json:"course_name"`
	CLOs       []CLOGroup `json:"CLOs"`
}

type CLOGroup struct {
	CLO         string       `json:"CLO"`
	Assessments []Assessment `json:"assessments"`
}

type Assessment struct {
	Name string                       `json:"name"`
	PLOs map[string]map[string]string `json:"PLOs"`
	SOs  map[string]map[string]string `json:"SOs"`
	POs  map[string]string            `json:"POs"`
}

type CourseLinkedOutcome struct {
	CourseCode string                       `json:"course_code"`
	CourseName string                       `json:"course_name"`
	Year       string                       `json:"year"`
	PLOs       map[string]map[string]string `json:"PLOs"`
	SOs        map[string]map[string]string `json:"SOs"`
	POs        map[string]string            `json:"POs"`
}

func ProcessCourseLinkedOutcome(rows []entity.FlatRow) []CourseLinkedOutcome {
	courseMap := make(map[string]*CourseLinkedOutcome)

	for _, row := range rows {
		courseKey := row.CourseCode + "_" + row.Semester

		// Initialize course
		if _, exists := courseMap[courseKey]; !exists {
			courseMap[courseKey] = &CourseLinkedOutcome{
				CourseCode: row.CourseCode,
				CourseName: row.CourseName,
				Year:       row.Semester,
				PLOs:       map[string]map[string]string{},
				SOs:        map[string]map[string]string{},
				POs:        map[string]string{},
			}
		}

		course := courseMap[courseKey]

		// Populate nested PLO map
		if row.PLOCode != "" && row.SPLOCode != "" {
			if _, exists := course.PLOs[row.PLOCode]; !exists {
				course.PLOs[row.PLOCode] = map[string]string{}
			}
			course.PLOs[row.PLOCode][row.SPLOCode] = "X"
		}

		// Populate nested SO map
		if row.SOCode != "" && row.SSOCode != "" {
			if _, exists := course.SOs[row.SOCode]; !exists {
				course.SOs[row.SOCode] = map[string]string{}
			}
			course.SOs[row.SOCode][row.SSOCode] = "X"
		}

		// Populate PO map
		if row.POCode != "" {
			course.POs[row.POCode] = "X"
		}
	}

	//sort the courseMap by CourseCode
	sortedKeys := make([]string, 0, len(courseMap))
	for k := range courseMap {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	output := make([]CourseLinkedOutcome, 0, len(courseMap))
	for _, v := range courseMap {
		output = append(output, *v)
	}
	return output
}

func WriteCourseLinkedOutcomes(outputs []CourseLinkedOutcome, filename string) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName(f.GetSheetName(0), sheet)

	// Create central style
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create style: %v", err)
	}

	// Collect unique outcomes
	uniquePLOs := map[string]map[string]bool{}
	uniqueSOs := map[string]map[string]bool{}
	uniquePOs := map[string]bool{}

	for _, output := range outputs {
		for plo, splos := range output.PLOs {
			if uniquePLOs[plo] == nil {
				uniquePLOs[plo] = map[string]bool{}
			}
			for splo := range splos {
				uniquePLOs[plo][splo] = true
			}
		}
		for so, ssos := range output.SOs {
			if uniqueSOs[so] == nil {
				uniqueSOs[so] = map[string]bool{}
			}
			for sso := range ssos {
				uniqueSOs[so][sso] = true
			}
		}
		for po := range output.POs {
			uniquePOs[po] = true
		}

	}

	// Write static headers
	staticHeaders := []string{"Course Code", "Course Name", "Year"}
	for i, h := range staticHeaders {
		col := i + 1
		cell1, _ := excelize.CoordinatesToCellName(col, 1)
		cell2, _ := excelize.CoordinatesToCellName(col, 2)
		if err := f.SetCellValue(sheet, cell1, h); err != nil {
			return err
		}
		if err := f.MergeCell(sheet, cell1, cell2); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheet, cell1, cell2, style); err != nil {
			return err
		}
	}

	// Write dynamic outcome headers
	colIndex := len(staticHeaders) + 1
	outcomeMap := map[string]int{}

	writeOutcomeHeaders := func(groups map[string]map[string]bool) error {
		for mainKey, subMap := range groups {
			startCol := colIndex
			for _, subKey := range getSortedKeys(subMap) {
				cellTop, _ := excelize.CoordinatesToCellName(colIndex, 1)
				cellBottom, _ := excelize.CoordinatesToCellName(colIndex, 2)
				if err := f.SetCellValue(sheet, cellTop, mainKey); err != nil {
					return err
				}
				if err := f.SetCellValue(sheet, cellBottom, subKey); err != nil {
					return err
				}
				if err := f.SetCellStyle(sheet, cellTop, cellBottom, style); err != nil {
					return err
				}
				outcomeMap[subKey] = colIndex
				colIndex++
			}
			startCell, _ := excelize.CoordinatesToCellName(startCol, 1)
			endCell, _ := excelize.CoordinatesToCellName(colIndex-1, 1)
			if err := f.MergeCell(sheet, startCell, endCell); err != nil {
				return err
			}
		}
		return nil
	}

	if err := writeOutcomeHeaders(uniquePLOs); err != nil {
		return err
	}
	if err := writeOutcomeHeaders(uniqueSOs); err != nil {
		return err
	}

	// Write POs
	startCol := colIndex
	for _, po := range getSortedKeysBoolMap(uniquePOs) {
		cellTop, _ := excelize.CoordinatesToCellName(colIndex, 1)
		cellBottom, _ := excelize.CoordinatesToCellName(colIndex, 2)
		if err := f.SetCellValue(sheet, cellTop, "PO"); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, cellBottom, po); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheet, cellTop, cellBottom, style); err != nil {
			return err
		}
		outcomeMap[po] = colIndex
		colIndex++
	}
	startPOCell, _ := excelize.CoordinatesToCellName(startCol, 1)
	endPOCell, _ := excelize.CoordinatesToCellName(colIndex-1, 1)
	if err := f.MergeCell(sheet, startPOCell, endPOCell); err != nil {
		return err
	}

	// Write data rows
	row := 3
	for _, output := range outputs {
		f.SetCellValue(sheet, getCell(1, row), output.CourseCode)
		f.SetCellValue(sheet, getCell(2, row), output.CourseName)
		f.SetCellValue(sheet, getCell(3, row), output.Year)
		for plo := range output.PLOs {
			for splo := range output.PLOs[plo] {
				if col, ok := outcomeMap[splo]; ok {
					f.SetCellValue(sheet, getCell(col, row), "X")
				}
			}
		}
		for so := range output.SOs {
			for sso := range output.SOs[so] {
				if col, ok := outcomeMap[sso]; ok {
					f.SetCellValue(sheet, getCell(col, row), "X")
				}
			}
		}
		for po := range output.POs {
			if col, ok := outcomeMap[po]; ok {
				f.SetCellValue(sheet, getCell(col, row), "X")
			}
		}
		row++
	}

	// Adjust column widths
	colWidths := make(map[int]int)
	for r := 1; r < row; r++ {
		for c := 1; c < colIndex; c++ {
			cell, _ := excelize.CoordinatesToCellName(c, r)
			val, err := f.GetCellValue(sheet, cell)
			if err != nil {
				continue
			}
			if len(val) > colWidths[c] {
				colWidths[c] = len(val)
			}
		}
	}

	for colNum, maxLen := range colWidths {
		colName, _ := excelize.ColumnNumberToName(colNum)
		width := float64(maxLen + 10) // Add padding
		if err := f.SetColWidth(sheet, colName, colName, width); err != nil {
			return fmt.Errorf("failed to set column width for %s: %v", colName, err)
		}
	}

	// Save file
	if err := f.SaveAs(filename); err != nil {
		return err
	}

	// Cleanup old files
	fileFolder := filepath.Dir(filename)
	if err := utils.DeleteOldFiles(fileFolder, 1); err != nil {
		return fmt.Errorf("cannot delete old files: %w", err)
	}

	return nil
}

func (u coursePortfolioUseCase) GetCourseOutcomesSuccessRate(programmeId string, fromSerm, toSerm int) (*entity.FileResponse, error) {
	output, err := u.CoursePortfolioRepository.GetCourseOutcomesSuccessRate(
		programmeId, fromSerm, toSerm,
	)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course_outcomes_success_rate %s", err)
	}

	// jsonData, err := json.MarshalIndent(output, "", "  ")
	// if err != nil {
	// 	return errs.New(errs.SameCode, "cannot marshal course_outcomes_success_rate %s", err)
	// }

	// fmt.Println(string(jsonData))

	fileDir := filepath.Join("output", "course_outcomes_success_rate")
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return nil, errs.New(errs.SameCode, "cannot create directory %s", err)
	}
	fileName := fmt.Sprintf("course_outcomes_success_rate_%s.xlsx", time.Now().Format("20060102150405"))
	filepath := filepath.Join(fileDir, fileName)

	err = WriteCourseOutcomesSuccessRate(output, filepath)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot write to excel %s", err)
	}

	return &entity.FileResponse{
		FileName: fileName,
		FilePath: filepath,
		FileType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}, nil
}

func WriteCourseOutcomesSuccessRate(outputs []entity.CourseOutcomeSuccessRate, filename string) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName(f.GetSheetName(0), sheet)

	// Create central style
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create style: %v", err)
	}

	// Collect unique outcomes
	uniquePLOs := map[string]map[string]bool{}
	uniqueSOs := map[string]map[string]bool{}
	uniquePOs := map[string]bool{}

	for _, output := range outputs {
		for plo, splos := range output.PLOs {
			if uniquePLOs[plo] == nil {
				uniquePLOs[plo] = map[string]bool{}
			}
			for splo := range splos {
				uniquePLOs[plo][splo] = true
			}
		}
		for so, ssos := range output.SOs {
			if uniqueSOs[so] == nil {
				uniqueSOs[so] = map[string]bool{}
			}
			for sso := range ssos {
				uniqueSOs[so][sso] = true
			}
		}
		for po := range output.POs {
			uniquePOs[po] = true
		}

	}

	// Write static headers
	staticHeaders := []string{"Course Code", "Course Name", "Semester"}
	for i, h := range staticHeaders {
		col := i + 1
		cell1, _ := excelize.CoordinatesToCellName(col, 1)
		cell2, _ := excelize.CoordinatesToCellName(col, 2)
		if err := f.SetCellValue(sheet, cell1, h); err != nil {
			return err
		}
		if err := f.MergeCell(sheet, cell1, cell2); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheet, cell1, cell2, style); err != nil {
			return err
		}
	}

	// Write dynamic outcome headers
	colIndex := len(staticHeaders) + 1
	outcomeMap := map[string]int{}

	writeOutcomeHeaders := func(groups map[string]map[string]bool) error {
		for mainKey, subMap := range groups {
			startCol := colIndex
			for _, subKey := range getSortedKeys(subMap) {
				cellTop, _ := excelize.CoordinatesToCellName(colIndex, 1)
				cellBottom, _ := excelize.CoordinatesToCellName(colIndex, 2)
				if err := f.SetCellValue(sheet, cellTop, mainKey); err != nil {
					return err
				}
				if err := f.SetCellValue(sheet, cellBottom, subKey); err != nil {
					return err
				}
				if err := f.SetCellStyle(sheet, cellTop, cellBottom, style); err != nil {
					return err
				}
				outcomeMap[subKey] = colIndex
				colIndex++
			}
			startCell, _ := excelize.CoordinatesToCellName(startCol, 1)
			endCell, _ := excelize.CoordinatesToCellName(colIndex-1, 1)
			if err := f.MergeCell(sheet, startCell, endCell); err != nil {
				return err
			}
		}
		return nil
	}

	if err := writeOutcomeHeaders(uniquePLOs); err != nil {
		return err
	}
	if err := writeOutcomeHeaders(uniqueSOs); err != nil {
		return err
	}

	// Write POs
	startCol := colIndex
	for _, po := range getSortedKeysBoolMap(uniquePOs) {
		cellTop, _ := excelize.CoordinatesToCellName(colIndex, 1)
		cellBottom, _ := excelize.CoordinatesToCellName(colIndex, 2)
		if err := f.SetCellValue(sheet, cellTop, "PO"); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, cellBottom, po); err != nil {
			return err
		}
		if err := f.SetCellStyle(sheet, cellTop, cellBottom, style); err != nil {
			return err
		}
		outcomeMap[po] = colIndex
		colIndex++
	}
	startPOCell, _ := excelize.CoordinatesToCellName(startCol, 1)
	endPOCell, _ := excelize.CoordinatesToCellName(colIndex-1, 1)
	if err := f.MergeCell(sheet, startPOCell, endPOCell); err != nil {
		return err
	}

	// Write data rows
	row := 3
	for _, output := range outputs {
		f.SetCellValue(sheet, getCell(1, row), output.CourseCode)
		f.SetCellValue(sheet, getCell(2, row), output.CourseName)
		f.SetCellValue(sheet, getCell(3, row), output.CourseSemester)
		for plo := range output.PLOs {
			for splo, value := range output.PLOs[plo] {
				if col, ok := outcomeMap[splo]; ok {
					f.SetCellValue(sheet, getCell(col, row), fmt.Sprintf("%.2f", value))
				}
			}
		}
		for so := range output.SOs {
			for sso, value := range output.SOs[so] {
				if col, ok := outcomeMap[sso]; ok {
					f.SetCellValue(sheet, getCell(col, row), fmt.Sprintf("%.2f", value))
				}
			}
		}
		for po, value := range output.POs {
			if col, ok := outcomeMap[po]; ok {
				f.SetCellValue(sheet, getCell(col, row), fmt.Sprintf("%.2f", value))
			}
		}
		row++
	}

	// Adjust column widths
	colWidths := make(map[int]int)
	for r := 1; r < row; r++ {
		for c := 1; c < colIndex; c++ {
			cell, _ := excelize.CoordinatesToCellName(c, r)
			val, err := f.GetCellValue(sheet, cell)
			if err != nil {
				continue
			}
			if len(val) > colWidths[c] {
				colWidths[c] = len(val)
			}
		}
	}

	for colNum, maxLen := range colWidths {
		colName, _ := excelize.ColumnNumberToName(colNum)
		width := float64(maxLen + 10) // Add padding
		if err := f.SetColWidth(sheet, colName, colName, width); err != nil {
			return fmt.Errorf("failed to set column width for %s: %v", colName, err)
		}
	}

	// Save file
	if err := f.SaveAs(filename); err != nil {
		return err
	}

	// Cleanup old files
	fileFolder := filepath.Dir(filename)
	if err := utils.DeleteOldFiles(fileFolder, 1); err != nil {
		return fmt.Errorf("cannot delete old files: %w", err)
	}

	return nil
}
