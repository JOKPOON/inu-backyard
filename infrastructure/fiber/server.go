package fiber

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/captcha"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/controller"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/middleware"
	"github.com/team-inu/inu-backyard/internal/config"
	"github.com/team-inu/inu-backyard/internal/utils/session"
	"github.com/team-inu/inu-backyard/internal/validator"
	"github.com/team-inu/inu-backyard/repository"
	"github.com/team-inu/inu-backyard/usecase"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type fiberServer struct {
	config    config.FiberServerConfig
	gorm      *gorm.DB
	turnstile *captcha.Turnstile
	logger    *zap.Logger
	session   session.SessionM

	studentRepository                entity.StudentRepository
	courseRepository                 entity.CourseRepository
	courseLearningOutcomeRepository  entity.CourseLearningOutcomeRepository
	studentOutcomeRepository         entity.StudentOutcomeRepository
	programLearningOutcomeRepository entity.ProgramLearningOutcomeRepository
	programOutcomeRepository         entity.ProgramOutcomeRepository
	facultyRepository                entity.FacultyRepository
	departmentRepository             entity.DepartmentRepository
	scoreRepository                  entity.ScoreRepository
	userRepository                   entity.UserRepository
	assignmentRepository             entity.AssignmentRepository
	programmeRepository              entity.ProgrammeRepository
	semesterRepository               entity.SemesterRepository
	enrollmentRepository             entity.EnrollmentRepository
	gradeRepository                  entity.GradeRepository
	sessionRepository                entity.SessionRepository
	coursePortfolioRepository        entity.CoursePortfolioRepository
	courseStreamRepository           entity.CourseStreamRepository
	importerRepository               repository.ImporterRepositoryGorm
	mailRepository                   entity.MailRepository
	surveyRepository                 entity.SurveyRepository

	studentUseCase                entity.StudentUseCase
	courseUseCase                 entity.CourseUseCase
	courseLearningOutcomeUseCase  entity.CourseLearningOutcomeUseCase
	studentOutcomeUseCase         entity.StudentOutcomeUseCase
	programLearningOutcomeUseCase entity.ProgramLearningOutcomeUseCase
	programOutcomeUseCase         entity.ProgramOutcomeUseCase
	facultyUseCase                entity.FacultyUseCase
	departmentUseCase             entity.DepartmentUseCase
	scoreUseCase                  entity.ScoreUseCase
	userUseCase                   entity.UserUseCase
	assignmentUseCase             entity.AssignmentUseCase
	programmeUseCase              entity.ProgrammeUseCase
	semesterUseCase               entity.SemesterUseCase
	enrollmentUseCase             entity.EnrollmentUseCase
	gradeUseCase                  entity.GradeUseCase
	sessionUseCase                entity.SessionUseCase
	authUseCase                   entity.AuthUseCase
	coursePortfolioUseCase        entity.CoursePortfolioUseCase
	predictionUseCase             entity.PredictionUseCase
	courseStreamUseCase           entity.CourseStreamsUseCase
	importerUseCase               usecase.ImporterUseCase
	surveyUseCase                 entity.SurveyUseCase

	mailUseCase entity.MailUseCase
}

func NewFiberServer(
	config config.FiberServerConfig,
	gorm *gorm.DB,
	turnstile *captcha.Turnstile,
	logger *zap.Logger,
	session session.SessionM,
) *fiberServer {
	return &fiberServer{
		config:    config,
		gorm:      gorm,
		turnstile: turnstile,
		logger:    logger,
		session:   session,
	}
}

func (f *fiberServer) Run() {
	f.initRepository()
	f.initUseCase()

	err := f.initController()
	if err != nil {
		panic(err)
	}
}

func (f *fiberServer) initRepository() {
	f.studentRepository = repository.NewStudentRepositoryGorm(f.gorm)
	f.courseRepository = repository.NewCourseRepositoryGorm(f.gorm)
	f.courseLearningOutcomeRepository = repository.NewCourseLearningOutcomeRepositoryGorm(f.gorm)
	f.studentOutcomeRepository = repository.NewStudentOutcomeRepositoryGorm(f.gorm)
	f.programLearningOutcomeRepository = repository.NewProgramLearningOutcomeRepositoryGorm(f.gorm)
	f.programOutcomeRepository = repository.NewProgramOutcomeRepositoryGorm(f.gorm)
	f.facultyRepository = repository.NewFacultyRepositoryGorm(f.gorm)
	f.departmentRepository = repository.NewDepartmentRepositoryGorm(f.gorm)
	f.scoreRepository = repository.NewScoreRepositoryGorm(f.gorm)
	f.userRepository = repository.NewUserRepositoryGorm(f.gorm)
	f.assignmentRepository = repository.NewAssignmentRepositoryGorm(f.gorm)
	f.programmeRepository = repository.NewProgrammeRepositoryGorm(f.gorm)
	f.semesterRepository = repository.NewSemesterRepositoryGorm(f.gorm)
	f.enrollmentRepository = repository.NewEnrollmentRepositoryGorm(f.gorm)
	f.gradeRepository = repository.NewGradeRepositoryGorm(f.gorm)
	f.sessionRepository = repository.NewSessionRepository(f.gorm)
	f.coursePortfolioRepository = repository.NewCoursePortfolioRepositoryGorm(f.gorm)
	f.courseStreamRepository = repository.NewCourseStreamRepository(f.gorm)
	f.importerRepository = repository.NewImporterRepositoryGorm(f.gorm)
	f.mailRepository = repository.NewMailRepository(f.session)
	f.surveyRepository = repository.NewSurveyRepositoryGorm(f.gorm)
}

func (f *fiberServer) initUseCase() {
	f.programmeUseCase = usecase.NewProgrammeUseCase(f.programmeRepository)
	f.facultyUseCase = usecase.NewFacultyUseCase(f.facultyRepository)
	f.departmentUseCase = usecase.NewDepartmentUseCase(f.departmentRepository)
	f.studentUseCase = usecase.NewStudentUseCase(f.studentRepository, f.departmentUseCase, f.programmeUseCase)

	f.programLearningOutcomeUseCase = usecase.NewProgramLearningOutcomeUseCase(f.programLearningOutcomeRepository, f.programmeUseCase)
	f.userUseCase = usecase.NewUserUseCase(f.userRepository)
	f.semesterUseCase = usecase.NewSemesterUseCase(f.semesterRepository)
	f.courseUseCase = usecase.NewCourseUseCase(f.courseRepository, f.semesterUseCase, f.userUseCase)
	f.enrollmentUseCase = usecase.NewEnrollmentUseCase(f.enrollmentRepository, f.studentUseCase, f.courseUseCase)
	f.gradeUseCase = usecase.NewGradeUseCase(f.gradeRepository, f.studentUseCase, f.semesterUseCase)
	f.sessionUseCase = usecase.NewSessionUseCase(f.sessionRepository, f.config.Client.Auth)
	f.mailUseCase = usecase.NewMailUseCase(f.mailRepository)
	f.authUseCase = usecase.NewAuthUseCase(f.sessionUseCase, f.userUseCase, f.mailUseCase)
	f.programOutcomeUseCase = usecase.NewProgramOutcomeUseCase(f.programOutcomeRepository, f.semesterUseCase)
	f.studentOutcomeUseCase = usecase.NewStudentOutcomeUseCase(f.studentOutcomeRepository, f.programmeUseCase)
	f.courseLearningOutcomeUseCase = usecase.NewCourseLearningOutcomeUseCase(f.courseLearningOutcomeRepository, f.courseUseCase, f.programmeUseCase, f.programOutcomeUseCase, f.programLearningOutcomeUseCase, f.studentOutcomeUseCase)

	f.assignmentUseCase = usecase.NewAssignmentUseCase(f.assignmentRepository, f.courseLearningOutcomeUseCase, f.courseUseCase)
	f.scoreUseCase = usecase.NewScoreUseCase(f.scoreRepository, f.enrollmentUseCase, f.assignmentUseCase, f.courseUseCase, f.userUseCase, f.studentUseCase)
	f.courseStreamUseCase = usecase.NewCourseStreamUseCase(f.courseStreamRepository, f.courseUseCase)
	f.coursePortfolioUseCase = usecase.NewCoursePortfolioUseCase(f.coursePortfolioRepository, f.courseUseCase, f.userUseCase, f.enrollmentUseCase, f.assignmentUseCase, f.scoreUseCase, f.studentUseCase, f.courseLearningOutcomeUseCase, f.courseStreamUseCase)
	f.importerUseCase = usecase.NewImporterUseCase(f.importerRepository, f.courseUseCase, f.enrollmentUseCase, f.assignmentUseCase, f.programOutcomeUseCase, f.programLearningOutcomeUseCase, f.courseLearningOutcomeUseCase, f.userUseCase)
	f.predictionUseCase = usecase.NewPredictionUseCase(f.config)
	f.surveyUseCase = usecase.NewSurveyUseCase(f.surveyRepository)
}

func (f *fiberServer) initController() error {
	app := fiber.New(fiber.Config{
		AppName:      "inu-backyard",
		ErrorHandler: errorHandler(f.logger),
	})

	app.Use(middleware.NewCorsMiddleware(f.config.Client.Cors.AllowOrigins))
	app.Use(middleware.NewLogger(fiberzap.Config{
		Logger: f.logger,
	}))

	validator := validator.NewPayloadValidator(&f.config.Client.Auth)

	authMiddleware := middleware.NewAuthMiddleware(validator, f.authUseCase)

	studentController := controller.NewStudentController(validator, f.studentUseCase)
	courseController := controller.NewCourseController(validator, f.courseUseCase, f.importerUseCase)
	courseLearningOutcomeController := controller.NewCourseLearningOutcomeController(validator, f.courseLearningOutcomeUseCase)
	studentOutcomeController := controller.NewStudentOutcomeController(validator, f.studentOutcomeUseCase)
	programLearningOutcomeController := controller.NewProgramLearningOutcomeController(validator, f.programLearningOutcomeUseCase)
	subProgramLearningOutcomeController := controller.NewSubProgramLearningOutcomeController(validator, f.programLearningOutcomeUseCase)
	subStudentOutcomeController := controller.NewSubStudentOutcomeController(validator, f.studentOutcomeUseCase)
	programOutcomeController := controller.NewProgramOutcomeController(validator, f.programOutcomeUseCase)
	facultyController := controller.NewFacultyController(validator, f.facultyUseCase)
	departmentController := controller.NewDepartmentController(validator, f.departmentUseCase)
	scoreController := controller.NewScoreController(validator, f.scoreUseCase)
	userController := controller.NewUserController(validator, f.userUseCase, f.authUseCase)
	assignmentController := controller.NewAssignmentController(validator, f.assignmentUseCase)
	programmeController := controller.NewProgrammeController(validator, f.programmeUseCase)
	semesterController := controller.NewSemesterController(validator, f.semesterUseCase)
	enrollmentController := controller.NewEnrollmentController(validator, f.enrollmentUseCase)
	gradeController := controller.NewGradeController(validator, f.gradeUseCase)
	predictionController := controller.NewPredictionController(validator, f.predictionUseCase)
	coursePortfolioController := controller.NewCoursePortfolioController(validator, f.coursePortfolioUseCase)
	courseStreamController := controller.NewCourseStreamController(validator, f.courseStreamUseCase)
	importerController := controller.NewImporterController(validator, f.importerUseCase)
	surveyController := controller.NewSurveyController(validator, f.surveyUseCase)
	authController := controller.NewAuthController(validator, f.config.Client.Auth, *f.turnstile, f.authUseCase, f.userUseCase)

	api := app.Group("/")

	api.Post("/importer", authMiddleware, importerController.Import)

	api.Get("/schools", authMiddleware, studentController.GetAllSchools)
	api.Get("/admissions", authMiddleware, studentController.GetAllAdmissions)

	// student route
	student := api.Group("/students", authMiddleware)

	student.Get("/", studentController.GetStudents)
	student.Post("/", studentController.Create)
	student.Post("/bulk", studentController.CreateMany)
	student.Get("/:studentId", studentController.GetById)
	student.Get("/:studentId/outcomes", coursePortfolioController.GetOutcomesByStudentId)
	student.Patch("/:studentId", studentController.Update)
	student.Delete("/:studentId", studentController.Delete)

	// course route
	course := api.Group("/courses", authMiddleware)

	course.Get("/", courseController.GetAll)
	course.Post("/", courseController.Create)

	course.Patch("/:courseId", courseController.Update)
	course.Delete("/:courseId", courseController.Delete)
	course.Get("/:courseId/students/clos", courseController.GetStudentsPassingCLOs)

	course.Get("/:courseId/clos", courseLearningOutcomeController.GetByCourseId)
	course.Get("/:courseId/clos/students", coursePortfolioController.GetCloPassingStudentsByCourseId)
	course.Get("/:courseId/students/outcomes", coursePortfolioController.GetStudentOutcomeStatusByCourseId)
	course.Get("/:courseId/result", coursePortfolioController.GetCourseResult)
	course.Get("/:courseId/enrollments", enrollmentController.GetByCourseId)
	course.Get("/:courseId/portfolio", coursePortfolioController.Generate)
	course.Patch("/:courseId/portfolio", coursePortfolioController.Update)
	course.Get("/:courseId/assignments", assignmentController.GetByCourseId)
	course.Get("/:courseId/assignment-groups", assignmentController.GetGroupByCourseId)
	course.Get("/:courseId/survey", surveyController.GetByCourseId)
	course.Get("/:courseId", courseController.GetById)

	// course learning outcome route
	clo := api.Group("/clos", authMiddleware)

	clo.Get("/", courseLearningOutcomeController.GetAll)
	clo.Post("/", courseLearningOutcomeController.Create)
	clo.Get("/:cloId", courseLearningOutcomeController.GetById)
	clo.Patch("/:cloId", courseLearningOutcomeController.Update)
	clo.Delete("/:cloId", courseLearningOutcomeController.Delete)

	// program outcome by course learning outcome route
	poByClo := clo.Group("/:cloId/pos", authMiddleware)

	poByClo.Post("/", courseLearningOutcomeController.CreateLinkProgramOutcome)
	poByClo.Delete("/:poId", courseLearningOutcomeController.DeleteLinkProgramOutcome)

	// sub program learning outcome by course learning outcome route
	sploByClo := clo.Group("/:cloId/splos", authMiddleware)

	sploByClo.Post("/", courseLearningOutcomeController.CreateLinkSubProgramLearningOutcome)
	sploByClo.Delete("/:sploId", courseLearningOutcomeController.DeleteLinkSubProgramLearningOutcome)

	// sub student outcome by course learning outcome route
	ssoByClo := clo.Group("/:cloId/ssos", authMiddleware)

	ssoByClo.Post("/", courseLearningOutcomeController.CreateLinkSubStudentOutcome)
	ssoByClo.Delete("/:ssoId", courseLearningOutcomeController.DeleteLinkSubStudentOutcome)

	// student outcome route
	so := api.Group("/sos", authMiddleware)

	so.Get("/", studentOutcomeController.GetAll)
	so.Post("/", studentOutcomeController.Create)
	so.Get("/:soId", studentOutcomeController.GetById)
	so.Patch("/:soId", studentOutcomeController.Update)
	so.Delete("/:soId", studentOutcomeController.Delete)

	// sub student outcome route
	sso := api.Group("ssos", authMiddleware)

	sso.Get("/", subStudentOutcomeController.GetAll)
	sso.Get("/:ssoId", subStudentOutcomeController.GetById)
	sso.Post("/", subStudentOutcomeController.Create)
	sso.Patch("/:ssoId", subStudentOutcomeController.Update)
	sso.Delete("/:ssoId", subStudentOutcomeController.Delete)

	// program learning outcome route
	plo := api.Group("/plos", authMiddleware)

	plo.Get("/", programLearningOutcomeController.GetAll)
	plo.Get("/courses", coursePortfolioController.GetAllProgramLearningOutcomeCourses)
	plo.Post("/", programLearningOutcomeController.Create)
	plo.Get("/:ploId", programLearningOutcomeController.GetById)
	plo.Patch("/:ploId", programLearningOutcomeController.Update)
	plo.Delete("/:ploId", programLearningOutcomeController.Delete)

	// sub program learning outcome route
	splo := api.Group("/splos", authMiddleware)

	splo.Get("/", subProgramLearningOutcomeController.GetAll)
	splo.Post("/", subProgramLearningOutcomeController.Create)
	splo.Get("/:sploId", subProgramLearningOutcomeController.GetById)
	splo.Patch("/:sploId", subProgramLearningOutcomeController.Update)
	splo.Delete("/:sploId", subProgramLearningOutcomeController.Delete)

	// program outcome route
	pos := api.Group("/pos", authMiddleware)

	pos.Get("/", programOutcomeController.GetAll)
	pos.Get("/courses", coursePortfolioController.GetAllProgramOutcomeCourses)
	pos.Post("/", programOutcomeController.Create)
	pos.Get("/:poId", programOutcomeController.GetById)
	pos.Patch("/:poId", programOutcomeController.Update)
	pos.Delete("/:poId", programOutcomeController.Delete)

	// faculty route
	faculty := api.Group("/faculties", authMiddleware)

	faculty.Get("/", facultyController.GetAll)
	faculty.Post("/", facultyController.Create)
	faculty.Get("/:facultyName", facultyController.GetById)
	faculty.Patch("/:facultyName", facultyController.Update)
	faculty.Delete("/:facultyName", facultyController.Delete)

	// department route
	department := api.Group("/departments", authMiddleware)

	department.Get("/", departmentController.GetAll)
	department.Post("/", departmentController.Create)
	department.Get("/:departmentName", departmentController.GetByName)
	department.Patch("/:departmentName", departmentController.Update)
	department.Delete("/:departmentName", departmentController.Delete)

	// score route
	score := api.Group("/scores", authMiddleware)

	score.Get("/", scoreController.GetAll)
	score.Post("/", scoreController.CreateMany)
	score.Get("/:scoreId", scoreController.GetById)
	score.Patch("/:scoreId", scoreController.Update)
	score.Delete("/:scoreId", scoreController.Delete)

	// user route
	user := api.Group("/users", authMiddleware)

	user.Get("/", userController.GetAll)
	user.Post("/", userController.Create)
	user.Get("/:userId", userController.GetById)
	user.Patch("/:userId", userController.Update)
	user.Delete("/:userId", userController.Delete)
	user.Post("/:userId/password", userController.ChangePassword)
	user.Post("/bulk", userController.CreateMany)

	user.Get("/:userId/course", courseController.GetByUserId)

	// assignment route
	assignment := api.Group("/assignments", authMiddleware)

	assignment.Post("/", assignmentController.Create)

	assignment.Get("/", assignmentController.GetAll)
	assignment.Get("/:assignmentId", assignmentController.GetById)
	assignment.Patch("/:assignmentId", assignmentController.Update)
	assignment.Delete("/:assignmentId", assignmentController.Delete)
	assignment.Get("/:assignmentId/scores", scoreController.GetByAssignmentId)

	assignmentGroup := api.Group("/assignment-groups", authMiddleware)
	assignmentGroup.Get("/", assignmentController.GetAllGroup)
	assignmentGroup.Post("/", assignmentController.CreateGroup)
	assignmentGroup.Patch("/:assignmentGroupId", assignmentController.UpdateGroup)
	assignmentGroup.Delete("/:assignmentGroupId", assignmentController.DeleteGroup)

	assignmentGroup.Get("/:assignmentGroupID/assignments", assignmentController.GetByGroupId)

	// clo by assignment route
	cloByAssignment := assignment.Group("/:assignmentId/clos/", authMiddleware)

	cloByAssignment.Post("/", assignmentController.CreateLinkCourseLearningOutcome)
	cloByAssignment.Delete("/:cloId", assignmentController.DeleteLinkCourseLearningOutcome)

	// programme route
	programme := api.Group("/programmes", authMiddleware)

	programme.Post("/", programmeController.Create)
	programme.Post("/:programmeId/link/po", programmeController.CreateLinkWithPO)
	programme.Post("/:programmeId/link/plo", programmeController.CreateLinkWithPLO)
	programme.Post("/:programmeId/link/so", programmeController.CreateLinkWithSO)

	// programme.Get("/", programmeController.GetByNameAndYear)
	// programme.Get("/", programmeController.GetByName)
	programme.Get("/", programmeController.GetAll)
	programme.Patch("/:programmeName", programmeController.Update)
	programme.Delete("/:programmeId", programmeController.Delete)

	programme.Get("/:programmeId/outcome", programmeController.GetAllCourseOutcomeLinked)
	programme.Get("/outcomes/po", programmeController.GetAllCourseLinkedPO)
	programme.Get("/outcomes/plo", programmeController.GetAllCourseLinkedPLO)
	programme.Get("/outcomes/so", programmeController.GetAllCourseLinkedSO)

	// enrollment route
	enrollment := api.Group("/enrollments", authMiddleware)

	enrollment.Get("/", enrollmentController.GetAll)
	enrollment.Post("/", enrollmentController.Create)
	enrollment.Get("/:enrollmentId", enrollmentController.GetById)
	enrollment.Patch("/:enrollmentId", enrollmentController.Update)
	enrollment.Delete("/:enrollmentId", enrollmentController.Delete)

	// semester route
	semester := api.Group("/semesters", authMiddleware)

	semester.Get("/", semesterController.GetAll)
	semester.Get("/:semesterId", semesterController.GetById)
	semester.Post("/", semesterController.Create)
	semester.Patch("/:semesterId", semesterController.Update)
	semester.Delete("/:semesterId", semesterController.Delete)

	// grade route
	grade := api.Group("/grades", authMiddleware)

	grade.Get("/", gradeController.GetAll)
	grade.Post("/", gradeController.CreateMany)
	grade.Get("/:gradeId", gradeController.GetById)
	grade.Patch("/:gradeId", gradeController.Update)
	grade.Delete("/:gradeId", gradeController.Delete)

	// course stream route
	courseStream := api.Group("/course-streams", authMiddleware)
	courseStream.Get("/", courseStreamController.Get)
	courseStream.Post("/", courseStreamController.Create)
	courseStream.Delete("/:courseStreamId", courseStreamController.Delete)

	// prediction
	prediction := api.Group("/prediction", authMiddleware)
	prediction.Post("/predict", predictionController.Predict)

	// survey
	survey := api.Group("/surveys", authMiddleware)

	survey.Get("/", surveyController.GetAll)
	survey.Get("/courses/outcome", surveyController.GetSurveysWithCourseAndOutcomes)
	survey.Post("/", surveyController.Create)
	survey.Get("/:surveyId", surveyController.GetById)
	survey.Get("/:surveyId/questions", surveyController.GetQuestionBySurveyId)
	survey.Patch("/:surveyId", surveyController.Update)
	survey.Delete("/:surveyId", surveyController.Delete)

	question := api.Group("/questions", authMiddleware)
	question.Post("/:surveyId", surveyController.CreateQuestion)
	question.Get("/:questionId", surveyController.GetQuestionById)
	question.Patch("/:questionId", surveyController.UpdateQuestion)
	question.Delete("/:questionId", surveyController.DeleteQuestion)

	// authentication route
	auth := app.Group("/auth")

	auth.Post("/login", authController.SignIn)
	auth.Get("/logout", authMiddleware, authController.SignOut)
	auth.Get("/me", authMiddleware, authController.Me)

	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Post("/reset-password", authController.ResetPassword)
	auth.Get("/:email", authController.GetSessionData)

	app.Get("/metrics", monitor.New())

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	err := app.ListenTLS(":443", "certs/localhost.crt", "certs/localhost.key")
	return err
}
