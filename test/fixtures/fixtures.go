package fixtures

import (
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/usecase/input"

	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
)

var (
	// 変数名の昇順
	AccountPassword    = "123456"
	AccountPasswordNew = "1Aa@1Aa@1Aa@"
	ApiKey             = "36cfd2a8-9f45-0766-77d3-7098c1336a32"
	AssertMessage      = "比較対象の２つの値は定義順に関係なく、一致する必要があります。"
	Email              = "testaccount_user122@example.com"
	GlobalOperatorId   = "GlobalOperatorId"
	GlobalOperatorID   = "GlobalOperatorId"
	GlobalPlantID      = "GlobalPlantId"
	GlobalPlantId      = "GlobalPlantId"
	InvalidIPAddress   = "invalid_ip_address"
	InvalidUUID        = "invalid_uuid"
	InvalidEnum        = "invalid_enum"
	IpAddress          = "127.0.0.1"
	OpenOperatorID     = "AAAA-123456"
	OpenPlantID        = "AAAA-123456"
	OperatorAccountID  = "aaa@bbb.com"
	OperatorAddress    = "東京都"
	OperatorID         = "e03cc699-7234-31ed-86be-cc18c92208e5"
	OperatorId         = "e03cc699-7234-31ed-86be-cc18c92208e5"
	OperatorName       = "A株式会社"
	PlantAddress       = "東京都"
	PlantID            = "eedf264e-cace-4414-8bd3-e10ce1c090e0"
	PlantId            = "eedf264e-cace-4414-8bd3-e10ce1c090e0"
	PlantName          = "A工場"
	Token              = "valid_token" // 実際には無効。有効なtokenを定義することはできないので、ダミー
	UID                = "uid"
)

func NewClaims() authentication.Claims {
	return authentication.Claims{
		OperatorID: OperatorId,
		Token: auth.Token{
			UID: UID,
		},
	}
}

func NewGetOperatorInput() traceability.GetOperatorInput {
	return traceability.GetOperatorInput{
		OperatorID: OperatorId,
	}
}

func NewInputLoginParam() input.LoginParam {
	return input.LoginParam{
		OperatorAccountID: OperatorAccountID,
		AccountPassword:   AccountPassword,
	}
}

func NewInputRefreshParam() input.RefreshParam {
	return input.RefreshParam{
		RefreshToken: Token,
	}
}

func NewInputChangePasswordParam() input.ChangePasswordParam {
	return input.ChangePasswordParam{
		UID:         UID,
		NewPassword: authentication.Password(AccountPasswordNew),
	}
}

func NewInputVerifyTokenParam() input.VerifyTokenParam {
	return input.VerifyTokenParam{
		IDToken: Token,
	}
}

func NewInputVerifyIDTokenParam() input.VerifyIDTokenParam {
	return input.VerifyIDTokenParam{
		IDToken: Token,
	}
}

func NewInputVerifyAPIKeyParam() input.VerifyAPIKeyParam {
	return input.VerifyAPIKeyParam{
		IPAddress: IpAddress,
		APIKey:    ApiKey,
	}
}

func NewGetOperatorsInput() traceability.GetOperatorsInput {
	return traceability.GetOperatorsInput{
		OperatorIDs: []uuid.UUID{uuid.MustParse(OperatorId)},
	}
}

func NewOperatorModel() traceability.OperatorModel {
	return traceability.OperatorModel{
		OperatorID:      uuid.MustParse(OperatorId),
		OperatorName:    OperatorName,
		OperatorAddress: OperatorAddress,
		OpenOperatorID:  OpenOperatorID,
		OperatorAttribute: traceability.OperatorAttribute{
			GlobalOperatorID: &GlobalOperatorID,
		},
	}
}

func NewOperatorModels(n int) []traceability.OperatorModel {
	operatorModels := make([]traceability.OperatorModel, n)
	for i := 0; i < n; i++ {
		operatorModels[i] = NewOperatorModel()
	}
	return operatorModels
}

func NewGetPlantModel() traceability.GetPlantModel {
	return traceability.GetPlantModel{
		OperatorID: uuid.MustParse(PlantId),
	}
}

func NewPutOperatorInput() traceability.PutOperatorInput {
	return traceability.PutOperatorInput{
		OperatorID:      OperatorId,
		OperatorName:    OperatorName,
		OperatorAddress: OperatorAddress,
		OpenOperatorID:  OpenOperatorID,
		OperatorAttributeInput: &traceability.OperatorAttributeInput{
			GlobalOperatorID: &GlobalOperatorId,
		},
	}
}

func NewPutOperatorInterface() interface{} {
	return map[string]interface{}{
		"operatorId":      OperatorId,
		"operatorName":    OperatorName,
		"operatorAddress": OperatorAddress,
		"openOperatorId":  OpenOperatorID,
		"operatorAttribute": map[string]interface{}{
			"globalOperatorId": GlobalOperatorId,
		},
	}
}

func NewChangePasswordParam() input.ChangePasswordParam {
	return input.ChangePasswordParam{
		UID:         UID,
		NewPassword: authentication.Password(AccountPasswordNew),
	}
}

func NewChangePasswordInterface() interface{} {
	return map[string]interface{}{
		"uid":         UID,
		"newPassword": AccountPasswordNew,
	}
}

func NewLoginParam() input.LoginParam {
	return input.LoginParam{
		OperatorAccountID: Email,
		AccountPassword:   AccountPassword,
	}
}

func NewLoginInterface() interface{} {
	return map[string]interface{}{
		"operatorAccountID": Email,
		"accountPassword":   AccountPassword,
	}
}

func NewRefreshParam() input.RefreshParam {
	return input.RefreshParam{
		RefreshToken: Token,
	}
}

func NewRefreshInterface() interface{} {
	return map[string]interface{}{
		"refreshToken": Token,
	}
}

func NewPutPlantInput() traceability.PutPlantInput {
	return traceability.PutPlantInput{
		PlantID:      &PlantID,
		OperatorID:   OperatorID,
		PlantName:    PlantName,
		PlantAddress: PlantAddress,
		OpenPlantID:  &OpenPlantID,
		PlantAttributeInput: &traceability.PlantAttributeInput{
			GlobalPlantID: &GlobalPlantID,
		},
	}
}

func NewPutPlantInterface() interface{} {
	return map[string]interface{}{
		"plantId":      PlantId,
		"operatorId":   OperatorId,
		"plantName":    PlantName,
		"plantAddress": PlantAddress,
		"openPlantId":  OpenPlantID,
		"plantAttribute": map[string]interface{}{
			"globalPlantId": GlobalPlantID,
		},
	}
}

func NewPlantModel() traceability.PlantModel {
	return traceability.PlantModel{
		PlantID:      uuid.MustParse(PlantId),
		OperatorID:   uuid.MustParse(OperatorId),
		PlantName:    PlantName,
		PlantAddress: PlantAddress,
		OpenPlantID:  OpenPlantID,
		PlantAttribute: traceability.PlantAttribute{
			GlobalPlantID: &GlobalPlantID,
		},
	}
}
