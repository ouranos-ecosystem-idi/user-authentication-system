package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/presentation/http/echo/handler"
	f "authenticator-backend/test/fixtures"
	mocks "authenticator-backend/test/mock"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// /////////////////////////////////////////////////////////////////////////////////
// GET /api/v1/authInfo/operator テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 400: バリデーションエラー：operatorIdのフォーマットエラーの場合
// [x] 1-2. 500: システムエラー: 事業所取得失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_GetPlant_Abnormal(tt *testing.T) {
	var method = "GET"
	var endPoint = "/api/v1/authInfo"
	var dataTarget = "plant"

	tests := []struct {
		name              string
		input             string
		modifyQueryParams func(q url.Values)
		receive           error
		expectError       string
	}{
		{
			name:  "1-1. 400: バリデーションエラー：operatorIdのフォーマットエラーの場合",
			input: "invalidUUID",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
			},
			expectError: "code=400, message={[auth] BadRequest Invalid or expired token",
		},
		{
			name:  "1-2. 500: システムエラー: 事業所取得失敗の場合",
			input: f.OperatorId,
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
			},
			receive:     fmt.Errorf("Access Error"),
			expectError: "code=500, message={[auth] InternalServerError Unexpected error occurred",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			operatorUUID, _ := uuid.Parse(f.OperatorId)
			input := traceability.GetPlantModel{
				OperatorID: operatorUUID,
			}

			q := make(url.Values)
			test.modifyQueryParams(q)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)
			c.Set("operatorID", test.input)

			plantUsecase := new(mocks.IPlantUsecase)
			plantHandler := handler.NewPlantHandler(plantUsecase)
			plantUsecase.On("ListPlants", input).Return(nil, test.receive)

			err := plantHandler.GetPlant(c)
			if assert.Error(t, err) {
				assert.ErrorContains(t, err, test.expectError)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// GET /api/v1/authInfo/plant 正常系テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 200: 正常系
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_GetPlant_Normal(tt *testing.T) {
	var method = "PUT"
	var endPoint = "/api/v1/authInfo/plant"
	var dataTarget = "plant"

	tests := []struct {
		name              string
		modifyQueryParams func(q url.Values)
		expectStatus      int
	}{
		{
			name: "1-1. 200: 正常系",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
			},
			expectStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			operatorUUID, _ := uuid.Parse(f.OperatorId)
			input := traceability.GetPlantModel{
				OperatorID: operatorUUID,
			}

			q := make(url.Values)
			test.modifyQueryParams(q)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)
			c.Set("operatorID", f.OperatorId)

			plantUsecase := new(mocks.IPlantUsecase)
			plantHandler := handler.NewPlantHandler(plantUsecase)
			plantUsecase.On("ListPlants", input).Return(nil, nil)

			if assert.NoError(t, plantHandler.GetPlant(c)) {
				assert.Equal(t, test.expectStatus, rec.Code)
				plantUsecase.AssertExpectations(t)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// PUT /api/v1/authInfo/plant 正常系テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 201: 正常系：新規作成
// [x] 1-2. 201: 正常系：globalPlantId値あり
// [x] 1-3. 201: 正常系：globalPlantId空文字指定
// [x] 1-4. 201: 正常系：globalPlantIdなし、空Object
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_PutPlant_Normal(tt *testing.T) {
	var method = "PUT"
	var endPoint = "/api/v1/authInfo/plant"
	var dataTarget = "plant"

	tests := []struct {
		name         string
		inputFunc    func() traceability.PutPlantInput
		expectStatus int
	}{
		{
			name: "1-1. 201: 正常系：新規作成",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantID = nil
				putPlantInput.PlantAttributeInput.GlobalPlantID = nil
				return putPlantInput
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "1-2. 201: 正常系：globalPlantId値あり",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantID = nil
				return putPlantInput
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "1-3. 201: 正常系：globalPlantId空文字指定",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantID = nil
				putPlantInput.PlantAttributeInput.GlobalPlantID = common.StringPtr("")
				return putPlantInput
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "1-4. 201: 正常系：globalPlantIdなし、空Object",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantID = nil
				putPlantInput.PlantAttributeInput = &traceability.PlantAttributeInput{}
				return putPlantInput
			},
			expectStatus: http.StatusCreated,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			inputJSON, _ := json.Marshal(test.inputFunc())

			putPlantModel, _ := test.inputFunc().ToModel()

			q := make(url.Values)
			q.Set("dataTarget", dataTarget)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)
			c.Set("operatorID", f.OperatorId)

			plantUsecase := new(mocks.IPlantUsecase)
			plantHandler := handler.NewPlantHandler(plantUsecase)
			plantUsecase.On("PutPlant", putPlantModel).Return(putPlantModel, nil)

			err := plantHandler.PutPlant(c)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectStatus, rec.Code)
				plantUsecase.AssertExpectations(t)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// PUT /api/v1/authInfo/plant 異常系テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 400: バリデーションエラー：operatorIdが含まれない場合
// [x] 2-2. 400: バリデーションエラー：operatorIdがUUID形式でない場合
// [x] 2-3. 400: バリデーションエラー：plantNameが含まれない場合
// [x] 2-4. 400: バリデーションエラー：plantNameが257文字以上の場合
// [x] 2-5. 400: バリデーションエラー：plantNameがstring形式でない場合
// [x] 2-6. 400: バリデーションエラー：plantIdの値がUUID形式ではない場合
// [x] 2-7. 400: バリデーションエラー：plantAddress含まれない場合
// [x] 2-8. 400: バリデーションエラー：plantAddressの値が257文字以上の場合
// [x] 2-9. 400: バリデーションエラー：plantAddressがstring形式でない場合
// [x] 2-10. 400: バリデーションエラー：openPlantIdが含まれない場合
// [x] 2-11. 400: バリデーションエラー：openPlantIdの値が27文字以上の場合
// [x] 2-12. 400: バリデーションエラー：openPlantIdの値の末尾6文字が数字でない場合
// [x] 2-13. 400: バリデーションエラー：plantAttributeが含まれない場合
// [x] 2-14. 400: バリデーションエラー：globalPlantIdの値が257文字以上の場合
// [x] 2-15. 400: バリデーションエラー：operatorIdとplantNameが含まれない場合
// [x] 2-16. 400: バリデーションエラー：operatorIDがjwtのoperatorIdと一致しない場合
// [x] 2-17. 400: バリデーションエラー：operatorIdがUUID形式でない
// [x] 2-18. 400: バリデーションエラー：operatorIdがstring形式でない
// [x] 2-19. 400: バリデーションエラー：GlobalPlantIDがUUID形式でない
// [x] 2-20. 500: システムエラー：更新エラーの場合
// [x] 2-21. 500: システムエラー：更新エラーの場合
// [x] 2-22. 400: バリデーションエラー：plantNameがstring形式でない場合
// [x] 2-23. 400: バリデーションエラー：plantNameが空文字の場合
// [x] 2-24. 400: バリデーションエラー：plantAddressがstring形式でない場合
// [x] 2-25. 400: バリデーションエラー：plantAddressが空文字の場合
// [x] 2-26. 400: バリデーションエラー：2-3と2-5が同時に発生する場合
// [x] 2-27. 400: バリデーションエラー：2-3と2-12が同時に発生する場合(operatorIdとplantAttributeが含まれない)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_PutPlant_Abnormal(tt *testing.T) {
	var method = "PUT"
	var endPoint = "/api/v1/authInfo/plant"
	var dataTarget = "plant"

	tests := []struct {
		name              string
		invalidOperatorId string
		inputFunc         func() traceability.PutPlantInput
		invalidInputFunc  func() interface{}
		receive           error
		expectError       string
		expectStatus      int
	}{
		{
			name: "2-1. 400: バリデーションエラー：operatorIdが含まれない場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.OperatorID = ""
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, operatorId: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-2. 400: バリデーションエラー：operatorIdの値がUUID形式でない場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.OperatorID = f.InvalidUUID
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, operatorId: invalid UUID.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-3. 400: バリデーションエラー：plantNameが含まれない場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantName = ""
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantName: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-4. 400: バリデーションエラー：plantNameが257文字以上の場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantName = "Jv0ceJYX9Pa9zQTtlYPQLqLyUUhhNZ5EQCL2JDj9jLfrrgFK8MzV7zkaPvVj1wVtq5ESQGAbXrOhElxsVJzBjSxMBhwOa7hJwBrEkJjmYV8njJma2Zq6OZ7z9lXXh3xt6rYY0mYLLWpPGorQTOSY4XOkvOHfcmusmBl8OaFWjrAIUo9XwYfN2wVF4bKS32uD5vfwAzU5mhWCNwlZqABU9skfSQW9aMmCxbPkFiTq3P9hN9x4FR4m2SqB1AMLbNGu4"
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantName: the length must be between 1 and 256.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-6. 400: バリデーションエラー：plantIdの値がUUID形式ではない場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantID = common.StringPtr(f.InvalidUUID)
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantId: invalid UUID.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-7. 400: バリデーションエラー：plantAddressが含まれない場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantAddress = ""
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantAddress: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-8. 400: バリデーションエラー：plantAddressの値が257文字以上の場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantAddress = "Jv0ceJYX9Pa9zQTtlYPQLqLyUUhhNZ5EQCL2JDj9jLfrrgFK8MzV7zkaPvVj1wVtq5ESQGAbXrOhElxsVJzBjSxMBhwOa7hJwBrEkJjmYV8njJma2Zq6OZ7z9lXXh3xt6rYY0mYLLWpPGorQTOSY4XOkvOHfcmusmBl8OaFWjrAIUo9XwYfN2wVF4bKS32uD5vfwAzU5mhWCNwlZqABU9skfSQW9aMmCxbPkFiTq3P9hN9x4FR4m2SqB1AMLbNGu4"
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantAddress: the length must be between 1 and 256.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-10. 400: バリデーションエラー：openPlantIdが含まれない場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.OpenPlantID = nil
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, openPlantId: is required.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-11. 400: バリデーションエラー：openPlantIdの値が27文字以上の場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.OpenPlantID = common.StringPtr("Jv0ceJYX9Pa9zQTtlYPQLqLyUUh")
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, openPlantId: the length must be no more than 26.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-12. 400: バリデーションエラー：openPlantIdの値の末尾6文字が数字でない場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.OpenPlantID = common.StringPtr("xxxxx-12345x")
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, openPlantId: the last 6 digits must always be numeric.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-13. 400: バリデーションエラー：plantAttributeが含まれない場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantAttributeInput = nil
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantAttribute: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-14. 400: バリデーションエラー：globalPlantIdの値が257文字以上の場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantAttributeInput.GlobalPlantID = common.StringPtr("Jv0ceJYX9Pa9zQTtlYPQLqLyUUhhNZ5EQCL2JDj9jLfrrgFK8MzV7zkaPvVj1wVtq5ESQGAbXrOhElxsVJzBjSxMBhwOa7hJwBrEkJjmYV8njJma2Zq6OZ7z9lXXh3xt6rYY0mYLLWpPGorQTOSY4XOkvOHfcmusmBl8OaFWjrAIUo9XwYfN2wVF4bKS32uD5vfwAzU5mhWCNwlZqABU9skfSQW9aMmCxbPkFiTq3P9hN9x4FR4m2SqB1AMLbNGu4")
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantAttribute: (globalPlantId: the length must be no more than 256.).",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-15. 400: バリデーションエラー：operatorIdとplantNameが含まれない場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.OperatorID = ""
				putPlantInput.PlantName = ""
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, operatorId: cannot be blank; plantName: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-16. 400: バリデーションエラー：operatorIDがjwtのoperatorIdと一致しない場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.OperatorID = "80762b76-cf76-4485-9a99-cbe609c677c8"
				return putPlantInput
			},
			expectError:  "code=403, message={[auth] AccessDenied You do not have the necessary privileges",
			expectStatus: http.StatusForbidden,
		},
		{
			name: "2-17. 400: バリデーションエラー：operatorIdがUUID形式でない",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.OperatorID = f.InvalidUUID
				return putPlantInput
			},
			invalidOperatorId: "invalidOp",
			expectError:       "code=400, message={[auth] BadRequest Invalid or expired token",
			expectStatus:      http.StatusBadRequest,
		},
		{
			name: "2-18. 400: バリデーションエラー：operatorIdがstring形式でない",
			invalidInputFunc: func() interface{} {
				putPlantInterface := f.NewPutPlantInterface()
				putPlantInterface.(map[string]interface{})["operatorId"] = 1
				return putPlantInterface
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, operatorId: Unmarshal type error: expected=string, got=number.",
			expectStatus: http.StatusBadRequest,
		},
		// デッドコードのため到達不可
		// {
		// 	name: "2-19. 400: バリデーションエラー：GlobalPlantIDがUUID形式でない場合",
		// 	modifyInput: func(i *traceability.PutPlantInput) {
		// 		i.PlantID = common.StringPtr("invalidUUID")
		// 		i.PlantAttributeInput = &traceability.PlantAttributeInput{
		// 		    GlobalPlantID: nil,
		// 		}
		// 	},
		// 	receive:     common.NewCustomError(common.CustomErrorCode500, common.Err500Unexpected, nil, common.HTTPErrorSourceAuth),
		// 	expectError: "code=400, message={[auth] BadRequest Validation failed",
		// },
		{
			name:         "2-20. 500: システムエラー：更新エラーの場合",
			inputFunc:    func() traceability.PutPlantInput { return f.NewPutPlantInput() },
			receive:      common.NewCustomError(common.CustomErrorCode500, common.Err500Unexpected, nil, common.HTTPErrorSourceAuth),
			expectError:  "code=500, message={[auth] InternalServerError Unexpected error occurred",
			expectStatus: http.StatusInternalServerError,
		},
		{
			name:         "2-21. 500: システムエラー：更新エラーの場合",
			inputFunc:    func() traceability.PutPlantInput { return f.NewPutPlantInput() },
			receive:      fmt.Errorf("Access Error"),
			expectError:  "code=500, message={[auth] InternalServerError Unexpected error occurred",
			expectStatus: http.StatusInternalServerError,
		},
		{
			name: "2-22. 400: バリデーションエラー：plantNameがstring形式でない",
			invalidInputFunc: func() interface{} {
				putPlantInterface := f.NewPutPlantInterface()
				putPlantInterface.(map[string]interface{})["plantName"] = 1
				return putPlantInterface
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantName: Unmarshal type error: expected=string, got=number.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-23. 400: バリデーションエラー：plantNameが空文字の場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantName = ""
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantName: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-24. 400: バリデーションエラー：plantAddressがstring形式でない",
			invalidInputFunc: func() interface{} {
				putPlantInterface := f.NewPutPlantInterface()
				putPlantInterface.(map[string]interface{})["plantAddress"] = 1
				return putPlantInterface
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantAddress: Unmarshal type error: expected=string, got=number.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-25. 400: バリデーションエラー：plantAddressが空文字の場合",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.PlantAddress = ""
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, plantAddress: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-26. 400: バリデーションエラー：2-3と2-5が同時に発生する場合 (operatorIdとplantNameが含まれない場合)",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.OperatorID = ""
				putPlantInput.PlantName = ""
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, operatorId: cannot be blank; plantName: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-27. 400: バリデーションエラー：2-3と2-12が同時に発生する場合 (operatorIdとplantAttributeが含まれない場合)",
			inputFunc: func() traceability.PutPlantInput {
				putPlantInput := f.NewPutPlantInput()
				putPlantInput.OperatorID = ""
				putPlantInput.PlantAttributeInput = nil
				return putPlantInput
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, operatorId: cannot be blank; plantAttribute: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var inputJSON []byte
			if test.invalidInputFunc != nil {
				inputJSON, _ = json.Marshal(test.invalidInputFunc())
			} else {
				inputJSON, _ = json.Marshal(test.inputFunc())
			}

			q := make(url.Values)
			q.Set("dataTarget", dataTarget)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint, strings.NewReader(string(inputJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)
			if test.invalidOperatorId != "" {
				c.Set("operatorID", test.invalidOperatorId)
			} else {
				c.Set("operatorID", f.OperatorId)
			}

			plantUsecase := new(mocks.IPlantUsecase)
			plantHandler := handler.NewPlantHandler(plantUsecase)
			plantUsecase.On("PutPlant", mock.Anything).Return(traceability.PlantModel{}, test.receive)

			err := plantHandler.PutPlant(c)
			e.HTTPErrorHandler(err, c)
			if assert.Error(t, err) {
				assert.Equal(t, test.expectStatus, rec.Code)
				assert.ErrorContains(t, err, test.expectError)
			}
		})
	}
}
