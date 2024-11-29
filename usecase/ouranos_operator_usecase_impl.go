package usecase

import (
	"errors"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/domain/repository"
	"authenticator-backend/extension/logger"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// operatorUsecase
// Summary: This is function which defines operatorUsecase.
type operatorUsecase struct {
	r repository.OuranosRepository
}

// NewOperatorUsecase
// Summary: This is function which create OperatorUsecase with OuranosRepository.
// input: r(repository.OuranosRepository) repository interface
// output: (IOperatorUsecase) usecase interface
func NewOperatorUsecase(r repository.OuranosRepository) IOperatorUsecase {
	return &operatorUsecase{r}
}

// GetOperator
// Summary: This is function which get operator.
// input: getOperatorInput(GetOperatorInput) GetOperatorInput object
// output: (OperatorModel) OperatorModel object
// output: (error) error object
func (u *operatorUsecase) GetOperator(getOperatorInput traceability.GetOperatorInput) (traceability.OperatorModel, error) {
	if getOperatorInput.OpenOperatorID != nil {
		// Obtain business operator from public business identifier
		e, err := u.r.GetOperatorByOpenOperatorID(*getOperatorInput.OpenOperatorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Set(nil).Warnf(err.Error())

				return traceability.OperatorModel{}, common.NewCustomError(common.CustomErrorCode404, err.Error(), nil, common.HTTPErrorSourceAuth)
			}
			logger.Set(nil).Errorf(err.Error())

			return traceability.OperatorModel{}, err
		}
		return e.ToModel(), nil
	} else {
		// Obtain business operator from the business identifier
		e, err := u.r.GetOperator(getOperatorInput.OperatorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Set(nil).Warnf(err.Error())

				return traceability.OperatorModel{}, common.NewCustomError(common.CustomErrorCode404, err.Error(), nil, common.HTTPErrorSourceAuth)
			}

			logger.Set(nil).Errorf(err.Error())

			return traceability.OperatorModel{}, err
		}
		return e.ToModel(), nil
	}
}

// GetOperators
// Summary: This is function which get list of operator.
// input: GetOperatorsInput(GetOperatorsInput) GetOperatorsInput object
// output: ([]OperatorModel) slice of OperatorModel object
// output: (error) error object
func (u *operatorUsecase) GetOperators(getOperatorsInput traceability.GetOperatorsInput) ([]traceability.OperatorModel, error) {
	operatorIDs := []string{}
	for _, operatorID := range getOperatorsInput.OperatorIDs {
		operatorIDs = append(operatorIDs, operatorID.String())
	}
	es, err := u.r.GetOperators(operatorIDs)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return []traceability.OperatorModel{}, err
	}
	return es.ToModels(), nil
}

// PutOperator
// Summary: This is function which put operator.
// input: operatorModel(OperatorModel) OperatorModel object
// output: (OperatorModel) OperatorModel object
// output: (error) error object
func (u *operatorUsecase) PutOperator(operatorModel traceability.OperatorModel) (traceability.OperatorModel, error) {
	e, err := u.r.GetOperator(operatorModel.OperatorID.String())
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return traceability.OperatorModel{}, err
	}

	if err := e.Update(operatorModel); err != nil {
		logger.Set(nil).Errorf(err.Error())

		return traceability.OperatorModel{}, err

	}

	e, err = u.r.PutOperator(e)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == common.PgErrorUniqueViolation {
				// Unique constraint error
				msg := common.DuplicateOperatorError(pgError, operatorModel.OperatorAttribute.GlobalOperatorID)
				logger.Set(nil).Warnf(msg)

				return traceability.OperatorModel{}, common.NewCustomError(common.CustomErrorCode400, msg, nil, common.HTTPErrorSourceAuth)
			}
			logger.Set(nil).Errorf(err.Error())

			return traceability.OperatorModel{}, err
		}
		logger.Set(nil).Errorf(err.Error())

		return traceability.OperatorModel{}, err
	}
	return e.ToModel(), nil
}
