package traceability

// NOTE: This file is used only for dataReset

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TradeEntityModel struct {
	TradeID              *uuid.UUID     `json:"tradeId" gorm:"type:uuid"`
	DownstreamOperatorID uuid.UUID      `json:"downstreamOperatorId" gorm:"type:uuid;not null"`
	UpstreamOperatorID   *uuid.UUID     `json:"upstreamOperatorId" gorm:"type:uuid;not null"`
	DownstreamTraceID    uuid.UUID      `json:"downstreamTraceId" gorm:"type:uuid;not null"`
	UpstreamTraceID      *uuid.UUID     `json:"upstreamTraceId" gorm:"type:uuid"`
	TradeDate            *string        `json:"tradeDate" gorm:"type:string"`
	DeletedAt            gorm.DeletedAt `json:"deletedAt"`
	CreatedAt            time.Time      `json:"createdAt" gorm:"<-:create "`
	CreatedUserID        string         `json:"createdUserId" gorm:"type:varchar(256);not null; <-:create"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	UpdatedUserID        string         `json:"updatedUserId" gorm:"type:varchar(256);not null"`
}

type TradeEntityModels []TradeEntityModel
