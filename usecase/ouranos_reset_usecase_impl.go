package usecase

import (
	"errors"

	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/domain/repository"
	"authenticator-backend/extension/logger"

	"gorm.io/gorm"
)

// resetUsecase
// Summary: This is the structure which defines the usecase for the reset.
type resetUsecase struct {
	OuranosRepository repository.OuranosRepository
	AuthRepository    repository.AuthRepository
}

// NewResetUsecase
// Summary: This is the function which creates the reset usecase.
// input: r(repository.OuranosRepository) Ouranos repository
// input: a(repository.AuthRepository) Auth repository
// output: (IResetUsecase) reset usecase
func NewResetUsecase(r repository.OuranosRepository, a repository.AuthRepository) IResetUsecase {
	return &resetUsecase{r, a}
}

// Reset
// Summary: This is the function which resets the datastore resources.
// input: apikey(string): apikey
// output: (error) error object
func (u *resetUsecase) Reset(apikey string) error {
	param := repository.APIKeyOperatorsParam{APIKey: &apikey}
	apikeyOperators, err := u.AuthRepository.ListAPIKeyOperators(param)
	if err != nil {
		logger.Set(nil).Error(err.Error())

		return err
	}

	operatorIds := apikeyOperators.GetOperatorIds()
	var resetCfps []traceability.CfpEntityModel
	var resetParts []traceability.PartsModel
	var resetTrades []traceability.TradeEntityModel
	var resetPlants []traceability.PlantEntityModel

	for _, operatorId := range operatorIds {
		getPartsModel := traceability.GetPartsModel{OperatorID: operatorId, Limit: 100}

		// get plants
		plants, err := u.OuranosRepository.ListPlantsByOperatorID(operatorId)
		// in case of 0 records, it is not an error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			if err != nil {
				logger.Set(nil).Errorf(err.Error())

				return err
			}
		}
		resetPlants = append(resetPlants, plants...)

		// get parts
		parts, err := u.OuranosRepository.ListParts(getPartsModel)
		// in case of 0 records, it is not an error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			if err != nil {
				logger.Set(nil).Errorf(err.Error())

				return err
			}
		}

		for _, v := range parts {
			resetParts = append(resetParts, *v)
		}

		// get cfpInformation
		var cfps []traceability.CfpEntityModel
		for _, part := range parts {
			cfp, err := u.OuranosRepository.GetCFPInformation(part.TraceID.String())
			// in case of 0 records, it is not an error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}

			if err != nil {
				logger.Set(nil).Errorf(err.Error())

				return err
			}
			cfps = append(cfps, cfp)
			resetCfps = append(resetCfps, cfps...)
		}

		// get trades
		trades, err := u.OuranosRepository.ListTradesByOperatorID(operatorId)
		// in case of 0 records, it is not an error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			if err != nil {
				logger.Set(nil).Errorf(err.Error())

				return err
			}
		}
		resetTrades = append(resetTrades, trades...)
	}

	// delete cfpCertificate
	for _, cfp := range resetCfps {
		if err := u.OuranosRepository.DeleteCFPCertificateByCFPID(cfp.CfpID.String()); err != nil {
			logger.Set(nil).Errorf(err.Error())

			return err
		}
	}

	// delete cfpInformation
	for _, cfp := range resetCfps {
		if err := u.OuranosRepository.DeleteCFPInformation(cfp.CfpID.String()); err != nil {
			logger.Set(nil).Errorf(err.Error())

			return err
		}
	}

	// delete partsStructure
	for _, part := range resetParts {
		if err := u.OuranosRepository.DeletePartsStructure(part.TraceID.String()); err != nil {
			logger.Set(nil).Errorf(err.Error())

			return err
		}
	}

	// delete tradeRequest
	for _, trade := range resetTrades {
		if trade.TradeID != nil {
			tradeID := trade.TradeID.String()
			if err := u.OuranosRepository.DeleteRequestStatusByTradeID(tradeID); err != nil {
				logger.Set(nil).Errorf(err.Error())

				return err
			}
		}
	}

	// delete trade
	for _, trade := range resetTrades {
		if err := u.OuranosRepository.DeleteTrade(trade.TradeID.String()); err != nil {
			logger.Set(nil).Errorf(err.Error())

			return err
		}
	}

	// delete parts
	for _, part := range resetParts {
		if err := u.OuranosRepository.DeleteParts(part.TraceID.String()); err != nil {
			logger.Set(nil).Errorf(err.Error())

			return err
		}
	}

	// delete plants
	for _, plant := range resetPlants {
		if err := u.OuranosRepository.DeletePlant(plant.PlantID.String()); err != nil {
			logger.Set(nil).Errorf(err.Error())

			return err
		}
	}

	return nil
}
