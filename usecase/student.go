package usecase

import (
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils"
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

func (u studentUseCase) CreateMany(students []entity.Student) error {
	departmentNames := []string{}
	programmeNames := []string{}
	studentIds := []string{}
	for _, student := range students {
		departmentNames = append(departmentNames, student.DepartmentName)
		programmeNames = append(programmeNames, student.ProgrammeName)
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

	deduplicateDepartmentNames := slice.DeduplicateValues(departmentNames)
	nonExistedDepartmentNames, err := u.departmentUseCase.FilterNonExisted(deduplicateDepartmentNames)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get non existed department while creating students")
	} else if len(nonExistedDepartmentNames) != 0 {
		return errs.New(errs.ErrCreateEnrollment, "there are non exist department %v", nonExistedDepartmentNames)
	}

	deduplicateProgrammeNames := slice.DeduplicateValues(programmeNames)
	nonExistedProgrammeNames, err := u.programmeUseCase.FilterNonExisted(deduplicateProgrammeNames)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get non existed programme while creating students")
	} else if len(nonExistedProgrammeNames) != 0 {
		return errs.New(errs.ErrCreateEnrollment, "there are non exist programme %v", nonExistedProgrammeNames)
	}

	err = u.studentRepo.CreateMany(students)
	if err != nil {
		return errs.New(errs.ErrCreateStudent, "cannot create students", err)
	}

	return nil
}

func (u studentUseCase) Update(id string, student *entity.Student) error {
	existStudent, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get student id %s to update", id, err)
	} else if existStudent == nil {
		return errs.New(errs.ErrSubPLONotFound, "cannot get student id %s to update", id)
	}

	err = u.studentRepo.Update(id, student)

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
