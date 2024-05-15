package fixtures

import (
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/domain/model/traceability"

	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
)

var (
	// 変数名の昇順
	ApiKey           = "36cfd2a8-9f45-0766-77d3-7098c1336a32"
	AssertMessage    = "比較対象の２つの値は定義順に関係なく、一致する必要があります。"
	Email            = "testaccount_user122@example.com"
	GlobalOperatorId = "GlobalOperatorId"
	GlobalOperatorID = "GlobalOperatorId"
	GlobalPlantID    = "GlobalPlantId"
	GlobalPlantId    = "GlobalPlantId"
	InvalidIPAddress = "invalid_ip_address"
	InvalidUUID      = "invalid_uuid"
	InvalidEnum      = "invalid_enum"
	IpAddress        = "127.0.0.1"
	OpenOperatorID   = "AAAA-123456"
	OpenPlantID      = "AAAA-BBBB"
	OperatorAddress  = "東京都"
	OperatorID       = "e03cc699-7234-31ed-86be-cc18c92208e5"
	OperatorId       = "e03cc699-7234-31ed-86be-cc18c92208e5"
	OperatorName     = "A株式会社"
	PlantAddress     = "東京都"
	PlantID          = "eedf264e-cace-4414-8bd3-e10ce1c090e0"
	PlantId          = "eedf264e-cace-4414-8bd3-e10ce1c090e0"
	PlantName        = "A工場"
	Token            = "valid_token" // 実際には無効。有効なtokenを定義することはできないので、ダミー
	UID              = "uid"

	// Input
	PutOperatorInput = traceability.PutOperatorInput{
		OperatorID:      OperatorId,
		OperatorName:    OperatorName,
		OperatorAddress: OperatorAddress,
		OpenOperatorID:  OpenOperatorID,
		OperatorAttributeInput: &traceability.OperatorAttributeInput{
			GlobalOperatorID: &GlobalOperatorId,
		},
	}
	PutPlantInput = traceability.PutPlantInput{
		OperatorID:   OperatorId,
		PlantName:    PlantName,
		PlantAddress: PlantAddress,
		OpenPlantID:  &OpenOperatorID,
		PlantAttributeInput: &traceability.PlantAttributeInput{
			GlobalPlantID: &GlobalPlantId,
		},
	}

	// Claims
	Claims = authentication.Claims{
		OperatorID: OperatorId,
		Token: auth.Token{
			UID: UID,
		},
	}
)

func NewGetOperatorInput(inputOpenOperatorId bool) traceability.GetOperatorInput {
	if inputOpenOperatorId {
		return traceability.GetOperatorInput{
			OpenOperatorID: &OpenOperatorID,
		}
	} else {
		return traceability.GetOperatorInput{
			OperatorID: OperatorId,
		}
	}
}

func NewGetOperatorsInput() traceability.GetOperatorsInput {
	return traceability.GetOperatorsInput{
		OperatorIDs: []uuid.UUID{uuid.MustParse(OperatorId)},
	}
}

func NewOperatorModel(operatorId string, openOperatorId string) traceability.OperatorModel {
	return traceability.OperatorModel{
		OperatorID:      uuid.MustParse(operatorId),
		OperatorName:    OperatorName,
		OperatorAddress: OperatorAddress,
		OpenOperatorID:  openOperatorId,
		OperatorAttribute: traceability.OperatorAttribute{
			GlobalOperatorID: &GlobalOperatorId,
		},
	}
}

func NewGetPlantModel() traceability.GetPlantModel {
	return traceability.GetPlantModel{
		OperatorID: uuid.MustParse(PlantId),
	}
}

func NewPutPlantInput() traceability.PutPlantInput {
	return traceability.PutPlantInput{
		OperatorID:   OperatorId,
		PlantName:    PlantName,
		PlantAddress: PlantAddress,
		OpenPlantID:  &OpenOperatorID,
		PlantAttributeInput: &traceability.PlantAttributeInput{
			GlobalPlantID: &GlobalPlantId,
		},
	}
}

func NewPlantModel(inputPlantId bool) traceability.PlantModel {
	if inputPlantId {
		return traceability.PlantModel{
			PlantID:      uuid.MustParse(PlantId),
			OperatorID:   uuid.MustParse(OperatorId),
			PlantName:    PlantName,
			PlantAddress: PlantAddress,
			OpenPlantID:  OpenPlantID,
			PlantAttribute: traceability.PlantAttribute{
				GlobalPlantID: &GlobalPlantId,
			},
		}
	} else {
		return traceability.PlantModel{
			PlantID:      uuid.Nil,
			OperatorID:   uuid.MustParse(OperatorId),
			PlantName:    PlantName,
			PlantAddress: PlantAddress,
			OpenPlantID:  OpenPlantID,
			PlantAttribute: traceability.PlantAttribute{
				GlobalPlantID: &GlobalPlantId,
			},
		}
	}
}
