package traceability

// NOTE: This file is used only for dataReset

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CfpEntityModel struct {
	CfpID              *uuid.UUID     `json:"cfpId" gorm:"type:uuid"`
	TraceID            uuid.UUID      `json:"traceId" gorm:"type:uuid;not null"`
	GhgEmission        *float64       `json:"ghgEmission" gorm:"type:number"`
	GhgDeclaredUnit    string         `json:"ghgDeclaredUnit" gorm:"type:string"`
	CfpCertificateList []string       `json:"cfpCertificateList" gorm:"-"`
	CfpType            string         `json:"cfpType" gorm:"type:string"`
	DqrType            string         `json:"dqrType" gorm:"type:string"`
	TeR                *float64       `json:"TeR" gorm:"type:number"`
	GeR                *float64       `json:"GeR" gorm:"type:number"`
	TiR                *float64       `json:"TiR" gorm:"type:number"`
	DeletedAt          gorm.DeletedAt `json:"deletedAt"`
	CreatedAt          time.Time      `json:"createdAt" gorm:"<-:create "`
	CreatedUserID      string         `json:"createdUserId" gorm:"type:varchar(256);not null; <-:create"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	UpdatedUserID      string         `json:"updatedUserId" gorm:"type:varchar(256);not null"`
}
