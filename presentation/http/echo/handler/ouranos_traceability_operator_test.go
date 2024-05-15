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
// [x] 1-1. 200: 正常系(operatorId指定あり)
// [x] 1-2. 200: 正常系(operatorId指定なし)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_GetOperator_Normal(tt *testing.T) {
	var method = "GET"
	var endPoint = "/api/v1/authInfo"
	var dataTarget = "operator"

	tests := []struct {
		name              string
		modifyQueryParams func(q *url.Values)
		expectStatus      int
	}{
		{
			name: "1-1. 200: 正常系(operatorId指定あり)",
			modifyQueryParams: func(q *url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("openOperatorId", f.OpenOperatorID)
			},
			expectStatus: http.StatusOK,
		},
		{
			name: "1-2. 200: 正常系(operatorId指定なし)",
			modifyQueryParams: func(q *url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("openOperatorId", "")
			},
			expectStatus: http.StatusOK,
		},
		{
			name: "1-3. 200: 正常系(operatorIds指定あり)",
			modifyQueryParams: func(q *url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("operatorIds", uuid.New().String())
			},
			expectStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			q := make(url.Values)
			test.modifyQueryParams(&q)
			operatorID := f.OperatorId
			operatorModel := traceability.OperatorModel{}
			operatorModels := traceability.OperatorModels{}
			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)
			c.Set("operatorID", f.OperatorId)

			operatorUsecase := new(mocks.IOperatorUsecase)
			operatorHandler := handler.NewOperatorHandler(operatorUsecase)
			openOperatorID := q.Get("openOperatorId")
			input := traceability.GetOperatorInput{
				OperatorID: operatorID,
			}
			if openOperatorID != "" {
				input.OpenOperatorID = &openOperatorID
			}

			if test.name == "1-1. 200: 正常系(operatorId指定あり)" {
				operatorUsecase.On("GetOperator", input).Return(operatorModel, nil)
			} else if test.name == "1-2. 200: 正常系(operatorId指定なし)" {
				operatorUsecase.On("GetOperator", input).Return(operatorModel, nil)
			} else if test.name == "1-3. 200: 正常系(operatorIds指定あり)" {
				operatorUsecase.On("GetOperators", mock.Anything).Return(operatorModels, nil)
			}
			err := operatorHandler.GetOperator(c)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectStatus, rec.Code)
				operatorUsecase.AssertExpectations(t)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// GET /api/v1/authInfo/operator テストケース
// /////////////////////////////////////////////////////////////////////////////////
// No.2 事業者情報取得APIのテストケース
// [x] 1-1. 400: バリデーションエラー：OpenOperatorIDが257文字以上の場合
// No.17 事業者情報一覧取得APIのテストケース
// [x] 1-1. 400: バリデーションエラー: operatorIdsが含まれない場合
// [x] 1-2. 400: バリデーションエラー: operatorIdsがUUID形式ではない場合
// [x] 1-3. 400: バリデーションエラー: operatorIdsが複数指定されているうちの1つがUUID形式ではない場合
// [x] 1-4. 400: バリデーションエラー: operatorIdsが101個以上指定されている場合
// [x] 1-5. 400: バリデーションエラー: operatorIdsとopenOperatorIdsが同時に指定されている場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_GetOperator(tt *testing.T) {
	var method = "GET"
	var endPoint = "/api/v1/authInfo"
	var dataTarget = "operator"

	tests := []struct {
		name              string
		modifyQueryParams func(q url.Values)
		receive           error
		expectError       string
	}{
		{
			name: "1-1. 400: バリデーションエラー：OpenOperatorIDが257文字以上の場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("openOperatorId", "Jv0ceJYX9Pa9zQTtlYPQLqLyUUhhNZ5EQCL2JDj9jLfrrgFK8MzV7zkaPvVj1wVtq5ESQGAbXrOhElxsVJzBjSxMBhwOa7hJwBrEkJjmYV8njJma2Zq6OZ7z9lXXh3xt6rYY0mYLLWpPGorQTOSY4XOkvOHfcmusmBl8OaFWjrAIUo9XwYfN2wVF4bKS32uD5vfwAzU5mhWCNwlZqABU9skfSQW9aMmCxbPkFiTq3P9hN9x4FR4m2SqB1AMLbNGu4")
			},
			expectError: "code=400, message={[auth] BadRequest Invalid request parameters, openOperatorId: Unexpected query parameter",
		},
		{
			name: "1-2. 400: バリデーションエラー: operatorIdsが含まれない場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("operatorIds", "")
			},
			expectError: "code=400, message={[auth] BadRequest Invalid request parameters, operatorIds: Unexpected query parameter",
		},
		{
			name: "1-3. 400: バリデーションエラー: operatorIdsがUUID形式ではない場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("operatorIds", "hoge")
			},
			expectError: "code=400, message={[auth] BadRequest Invalid request parameters, operatorIds: Unexpected query parameter",
		},
		{
			name: "1-4. 400: バリデーションエラー: operatorIdsが複数指定されているうちの1つがUUID形式ではない場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
				operatorIds := []string{"hoge", uuid.New().String()}
				q.Set("operatorIds", strings.Join(operatorIds, ","))
			},
			expectError: "code=400, message={[auth] BadRequest Invalid request parameters, operatorIds: Unexpected query parameter",
		},
		{
			name: "1-5. 400: バリデーションエラー: operatorIdsが101個以上指定されている場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
				operatorIds := []string{}
				for i := 0; i < 101; i++ {
					operatorIds = append(operatorIds, uuid.New().String())
				}
				q.Set("operatorIds", strings.Join(operatorIds, ","))
			},
			expectError: "code=400, message={[auth] BadRequest Invalid request parameters, operatorIds: Unexpected query parameter",
		},
		{
			name: "1-6. 400: バリデーションエラー: operatorIdsとopenOperatorIdsが同時に指定されている場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("operatorIds", uuid.New().String())
				q.Set("openOperatorId", uuid.New().String())
			},
			expectError: "code=400, message={[auth] BadRequest Invalid request parameters, only one of operatorIds and openOperatorId can be set.",
		},
		{
			name: "1-7. 500: システムエラー: 事業者取得に失敗した場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("operatorIds", uuid.New().String())
			},
			receive:     fmt.Errorf("Access Error"),
			expectError: "code=500, message={[auth] InternalServerError Unexpected error occurred",
		},
		{
			name: "1-8. 500: システムエラー: 事業者取得に失敗した場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("openOperatorId", uuid.New().String())
			},
			receive:     common.NewCustomError(common.CustomErrorCode500, common.Err500Unexpected, nil, common.HTTPErrorSourceAuth),
			expectError: "code=500, message={[auth] InternalServerError Unexpected error occurred",
		},
		{
			name: "1-9. 500: システムエラー: 事業者取得に失敗した場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("openOperatorId", uuid.New().String())
			},
			receive:     fmt.Errorf("Access Error"),
			expectError: "code=500, message={[auth] InternalServerError Unexpected error occurred",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			q := make(url.Values)
			test.modifyQueryParams(q)
			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)
			c.Set("operatorID", f.OperatorId)

			operatorUsecase := new(mocks.IOperatorUsecase)
			operatorHandler := handler.NewOperatorHandler(operatorUsecase)
			operatorUsecase.On("GetOperator", mock.Anything).Return(traceability.OperatorModel{}, test.receive)
			operatorUsecase.On("GetOperators", mock.Anything).Return(traceability.OperatorModels{}, test.receive)

			err := operatorHandler.GetOperator(c)
			if assert.Error(t, err) {
				assert.ErrorContains(t, err, test.expectError)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// GET /api/v1/authInfo/operator/operatorIds テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 200: 正常系(operatorIds指定、数は1)
// [x] 2-2. 200: 正常系(operatorIds指定、数は2)
// [x] 2-3. 200: 正常系(operatorIds指定、数は100)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_GetOperators_Normal(tt *testing.T) {
	var method = "GET"
	var endPoint = "/api/v1/authInfo"
	var dataTarget = "operator"

	tests := []struct {
		name              string
		modifyQueryParams func(q *url.Values)
		expectStatus      int
	}{
		{
			name: "2-1. 200: 正常系(operatorIds指定、数は1)",
			modifyQueryParams: func(q *url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("operatorIds", f.OperatorId)
			},
			expectStatus: http.StatusOK,
		},
		{
			name: "2-2. 200: 正常系(operatorIds指定、数は2)",
			modifyQueryParams: func(q *url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("operatorIds", common.GenerateUUIDString(2))
			},
			expectStatus: http.StatusOK,
		},
		{
			name: "2-3. 200: 正常系(operatorIds指定、数は100)",
			modifyQueryParams: func(q *url.Values) {
				q.Set("dataTarget", dataTarget)
				q.Set("operatorIds", common.GenerateUUIDString(100))
			},
			expectStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			q := make(url.Values)
			test.modifyQueryParams(&q)
			operatorModels := traceability.OperatorModels{}

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)
			c.Set("operatorID", f.OperatorId)
			c.Set("operatorIds", f.OperatorId)

			input := traceability.GetOperatorsInput{}
			// operatorIdsを区切り文字で分割して配列に格納
			operatorIds := strings.Split(q.Get("operatorIds"), ",")
			if len(operatorIds) > 0 {
				// UUIDの配列
				uuids := make([]uuid.UUID, len(operatorIds))

				// 文字列をUUIDに変換して配列に格納
				for i, str := range operatorIds {
					parsedUUID, err := uuid.Parse(str)
					if err != nil {
						fmt.Println("Error parsing UUID:", err)
						return
					}
					uuids[i] = parsedUUID
				}
				input.OperatorIDs = uuids
			}

			operatorUsecase := new(mocks.IOperatorUsecase)
			operatorHandler := handler.NewOperatorHandler(operatorUsecase)

			operatorUsecase.On("GetOperators", input).Return(operatorModels, nil)
			err := operatorHandler.GetOperator(c)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectStatus, rec.Code)
				operatorUsecase.AssertExpectations(t)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// PUT /api/v1/authInfo/operator 正常系テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 201: 正常系(globalOperatorId値あり)
// [x] 1-2. 201: 正常系(globalOperatorId空文字指定)
// [x] 1-3. 201: 正常系(globalOperatorIdなし、空Object)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_PutOperator_Normal(tt *testing.T) {
	var method = "PUT"
	var endPoint = "/api/v1/authInfo"
	var dataTarget = "operator"

	tests := []struct {
		name         string
		modifyInput  func(i *traceability.PutOperatorInput)
		expectStatus int
	}{
		{
			name: "1-1. 201: 正常系(globalOperatorId値あり)",
			modifyInput: func(i *traceability.PutOperatorInput) {
				i.OperatorAttributeInput.GlobalOperatorID = &f.GlobalOperatorId
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "1-2. 201: 正常系(globalOperatorId空文字指定)",
			modifyInput: func(i *traceability.PutOperatorInput) {
				i.OperatorAttributeInput.GlobalOperatorID = common.StringPtr("")
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "1-3. 201: 正常系(globalOperatorIdなし、空Object)",
			modifyInput: func(i *traceability.PutOperatorInput) {
				i.OperatorAttributeInput = &traceability.OperatorAttributeInput{}
			},
			expectStatus: http.StatusCreated,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				input := f.PutOperatorInput
				test.modifyInput(&input)
				inputJSON, _ := json.Marshal(input)
				operatorModel, _ := input.ToModel()

				q := make(url.Values)
				q.Set("dataTarget", dataTarget)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorId)

				operatorUsecase := new(mocks.IOperatorUsecase)
				operatorHandler := handler.NewOperatorHandler(operatorUsecase)
				operatorUsecase.On("PutOperator", operatorModel).Return(operatorModel, nil)

				err := operatorHandler.PutOperator(c)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expectStatus, rec.Code)
					operatorUsecase.AssertExpectations(t)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// PUT /api/v1/authInfo/operator テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 400: バリデーションエラー：dataTargetが含まれない場合
// [x] 1-2. 400: バリデーションエラー：dataTargetがoperator以外の場合
// [x] 1-3. 400: バリデーションエラー：OperatorIDが含まれない場合
// [x] 1-4. 400: バリデーションエラー：OperatorIDがUUID形式ではない場合
// [x] 1-5. 400: バリデーションエラー：OperatorNameが含まれない場合
// [x] 1-6. 400: バリデーションエラー：OperatorNameが256文字以上の場合
// [x] 1-7. 400: バリデーションエラー：OperatorAddressが含まれない場合
// [x] 1-8. 400: バリデーションエラー：OperatorAddressが256文字以上の場合
// [x] 1-9. 400: バリデーションエラー：OpenOperatorIDが含まれない場合
// [x] 1-10. 400: バリデーションエラー：OpenOperatorIDが20文字以上の場合
// [x] 1-11. 400: バリデーションエラー：OperatorAttribute.GlobalOperatorIDが256文字以上の場合
// [x] 1-12. 400: バリデーションエラー：OperatorAttributeが含まれない場合
// [x] 1-13. 400: バリデーションエラー：operatorNameがstring形式でない場合
// [x] 1-14. 400: バリデーションエラー：OperatorIDがUUID形式でない場合
// [x] 1-15. 403: バリデーションエラー：OperatorIDがjwtのOperatorIdと一致しない場合
// [x] 1-16. 500: システムエラー：更新エラーの場合
// [x] 1-17. 500: システムエラー：更新エラーの場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_PutOperator(tt *testing.T) {
	var method = "PUT"
	var endPoint = "/api/v1/authInfo"
	var dataTarget = "operator"

	tests := []struct {
		name            string
		modifyInput     func(input *traceability.PutOperatorInput)
		invalidInput    any
		invalidOperator string
		receive         error
		expectError     string
	}{
		{
			name: "1-3. 400: バリデーションエラー：OperatorIDが含まれない場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorID = ""
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorId: cannot be blank.",
		},
		{
			name: "1-4. 400: バリデーションエラー：OperatorIDがUUID形式ではない場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorID = "hoge"
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorId: invalid UUID.",
		},
		{
			name: "1-5. 400: バリデーションエラー：OperatorNameが含まれない場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorName = ""
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorName: cannot be blank.",
		},
		{
			name: "1-6. 400: バリデーションエラー：OperatorNameが256文字以上の場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorName = "Jv0ceJYX9Pa9zQTtlYPQLqLyUUhhNZ5EQCL2JDj9jLfrrgFK8MzV7zkaPvVj1wVtq5ESQGAbXrOhElxsVJzBjSxMBhwOa7hJwBrEkJjmYV8njJma2Zq6OZ7z9lXXh3xt6rYY0mYLLWpPGorQTOSY4XOkvOHfcmusmBl8OaFWjrAIUo9XwYfN2wVF4bKS32uD5vfwAzU5mhWCNwlZqABU9skfSQW9aMmCxbPkFiTq3P9hN9x4FR4m2SqB1AMLbNGu4"
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorName: the length must be between 1 and 256.",
		},
		{
			name: "1-7. 400: バリデーションエラー：OperatorAddressが含まれない場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorAddress = ""
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorAddress: cannot be blank.",
		},
		{
			name: "1-8. 400: バリデーションエラー：OperatorAddressが256文字以上の場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorAddress = "Jv0ceJYX9Pa9zQTtlYPQLqLyUUhhNZ5EQCL2JDj9jLfrrgFK8MzV7zkaPvVj1wVtq5ESQGAbXrOhElxsVJzBjSxMBhwOa7hJwBrEkJjmYV8njJma2Zq6OZ7z9lXXh3xt6rYY0mYLLWpPGorQTOSY4XOkvOHfcmusmBl8OaFWjrAIUo9XwYfN2wVF4bKS32uD5vfwAzU5mhWCNwlZqABU9skfSQW9aMmCxbPkFiTq3P9hN9x4FR4m2SqB1AMLbNGu4"
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorAddress: the length must be between 1 and 256.",
		},
		{
			name: "1-9. 400: バリデーションエラー：OpenOperatorIDが含まれない場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OpenOperatorID = ""
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, openOperatorId: cannot be blank.",
		},
		{
			name: "1-10. 400: バリデーションエラー：OpenOperatorIDが20文字以上の場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OpenOperatorID = "Jv0ceJYX9Pa9zQTtlYPQLqLyUUhhNZ5EQCL2JDj9jLfrrgFK8MzV7zkaPvVj1wVtq5ESQGAbXrOhElxsVJzBjSxMBhwOa7hJwBrEkJjmYV8njJma2Zq6OZ7z9lXXh3xt6rYY0mYLLWpPGorQTOSY4XOkvOHfcmusmBl8OaFWjrAIUo9XwYfN2wVF4bKS32uD5vfwAzU5mhWCNwlZqABU9skfSQW9aMmCxbPkFiTq3P9hN9x4FR4m2SqB1AMLbNGu4"
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, openOperatorId: the length must be between 1 and 20.",
		},
		{
			name: "1-11. 400: バリデーションエラー：OperatorAttribute.GlobalOperatorIDが256文字以上の場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorAttributeInput.GlobalOperatorID = common.StringPtr("Jv0ceJYX9Pa9zQTtlYPQLqLyUUhhNZ5EQCL2JDj9jLfrrgFK8MzV7zkaPvVj1wVtq5ESQGAbXrOhElxsVJzBjSxMBhwOa7hJwBrEkJjmYV8njJma2Zq6OZ7z9lXXh3xt6rYY0mYLLWpPGorQTOSY4XOkvOHfcmusmBl8OaFWjrAIUo9XwYfN2wVF4bKS32uD5vfwAzU5mhWCNwlZqABU9skfSQW9aMmCxbPkFiTq3P9hN9x4FR4m2SqB1AMLbNGu4")
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorAttribute: (globalOperatorId: the length must be no more than 256.).",
		},
		{
			name: "1-12. 400: バリデーションエラー：OperatorAttributeが含まれない場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorAttributeInput = nil
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorAttribute: cannot be blank.",
		},
		{
			name: "1-13. 400: バリデーションエラー：operatorNameがstring形式でない場合",
			invalidInput: struct {
				OperatorID      int
				OperatorName    int
				OperatorAddress int
				OpenOperatorID  int
			}{
				1,
				1,
				1,
				1,
			},
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorId: Unmarshal type error: expected=string, got=number.",
		},
		{
			name: "1-14. 400: バリデーションエラー：OperatorIDがUUID形式でない場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorAttributeInput.GlobalOperatorID = &f.GlobalOperatorId
			},
			invalidOperator: "invalidOperator",
			expectError:     "code=400, message={[auth] BadRequest Invalid or expired token",
		},
		{
			name: "1-15. 403: バリデーションエラー：OperatorIDがjwtのOperatorIdと一致しない場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorAttributeInput.GlobalOperatorID = &f.GlobalOperatorId
			},
			invalidOperator: "e03cc699-7234-31ed-86be-cc18c92208e6",
			expectError:     "code=403, message={[auth] AccessDenied You do not have the necessary privileges",
		},
		{
			name: "1-16. 500: システムエラー：更新エラーの場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorAttributeInput.GlobalOperatorID = &f.GlobalOperatorId
			},
			receive:     common.NewCustomError(common.CustomErrorCode500, common.Err500Unexpected, nil, common.HTTPErrorSourceAuth),
			expectError: "code=500, message={[auth] InternalServerError Unexpected error occurred",
		},
		{
			name: "1-17. 500: システムエラー：更新エラーの場合",
			modifyInput: func(input *traceability.PutOperatorInput) {
				input.OperatorAttributeInput.GlobalOperatorID = &f.GlobalOperatorId
			},
			receive:     fmt.Errorf("Access Error"),
			expectError: "code=500, message={[auth] InternalServerError Unexpected error occurred",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				var inputJSON []byte
				input := f.PutOperatorInput
				if test.invalidInput != nil {
					inputJSON, _ = json.Marshal(test.invalidInput)
				} else {
					test.modifyInput(&input)
					inputJSON, _ = json.Marshal(input)
				}
				q := make(url.Values)
				q.Set("dataTarget", dataTarget)
				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				if test.invalidOperator != "" {
					c.Set("operatorID", test.invalidOperator)
				} else {
					c.Set("operatorID", f.OperatorId)
				}

				operatorUsecase := new(mocks.IOperatorUsecase)
				operatorHandler := handler.NewOperatorHandler(operatorUsecase)
				operatorUsecase.On("PutOperator", mock.Anything).Return(traceability.OperatorModel{}, test.receive)

				err := operatorHandler.PutOperator(c)
				if assert.Error(t, err) {
					assert.ErrorContains(t, err, test.expectError)
				}
			},
		)
	}
}
