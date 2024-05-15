package usecase

import (
	"errors"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/domain/repository"
	"authenticator-backend/extension/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// plantUsecase
// Summary: This is function which defines plantUsecase.
type plantUsecase struct {
	r repository.OuranosRepository
}

// NewPlantUsecase
// Summary: This is function which create PlantUsecase with OuranosRepository.
// input: r(repository.OuranosRepository) repository interface
// output: (IPlantUsecase) usecase interface
func NewPlantUsecase(r repository.OuranosRepository) IPlantUsecase {
	return &plantUsecase{r}
}

// ListPlants
// Summary: This is function which get list of plant.
// input: getPlantModel(traceability.GetPlantModel) GetPlantModel object
// output: ([]PlantModel) list of PlantModel
// output: (error) error object
func (u *plantUsecase) ListPlants(getPlantModel traceability.GetPlantModel) ([]traceability.PlantModel, error) {
	es, err := u.r.ListPlantsByOperatorID(getPlantModel.OperatorID.String())
	if err != nil {
		return nil, err
	}
	ms := es.ToModels()
	return ms, nil
}

// PutPlant
// Summary: This is function which put plant.
// input: plantModel(PlantModel) PlantModel object
// output: (PlantModel) PlantModel object
// output: (error) error object
func (u *plantUsecase) PutPlant(plantModel traceability.PlantModel) (traceability.PlantModel, error) {
	if plantModel.PlantID == uuid.Nil {
		// If pm.PlantID does not exist, create a new one.
		e := traceability.NewPlantEntityModel(
			plantModel.OperatorID,
			plantModel.PlantName,
			plantModel.PlantAddress,
			plantModel.OpenPlantID,
			plantModel.PlantAttribute.GlobalPlantID,
		)
		e, err := u.r.CreatePlant(e)
		if err != nil {
			var pgError *pgconn.PgError
			if errors.As(err, &pgError) {
				if pgError.Code == common.PgErrorUniqueViolation {
					// Unique constraint error
					msg := common.DuplicatePlantError(pgError, plantModel.OpenPlantID, plantModel.PlantAttribute.GlobalPlantID)
					logger.Set(nil).Warnf(err.Error())

					return traceability.PlantModel{}, common.NewCustomError(common.CustomErrorCode400, msg, nil, common.HTTPErrorSourceAuth)
				}
				logger.Set(nil).Errorf(err.Error())

				return traceability.PlantModel{}, err
			}
			logger.Set(nil).Errorf(err.Error())

			return traceability.PlantModel{}, err
		}
		return e.ToModel(), nil
	} else {
		// Update if there is pm.PlantID
		e, err := u.r.GetPlant(plantModel.OperatorID.String(), plantModel.PlantID.String())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Set(nil).Warnf(err.Error())

				return traceability.PlantModel{}, common.NewCustomError(common.CustomErrorCode404, err.Error(), nil, common.HTTPErrorSourceAuth)
			}
			logger.Set(nil).Errorf(err.Error())

			return traceability.PlantModel{}, err
		}
		e.Update(
			plantModel.OperatorID,
			plantModel.PlantName,
			plantModel.PlantAddress,
			plantModel.OpenPlantID,
			plantModel.PlantAttribute.GlobalPlantID,
		)
		e, err = u.r.UpdatePlant(e)
		if err != nil {
			var pgError *pgconn.PgError
			if errors.As(err, &pgError) {
				if pgError.Code == common.PgErrorUniqueViolation {
					// Unique constraint error
					msg := common.DuplicatePlantError(pgError, e.OpenPlantID, plantModel.PlantAttribute.GlobalPlantID)
					logger.Set(nil).Warnf(msg)

					return traceability.PlantModel{}, common.NewCustomError(common.CustomErrorCode400, msg, nil, common.HTTPErrorSourceAuth)
				}
				logger.Set(nil).Errorf(err.Error())

				return traceability.PlantModel{}, err
			}
			logger.Set(nil).Errorf(err.Error())

			return traceability.PlantModel{}, err
		}
		return e.ToModel(), nil
	}
}
