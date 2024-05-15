package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/domain/model/traceability"
	f "authenticator-backend/test/fixtures"
	mocks "authenticator-backend/test/mock"
	"authenticator-backend/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// /////////////////////////////////////////////////////////////////////////////////
// Post reset テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 200: 正常系
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectUsecase_Reset(tt *testing.T) {

	apikeyOperator := authentication.APIKeyOperators{
		{
			APIKey:     "apikey",
			OperatorID: f.OperatorId,
		},
	}

	plantId := uuid.MustParse(f.PlantId)
	operatorId := uuid.MustParse(f.OperatorId)
	tradeId := uuid.MustParse("97a72868-63e3-43fb-9997-488af61d3be7")
	traceId := uuid.MustParse("d17833fe-22b7-4a4a-b097-bc3f2150c9a6")
	dqrValue := 0.1
	ghgEmission := 1.12345
	partsName := "PartsA-002123"
	supportPartsName := "modelA"
	ghgDeclaredUnit := "kgCO2e/liter"
	cfpCertificateList := []string{"https://www.example1.com/1"}
	cfpId := uuid.MustParse("d9a38406-cae2-4679-b052-15a75f5531f5")
	var amountRequiredUnit traceability.AmountRequiredUnit = "liter"
	parts := traceability.PartsModels{
		{
			OperatorID:         operatorId,
			TraceID:            traceId,
			PlantID:            &plantId,
			PartsName:          partsName,
			SupportPartsName:   &supportPartsName,
			TerminatedFlag:     false,
			AmountRequired:     nil,
			AmountRequiredUnit: &amountRequiredUnit,
		},
	}

	cfp := traceability.CfpEntityModel{
		CfpID:              &cfpId,
		TraceID:            traceId,
		GhgEmission:        &ghgEmission,
		GhgDeclaredUnit:    ghgDeclaredUnit,
		CfpCertificateList: cfpCertificateList,
		CfpType:            "preProduction",
		DqrType:            "preProcessing",
		TeR:                &dqrValue,
		GeR:                &dqrValue,
		TiR:                &dqrValue,
		DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt:     time.Now(),
		CreatedUserID: f.OperatorId,
		UpdatedAt:     time.Now(),
		UpdatedUserID: f.OperatorId,
	}

	trade := traceability.TradeEntityModels{
		{
			TradeID:              &tradeId,
			DownstreamOperatorID: operatorId,
			UpstreamOperatorID:   &operatorId,
			DownstreamTraceID:    traceId,
			UpstreamTraceID:      &traceId,
			TradeDate:            common.StringPtr("2023-09-25T14:30:00Z"),
			DeletedAt: gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			},
			CreatedAt:     time.Now(),
			CreatedUserID: f.OperatorId,
			UpdatedAt:     time.Now(),
			UpdatedUserID: f.OperatorId,
		},
	}

	plants := traceability.PlantEntityModels{
		{
			PlantID:       plantId,
			OperatorID:    operatorId,
			PlantName:     f.PlantName,
			PlantAddress:  f.PlantAddress,
			OpenPlantID:   f.OpenPlantID,
			GlobalPlantID: &f.GlobalPlantId,
			DeletedAt: gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			},
			CreatedAt:     time.Now(),
			CreatedUserID: f.OperatorId,
			UpdatedAt:     time.Now(),
			UpdatedUserID: f.OperatorId,
		},
	}
	tests := []struct {
		name           string
		input          string
		receive_apikey authentication.APIKeyOperators
	}{
		{
			name:           "1-1. 200: 正常系",
			input:          "apikey",
			receive_apikey: apikeyOperator,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				ouranosRepositoryMock := new(mocks.OuranosRepository)
				authRepositoryMock := new(mocks.AuthRepository)
				authRepositoryMock.On("ListAPIKeyOperators", mock.Anything).Return(test.receive_apikey, nil)
				ouranosRepositoryMock.On("ListPlantsByOperatorID", mock.Anything).Return(plants, nil)
				ouranosRepositoryMock.On("ListParts", mock.Anything).Return(parts, nil)
				ouranosRepositoryMock.On("GetCFPInformation", mock.Anything).Return(cfp, nil)
				ouranosRepositoryMock.On("ListTradesByOperatorID", mock.Anything).Return(trade, nil)
				ouranosRepositoryMock.On("DeleteCFPCertificateByCFPID", mock.Anything).Return(nil)
				ouranosRepositoryMock.On("DeleteCFPInformation", mock.Anything).Return(nil)
				ouranosRepositoryMock.On("DeletePartsStructure", mock.Anything).Return(nil)
				ouranosRepositoryMock.On("DeleteRequestStatusByTradeID", mock.Anything).Return(nil)
				ouranosRepositoryMock.On("DeleteTrade", mock.Anything).Return(nil)
				ouranosRepositoryMock.On("DeleteParts", mock.Anything).Return(nil)
				ouranosRepositoryMock.On("DeletePlant", mock.Anything).Return(nil)

				usecase := usecase.NewResetUsecase(ouranosRepositoryMock, authRepositoryMock)

				err := usecase.Reset(test.input)
				assert.NoError(t, err)
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Post reset テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 500: データ取得エラー
// [x] 2-2. 500: データ取得エラー
// [x] 2-3. 500: データ取得エラー
// [x] 2-4. 500: データ取得エラー
// [x] 2-5. 500: データ取得エラー
// [x] 2-6. 500: データ取得エラー
// [x] 2-7. 500: データ取得エラー
// [x] 2-8. 500: データ取得エラー
// [x] 2-9. 500: データ取得エラー
// [x] 2-10. 500: データ取得エラー
// [x] 2-11. 500: データ取得エラー
// [x] 2-12. 500: データ取得エラー
// [x] 2-13. 500: データ取得エラー(CFPなし)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectUsecase_Reset_Abnormal(tt *testing.T) {

	dsResGetError := fmt.Errorf("DB AccessError")

	apikeyOperator := authentication.APIKeyOperators{
		{
			APIKey:     "apikey",
			OperatorID: f.OperatorId,
		},
	}

	plantId := uuid.MustParse(f.PlantId)
	operatorId := uuid.MustParse(f.OperatorId)
	tradeId := uuid.MustParse("97a72868-63e3-43fb-9997-488af61d3be7")
	traceId := uuid.MustParse("d17833fe-22b7-4a4a-b097-bc3f2150c9a6")
	dqrValue := 0.1
	ghgEmission := 1.12345
	partsName := "PartsA-002123"
	supportPartsName := "modelA"
	ghgDeclaredUnit := "kgCO2e/liter"
	cfpCertificateList := []string{"https://www.example1.com/1"}
	cfpId := uuid.MustParse("d9a38406-cae2-4679-b052-15a75f5531f5")
	var amountRequiredUnit traceability.AmountRequiredUnit = "liter"
	parts := traceability.PartsModels{
		{
			OperatorID:         operatorId,
			TraceID:            traceId,
			PlantID:            &plantId,
			PartsName:          partsName,
			SupportPartsName:   &supportPartsName,
			TerminatedFlag:     false,
			AmountRequired:     nil,
			AmountRequiredUnit: &amountRequiredUnit,
		},
	}

	cfp := traceability.CfpEntityModel{
		CfpID:              &cfpId,
		TraceID:            traceId,
		GhgEmission:        &ghgEmission,
		GhgDeclaredUnit:    ghgDeclaredUnit,
		CfpCertificateList: cfpCertificateList,
		CfpType:            "preProduction",
		DqrType:            "preProcessing",
		TeR:                &dqrValue,
		GeR:                &dqrValue,
		TiR:                &dqrValue,
		DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt:     time.Now(),
		CreatedUserID: f.OperatorId,
		UpdatedAt:     time.Now(),
		UpdatedUserID: f.OperatorId,
	}

	trade := traceability.TradeEntityModels{
		{
			TradeID:              &tradeId,
			DownstreamOperatorID: operatorId,
			UpstreamOperatorID:   &operatorId,
			DownstreamTraceID:    traceId,
			UpstreamTraceID:      &traceId,
			TradeDate:            common.StringPtr("2023-09-25T14:30:00Z"),
			DeletedAt: gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			},
			CreatedAt:     time.Now(),
			CreatedUserID: f.OperatorId,
			UpdatedAt:     time.Now(),
			UpdatedUserID: f.OperatorId,
		},
	}

	plants := traceability.PlantEntityModels{
		{
			PlantID:       plantId,
			OperatorID:    operatorId,
			PlantName:     f.PlantName,
			PlantAddress:  f.PlantAddress,
			OpenPlantID:   f.OpenPlantID,
			GlobalPlantID: &f.GlobalPlantId,
			DeletedAt: gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			},
			CreatedAt:     time.Now(),
			CreatedUserID: f.OperatorId,
			UpdatedAt:     time.Now(),
			UpdatedUserID: f.OperatorId,
		},
	}
	tests := []struct {
		name           string
		input          string
		receive_apikey authentication.APIKeyOperators
		error1         error
		error2         error
		error3         error
		error4         error
		error5         error
		error6         error
		error7         error
		error8         error
		error9         error
		error10        error
		error11        error
		error12        error
	}{
		{
			name:           "2-1. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         dsResGetError,
			error2:         nil,
			error3:         nil,
			error4:         nil,
			error5:         nil,
			error6:         nil,
			error7:         nil,
			error8:         nil,
			error9:         nil,
			error10:        nil,
			error11:        nil,
			error12:        nil,
		},
		{
			name:           "2-2. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         dsResGetError,
			error3:         nil,
			error4:         nil,
			error5:         nil,
			error6:         nil,
			error7:         nil,
			error8:         nil,
			error9:         nil,
			error10:        nil,
			error11:        nil,
			error12:        nil,
		},
		{
			name:           "2-3. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         dsResGetError,
			error4:         nil,
			error5:         nil,
			error6:         nil,
			error7:         nil,
			error8:         nil,
			error9:         nil,
			error10:        nil,
			error11:        nil,
			error12:        nil,
		},
		{
			name:           "2-4. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         nil,
			error4:         dsResGetError,
			error5:         nil,
			error6:         nil,
			error7:         nil,
			error8:         nil,
			error9:         nil,
			error10:        nil,
			error11:        nil,
			error12:        nil,
		},
		{
			name:           "2-5. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         nil,
			error4:         nil,
			error5:         dsResGetError,
			error6:         nil,
			error7:         nil,
			error8:         nil,
			error9:         nil,
			error10:        nil,
			error11:        nil,
			error12:        nil,
		},
		{
			name:           "2-6. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         nil,
			error4:         nil,
			error5:         nil,
			error6:         dsResGetError,
			error7:         nil,
			error8:         nil,
			error9:         nil,
			error10:        nil,
			error11:        nil,
			error12:        nil,
		},
		{
			name:           "2-7. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         nil,
			error4:         nil,
			error5:         nil,
			error6:         nil,
			error7:         dsResGetError,
			error8:         nil,
			error9:         nil,
			error10:        nil,
			error11:        nil,
			error12:        nil,
		},
		{
			name:           "2-8. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         nil,
			error4:         nil,
			error5:         nil,
			error6:         nil,
			error7:         nil,
			error8:         dsResGetError,
			error9:         nil,
			error10:        nil,
			error11:        nil,
			error12:        nil,
		},
		{
			name:           "2-9. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         nil,
			error4:         nil,
			error5:         nil,
			error6:         nil,
			error7:         nil,
			error8:         nil,
			error9:         dsResGetError,
			error10:        nil,
			error11:        nil,
			error12:        nil,
		},
		{
			name:           "2-10. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         nil,
			error4:         nil,
			error5:         nil,
			error6:         nil,
			error7:         nil,
			error8:         nil,
			error9:         nil,
			error10:        dsResGetError,
			error11:        nil,
			error12:        nil,
		},
		{
			name:           "2-11. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         nil,
			error4:         nil,
			error5:         nil,
			error6:         nil,
			error7:         nil,
			error8:         nil,
			error9:         nil,
			error10:        nil,
			error11:        dsResGetError,
			error12:        nil,
		},
		{
			name:           "2-12. 500: データ取得エラー",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         nil,
			error4:         nil,
			error5:         nil,
			error6:         nil,
			error7:         nil,
			error8:         nil,
			error9:         nil,
			error10:        nil,
			error11:        nil,
			error12:        dsResGetError,
		},
		{
			name:           "2-13. 500: データ取得エラー(CFPなし)",
			input:          "apikey",
			receive_apikey: apikeyOperator,
			error1:         nil,
			error2:         nil,
			error3:         nil,
			error4:         gorm.ErrRecordNotFound,
			error5:         nil,
			error6:         nil,
			error7:         nil,
			error8:         nil,
			error9:         nil,
			error10:        nil,
			error11:        nil,
			error12:        dsResGetError,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				ouranosRepositoryMock := new(mocks.OuranosRepository)
				authRepositoryMock := new(mocks.AuthRepository)
				authRepositoryMock.On("ListAPIKeyOperators", mock.Anything).Return(test.receive_apikey, test.error1)
				ouranosRepositoryMock.On("ListPlantsByOperatorID", mock.Anything).Return(plants, test.error2)
				ouranosRepositoryMock.On("ListParts", mock.Anything).Return(parts, test.error3)
				ouranosRepositoryMock.On("GetCFPInformation", mock.Anything).Return(cfp, test.error4)
				ouranosRepositoryMock.On("ListTradesByOperatorID", mock.Anything).Return(trade, test.error5)
				ouranosRepositoryMock.On("DeleteCFPCertificateByCFPID", mock.Anything).Return(test.error6)
				ouranosRepositoryMock.On("DeleteCFPInformation", mock.Anything).Return(test.error7)
				ouranosRepositoryMock.On("DeletePartsStructure", mock.Anything).Return(test.error8)
				ouranosRepositoryMock.On("DeleteRequestStatusByTradeID", mock.Anything).Return(test.error9)
				ouranosRepositoryMock.On("DeleteTrade", mock.Anything).Return(test.error10)
				ouranosRepositoryMock.On("DeleteParts", mock.Anything).Return(test.error11)
				ouranosRepositoryMock.On("DeletePlant", mock.Anything).Return(test.error12)

				usecase := usecase.NewResetUsecase(ouranosRepositoryMock, authRepositoryMock)

				err := usecase.Reset(test.input)
				assert.Error(t, err)
			},
		)
	}
}
