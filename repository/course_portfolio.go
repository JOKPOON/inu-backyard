package repository

import (
	"fmt"
	"sort"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type coursePortfolioRepositoryGorm struct {
	gorm *gorm.DB
}

func NewCoursePortfolioRepositoryGorm(gorm *gorm.DB) entity.CoursePortfolioRepository {
	return &coursePortfolioRepositoryGorm{gorm: gorm}
}

type TabeeSelector string

const (
	TabeeSelectorAssignment         TabeeSelector = "result_student_passing_assignment_percentage"
	TabeeSelectorPo                 TabeeSelector = "student_passing_po_percentage"
	TabeeSelectorCloPercentage      TabeeSelector = "student_passing_clo_percentage"
	TabeeSelectorCloPassingStudents TabeeSelector = "student_passing_clo_with_information"
	TabeeSelectorPloPassingStudents TabeeSelector = "student_passing_plo_with_information"
	TabeeSelectorPoPassingStudents  TabeeSelector = "student_passing_po_with_information"
	TabeeSelectorAllPloCourses      TabeeSelector = "plo_with_course_information"
	TabeeSelectorAllPoCourses       TabeeSelector = "po_with_course_information"
)

func (r coursePortfolioRepositoryGorm) EvaluatePassingAssignmentPercentage(courseId string) ([]entity.AssignmentPercentage, error) {
	var res = []entity.AssignmentPercentage{}

	err := r.evaluateTabeeOutcomes(courseId, TabeeSelectorAssignment, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot query to evaluate assignment percentage: %w", err)
	}

	return res, nil
}

func (r coursePortfolioRepositoryGorm) EvaluatePassingPoPercentage(courseId string) ([]entity.PoPercentage, error) {
	var res = []entity.PoPercentage{}

	err := r.evaluateTabeeOutcomes(courseId, TabeeSelectorPo, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot query to evaluate program outcome percentage: %w", err)
	}

	return res, nil
}

func (r coursePortfolioRepositoryGorm) EvaluatePassingCloPercentage(courseId string) ([]entity.CloPercentage, error) {
	var res = []entity.CloPercentage{}

	err := r.evaluateTabeeOutcomes(courseId, TabeeSelectorCloPercentage, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot query to evaluate course learning outcome percentage: %w", err)
	}

	return res, nil
}

func (r coursePortfolioRepositoryGorm) EvaluatePassingCloStudents(courseId string) ([]entity.CloPassingStudentGorm, error) {
	var res = []entity.CloPassingStudentGorm{}

	err := r.evaluateTabeeOutcomes(courseId, TabeeSelectorCloPassingStudents, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot query to evaluate course learning outcome passing students: %w", err)
	}

	return res, nil
}

func (r coursePortfolioRepositoryGorm) EvaluatePassingPloStudents(courseId string) ([]entity.PloPassingStudentGorm, error) {
	var res = []entity.PloPassingStudentGorm{}

	err := r.evaluateTabeeOutcomes(courseId, TabeeSelectorPloPassingStudents, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot query to evaluate program learning outcome passing students: %w", err)
	}

	return res, nil
}

func (r coursePortfolioRepositoryGorm) EvaluatePassingPoStudents(courseId string) ([]entity.PoPassingStudentGorm, error) {
	var res = []entity.PoPassingStudentGorm{}

	err := r.evaluateTabeeOutcomes(courseId, TabeeSelectorPoPassingStudents, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot query to evaluate program outcome passing students: %w", err)
	}

	return res, nil
}

func (r coursePortfolioRepositoryGorm) evaluateTabeeOutcomes(courseId string, selector TabeeSelector, x interface{}) error {
	template := `
		WITH
			target_course AS (
				SELECT expected_passing_clo_percentage, id
				FROM course
				WHERE id = ?
			),
			clos AS (
				SELECT
					target_course.id AS course_id,
					course_learning_outcome.id,
					expected_passing_assignment_percentage,
					program_outcome_id
				FROM
					course_learning_outcome
				JOIN target_course ON target_course.id = course_learning_outcome.course_id
			),
			assignments AS (
				SELECT
					assignment.name,
					assignment.max_score,
					assignment.expected_score_percentage,
					clos.expected_passing_assignment_percentage,
					clos.id AS c_id,
					assignment.id AS a_id,
					course_id
				FROM clos
				JOIN clo_assignment AS ca ON ca.course_learning_outcome_id = clos.id
				JOIN assignment ON ca.assignment_id = assignment.id
				WHERE assignment.is_included_in_clo IS True
			),
			scores AS (
				SELECT assignments.*, score.*
				FROM assignments
				JOIN score ON score.assignment_id = a_id
                JOIN enrollment ON enrollment.course_id = assignments.course_id AND enrollment.student_id = score.student_id
        		WHERE status != 'WITHDRAW'
			),
			student_passing_assignment AS (
				SELECT
					scores.score >= scores.expected_score_percentage / 100 * scores.max_score AS pass,
					scores.student_id,
					scores.a_id AS a_id,
					scores.c_id AS c_id
				FROM
					scores
			),
			total_assignment_pass AS (
				SELECT
					SUM(pass) AS pass_count,
					a_id,
					c_id
				FROM
					student_passing_assignment
				GROUP BY
					a_id, c_id
			),
			student_count_by_assignment AS (
				SELECT
					COUNT(*) AS count,
					a_id,
					c_id
				FROM
					student_passing_assignment
				GROUP BY
					a_id, c_id
			),
			student_passing_assignment_percentage AS (
				SELECT
					total_assignment_pass.pass_count / student_count_by_assignment.count * 100 AS passing_percentage,
					total_assignment_pass.a_id,
					total_assignment_pass.c_id
				FROM
					total_assignment_pass
					JOIN student_count_by_assignment ON total_assignment_pass.a_id = student_count_by_assignment.a_id
						AND total_assignment_pass.c_id = student_count_by_assignment.c_id
			),
			student_assignment_pass_count AS (
				SELECT
					SUM(pass) AS pass_count,
					c_id,
					student_id
				FROM
					student_passing_assignment
				GROUP BY
					c_id, student_id
			),
			student_count AS (
				SELECT COUNT(*) AS count, c_id FROM student_assignment_pass_count GROUP BY c_id
			),
			assignments_count AS (
				SELECT COUNT(*) AS count , c_id FROM assignments GROUP BY c_id
			),
			student_passing_clo AS (
				SELECT
					student_assignment_pass_count.pass_count >= (clos.expected_passing_assignment_percentage / 100 * assignments_count.count)
					AS pass,
					clos.program_outcome_id,
					clos.id AS clo_id,
					student_assignment_pass_count.student_id
				FROM
					clos
					JOIN assignments_count ON clos.id = assignments_count.c_id
					JOIN student_assignment_pass_count ON clos.id = student_assignment_pass_count.c_id
			),
			total_clo_pass AS (
				SELECT SUM(pass) AS count, clo_id FROM student_passing_clo GROUP BY clo_id
			),
			student_passing_clo_percentage AS (
				SELECT
					total_clo_pass.count / student_count.count * 100 AS passing_percentage, total_clo_pass.clo_id
				FROM
					total_clo_pass
					JOIN student_count ON total_clo_pass.clo_id = student_count.c_id
			),
			student_clo_passing_count_by_po AS (
				SELECT
					SUM(pass) AS pass_count,
					student_id,
					program_outcome_id
				FROM
					student_passing_clo
				GROUP BY
					program_outcome_id, student_id
			),
			clo_count_by_po AS (
				SELECT
					COUNT(*) AS clo_count,
					program_outcome_id AS p_id
				FROM
					clos
				GROUP BY
					program_outcome_id
			),
			student_passing_po AS (
				SELECT
					(pass_count >= target_course.expected_passing_clo_percentage / 100 * clo_count_by_po.clo_count) AS pass,
					clo_count_by_po.p_id,
					student_clo_passing_count_by_po.student_id
				FROM
					clo_count_by_po
					JOIN student_clo_passing_count_by_po ON clo_count_by_po.p_id = student_clo_passing_count_by_po.program_outcome_id,
					target_course
			),
			total_po_pass AS (
				SELECT
					SUM(pass) AS count,
					p_id
				FROM
					student_passing_po
				GROUP BY
					p_id
			),
			student_count_by_po AS (
				SELECT
					COUNT(*) AS count,
					program_outcome_id
				FROM
					student_clo_passing_count_by_po
				GROUP BY
					program_outcome_id
			),
			student_passing_po_percentage AS (
				SELECT
					total_po_pass.count / student_count_by_po.count * 100 AS passing_percentage,
					total_po_pass.p_id
				FROM
					total_po_pass
					JOIN student_count_by_po ON student_count_by_po.program_outcome_id = total_po_pass.p_id
			),
			plos AS (
				SELECT
					clos.id AS c_id,
					sub_program_learning_outcome.id AS splo_id,
					sub_program_learning_outcome.program_learning_outcome_id AS plo_id
				FROM
					clos
					JOIN clo_subplo ON clos.id = clo_subplo.course_learning_outcome_id
					JOIN sub_program_learning_outcome ON clo_subplo.sub_program_learning_outcome_id = sub_program_learning_outcome.id
			),
			distinct_plos AS (
				SELECT
					DISTINCT
					c_id,
					plo_id
				FROM
					plos
			),
			student_passing_clo_with_plo AS (
				SELECT
					pass,
					c_id,
					student_id,
					plo_id
				FROM
					student_passing_clo
					JOIN distinct_plos ON student_passing_clo.clo_id = distinct_plos.c_id
			),
			student_clo_passing_count_by_plo AS (
				SELECT
					SUM(pass) AS pass_count,
					plo_id,
					student_id
				FROM
					student_passing_clo_with_plo
				GROUP BY
					plo_id, student_id
			),
			clo_count_by_plo AS (
				SELECT
					COUNT(*) AS clo_count,
					plo_id
				FROM
					distinct_plos
				GROUP BY
					plo_id
			),
			student_passing_plo AS (
				SELECT
					(pass_count >= target_course.expected_passing_clo_percentage / 100 * clo_count_by_plo.clo_count) AS pass,
					clo_count_by_plo.plo_id,
					student_clo_passing_count_by_plo.student_id
				FROM
					clo_count_by_plo
					JOIN student_clo_passing_count_by_plo ON clo_count_by_plo.plo_id = student_clo_passing_count_by_plo.plo_id,
					target_course
			),
			total_plo_pass AS (
				SELECT
					SUM(pass) AS count,
					plo_id
				FROM
					student_passing_plo
				GROUP BY
					plo_id
			),
			student_count_by_plo AS (
				SELECT
					COUNT(*) AS count,
					plo_id
				FROM
					student_clo_passing_count_by_plo
				GROUP BY
					plo_id
			),
			student_passing_plo_percentage AS (
				SELECT
					total_plo_pass.count / student_count_by_plo.count * 100 AS passing_percentage,
					total_plo_pass.plo_id
				FROM
					total_plo_pass
					JOIN student_count_by_plo ON student_count_by_plo.plo_id = total_plo_pass.plo_id
			),
			result_student_passing_assignment_percentage AS (
                SELECT assignments.name, assignments.expected_score_percentage, student_passing_assignment_percentage.*
                FROM assignments
                JOIN student_passing_assignment_percentage ON assignments.a_id = student_passing_assignment_percentage.a_id AND assignments.c_id = student_passing_assignment_percentage.c_id
            ),
			student_passing_clo_with_information AS (
				SELECT student.first_name, student.last_name, student_passing_clo.student_id, student_passing_clo.pass, student_passing_clo.clo_id, course_learning_outcome.code, course_learning_outcome.description
				FROM student_passing_clo
				JOIN student ON student_passing_clo.student_id = student.id
				JOIN course_learning_outcome ON course_learning_outcome.id = student_passing_clo.clo_id
			),
			student_passing_plo_with_information AS (
				SELECT program_learning_outcome.code, program_learning_outcome.description_thai, program_learning_outcome.program_year, student_passing_plo.pass, student_passing_plo.plo_id, student_passing_plo.student_id
				FROM student_passing_plo
				JOIN program_learning_outcome ON student_passing_plo.plo_id = program_learning_outcome.id
			),
			student_passing_po_with_information AS (
				SELECT program_outcome.code, program_outcome.name, student_passing_po.pass, student_passing_po.p_id, student_passing_po.student_id
				FROM student_passing_po
				JOIN program_outcome ON student_passing_po.p_id = program_outcome.id
			)
		SELECT *
		FROM %s;
	`

	query := fmt.Sprintf(template, selector)

	err := r.gorm.Raw(query, courseId).Scan(x).Error
	if err != nil {
		return fmt.Errorf("cannot query to evaluate tabee outcomes: %w", err)
	}

	return nil
}

func (r coursePortfolioRepositoryGorm) EvaluateAllPloCourses() ([]entity.PloCoursesGorm, error) {
	var res = []entity.PloCoursesGorm{}

	err := r.evaluateOutcomesAllCourses(TabeeSelectorAllPloCourses, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot query to evaluate all program learning outcome courses: %w", err)
	}

	return res, nil
}

func (r coursePortfolioRepositoryGorm) EvaluateAllPoCourses() ([]entity.PoCoursesGorm, error) {
	var res = []entity.PoCoursesGorm{}

	err := r.evaluateOutcomesAllCourses(TabeeSelectorAllPoCourses, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot query to evaluate all program outcome courses: %w", err)
	}

	return res, nil
}

func (r coursePortfolioRepositoryGorm) evaluateOutcomesAllCourses(selector TabeeSelector, x interface{}) error {
	template := `
		SELECT *
		FROM %s;
	`

	query := fmt.Sprintf(template, selector)

	err := r.gorm.Raw(query).Scan(x).Error
	if err != nil {
		return fmt.Errorf("cannot query to evaluate outcomes: %w", err)
	}

	return nil
}

func (r coursePortfolioRepositoryGorm) UpdateCoursePortfolio(courseId string, data datatypes.JSON) error {
	completed := true

	err := r.gorm.Model(&entity.Course{}).Where("id = ?", courseId).Updates(&entity.Course{
		PortfolioData:        data,
		IsPortfolioCompleted: completed,
	}).Error
	if err != nil {
		return fmt.Errorf("cannot update course: %w", err)
	}

	return nil
}

func (r coursePortfolioRepositoryGorm) EvaluateProgramLearningOutcomesByStudentId(studentId string) ([]entity.StudentPlosGorm, error) {
	var res = []entity.StudentPlosGorm{}

	err := r.evaluateOutcomesByStudentId(studentId, TabeeSelectorPloPassingStudents, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot query to evaluate student program learning outcomes: %w", err)
	}

	return res, nil
}

func (r coursePortfolioRepositoryGorm) EvaluateProgramOutcomesByStudentId(studentId string) ([]entity.StudentPosGorm, error) {
	var res = []entity.StudentPosGorm{}

	err := r.evaluateOutcomesByStudentId(studentId, TabeeSelectorPoPassingStudents, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot query to evaluate student program outcomes: %w", err)
	}

	return res, nil
}

func (r coursePortfolioRepositoryGorm) evaluateOutcomesByStudentId(studentId string, selector TabeeSelector, x interface{}) error {
	template := `
		WITH
			courses AS (
				SELECT expected_passing_clo_percentage, id
				FROM course
			),
			clos AS (
				SELECT
                	courses.id AS course_id,
					course_learning_outcome.id,
					expected_passing_assignment_percentage,
					program_outcome_id
				FROM
					course_learning_outcome
				JOIN courses ON courses.id = course_learning_outcome.course_id
			),
			assignments AS (
				SELECT
					assignment.name,
					assignment.max_score,
					assignment.expected_score_percentage,
					clos.expected_passing_assignment_percentage,
					clos.id AS c_id,
					assignment.id AS a_id,
                	course_id
				FROM clos
				JOIN clo_assignment AS ca ON ca.course_learning_outcome_id = clos.id
				JOIN assignment ON ca.assignment_id = assignment.id
				WHERE assignment.is_included_in_clo IS True
			),
			scores AS (
				SELECT assignments.*, score.*
				FROM assignments
				JOIN score ON score.assignment_id = a_id
                JOIN enrollment ON enrollment.course_id = assignments.course_id AND enrollment.student_id = score.student_id
        		WHERE status != 'WITHDRAW' AND score.student_id = ?
			),
			student_passing_assignment AS (
				SELECT
					scores.score >= scores.expected_score_percentage / 100 * scores.max_score AS pass,
					scores.student_id,
					scores.a_id AS a_id,
					scores.c_id AS c_id,
                	course_id
				FROM
					scores
			),
			total_assignment_pass AS (
				SELECT
					SUM(pass) AS pass_count,
					a_id,
					c_id,
                	course_id
				FROM
					student_passing_assignment
				GROUP BY
					a_id, c_id
			),
			student_count_by_assignment AS (
				SELECT
					COUNT(*) AS count,
					a_id,
					c_id,
                	course_id
				FROM
					student_passing_assignment
				GROUP BY
					a_id, c_id
			),
			student_passing_assignment_percentage AS (
				SELECT
					total_assignment_pass.pass_count / student_count_by_assignment.count * 100 AS passing_percentage,
					total_assignment_pass.a_id,
					total_assignment_pass.c_id,
                	total_assignment_pass.course_id
				FROM
					total_assignment_pass
					JOIN student_count_by_assignment ON total_assignment_pass.a_id = student_count_by_assignment.a_id
						AND total_assignment_pass.c_id = student_count_by_assignment.c_id
			),
			student_assignment_pass_count AS (
				SELECT
					SUM(pass) AS pass_count,
					c_id,
					student_id,
                	course_id
				FROM
					student_passing_assignment
				GROUP BY
					c_id, student_id
			),
			student_count AS (
				SELECT COUNT(*) AS count, c_id, course_id FROM student_assignment_pass_count GROUP BY c_id
			),
			assignments_count AS (
				SELECT COUNT(*) AS count , c_id, course_id FROM assignments GROUP BY c_id
			),
			student_passing_clo AS (
				SELECT
					student_assignment_pass_count.pass_count >= (clos.expected_passing_assignment_percentage / 100 * assignments_count.count)
					AS pass,
					clos.program_outcome_id,
					clos.id AS clo_id,
					student_assignment_pass_count.student_id,
                	clos.course_id
				FROM
					clos
					JOIN assignments_count ON clos.id = assignments_count.c_id
					JOIN student_assignment_pass_count ON clos.id = student_assignment_pass_count.c_id
			),
			total_clo_pass AS (
				SELECT SUM(pass) AS count, clo_id, course_id FROM student_passing_clo GROUP BY clo_id
			),
			student_passing_clo_percentage AS (
				SELECT
					total_clo_pass.count / student_count.count * 100 AS passing_percentage, total_clo_pass.clo_id, total_clo_pass.course_id
				FROM
					total_clo_pass
					JOIN student_count ON total_clo_pass.clo_id = student_count.c_id
			),
			student_clo_passing_count_by_po AS (
				SELECT
					SUM(pass) AS pass_count,
					student_id,
					program_outcome_id,
                	course_id
				FROM
					student_passing_clo
				GROUP BY
					course_id, program_outcome_id, student_id
			),
			clo_count_by_po AS (
				SELECT
					COUNT(*) AS clo_count,
					program_outcome_id AS p_id,
                	course_id
				FROM
					clos
				GROUP BY
					course_id, program_outcome_id
			),
			student_passing_po AS (
				SELECT
					(pass_count >= courses.expected_passing_clo_percentage / 100 * clo_count_by_po.clo_count) AS pass,
					clo_count_by_po.p_id,
					student_clo_passing_count_by_po.student_id,
                	clo_count_by_po.course_id
				FROM
					clo_count_by_po
					JOIN student_clo_passing_count_by_po ON clo_count_by_po.p_id = student_clo_passing_count_by_po.program_outcome_id
					JOIN courses ON courses.id = clo_count_by_po.course_id AND courses.id = student_clo_passing_count_by_po.course_id
			),
			total_po_pass AS (
				SELECT
					SUM(pass) AS count,
					p_id,
                	course_id
				FROM
					student_passing_po
				GROUP BY
					p_id, course_id
			),
			student_count_by_po AS (
				SELECT
					COUNT(*) AS count,
					program_outcome_id,
                	course_id
				FROM
					student_clo_passing_count_by_po
				GROUP BY
					course_id, program_outcome_id
			),
			student_passing_po_percentage AS (
				SELECT
					total_po_pass.count / student_count_by_po.count * 100 AS passing_percentage,
					total_po_pass.p_id,
                	total_po_pass.course_id
				FROM
					total_po_pass
					JOIN student_count_by_po ON student_count_by_po.program_outcome_id = total_po_pass.p_id
                		AND student_count_by_po.course_id = total_po_pass.course_id
			),
			plos AS (
				SELECT
					clos.id AS c_id,
					sub_program_learning_outcome.id AS splo_id,
					sub_program_learning_outcome.program_learning_outcome_id AS plo_id,
                	course_id
				FROM
					clos
					JOIN clo_subplo ON clos.id = clo_subplo.course_learning_outcome_id
					JOIN sub_program_learning_outcome ON clo_subplo.sub_program_learning_outcome_id = sub_program_learning_outcome.id
			),
			distinct_plos AS (
				SELECT
					DISTINCT
					c_id,
					plo_id,
                	course_id
				FROM
					plos
			),
			student_passing_clo_with_plo AS (
				SELECT
					pass,
					c_id,
					student_id,
					plo_id,
                	student_passing_clo.course_id
				FROM
					student_passing_clo
					JOIN distinct_plos ON student_passing_clo.clo_id = distinct_plos.c_id
			),
			student_clo_passing_count_by_plo AS (
				SELECT
					SUM(pass) AS pass_count,
					plo_id,
					student_id,
                	course_id
				FROM
					student_passing_clo_with_plo
				GROUP BY
					course_id, plo_id, student_id
			),
			clo_count_by_plo AS (
				SELECT
					COUNT(*) AS clo_count,
					plo_id,
                	course_id
				FROM
					distinct_plos
				GROUP BY
					course_id, plo_id
			),
			student_passing_plo AS (
				SELECT
					(pass_count >= courses.expected_passing_clo_percentage / 100 * clo_count_by_plo.clo_count) AS pass,
					clo_count_by_plo.plo_id,
					student_clo_passing_count_by_plo.student_id,
                	clo_count_by_plo.course_id
				FROM
					clo_count_by_plo
					JOIN student_clo_passing_count_by_plo ON clo_count_by_plo.plo_id = student_clo_passing_count_by_plo.plo_id
					JOIN courses ON courses.id = clo_count_by_plo.course_id AND clo_count_by_plo.course_id = student_clo_passing_count_by_plo.course_id
				WHERE
                	student_id IS NOT NULL
			),
			total_plo_pass AS (
				SELECT
					SUM(pass) AS count,
					plo_id,
                	course_id
				FROM
					student_passing_plo
				GROUP BY
					course_id, plo_id
			),
			student_count_by_plo AS (
				SELECT
					COUNT(*) AS count,
					plo_id,
                	course_id
				FROM
					student_clo_passing_count_by_plo
				GROUP BY
					course_id, plo_id
			),
			student_passing_plo_percentage AS (
				SELECT
					total_plo_pass.count / student_count_by_plo.count * 100 AS passing_percentage,
					total_plo_pass.plo_id,
                	total_plo_pass.course_id
				FROM
					total_plo_pass
					JOIN student_count_by_plo ON student_count_by_plo.plo_id = total_plo_pass.plo_id
                		AND total_plo_pass.course_id = student_count_by_plo.course_id
			), student_passing_plo_with_information AS (
				SELECT
                	program_learning_outcome.code AS plo_code,
                	program_learning_outcome.description_thai,
                	program_learning_outcome.program_year,
                	student_passing_plo.pass,
                	student_passing_plo.plo_id,
                	student_passing_plo.student_id,
                	course_id,
                	course.name AS course_name,
                	course.code AS course_code,
                	semester.year,
                	semester.semester_sequence
				FROM student_passing_plo
				JOIN program_learning_outcome ON student_passing_plo.plo_id = program_learning_outcome.id
                JOIN course ON course.id = student_passing_plo.course_id
                JOIN semester ON semester.id = course.semester_id
			),
			student_passing_po_with_information AS (
				SELECT
                	program_outcome.code AS po_code,
                	program_outcome.name AS po_name,
                	student_passing_po.pass,
                	student_passing_po.p_id,
                	student_passing_po.student_id,
                	course_id,
                	course.name AS course_name,
                	course.code AS course_code,
                	semester.year,
                	semester.semester_sequence
				FROM student_passing_po
				JOIN program_outcome ON student_passing_po.p_id = program_outcome.id
                JOIN course ON course.id = student_passing_po.course_id
                JOIN semester ON semester.id = course.semester_id
			)
		SELECT *
		FROM %s;
	`

	query := fmt.Sprintf(template, selector)

	err := r.gorm.Raw(query, studentId).Scan(x).Error
	if err != nil {
		return fmt.Errorf("cannot query to evaluate student outcomes: %w", err)
	}

	return nil
}

func cacheOutcomes(gorm *gorm.DB, selector TabeeSelector) {
	template := `
	CREATE TABLE %s AS
	WITH
		courses AS (
			SELECT expected_passing_clo_percentage, id
			FROM course
		),
		clos AS (
			SELECT
				courses.id AS course_id,
				course_learning_outcome.id,
				expected_passing_assignment_percentage,
				program_outcome_id
			FROM
				course_learning_outcome
			JOIN courses ON courses.id = course_learning_outcome.course_id
		),
		assignments AS (
			SELECT
				assignment.name,
				assignment.max_score,
				assignment.expected_score_percentage,
				clos.expected_passing_assignment_percentage,
				clos.id AS c_id,
				assignment.id AS a_id,
				course_id
			FROM clos
			JOIN clo_assignment AS ca ON ca.course_learning_outcome_id = clos.id
			JOIN assignment ON ca.assignment_id = assignment.id
			WHERE assignment.is_included_in_clo IS True
		),
		scores AS (
			SELECT assignments.*, score.*
			FROM assignments
			JOIN score ON score.assignment_id = a_id
			JOIN enrollment ON enrollment.course_id = assignments.course_id AND enrollment.student_id = score.student_id
			WHERE status != 'WITHDRAW'
		),
		student_passing_assignment AS (
			SELECT
				scores.score >= scores.expected_score_percentage / 100 * scores.max_score AS pass,
				scores.student_id,
				scores.a_id AS a_id,
				scores.c_id AS c_id,
				course_id
			FROM
				scores
		),
		total_assignment_pass AS (
			SELECT
				SUM(pass) AS pass_count,
				a_id,
				c_id,
				course_id
			FROM
				student_passing_assignment
			GROUP BY
				a_id, c_id
		),
		student_count_by_assignment AS (
			SELECT
				COUNT(*) AS count,
				a_id,
				c_id,
				course_id
			FROM
				student_passing_assignment
			GROUP BY
				a_id, c_id
		),
		student_passing_assignment_percentage AS (
			SELECT
				total_assignment_pass.pass_count / student_count_by_assignment.count * 100 AS passing_percentage,
				total_assignment_pass.a_id,
				total_assignment_pass.c_id,
				total_assignment_pass.course_id
			FROM
				total_assignment_pass
				JOIN student_count_by_assignment ON total_assignment_pass.a_id = student_count_by_assignment.a_id
					AND total_assignment_pass.c_id = student_count_by_assignment.c_id
		),
		student_assignment_pass_count AS (
			SELECT
				SUM(pass) AS pass_count,
				c_id,
				student_id,
				course_id
			FROM
				student_passing_assignment
			GROUP BY
				c_id, student_id
		),
		student_count AS (
			SELECT COUNT(*) AS count, c_id, course_id FROM student_assignment_pass_count GROUP BY c_id
		),
		assignments_count AS (
			SELECT COUNT(*) AS count , c_id, course_id FROM assignments GROUP BY c_id
		),
		student_passing_clo AS (
			SELECT
				student_assignment_pass_count.pass_count >= (clos.expected_passing_assignment_percentage / 100 * assignments_count.count)
				AS pass,
				clos.program_outcome_id,
				clos.id AS clo_id,
				student_assignment_pass_count.student_id,
				clos.course_id
			FROM
				clos
				JOIN assignments_count ON clos.id = assignments_count.c_id
				JOIN student_assignment_pass_count ON clos.id = student_assignment_pass_count.c_id
		),
		total_clo_pass AS (
			SELECT SUM(pass) AS count, clo_id, course_id FROM student_passing_clo GROUP BY clo_id
		),
		student_passing_clo_percentage AS (
			SELECT
				total_clo_pass.count / student_count.count * 100 AS passing_percentage, total_clo_pass.clo_id, total_clo_pass.course_id
			FROM
				total_clo_pass
				JOIN student_count ON total_clo_pass.clo_id = student_count.c_id
		),
		student_clo_passing_count_by_po AS (
			SELECT
				SUM(pass) AS pass_count,
				student_id,
				program_outcome_id,
				course_id
			FROM
				student_passing_clo
			GROUP BY
				course_id, program_outcome_id, student_id
		),
		clo_count_by_po AS (
			SELECT
				COUNT(*) AS clo_count,
				program_outcome_id AS p_id,
				course_id
			FROM
				clos
			GROUP BY
				course_id, program_outcome_id
		),
		student_passing_po AS (
			SELECT
				(pass_count >= courses.expected_passing_clo_percentage / 100 * clo_count_by_po.clo_count) AS pass,
				clo_count_by_po.p_id,
				student_clo_passing_count_by_po.student_id,
				clo_count_by_po.course_id
			FROM
				clo_count_by_po
				JOIN student_clo_passing_count_by_po ON clo_count_by_po.p_id = student_clo_passing_count_by_po.program_outcome_id
				JOIN courses ON courses.id = clo_count_by_po.course_id AND courses.id = student_clo_passing_count_by_po.course_id
		),
		total_po_pass AS (
			SELECT
				SUM(pass) AS count,
				p_id,
				course_id
			FROM
				student_passing_po
			GROUP BY
				p_id, course_id
		),
		student_count_by_po AS (
			SELECT
				COUNT(*) AS count,
				program_outcome_id,
				course_id
			FROM
				student_clo_passing_count_by_po
			GROUP BY
				course_id, program_outcome_id
		),
		student_passing_po_percentage AS (
			SELECT
				total_po_pass.count / student_count_by_po.count * 100 AS passing_percentage,
				total_po_pass.p_id,
				total_po_pass.course_id
			FROM
				total_po_pass
				JOIN student_count_by_po ON student_count_by_po.program_outcome_id = total_po_pass.p_id
					AND student_count_by_po.course_id = total_po_pass.course_id
		),
		plos AS (
			SELECT
				clos.id AS c_id,
				sub_program_learning_outcome.id AS splo_id,
				sub_program_learning_outcome.program_learning_outcome_id AS plo_id,
				course_id
			FROM
				clos
				JOIN clo_subplo ON clos.id = clo_subplo.course_learning_outcome_id
				RIGHT JOIN sub_program_learning_outcome ON clo_subplo.sub_program_learning_outcome_id = sub_program_learning_outcome.id
		),
		distinct_plos AS (
			SELECT
				DISTINCT
				c_id,
				plo_id,
				course_id
			FROM
				plos
		),
		student_passing_clo_with_plo AS (
			SELECT
				pass,
				c_id,
				student_id,
				plo_id,
				student_passing_clo.course_id
			FROM
				student_passing_clo
				RIGHT JOIN distinct_plos ON student_passing_clo.clo_id = distinct_plos.c_id
		),
		student_clo_passing_count_by_plo AS (
			SELECT
				SUM(pass) AS pass_count,
				plo_id,
				student_id,
				course_id
			FROM
				student_passing_clo_with_plo
			GROUP BY
				course_id, plo_id, student_id
		),
		clo_count_by_plo AS (
			SELECT
				COUNT(*) AS clo_count,
				plo_id,
				course_id
			FROM
				distinct_plos
			GROUP BY
				course_id, plo_id
		),
		student_passing_plo AS (
			SELECT
				(pass_count >= courses.expected_passing_clo_percentage / 100 * clo_count_by_plo.clo_count) AS pass,
				clo_count_by_plo.plo_id,
				student_clo_passing_count_by_plo.student_id,
				clo_count_by_plo.course_id
			FROM
				clo_count_by_plo
				JOIN student_clo_passing_count_by_plo ON clo_count_by_plo.plo_id = student_clo_passing_count_by_plo.plo_id
				LEFT JOIN courses ON courses.id = clo_count_by_plo.course_id AND clo_count_by_plo.course_id = student_clo_passing_count_by_plo.course_id
		),
		total_plo_pass AS (
			SELECT
				SUM(pass) AS count,
				plo_id,
				course_id
			FROM
				student_passing_plo
			GROUP BY
				course_id, plo_id
		),
		student_count_by_plo AS (
			SELECT
				COUNT(*) AS count,
				plo_id,
				course_id
			FROM
				student_clo_passing_count_by_plo
			GROUP BY
				course_id, plo_id
		),
		student_passing_plo_percentage AS (
			SELECT
				total_plo_pass.count / student_count_by_plo.count * 100 AS passing_percentage,
				total_plo_pass.plo_id,
				total_plo_pass.course_id
			FROM
				total_plo_pass
				LEFT JOIN student_count_by_plo ON student_count_by_plo.plo_id = total_plo_pass.plo_id
					AND total_plo_pass.course_id = student_count_by_plo.course_id
		),
		plo_with_course_information AS (
			SELECT
				passing_percentage,
				plo_id,
				course_id,
				name,
				course.code,
				semester.year,
				semester.semester_sequence
			FROM
				student_passing_plo_percentage
				LEFT JOIN course ON course.id = student_passing_plo_percentage.course_id
				LEFT JOIN semester ON semester.id = course.semester_id
				RIGHT JOIN program_learning_outcome ON program_learning_outcome.id = student_passing_plo_percentage.plo_id
		),
		po_with_course_information AS (
			SELECT
				passing_percentage,
				program_outcome.id AS p_id,
				course_id,
				course.name,
				course.code,
				semester.year,
				semester.semester_sequence
			FROM
				student_passing_po_percentage
				JOIN course ON course.id = student_passing_po_percentage.course_id
				JOIN semester ON semester.id = course.semester_id
				RIGHT JOIN program_outcome ON program_outcome.id = student_passing_po_percentage.p_id
		)
		SELECT *
		FROM %s;
	`

	drop := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, selector)
	query := fmt.Sprintf(template, selector, selector)

	gorm.Exec(drop)
	gorm.Exec(query)
}

func (r coursePortfolioRepositoryGorm) GetCourseCloAssessment(programmeId string, fromSerm, toSerm int) ([]entity.FlatRow, error) {
	query := `
		SELECT
			c.code AS course_code,
			s.year AS semester,
			c.name AS course_name,
			clo.id AS clo_id,
			clo.code AS clo_description,
			a.id AS assessment_id,
			a.name AS assessment_name,
			splo.code AS splo_code,
			plo.code AS plo_code,
			sso.code AS sso_code,
			so.code AS so_code,
			po.code AS po_code
		FROM
			course c
		JOIN semester s ON
			s.id = c.semester_id
		JOIN course_learning_outcome clo ON
			clo.course_id = c.id
		LEFT JOIN clo_assignment ca ON
			ca.course_learning_outcome_id = clo.id
		LEFT JOIN assignment a ON
			a.id = ca.assignment_id
			-- Join to SPLOs and their parent PLOs
		LEFT JOIN clo_subplo csp ON
			csp.course_learning_outcome_id = clo.id
		LEFT JOIN sub_program_learning_outcome splo ON
			splo.id = csp.sub_program_learning_outcome_id
		LEFT JOIN program_learning_outcome plo ON
			plo.id = splo.program_learning_outcome_id
			-- Join to SSOs and their parent SOs
		LEFT JOIN clo_subso csso ON
			csso.course_learning_outcome_id = clo.id
		LEFT JOIN sub_student_outcome sso ON
			sso.id = csso.sub_student_outcome_id
		LEFT JOIN student_outcome so ON
			so.id = sso.student_outcome_id
			-- Join to POs
		LEFT JOIN clo_po cpo ON
			cpo.course_learning_outcome_id = clo.id
		LEFT JOIN program_outcome po ON
			po.id = cpo.program_outcome_id
		WHERE
			c.programme_id = ? AND s.year BETWEEN ? AND ?
		ORDER BY
			c.code,
			s.year,
			clo.code,
			a.name;
	`

	var rows []entity.FlatRow
	tx := r.gorm.Raw(query, programmeId, fromSerm, toSerm).Scan(&rows)
	if tx.Error != nil {
		return nil, fmt.Errorf("cannot query to get course linked outcomes: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, fmt.Errorf("no data found for the given parameters")
	}

	return rows, nil
}

func (r coursePortfolioRepositoryGorm) GetCourseLinkedOutcomes(programmeId string, fromSerm, toSerm int) ([]entity.FlatRow, error) {
	query := `
			SELECT
				c.code AS course_code,
				s.year AS semester,
				c.name AS course_name,
				clo.id AS clo_id,
				clo.code AS clo_description,
				splo.code AS splo_code,
				plo.code AS plo_code,
				sso.code AS sso_code,
				so.code AS so_code,
				po.code AS po_code
			FROM
				course c
			JOIN semester s ON
				s.id = c.semester_id
			JOIN course_learning_outcome clo ON
				clo.course_id = c.id
				-- Join to SPLOs and their parent PLOs
			LEFT JOIN clo_subplo csp ON
				csp.course_learning_outcome_id = clo.id
			LEFT JOIN sub_program_learning_outcome splo ON
				splo.id = csp.sub_program_learning_outcome_id
			LEFT JOIN program_learning_outcome plo ON
				plo.id = splo.program_learning_outcome_id
				-- Join to SSOs and their parent SOs
			LEFT JOIN clo_subso csso ON
				csso.course_learning_outcome_id = clo.id
			LEFT JOIN sub_student_outcome sso ON
				sso.id = csso.sub_student_outcome_id
			LEFT JOIN student_outcome so ON
				so.id = sso.student_outcome_id
				-- Join to POs
			LEFT JOIN clo_po cpo ON
				cpo.course_learning_outcome_id = clo.id
			LEFT JOIN program_outcome po ON
				po.id = cpo.program_outcome_id
			WHERE
				c.programme_id = ? AND s.year BETWEEN ? AND ?
			ORDER BY
				c.code,
				s.year,
				clo.code;
	`

	var rows []entity.FlatRow
	tx := r.gorm.Raw(query, programmeId, fromSerm, toSerm).Scan(&rows)
	if tx.Error != nil {
		return nil, fmt.Errorf("cannot query to get course linked outcomes: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, fmt.Errorf("no data found for the given parameters")
	}
	fmt.Printf("rows: %v", rows)

	return rows, nil
}

type JoinedScore struct {
	StudentID                           string
	Score                               float64
	AssignmentID                        string
	AssignmentName                      string
	MaxScore                            float64
	ExpectedScorePercentage             float64
	ExpectedPassingStudentPercentage    float64
	CLOID                               string
	CLOCode                             string
	ExpectedPassingAssignmentPercentage float64
	ExpectedPassingCLOPercentage        float64
	CourseID                            string
	POID                                string
	POCode                              string
	PLOID                               string
	PLOCode                             string
	SPLOID                              string
	SPLOCode                            string
	SOID                                string
	SOCode                              string
	SSOID                               string
	SSOCode                             string
}

func (r coursePortfolioRepositoryGorm) GetCourseOutcomesSuccessRate(programmeId string, fromSerm, toSerm int) ([]entity.CourseOutcomeSuccessRate, error) {
	var joinedScores []JoinedScore
	type Courses struct {
		Id       string `json:"id"`
		Code     string `json:"code"`
		Name     string `json:"name"`
		Semester string `json:"semester"`
	}
	var courses []Courses
	db := r.gorm.Raw(`
		SELECT
			c.id,
			c.code,
			c.name,
			CONCAT(s.semester_sequence, '/', s.year) AS semester
		FROM
			course c
		LEFT JOIN semester s ON
			c.semester_id = s.id
		WHERE
			c.programme_id = ? AND s.year BETWEEN ? AND ?
		ORDER BY
			c.code,
			s.year ASC,
			s.semester_sequence ASC
		`,
		programmeId, fromSerm, toSerm,
	).Scan(&courses)
	if db.Error != nil {
		return nil, fmt.Errorf("cannot query to get courses: %w", db.Error)
	}

	coursesOutcomeSuccessRate := make(map[string]entity.CourseOutcomeSuccessRate)

	fmt.Println("Courses: ", courses)
	for _, course := range courses {
		coursesOutcomeSuccessRate[course.Id] = entity.CourseOutcomeSuccessRate{
			CourseId:       course.Id,
			CourseCode:     course.Code,
			CourseName:     course.Name,
			CourseSemester: course.Semester,
			PLOs:           make(map[string]map[string]float64),
			SOs:            make(map[string]map[string]float64),
			POs:            make(map[string]float64),
		}

		r.gorm.Raw(`
		SELECT
			s.student_id,
			s.score,
			a.id AS assignment_id,
			a.name AS assignment_name,
			a.max_score,
			a.expected_score_percentage,
			a.expected_passing_student_percentage,
			clo.id AS clo_id,
			clo.code AS clo_code,
			clo.expected_passing_assignment_percentage,
			c.expected_passing_clo_percentage,
			c.id AS course_id,
			po.id AS po_id,
			po.code AS po_code,
			plo.id AS plo_id,
			plo.code AS plo_code,
			splo.id AS splo_id,
			splo.code AS splo_code,
			so.id AS so_id,
			so.code AS so_code,
			sso.id AS sso_id,
			sso.code AS sso_code
		FROM
			score s
		LEFT JOIN assignment a ON
			s.assignment_id = a.id
		LEFT JOIN clo_assignment clo_a ON
			a.id = clo_a.assignment_id
		LEFT JOIN course_learning_outcome clo ON
			clo_a.course_learning_outcome_id = clo.id
		LEFT JOIN course c ON
			clo.course_id = c.id
		LEFT JOIN clo_po ON clo.id = clo_po.course_learning_outcome_id
		LEFT JOIN program_outcome po ON
			clo_po.program_outcome_id = po.id
		LEFT JOIN clo_subplo ON
			clo_subplo.course_learning_outcome_id = clo.id
		LEFT JOIN sub_program_learning_outcome splo ON
    		clo_subplo.sub_program_learning_outcome_id = splo.id
		LEFT JOIN program_learning_outcome plo ON
			splo.program_learning_outcome_id = plo.id
		LEFT JOIN clo_subso ON
			clo_subso.course_learning_outcome_id = clo.id
		LEFT JOIN sub_student_outcome sso ON
    		clo_subso.sub_student_outcome_id = sso.id
		LEFT JOIN student_outcome so ON
			sso.student_outcome_id = so.id
		WHERE
			a.is_included_in_clo = TRUE
			AND c.id = ?
		ORDER BY
			s.student_id,
			a.id,
			clo.id
	`, course.Id).Scan(&joinedScores)

		type StudentStats struct {
			StudentID        string
			PassedAssignment map[string]bool
			PassedCLOs       map[string]bool
			PassedPOs        map[string]bool
			PassedSPLOs      map[string]bool
			PassedSSOs       map[string]bool
		}

		type CLOStats struct {
			CLOID                               string
			CLOCode                             string
			ExpectedPassingAssignmentPercentage float64
			Assignments                         map[string]bool
		}

		type POStats struct {
			POID                         string
			POCode                       string
			PassedPercentage             float64
			ExpectedPassingCloPercentage float64
			CLOs                         map[string]bool
		}

		type SPLOStats struct {
			SPLOID                       string
			SPLOCode                     string
			PLOID                        string
			PLOCode                      string
			PassedPercentage             float64
			ExpectedPassingCloPercentage float64
			CLOs                         map[string]bool
		}

		type SSOStats struct {
			SOID                         string
			SOCode                       string
			SSOID                        string
			SSOCode                      string
			PassedPercentage             float64
			ExpectedPassingCloPercentage float64
			CLOs                         map[string]bool
		}

		studentStats := make(map[string]*StudentStats)
		cloStats := make(map[string]*CLOStats)
		poStats := make(map[string]*POStats)
		sploStats := make(map[string]*SPLOStats)
		ssoStats := make(map[string]*SSOStats)

		for _, row := range joinedScores {
			if course.Id != row.CourseID {
				continue
			}
			if _, ok := studentStats[row.StudentID]; !ok {
				studentStats[row.StudentID] = &StudentStats{
					StudentID:        row.StudentID,
					PassedAssignment: make(map[string]bool),
					PassedCLOs:       make(map[string]bool),
					PassedPOs:        make(map[string]bool),
					PassedSPLOs:      make(map[string]bool),
					PassedSSOs:       make(map[string]bool),
				}
			}
			studentStats[row.StudentID].PassedAssignment[row.AssignmentID] = (row.Score/row.MaxScore)*100 >= row.ExpectedScorePercentage

			if _, ok := cloStats[row.CLOID]; !ok {
				cloStats[row.CLOID] = &CLOStats{
					CLOID:                               row.CLOID,
					CLOCode:                             row.CLOCode,
					ExpectedPassingAssignmentPercentage: row.ExpectedPassingAssignmentPercentage,
					Assignments:                         make(map[string]bool),
				}
			}
			if _, ok := cloStats[row.CLOID].Assignments[row.AssignmentID]; !ok {
				cloStats[row.CLOID].Assignments[row.AssignmentID] = true
			}

			if _, ok := poStats[row.POID]; !ok {
				poStats[row.POID] = &POStats{
					POID:                         row.POID,
					POCode:                       row.POCode,
					ExpectedPassingCloPercentage: row.ExpectedPassingCLOPercentage,
					CLOs:                         make(map[string]bool),
				}
			}
			if _, ok := poStats[row.POID].CLOs[row.CLOID]; !ok {
				poStats[row.POID].CLOs[row.CLOID] = true
			}

			if _, ok := sploStats[row.SPLOID]; !ok {
				sploStats[row.SPLOID] = &SPLOStats{
					PLOID:                        row.PLOID,
					PLOCode:                      row.PLOCode,
					SPLOID:                       row.SPLOID,
					SPLOCode:                     row.SPLOCode,
					ExpectedPassingCloPercentage: row.ExpectedPassingCLOPercentage,
					CLOs:                         make(map[string]bool),
				}
			}
			if _, ok := sploStats[row.SPLOID].CLOs[row.CLOID]; !ok {
				sploStats[row.SPLOID].CLOs[row.CLOID] = true
			}

			if _, ok := ssoStats[row.SSOID]; !ok {
				ssoStats[row.SSOID] = &SSOStats{
					SOID:                         row.SOID,
					SOCode:                       row.SOCode,
					SSOID:                        row.SSOID,
					SSOCode:                      row.SSOCode,
					ExpectedPassingCloPercentage: row.ExpectedPassingCLOPercentage,
					CLOs:                         make(map[string]bool),
				}
			}
			if _, ok := ssoStats[row.SSOID].CLOs[row.CLOID]; !ok {
				ssoStats[row.SSOID].CLOs[row.CLOID] = true
			}
		}

		for _, studentStat := range studentStats {
			for _, cloStat := range cloStats {
				passedAssignmentCount := 0
				for assignmentID := range cloStat.Assignments {
					if studentStat.PassedAssignment[assignmentID] {
						passedAssignmentCount++
					}
				}

				if (float64(passedAssignmentCount) / float64(len(cloStat.Assignments)) * 100) >= cloStat.ExpectedPassingAssignmentPercentage {
					studentStat.PassedCLOs[cloStat.CLOID] = true
				} else {
					studentStat.PassedCLOs[cloStat.CLOID] = false
				}
			}

		}

		for _, studentStat := range studentStats {
			for _, poStat := range poStats {
				passedCLOCount := 0
				for cloID := range poStat.CLOs {
					if studentStat.PassedCLOs[cloID] {
						passedCLOCount++
					}
				}
				if (float64(passedCLOCount) / float64(len(poStat.CLOs)) * 100) >= poStat.ExpectedPassingCloPercentage {
					studentStat.PassedPOs[poStat.POID] = true
				} else {
					studentStat.PassedPOs[poStat.POID] = false
				}
			}

			for _, sploStat := range sploStats {
				passedCLOCount := 0
				for cloID := range sploStat.CLOs {
					if studentStat.PassedCLOs[cloID] {
						passedCLOCount++
					}
				}
				if (float64(passedCLOCount) / float64(len(sploStat.CLOs)) * 100) >= sploStat.ExpectedPassingCloPercentage {
					studentStat.PassedSPLOs[sploStat.SPLOID] = true
				} else {
					studentStat.PassedSPLOs[sploStat.SPLOID] = false
				}
			}

			for _, ssoStat := range ssoStats {
				passedCLOCount := 0
				for cloID := range ssoStat.CLOs {
					if studentStat.PassedCLOs[cloID] {
						passedCLOCount++
					}
				}
				if (float64(passedCLOCount) / float64(len(ssoStat.CLOs)) * 100) >= ssoStat.ExpectedPassingCloPercentage {
					studentStat.PassedSSOs[ssoStat.SSOID] = true
				} else {
					studentStat.PassedSSOs[ssoStat.SSOID] = false
				}
			}
		}

		fmt.Printf("Course %s (%s , %s) passing rate:\n", course.Code, course.Name, course.Semester)
		for _, poStat := range poStats {
			passPOCount := 0
			for _, studentStat := range studentStats {
				if studentStat.PassedPOs[poStat.POID] {
					passPOCount++
				}
			}
			poStat.PassedPercentage = (float64(passPOCount) / float64(len(studentStats))) * 100

			coursesOutcomeSuccessRate[course.Id].POs[poStat.POCode] = poStat.PassedPercentage

			fmt.Printf("PO %s: %.2f%%\n", poStat.POCode, poStat.PassedPercentage)
		}
		for _, sploStat := range sploStats {
			passSPLOCount := 0
			for _, studentStat := range studentStats {
				if studentStat.PassedSPLOs[sploStat.SPLOID] {
					passSPLOCount++
				}
			}
			sploStat.PassedPercentage = (float64(passSPLOCount) / float64(len(studentStats))) * 100

			if _, ok := coursesOutcomeSuccessRate[course.Id].PLOs[sploStat.PLOID]; !ok {
				coursesOutcomeSuccessRate[course.Id].PLOs[sploStat.PLOCode] = make(map[string]float64)
			}
			coursesOutcomeSuccessRate[course.Id].PLOs[sploStat.PLOCode][sploStat.SPLOCode] = sploStat.PassedPercentage

			fmt.Printf("SPLO %s: %.2f%%\n", sploStat.SPLOCode, sploStat.PassedPercentage)

		}
		for _, ssoStat := range ssoStats {
			passSSOCount := 0
			for _, studentStat := range studentStats {
				if studentStat.PassedSSOs[ssoStat.SSOID] {
					passSSOCount++
				}
			}
			ssoStat.PassedPercentage = (float64(passSSOCount) / float64(len(studentStats))) * 100

			if _, ok := coursesOutcomeSuccessRate[course.Id].SOs[ssoStat.SOCode]; !ok {
				coursesOutcomeSuccessRate[course.Id].SOs[ssoStat.SOCode] = make(map[string]float64)
			}
			coursesOutcomeSuccessRate[course.Id].SOs[ssoStat.SOCode][ssoStat.SSOCode] = ssoStat.PassedPercentage

			fmt.Printf("SSO %s: %.2f%%\n", ssoStat.SSOCode, ssoStat.PassedPercentage)
		}
	}

	coursesOutcomeSuccessRateList := make([]entity.CourseOutcomeSuccessRate, 0, len(coursesOutcomeSuccessRate))
	for _, course := range coursesOutcomeSuccessRate {
		coursesOutcomeSuccessRateList = append(coursesOutcomeSuccessRateList, course)
	}
	sort.Slice(coursesOutcomeSuccessRateList, func(i, j int) bool {
		if coursesOutcomeSuccessRateList[i].CourseCode == coursesOutcomeSuccessRateList[j].CourseCode {
			return coursesOutcomeSuccessRateList[i].CourseSemester < coursesOutcomeSuccessRateList[j].CourseSemester
		}
		return coursesOutcomeSuccessRateList[i].CourseCode < coursesOutcomeSuccessRateList[j].CourseCode
	})

	for _, course := range coursesOutcomeSuccessRateList {
		fmt.Printf("Course %s (%s , %s) passing rate:\n", course.CourseCode, course.CourseName, course.CourseSemester)
		for ploID, splo := range course.PLOs {
			fmt.Printf("PLO %s: ", ploID)
			for sploID, rate := range splo {
				fmt.Printf("SPLO %s: %.2f%% ", sploID, rate)
			}
			fmt.Println()
		}
		for soID, sso := range course.SOs {
			fmt.Printf("SO %s: ", soID)
			for ssoID, rate := range sso {
				fmt.Printf("SSO %s: %.2f%% ", ssoID, rate)
			}
			fmt.Println()
		}
		for poID, rate := range course.POs {
			fmt.Printf("PO %s: %.2f%% ", poID, rate)
		}
		fmt.Println()
	}

	return coursesOutcomeSuccessRateList, nil
}

func (r coursePortfolioRepositoryGorm) GetCourseOutcomes(courseId string) (*entity.CoursePortfolioOutcome, error) {
	var joinedScores []JoinedScore

	coursesOutcomeSuccessRate := entity.CourseOutcomeSuccessRate{
		PLOs: make(map[string]map[string]float64),
		SOs:  make(map[string]map[string]float64),
		POs:  make(map[string]float64),
	}

	r.gorm.Raw(`
		SELECT
			s.student_id,
			s.score,
			a.id AS assignment_id,
			a.name AS assignment_name,
			a.max_score,
			a.expected_score_percentage,
			a.expected_passing_student_percentage,
			clo.id AS clo_id,
			clo.code AS clo_code,
			clo.expected_passing_assignment_percentage,
			c.expected_passing_clo_percentage,
			po.id AS po_id,
			po.code AS po_code,
			plo.id AS plo_id,
			plo.code AS plo_code,
			splo.id AS splo_id,
			splo.code AS splo_code,
			so.id AS so_id,
			so.code AS so_code,
			sso.id AS sso_id,
			sso.code AS sso_code
		FROM
			score s
		LEFT JOIN enrollment e ON
			s.student_id = e.student_id
		LEFT JOIN assignment a ON
			s.assignment_id = a.id
		LEFT JOIN clo_assignment clo_a ON
			a.id = clo_a.assignment_id
		LEFT JOIN course_learning_outcome clo ON
			clo_a.course_learning_outcome_id = clo.id
		LEFT JOIN course c ON
			clo.course_id = c.id
		LEFT JOIN clo_po ON clo.id = clo_po.course_learning_outcome_id
		LEFT JOIN program_outcome po ON
			clo_po.program_outcome_id = po.id
		LEFT JOIN clo_subplo ON
			clo_subplo.course_learning_outcome_id = clo.id
		LEFT JOIN sub_program_learning_outcome splo ON
    		clo_subplo.sub_program_learning_outcome_id = splo.id
		LEFT JOIN program_learning_outcome plo ON
			splo.program_learning_outcome_id = plo.id
		LEFT JOIN clo_subso ON
			clo_subso.course_learning_outcome_id = clo.id
		LEFT JOIN sub_student_outcome sso ON
    		clo_subso.sub_student_outcome_id = sso.id
		LEFT JOIN student_outcome so ON
			sso.student_outcome_id = so.id
		WHERE
			a.is_included_in_clo = TRUE
			AND e.status = 'ENROLL'
			AND c.id = ?
		ORDER BY
			s.student_id,
			a.id,
			clo.id
	`, courseId).Scan(&joinedScores)

	type Metrics struct {
		IsPass   bool
		Expected float64
		Actual   float64
	}

	type StudentStats struct {
		StudentID  string
		Assignment map[string]Metrics
		CLOs       map[string]Metrics
		POs        map[string]Metrics
		SPLOs      map[string]Metrics
		SSOs       map[string]Metrics
	}

	type AssignmentStats struct {
		AssignmentID                        string
		AssignmentName                      string
		PassedPercentage                    float64
		ExpectedPassingAssignmentPercentage float64
	}

	type CLOStats struct {
		CLOID                               string
		CLOCode                             string
		PassedPercentage                    float64
		ExpectedPassingAssignmentPercentage float64
		Assignments                         map[string]Metrics
	}

	type POStats struct {
		POID                         string
		POCode                       string
		PassedPercentage             float64
		ExpectedPassingCloPercentage float64
		CLOs                         map[string]Metrics
	}

	type SPLOStats struct {
		SPLOID                       string
		SPLOCode                     string
		PLOID                        string
		PLOCode                      string
		PassedPercentage             float64
		ExpectedPassingCloPercentage float64
		CLOs                         map[string]Metrics
	}

	type SSOStats struct {
		SOID                         string
		SOCode                       string
		SSOID                        string
		SSOCode                      string
		PassedPercentage             float64
		ExpectedPassingCloPercentage float64
		CLOs                         map[string]Metrics
	}

	studentStats := make(map[string]*StudentStats)
	cloStats := make(map[string]*CLOStats)
	poStats := make(map[string]*POStats)
	sploStats := make(map[string]*SPLOStats)
	ssoStats := make(map[string]*SSOStats)
	assignmentStats := make(map[string]*AssignmentStats)

	for _, row := range joinedScores {
		if _, ok := studentStats[row.StudentID]; !ok {
			studentStats[row.StudentID] = &StudentStats{
				StudentID:  row.StudentID,
				Assignment: make(map[string]Metrics),
				CLOs:       make(map[string]Metrics),
				POs:        make(map[string]Metrics),
				SPLOs:      make(map[string]Metrics),
				SSOs:       make(map[string]Metrics),
			}
		} else {
			continue
		}
	}

	for _, row := range joinedScores {
		if _, ok := assignmentStats[row.AssignmentID]; !ok {
			assignmentStats[row.AssignmentID] = &AssignmentStats{
				AssignmentID:   row.AssignmentID,
				AssignmentName: row.AssignmentName,
			}
			assignmentStats[row.AssignmentID].ExpectedPassingAssignmentPercentage = row.ExpectedPassingAssignmentPercentage
		} else {
			continue
		}
	}

	for _, row := range joinedScores {
		if _, ok := cloStats[row.CLOID]; !ok {
			cloStats[row.CLOID] = &CLOStats{
				CLOID:                               row.CLOID,
				CLOCode:                             row.CLOCode,
				ExpectedPassingAssignmentPercentage: row.ExpectedPassingAssignmentPercentage,
				Assignments:                         make(map[string]Metrics),
			}
		}
		if _, ok := cloStats[row.CLOID].Assignments[row.AssignmentID]; !ok {
			cloStats[row.CLOID].Assignments[row.AssignmentID] = Metrics{}
		}

		if _, ok := poStats[row.POID]; !ok {
			poStats[row.POID] = &POStats{
				POID:                         row.POID,
				POCode:                       row.POCode,
				ExpectedPassingCloPercentage: row.ExpectedPassingCLOPercentage,
				CLOs:                         make(map[string]Metrics),
			}
		}
		if _, ok := poStats[row.POID].CLOs[row.CLOID]; !ok {
			poStats[row.POID].CLOs[row.CLOID] = Metrics{}
		}

		if _, ok := sploStats[row.SPLOID]; !ok {
			sploStats[row.SPLOID] = &SPLOStats{
				PLOID:                        row.PLOID,
				PLOCode:                      row.PLOCode,
				SPLOID:                       row.SPLOID,
				SPLOCode:                     row.SPLOCode,
				ExpectedPassingCloPercentage: row.ExpectedPassingCLOPercentage,
				CLOs:                         make(map[string]Metrics),
			}
		}
		if _, ok := sploStats[row.SPLOID].CLOs[row.CLOID]; !ok {
			sploStats[row.SPLOID].CLOs[row.CLOID] = Metrics{}
		}

		if _, ok := ssoStats[row.SSOID]; !ok {
			ssoStats[row.SSOID] = &SSOStats{
				SOID:                         row.SOID,
				SOCode:                       row.SOCode,
				SSOID:                        row.SSOID,
				SSOCode:                      row.SSOCode,
				ExpectedPassingCloPercentage: row.ExpectedPassingCLOPercentage,
				CLOs:                         make(map[string]Metrics),
			}
		}
		if _, ok := ssoStats[row.SSOID].CLOs[row.CLOID]; !ok {
			ssoStats[row.SSOID].CLOs[row.CLOID] = Metrics{}
		}
	}

	for _, row := range joinedScores {
		if _, ok := studentStats[row.StudentID].Assignment[row.AssignmentID]; !ok {
			metrics := Metrics{
				IsPass:   (row.Score/row.MaxScore)*100 >= row.ExpectedScorePercentage,
				Expected: row.ExpectedScorePercentage,
				Actual:   (row.Score / row.MaxScore) * 100,
			}
			studentStats[row.StudentID].Assignment[row.AssignmentID] = metrics
		} else {
			continue
		}
	}

	for _, studentStat := range studentStats {
		for _, cloStat := range cloStats {
			passedAssignmentCount := 0
			for assignmentID := range cloStat.Assignments {
				if studentStat.Assignment[assignmentID].IsPass {
					passedAssignmentCount++
				}
			}

			if (float64(passedAssignmentCount) / float64(len(cloStat.Assignments)) * 100) >= cloStat.ExpectedPassingAssignmentPercentage {
				studentStat.CLOs[cloStat.CLOID] = Metrics{
					IsPass:   true,
					Expected: cloStat.ExpectedPassingAssignmentPercentage,
					Actual:   (float64(passedAssignmentCount) / float64(len(cloStat.Assignments))) * 100,
				}
			} else {
				studentStat.CLOs[cloStat.CLOID] = Metrics{
					IsPass:   false,
					Expected: cloStat.ExpectedPassingAssignmentPercentage,
					Actual:   (float64(passedAssignmentCount) / float64(len(cloStat.Assignments))) * 100,
				}
			}
		}
	}

	for _, studentStat := range studentStats {
		for _, poStat := range poStats {
			passedCLOCount := 0
			for cloID := range poStat.CLOs {
				if studentStat.CLOs[cloID].IsPass {
					passedCLOCount++
				}
			}
			if (float64(passedCLOCount) / float64(len(poStat.CLOs)) * 100) >= poStat.ExpectedPassingCloPercentage {
				studentStat.POs[poStat.POID] = Metrics{
					IsPass:   true,
					Expected: poStat.ExpectedPassingCloPercentage,
					Actual:   (float64(passedCLOCount) / float64(len(poStat.CLOs))) * 100,
				}
			} else {
				studentStat.POs[poStat.POID] = Metrics{
					IsPass:   false,
					Expected: poStat.ExpectedPassingCloPercentage,
					Actual:   (float64(passedCLOCount) / float64(len(poStat.CLOs))) * 100,
				}
			}
		}

		for _, sploStat := range sploStats {
			passedCLOCount := 0
			for cloID := range sploStat.CLOs {
				if studentStat.CLOs[cloID].IsPass {
					passedCLOCount++
				}
			}
			if (float64(passedCLOCount) / float64(len(sploStat.CLOs)) * 100) >= sploStat.ExpectedPassingCloPercentage {
				studentStat.SPLOs[sploStat.SPLOID] = Metrics{
					IsPass:   true,
					Expected: sploStat.ExpectedPassingCloPercentage,
					Actual:   (float64(passedCLOCount) / float64(len(sploStat.CLOs))) * 100,
				}
			} else {
				studentStat.SPLOs[sploStat.SPLOID] = Metrics{
					IsPass:   false,
					Expected: sploStat.ExpectedPassingCloPercentage,
					Actual:   (float64(passedCLOCount) / float64(len(sploStat.CLOs))) * 100,
				}
			}
		}

		for _, ssoStat := range ssoStats {
			passedCLOCount := 0
			for cloID := range ssoStat.CLOs {
				if studentStat.CLOs[cloID].IsPass {
					passedCLOCount++
				}
			}
			if (float64(passedCLOCount) / float64(len(ssoStat.CLOs)) * 100) >= ssoStat.ExpectedPassingCloPercentage {
				studentStat.SSOs[ssoStat.SSOID] = Metrics{
					IsPass:   true,
					Expected: ssoStat.ExpectedPassingCloPercentage,
					Actual:   (float64(passedCLOCount) / float64(len(ssoStat.CLOs))) * 100,
				}
			} else {
				studentStat.SSOs[ssoStat.SSOID] = Metrics{
					IsPass:   false,
					Expected: ssoStat.ExpectedPassingCloPercentage,
					Actual:   (float64(passedCLOCount) / float64(len(ssoStat.CLOs))) * 100,
				}
			}
		}
	}

	for _, assignmentStat := range assignmentStats {
		passedAssignmentCount := 0
		for _, studentStat := range studentStats {
			if studentStat.Assignment[assignmentStat.AssignmentID].IsPass {
				passedAssignmentCount++
			}
		}
		assignmentStat.PassedPercentage = (float64(passedAssignmentCount) / float64(len(studentStats))) * 100
		fmt.Printf("Assignment %s: %.2f%% actual, %.2f%% expected\n", assignmentStat.AssignmentName, assignmentStat.PassedPercentage, assignmentStat.ExpectedPassingAssignmentPercentage)
	}

	fmt.Println("CLOs passing rate:")
	for _, cloStat := range cloStats {
		passedAssignmentCount := 0
		for _, studentStat := range studentStats {
			if studentStat.CLOs[cloStat.CLOID].IsPass {
				passedAssignmentCount++
			}
		}

		cloStat.PassedPercentage = (float64(passedAssignmentCount) / float64(len(studentStats))) * 100

		fmt.Printf("CLO %s: %.2f%% actual, %.2f%% expected\n", cloStat.CLOCode, cloStat.PassedPercentage, cloStat.ExpectedPassingAssignmentPercentage)
	}

	for _, poStat := range poStats {
		passPOCount := 0
		for _, studentStat := range studentStats {
			if studentStat.POs[poStat.POID].IsPass {
				passPOCount++
			}
		}
		poStat.PassedPercentage = (float64(passPOCount) / float64(len(studentStats))) * 100

		coursesOutcomeSuccessRate.POs[poStat.POCode] = poStat.PassedPercentage

		fmt.Printf("PO %s: %.2f%% actual, %.2f%% expected\n", poStat.POCode, poStat.PassedPercentage, poStat.ExpectedPassingCloPercentage)
	}
	for _, sploStat := range sploStats {
		passSPLOCount := 0
		for _, studentStat := range studentStats {
			if studentStat.SPLOs[sploStat.SPLOID].IsPass {
				passSPLOCount++
			}
		}
		sploStat.PassedPercentage = (float64(passSPLOCount) / float64(len(studentStats))) * 100

		if _, ok := coursesOutcomeSuccessRate.PLOs[sploStat.PLOID]; !ok {
			coursesOutcomeSuccessRate.PLOs[sploStat.PLOCode] = make(map[string]float64)
		}
		coursesOutcomeSuccessRate.PLOs[sploStat.PLOCode][sploStat.SPLOCode] = sploStat.PassedPercentage

		fmt.Printf("SPLO %s: %.2f%% actual, %.2f%% expected\n", sploStat.SPLOCode, sploStat.PassedPercentage, sploStat.ExpectedPassingCloPercentage)

	}
	for _, ssoStat := range ssoStats {
		passSSOCount := 0
		for _, studentStat := range studentStats {
			if studentStat.SSOs[ssoStat.SSOID].IsPass {
				passSSOCount++
			}
		}
		ssoStat.PassedPercentage = (float64(passSSOCount) / float64(len(studentStats))) * 100

		if _, ok := coursesOutcomeSuccessRate.SOs[ssoStat.SOCode]; !ok {
			coursesOutcomeSuccessRate.SOs[ssoStat.SOCode] = make(map[string]float64)
		}
		coursesOutcomeSuccessRate.SOs[ssoStat.SOCode][ssoStat.SSOCode] = ssoStat.PassedPercentage

		fmt.Printf("SSO %s: %.2f%% actual, %.2f%% expected\n", ssoStat.SSOCode, ssoStat.PassedPercentage, ssoStat.ExpectedPassingCloPercentage)
	}

	closPassingRate := map[string]entity.CloPassingRate{}
	for cloId, cloStat := range cloStats {
		assignments := map[string]entity.AssignmentPassingRate{}
		for assignmentID := range cloStat.Assignments {
			assignments[assignmentID] = entity.AssignmentPassingRate{
				AssignmentID:                        assignmentID,
				AssignmentName:                      assignmentStats[assignmentID].AssignmentName,
				PassedPercentage:                    assignmentStats[assignmentID].PassedPercentage,
				ExpectedPassingAssignmentPercentage: assignmentStats[assignmentID].ExpectedPassingAssignmentPercentage,
			}
		}
		closPassingRate[cloId] = entity.CloPassingRate{
			CLOID:                               cloStat.CLOID,
			CLOCode:                             cloStat.CLOCode,
			PassedPercentage:                    cloStat.PassedPercentage,
			ExpectedPassingAssignmentPercentage: cloStat.ExpectedPassingAssignmentPercentage,
			Assignments:                         assignments,
		}
	}

	poPassingRate := map[string]entity.PoPassingRate{}
	for _, poStat := range poStats {
		clos := make(map[string]entity.CloPassingRate)
		for cloId := range poStat.CLOs {
			clos[cloId] = closPassingRate[cloId]
		}
		poPassingRate[poStat.POID] = entity.PoPassingRate{
			POID:                         poStat.POID,
			POCode:                       poStat.POCode,
			PassedPercentage:             poStat.PassedPercentage,
			ExpectedPassingCloPercentage: poStat.ExpectedPassingCloPercentage,
			CLOPassingRate:               clos,
		}
	}

	sploPassingRate := map[string]entity.SploPassingRate{}
	for _, sploStat := range sploStats {
		clos := make(map[string]entity.CloPassingRate)
		for cloId := range sploStat.CLOs {
			clos[cloId] = closPassingRate[cloId]
		}
		sploPassingRate[sploStat.SPLOID] = entity.SploPassingRate{
			SPLOID:                       sploStat.SPLOID,
			SPLOCode:                     sploStat.SPLOCode,
			PLOID:                        sploStat.PLOID,
			PLOCode:                      sploStat.PLOCode,
			PassedPercentage:             sploStat.PassedPercentage,
			ExpectedPassingCloPercentage: sploStat.ExpectedPassingCloPercentage,
			CLOPassingRate:               clos,
		}
	}

	ssoPassingRate := map[string]entity.SsoPassingRate{}
	for _, ssoStat := range ssoStats {
		clos := make(map[string]entity.CloPassingRate)
		for cloId := range ssoStat.CLOs {
			clos[cloId] = closPassingRate[cloId]
		}
		ssoPassingRate[ssoStat.SSOID] = entity.SsoPassingRate{
			SOID:                         ssoStat.SOID,
			SOCode:                       ssoStat.SOCode,
			SSOID:                        ssoStat.SSOID,
			SSOCode:                      ssoStat.SSOCode,
			PassedPercentage:             ssoStat.PassedPercentage,
			ExpectedPassingCloPercentage: ssoStat.ExpectedPassingCloPercentage,
			CLOPassingRate:               clos,
		}
	}

	return &entity.CoursePortfolioOutcome{
		CLOs: closPassingRate,
		POs:  poPassingRate,
		PLOs: sploPassingRate,
		SOs:  ssoPassingRate,
	}, nil
}
