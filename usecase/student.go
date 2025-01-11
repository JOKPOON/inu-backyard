package usecase

import (
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils/slice"
)

type studentUseCase struct {
	studentRepo       entity.StudentRepository
	departmentUseCase entity.DepartmentUseCase
	programmeUseCase  entity.ProgrammeUseCase
}

func NewStudentUseCase(
	studentRepo entity.StudentRepository,
	departmentUseCase entity.DepartmentUseCase,
	programmeUseCase entity.ProgrammeUseCase,
) entity.StudentUseCase {
	return &studentUseCase{
		studentRepo:       studentRepo,
		departmentUseCase: departmentUseCase,
		programmeUseCase:  programmeUseCase,
	}
}

func (u studentUseCase) GetById(id string) (*entity.Student, error) {
	student, err := u.studentRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot get student by id %s", id, err)
	}

	return student, nil
}

func (u studentUseCase) GetAll() ([]entity.Student, error) {
	students, err := u.studentRepo.GetAll()
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot get all students", err)
	}

	return students, nil
}

func (u studentUseCase) GetByParams(params *entity.Student, limit int, offset int) ([]entity.Student, error) {
	students, err := u.studentRepo.GetByParams(params, limit, offset)

	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot get student by params", err)
	}

	return students, nil
}

func (u studentUseCase) CreateMany(students []entity.CreateStudentPayload) error {
	departmentNames := []string{}
	programmeIds := []string{}
	studentIds := []string{}
	for _, student := range students {
		departmentNames = append(departmentNames, student.DepartmentName)
		programmeIds = append(programmeIds, student.ProgrammeId)
		studentIds = append(studentIds, student.Id)
	}

	duplicateStudentIds := slice.GetDuplicateValue(studentIds)
	if len(duplicateStudentIds) != 0 {
		return errs.New(errs.ErrCreateStudent, "there are duplicate student ids in the payload %v", duplicateStudentIds)
	}

	exitedStudentIds, err := u.FilterExisted(studentIds)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get existed student ids while creating students", err)
	} else if len(exitedStudentIds) > 0 {
		return errs.New(errs.ErrCreateStudent, "there are existed student id in the database %v", exitedStudentIds)
	}

	deduplicateDepartmentNames := slice.RemoveDuplicates(departmentNames)
	nonExistedDepartmentNames, err := u.departmentUseCase.FilterNonExisted(deduplicateDepartmentNames)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get non existed department while creating students")
	} else if len(nonExistedDepartmentNames) != 0 {
		return errs.New(errs.ErrCreateEnrollment, "there are non exist department %v", nonExistedDepartmentNames)
	}

	deduplicateProgrammeNames := slice.RemoveDuplicates(programmeIds)
	for _, id := range deduplicateProgrammeNames {
		programme, err := u.programmeUseCase.GetById(id)
		if err != nil {
			return errs.New(errs.SameCode, "cannot get programme id %s while creating students", id, err)
		}

		if programme == nil {
			return errs.New(errs.ErrCreateStudent, "programme id %s not found while creating students", id)
		}
	}

	studentsToCreate := make([]entity.Student, 0, len(students))
	for _, student := range students {
		studentsToCreate = append(studentsToCreate, entity.Student{
			Id:             student.Id,
			FirstName:      student.FirstName,
			LastName:       student.LastName,
			Email:          student.Email,
			ProgrammeId:    student.ProgrammeId,
			DepartmentName: student.DepartmentName,
			GPAX:           *student.GPAX,
			MathGPA:        *student.MathGPA,
			EngGPA:         *student.EngGPA,
			SciGPA:         *student.SciGPA,
			School:         student.School,
			City:           student.City,
			Year:           student.Year,
			Admission:      student.Admission,
			Remark:         student.Remark,
		})
	}

	err = u.studentRepo.CreateMany(studentsToCreate)
	if err != nil {
		return errs.New(errs.ErrCreateStudent, "cannot create students", err)
	}

	return nil
}

func (u studentUseCase) Update(id string, student *entity.UpdateStudentPayload) error {
	existStudent, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get student id %s to update", id, err)
	} else if existStudent == nil {
		return errs.New(errs.ErrSubPLONotFound, "cannot get student id %s to update", id)
	}

	err = u.studentRepo.Update(id,
		&entity.Student{
			Id:             student.Id,
			FirstName:      student.FirstName,
			LastName:       student.LastName,
			Email:          student.Email,
			ProgrammeId:    student.ProgrammeId,
			DepartmentName: student.DepartmentName,
			GPAX:           *student.GPAX,
			MathGPA:        *student.MathGPA,
			EngGPA:         *student.EngGPA,
			SciGPA:         *student.SciGPA,
			School:         student.School,
			City:           student.City,
			Year:           student.Year,
			Admission:      student.Admission,
			Remark:         *student.Remark,
		},
	)

	if err != nil {
		return errs.New(errs.ErrUpdateStudent, "cannot update student by id %s", student.Id, err)
	}

	return nil
}

func (u studentUseCase) Delete(id string) error {
	existStudent, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get student id %s to delete", id, err)
	} else if existStudent == nil {
		return errs.New(errs.ErrSubPLONotFound, "cannot get student id %s to delete", id)
	}

	err = u.studentRepo.Delete(id)
	if err != nil {
		return errs.New(errs.ErrDeleteSubPLO, "cannot delete student", err)
	}

	return nil
}

func (u studentUseCase) FilterExisted(studentIds []string) ([]string, error) {
	existedIds, err := u.studentRepo.FilterExisted(studentIds)
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot query students", err)
	}

	return existedIds, nil
}

func (u studentUseCase) FilterNonExisted(studentIds []string) ([]string, error) {
	existedIds, err := u.studentRepo.FilterExisted(studentIds)
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot query students", err)
	}

	nonExistedIds := slice.Subtraction(studentIds, existedIds)

	return nonExistedIds, nil
}

func (u studentUseCase) GetAllSchools() ([]string, error) {
	schools, err := u.studentRepo.GetAllSchools()
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot query schools student", err)
	}

	return schools, nil
}
func (u studentUseCase) GetAllAdmissions() ([]string, error) {
	admissions, err := u.studentRepo.GetAllAdmissions()
	if err != nil {
		return nil, errs.New(errs.ErrQueryStudent, "cannot query admissions student", err)
	}

	return admissions, nil
}
