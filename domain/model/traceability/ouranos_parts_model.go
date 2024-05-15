package traceability

// NOTE: This file is used only for dataReset

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PartsModel struct {
	TraceID            uuid.UUID           `json:"traceId"`
	OperatorID         uuid.UUID           `json:"operatorId"`
	PlantID            *uuid.UUID          `json:"plantId"`
	PartsName          string              `json:"partsName"`
	SupportPartsName   *string             `json:"supportPartsName"`
	TerminatedFlag     bool                `json:"terminatedFlag"`
	AmountRequired     *float64            `json:"amountRequired"`
	AmountRequiredUnit *AmountRequiredUnit `json:"amountRequiredUnit"`
}

type PartsModels []*PartsModel

type AmountRequiredUnit string

type PartsModelEntity struct {
	TraceID            uuid.UUID      `json:"traceId" gorm:"type:uuid"`
	OperatorID         uuid.UUID      `json:"operatorId" gorm:"type:uuid;not null"`
	PlantID            uuid.UUID      `json:"plantId" gorm:"type:uuid;not null"`
	PartsName          string         `json:"partsName" gorm:"type:string;not null"`
	SupportPartsName   *string        `json:"supportPartsName" gorm:"type:string"`
	TerminatedFlag     bool           `json:"terminatedFlag" gorm:"type:boolean;not null"`
	AmountRequired     *float64       `json:"amountRequired" gorm:"type:float(64)"`
	AmountRequiredUnit *string        `json:"amountRequiredUnit" gorm:"type:string"`
	DeletedAt          gorm.DeletedAt `json:"deletedAt"`
	CreatedAt          time.Time      `json:"createdAt" gorm:"<-:create "`
	CreatedUserID      string         `json:"createdUserId" gorm:"type:varchar(256);not null; <-:create"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	UpdatedUserID      string         `json:"updatedUserId" gorm:"type:varchar(256);not null"`
}

type PartsModelEntities []PartsModelEntity

type GetPartsModel struct {
	OperatorID string
	TraceID    *string
	PartsName  *string
	PlantID    *string
	ParentFlag *bool
	Limit      int
	After      *uuid.UUID
}
