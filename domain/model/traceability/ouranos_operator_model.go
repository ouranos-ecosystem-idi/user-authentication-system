package traceability

import (
	"fmt"
	"time"

	"authenticator-backend/domain/common"
	"authenticator-backend/extension/logger"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OperatorModel
// Summary: This is structure which defines OperatorModel.
// Service: Dataspace
// Router: [PUT] /api/v1/authInfo?dataTarget=operator
// Usage: output
type OperatorModel struct {
	OperatorID        uuid.UUID         `json:"operatorId"`
	OperatorName      string            `json:"operatorName"`
	OperatorAddress   string            `json:"operatorAddress"`
	OpenOperatorID    string            `json:"openOperatorId"`
	OperatorAttribute OperatorAttribute `json:"operatorAttribute"`
}

// OperatorAttribute
// Summary: This is structure which defines OperatorAttribute.
type OperatorAttribute struct {
	GlobalOperatorID *string `json:"globalOperatorId"`
}

// OperatorModels
// Summary: This is structure which defines OperatorModels.
// Service: Dataspace
// Router: [GET] /api/v1/authInfo?dataTarget=operator
// Usage: output
type OperatorModels []*OperatorModel

// PutOperatorInput
// Summary: This is structure which defines PutOperatorInput.
// Service: Dataspace
// Router: [PUT] /api/v1/authInfo?dataTarget=operator
// Usage: input
type PutOperatorInput struct {
	OperatorID             string                  `json:"operatorId"`
	OperatorName           string                  `json:"operatorName"`
	OperatorAddress        string                  `json:"operatorAddress"`
	OpenOperatorID         string                  `json:"openOperatorId"`
	OperatorAttributeInput *OperatorAttributeInput `json:"operatorAttribute"`
}

// OperatorAttributeInput
// Summary: This is structure which defines OperatorAttributeInput.
type OperatorAttributeInput struct {
	GlobalOperatorID *string `json:"globalOperatorId"`
}

// OperatorEntityModel
// Summary: This is structure which defines OperatorEntityModel.
// DBName: operators
type OperatorEntityModel struct {
	OperatorID        uuid.UUID      `json:"operatorId" gorm:"type:uuid;not null"`
	OperatorName      string         `json:"operatorName" gorm:"type:string"`
	OperatorAddress   string         `json:"operatorAddress" gorm:"type:string"`
	OpenOperatorID    string         `json:"openOperatorId" gorm:"type:string"`
	GlobalOperatorID  *string        `json:"globalOperatorId" gorm:"type:string"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt"`
	CreatedAt         time.Time      `json:"createdAt" gorm:"<-:create "`
	CreatedOperatorID string         `json:"createdUserId" gorm:"type:varchar(256);not null; <-:create"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	UpdatedOperatorID string         `json:"updatedUserId" gorm:"type:varchar(256);not null"`
}

// OperatorEntityModels
// Summary: This is structure which defines OperatorEntityModels.
type OperatorEntityModels []*OperatorEntityModel

// GetOperatorInput
// Summary: This is structure which defines GetOperatorInput.
// Service: Dataspace
// Router: [GET] /api/v1/authInfo?dataTarget=operator
// Usage: input
type GetOperatorInput struct {
	OperatorID     string  `json:"operatorId"`
	OpenOperatorID *string `json:"openOperatorId"`
}

// GetOperatorsInput
// Summary: This is structure which defines GetOperatorsInput.
// Service: Dataspace
// Router: [GET] /api/v1/authInfo?dataTarget=operator
// Usage: input
type GetOperatorsInput struct {
	OperatorIDs []uuid.UUID `json:"operatorIds"`
}

// validate
// Summary: This is function which validate value of PutOperatorInput.
// output: (error) error object
func (i PutOperatorInput) validate() error {
	errors := []error{}
	err := validation.ValidateStruct(&i,
		validation.Field(
			&i.OperatorID,
			validation.By(common.StringUUIDValid),
		),
		validation.Field(
			&i.OperatorName,
			validation.Required,
			validation.RuneLength(1, 256),
		),
		validation.Field(
			&i.OperatorAddress,
			validation.Required,
			validation.RuneLength(1, 256),
		),
		validation.Field(
			&i.OpenOperatorID,
			validation.Required,
			validation.RuneLength(1, 20),
		),
		validation.Field(
			&i.OperatorAttributeInput,
			validation.Required,
		),
	)
	if err != nil {
		errors = append(errors, err)
	}

	var attributeErr error
	if i.OperatorAttributeInput != nil {
		if attributeErr = i.OperatorAttributeInput.validate(); attributeErr != nil {
			attributeErr = fmt.Errorf(common.ValidateStructureError("operatorAttribute", attributeErr))
			errors = append(errors, attributeErr)
		}
	}

	if len(errors) > 0 {
		if attributeErr != nil {
			return common.JoinErrors(errors)
		} else {
			return err
		}
	}

	return nil
}

// validate
// Summary: This is function which validate value of OperatorAttributeInput.
// output: (error) error object
func (i *OperatorAttributeInput) validate() error {
	return validation.ValidateStruct(i,
		validation.Field(
			&i.GlobalOperatorID,
			validation.RuneLength(0, 256),
		),
	)
}

// Validate
// Summary: This is function which validate value of PutOperatorInput.
// output: (error) error object
func (i PutOperatorInput) Validate() error {
	if err := i.validate(); err != nil {
		logger.Set(nil).Warnf(err.Error())

		return err
	}

	return nil
}

// ToModel
// Summary: This is function which convert PutOperatorInput to OperatorModel.
// output: (OperatorModel) OperatorModel object
// output: (error) error object
func (i PutOperatorInput) ToModel() (OperatorModel, error) {
	operatorID, err := uuid.Parse(i.OperatorID)
	if err != nil {
		logger.Set(nil).Warnf(err.Error())

		return OperatorModel{}, fmt.Errorf(common.InvalidUUIDError("operatorId"))
	}

	operatorAttribute := OperatorAttribute{
		GlobalOperatorID: i.OperatorAttributeInput.GlobalOperatorID,
	}

	m := OperatorModel{
		OperatorID:        operatorID,
		OperatorName:      i.OperatorName,
		OperatorAddress:   i.OperatorAddress,
		OpenOperatorID:    i.OpenOperatorID,
		OperatorAttribute: operatorAttribute,
	}
	return m, nil
}

// ToModel
// Summary: This is function which convert OperatorEntityModel to OperatorModel.
// output: (OperatorModel) OperatorModel object
func (e *OperatorEntityModel) ToModel() OperatorModel {
	OperatorAttribute := OperatorAttribute{
		GlobalOperatorID: e.GlobalOperatorID,
	}

	return OperatorModel{
		OperatorID:        e.OperatorID,
		OperatorName:      e.OperatorName,
		OperatorAddress:   e.OperatorAddress,
		OpenOperatorID:    e.OpenOperatorID,
		OperatorAttribute: OperatorAttribute,
	}
}

// ToModels
// Summary: This is function which convert OperatorEntityModels to []OperatorModel.
// output: ([]OperatorModel) slice of OperatorModel object
func (es OperatorEntityModels) ToModels() []OperatorModel {
	ms := make([]OperatorModel, len(es))
	for i, e := range es {
		m := e.ToModel()
		ms[i] = m
	}
	return ms
}

// Update
// Summary: This is function which update value of Operator.
// input: operatorModel(OperatorModel) OperatorModel object
// output: (error) error object
func (e *OperatorEntityModel) Update(operatorModel OperatorModel) error {
	if e.OpenOperatorID != operatorModel.OpenOperatorID {
		err := fmt.Errorf(common.FieldIsImutable("openOperatorID"))
		logger.Set(nil).Warnf(err.Error())

		return common.NewCustomError(common.CustomErrorCode400, err.Error(), nil, common.HTTPErrorSourceAuth)
	}

	e.OperatorName = operatorModel.OperatorName
	e.OperatorAddress = operatorModel.OperatorAddress
	e.GlobalOperatorID = operatorModel.OperatorAttribute.GlobalOperatorID
	e.UpdatedAt = time.Now()

	return nil
}

// Validate
// Summary: This is function which validate value of GetOperatorInput.
// output: (error) error object
func (i GetOperatorInput) Validate() error {
	if err := i.validate(); err != nil {
		logger.Set(nil).Warnf(err.Error())

		return err
	}
	return nil
}

// validate
// Summary: This is function which validate value of GetOperatorInput.
// output: (error) error object
func (i GetOperatorInput) validate() error {
	// Corporate number is 13 digits, verified only if not nil
	return validation.ValidateStruct(&i,
		validation.Field(
			&i.OperatorID,
			validation.By(common.StringUUIDValid),
		),
		validation.Field(
			&i.OpenOperatorID,
			validation.When(i.OpenOperatorID != nil, validation.RuneLength(1, 256)),
		),
	)
}
