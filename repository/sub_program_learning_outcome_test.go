package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

func setupProgrammeAndPLOs(db *gorm.DB) (*entity.Programme, *entity.ProgramLearningOutcome) {
	pro := &entity.Programme{
		Id: "1",
	}
	db.Create(pro)

	plo := &entity.ProgramLearningOutcome{Id: "1", ProgrammeId: pro.Id}
	db.Create(plo)

	return pro, plo
}

func setupSubPLOs(db *gorm.DB, ploId string) []entity.SubProgramLearningOutcome {
	splos := []entity.SubProgramLearningOutcome{
		{Id: "1", Code: "SPL01", ProgramLearningOutcomeId: ploId},
		{Id: "2", Code: "SPL02", ProgramLearningOutcomeId: ploId},
	}
	db.Create(&splos)
	return splos
}

func TestProgramLearningOutcomeRepository(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	repo := NewProgramLearningOutcomeRepositoryGorm(db)

	t.Run("TestCreateSubPLO", func(t *testing.T) {
		_, plo := setupProgrammeAndPLOs(db)

		splos := []entity.SubProgramLearningOutcome{
			{Id: "1", Code: "SPL01", ProgramLearningOutcomeId: plo.Id},
			{Id: "2", Code: "SPL02", ProgramLearningOutcomeId: plo.Id},
		}

		err := repo.CreateSubPLO(splos)
		assert.Nil(t, err, "Expected no error while creating subPLOs, got %v", err)

		var results []entity.SubProgramLearningOutcome
		db.Find(&results)
		assert.Len(t, results, 2, "Expected to find 2 created subPLOs")
	})

	t.Run("TestUpdateSubPLO", func(t *testing.T) {
		_, plo := setupProgrammeAndPLOs(db)
		splo := entity.SubProgramLearningOutcome{Id: "1", Code: "SPL01", ProgramLearningOutcomeId: plo.Id}
		db.Create(&splo)

		update := &entity.SubProgramLearningOutcome{Code: "UPDATED_CODE", ProgramLearningOutcomeId: plo.Id}
		err := repo.UpdateSubPLO("1", update)
		assert.Nil(t, err, "Expected no error while updating subPLO, got %v", err)

		var result entity.SubProgramLearningOutcome
		db.First(&result, "id = ?", "1")
		assert.Equal(t, "UPDATED_CODE", result.Code, "Expected the code to be updated")

		err = repo.UpdateSubPLO("non-existing", update)
		assert.NotNil(t, err, "Expected an error when updating non-existing subPLO")
	})

	t.Run("TestGetAllSubPlo", func(t *testing.T) {
		_, plo := setupProgrammeAndPLOs(db)
		setupSubPLOs(db, plo.Id)

		result, err := repo.GetAllSubPlo()
		assert.Nil(t, err, "Expected no error while getting all subPLOs, got %v", err)
		assert.Len(t, result, 2, "Expected to find 2 subPLOs")
	})

	t.Run("TestGetSubPloByPloId", func(t *testing.T) {
		_, plo := setupProgrammeAndPLOs(db)
		setupSubPLOs(db, plo.Id)

		result, err := repo.GetSubPloByPloId("1")
		assert.Nil(t, err, "Expected no error while getting subPLOs by PLO ID, got %v", err)
		assert.Len(t, result, 2, "Expected to find 2 subPLOs for PLO ID 1")

		result, err = repo.GetSubPloByPloId("non-existing")
		assert.Nil(t, err, "Expected no error when querying with non-existing PLO ID, got %v", err)
		assert.Empty(t, result, "Expected no subPLOs for non-existing PLO ID")
	})

	t.Run("TestGetSubPloByCode", func(t *testing.T) {
		_, plo := setupProgrammeAndPLOs(db)
		splo := entity.SubProgramLearningOutcome{Id: "1", Code: "SPL01", ProgramLearningOutcomeId: plo.Id}
		db.Create(&splo)

		result, err := repo.GetSubPloByCode("SPL01", "CS", 2024)
		assert.Nil(t, err, "Expected no error while getting subPLO by code, got %v", err)
		assert.NotNil(t, result, "Expected to find a subPLO")
		assert.Equal(t, splo.Code, result.Code, "Expected subPLO code to match")

		result, err = repo.GetSubPloByCode("non-existing", "CS", 2024)
		assert.Nil(t, err, "Expected no error for non-existing subPLO code, got %v", err)
		assert.Nil(t, result, "Expected no subPLO for non-existing code")
	})

	t.Run("TestFilterExistedSubPLO", func(t *testing.T) {
		_, plo := setupProgrammeAndPLOs(db)
		setupSubPLOs(db, plo.Id)

		existingIds, err := repo.FilterExistedSubPLO([]string{"1", "3"})
		assert.Nil(t, err, "Expected no error while filtering existing subPLOs, got %v", err)
		assert.Equal(t, []string{"1"}, existingIds, "Expected to find only existing subPLO ids")

		existingIds, err = repo.FilterExistedSubPLO([]string{})
		assert.Nil(t, err, "Expected no error with empty input, got %v", err)
		assert.Empty(t, existingIds, "Expected no results with empty input")
	})

	t.Run("TestDeleteSubPLO", func(t *testing.T) {
		_, plo := setupProgrammeAndPLOs(db)
		splo := entity.SubProgramLearningOutcome{Id: "1", Code: "SPL01", ProgramLearningOutcomeId: plo.Id}
		db.Create(&splo)

		err := repo.DeleteSubPLO("1")
		assert.Nil(t, err, "Expected no error while deleting subPLO, got %v", err)

		var result entity.SubProgramLearningOutcome
		err = db.First(&result, "id = ?", "1").Error
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Expected subPLO to be deleted")

		err = repo.DeleteSubPLO("non-existing")
		assert.NotNil(t, err, "Expected an error when deleting non-existing subPLO")
	})
}
