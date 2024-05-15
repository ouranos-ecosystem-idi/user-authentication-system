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

// PlantModel
// Summary: This is structure which defines PlantModel.
// Service: Dataspace
// Router: [PUT] /api/v1/authInfo?dataTarget=plant
// Usage: output
type PlantModel struct {
	PlantID        uuid.UUID      `json:"plantId"`
	OperatorID     uuid.UUID      `json:"operatorId"`
	PlantName      string         `json:"plantName"`
	PlantAddress   string         `json:"plantAddress"`
	OpenPlantID    string         `json:"openPlantId"`
	PlantAttribute PlantAttribute `json:"plantAttribute"`
}

// PlantAttribute
// Summary: This is structure which defines PlantAttribute.
type PlantAttribute struct {
	GlobalPlantID *string `json:"globalPlantId"`
}

// PlantEntityModel
// Summary: This is structure which defines PlantEntityModel.
// DBName: plants
type PlantEntityModel struct {
	PlantID       uuid.UUID      `json:"plantId" gorm:"type:uuid"`
	OperatorID    uuid.UUID      `json:"operatorId" gorm:"type:uuid;not null"`
	PlantName     string         `json:"plantName" gorm:"type:string"`
	PlantAddress  string         `json:"plantAddress" gorm:"type:string"`
	OpenPlantID   string         `json:"openPlantId" gorm:"type:string"`
	GlobalPlantID *string        `json:"globalPlantId" gorm:"type:string"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"<-:create "`
	CreatedUserID string         `json:"createdUserId" gorm:"type:varchar(256);not null; <-:create"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	UpdatedUserID string         `json:"updatedUserId" gorm:"type:varchar(256);not null"`
}

// PlantModels
// Summary: This is structure which defines PlantModels.
// Service: Dataspace
// Router: [GET] /api/v1/authInfo?dataTarget=plant
// Usage: output
type PlantModels []PlantModel

// PlantEntityModels
// Summary: This is structure which defines PlantEntityModels.
type PlantEntityModels []PlantEntityModel

// GetPlantModel
// Summary: This is structure which defines GetPlantModel.
// Service: Dataspace
// Router: [GET] /api/v1/authInfo?dataTarget=plant
// Usage: input
type GetPlantModel struct {
	OperatorID uuid.UUID
}

// NewPlantEntityModel
// Summary: This is function which make PlantEntityModel.
// input: operatorID(uuid.UUID) UUID of operatorID
// input: plantName(string) value of plantName
// input: plantAddress(string) value of plantAddress
// input: openPlantID(string) value of openPlantID
// input: globalPlantID(string) pointer of globalPlantID
// output: (PlantEntityModel) PlantEntityModel object
func NewPlantEntityModel(
	operatorID uuid.UUID,
	plantName string,
	plantAddress string,
	openPlantID string,
	globalPlantID *string,
) PlantEntityModel {
	t := time.Now()
	e := PlantEntityModel{
		PlantID:       uuid.New(),
		OperatorID:    operatorID,
		PlantName:     plantName,
		PlantAddress:  plantAddress,
		OpenPlantID:   openPlantID,
		GlobalPlantID: globalPlantID,
		CreatedAt:     t,
		DeletedAt:     gorm.DeletedAt{},
		CreatedUserID: operatorID.String(),
		UpdatedAt:     t,
		UpdatedUserID: operatorID.String(),
	}
	return e
}

// PutPlantInput
// Summary: This is structure which defines PutPlantInput.
// Service: Dataspace
// Router: [PUT] /api/v1/authInfo?dataTarget=plant
// Usage: input
type PutPlantInput struct {
	PlantID             *string              `json:"plantId"`
	OperatorID          string               `json:"operatorId"`
	PlantName           string               `json:"plantName"`
	PlantAddress        string               `json:"plantAddress"`
	OpenPlantID         *string              `json:"openPlantId"`
	PlantAttributeInput *PlantAttributeInput `json:"plantAttribute"`
}

// PlantAttributeInput
// Summary: This is structure which defines PlantAttributeInput.
type PlantAttributeInput struct {
	GlobalPlantID *string `json:"globalPlantId"`
}

// Validate
// Summary: This is function which validate value of PutPlantInput.
// output: (error) error object
func (i PutPlantInput) Validate() error {
	if err := i.validate(); err != nil {
		logger.Set(nil).Warnf(err.Error())

		return err
	}

	return nil
}

// validate
// Summary: This is function which validate value of PutPlantInput.
// output: (error) error object
func (i PutPlantInput) validate() error {
	errors := []error{}
	err := validation.ValidateStruct(&i,
		validation.Field(
			&i.PlantID,
			validation.By(common.StringPtrNilOrUUIDValid),
		),
		validation.Field(
			&i.OperatorID,
			validation.By(common.StringUUIDValid),
		),
		validation.Field(
			&i.PlantName,
			validation.Required,
			validation.RuneLength(1, 256),
		),
		validation.Field(
			&i.PlantAddress,
			validation.Required,
			validation.RuneLength(1, 256),
		),
		validation.Field(
			&i.OpenPlantID,
			validation.RuneLength(0, 26),
			validation.NotNil,
			validation.By(common.StringPtrLast6CharsNumeric),
		),
		validation.Field(
			&i.PlantAttributeInput,
			validation.Required,
		),
	)

	if err != nil {
		errors = append(errors, err)
	}

	var attributeErr error
	if i.PlantAttributeInput != nil {
		if attributeErr = i.PlantAttributeInput.validate(); attributeErr != nil {
			attributeErr = fmt.Errorf("plantAttribute: (%v)", attributeErr)
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
// Summary: This is function which validate value of PlantAttributeInput.
// output: (error) error object
func (i PlantAttributeInput) validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(
			&i.GlobalPlantID,
			validation.RuneLength(0, 256),
		),
	)
}

// ToModel
// Summary: This is function which convert PutPlantInput to PlantModel.
// output: (PlantModel) PlantModel object
// output: (error) error object
func (i PutPlantInput) ToModel() (PlantModel, error) {
	var m PlantModel

	if i.PlantID != nil {
		plantID, err := uuid.Parse(*i.PlantID)
		if err != nil {
			logger.Set(nil).Warnf(err.Error())

			return PlantModel{}, fmt.Errorf(common.InvalidUUIDError("plantId"))
		}
		m.PlantID = plantID
	}

	operatorID, err := uuid.Parse(i.OperatorID)
	if err != nil {
		logger.Set(nil).Warnf(err.Error())

		return PlantModel{}, fmt.Errorf(common.InvalidUUIDError("operatorId"))
	}
	m.OperatorID = operatorID

	m.PlantName = i.PlantName
	m.PlantAddress = i.PlantAddress
	m.OpenPlantID = *i.OpenPlantID
	PlantAttribute := PlantAttribute{
		GlobalPlantID: i.PlantAttributeInput.GlobalPlantID,
	}
	m.PlantAttribute = PlantAttribute

	return m, nil
}

// Update
// Summary: This is function which update value of Plant.
// input: operatorID(uuid.UUID) UUID of operatorID
// input: plantName(string) value of plantName
// input: plantAddress(string) value of plantAddress
// input: openPlantID(string) value of openPlantID
// input: globalPlantID(string) pointer of globalPlantID
func (e *PlantEntityModel) Update(
	operatorID uuid.UUID,
	plantName string,
	plantAddress string,
	openPlantID string,
	globalPlantID *string,
) {
	e.PlantName = plantName
	e.PlantAddress = plantAddress
	e.OpenPlantID = openPlantID
	e.GlobalPlantID = globalPlantID
	e.UpdatedAt = time.Now()
	e.UpdatedUserID = operatorID.String()
}

func (e PlantEntityModel) ToModel() PlantModel {
	plantAttribute := PlantAttribute{
		GlobalPlantID: e.GlobalPlantID,
	}
	return PlantModel{
		PlantID:        e.PlantID,
		OperatorID:     e.OperatorID,
		PlantName:      e.PlantName,
		PlantAddress:   e.PlantAddress,
		OpenPlantID:    e.OpenPlantID,
		PlantAttribute: plantAttribute,
	}
}

// Update
// Summary: This is function which convert PlantEntityModels to array of PlantModel.
// output: globalPlantID(PlantModel) Array of PlantModel
func (es PlantEntityModels) ToModels() []PlantModel {
	ms := []PlantModel{}
	for _, e := range es {
		ms = append(ms, e.ToModel())
	}
	return ms
}
