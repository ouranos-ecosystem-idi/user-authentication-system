package common

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgconn"
)

type HTTPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

// 下記はSwagger用のModel
type HTTP400Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type HTTP401Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type HTTP403Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type HTTP404Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type HTTP500Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type HTTP503Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

var (
	// 400 Error Messages
	Err400InvalidRequest  = "Invalid request parameters"
	Err400InvalidJSON     = "Invalid JSON format"
	Err400RequestTooLarge = "Request payload too large"
	Err400Validation      = "Validation failed"
	// 401 Error Messages
	Err401InvalidCredentials = "Invalid credentials"
	Err401Authentication     = "Authentication required"
	Err401InvalidToken       = "Invalid or expired token"
	// 403 Error Messages
	Err403AccessDenied          = "You do not have the necessary privileges"
	Err403InvalidKey            = "Invalid key"
	Err403IPNotAuthorizedForKey = "IP address not authorized for this API key"
	// 404 Error Messages
	Err404ResourceNotFound = "Resource Not Found"
	Err404ItemNotFound     = "Item or record Not Found"
	Err404EndpointNotFound = "Endpoint Not Found"
	// 500 Error Messages
	Err500Unexpected = "Unexpected error occurred"
	// 503 Error Messages
	Err503OuterService = "Unexpected error occurred in outer service"
)

// HTTPErrorSource
// Summary: This is enum which defines HTTPErrorSource.
type HTTPErrorSource string

// ToString
// Summary: This is the function to convert HTTPErrorSource to string.
// output: (string) converted to string
func (t HTTPErrorSource) ToString() string {
	return string(t)
}

var (
	HTTPErrorSourceDataspace HTTPErrorSource = "dataspace"
	HTTPErrorSourceAuth      HTTPErrorSource = "auth"
)

// HTTPErrorGenerate
// Summary: This is the function to generate HTTPError.
// input: httpStatusCode(int) http status code
// input: source(HTTPErrorSource) source of error
// input: errorMsg(string) error message
// input: operatorID(string) ID of the operator
// input: dataTarget(string) target of the data
// input: method(string) method of the request
// input: errorDetails(...string) error details
// output: (int) http status code
// output: (HTTPError) HTTPError object
func HTTPErrorGenerate(
	httpStatusCode int,
	source HTTPErrorSource,
	errorMsg string,
	operatorID string,
	dataTarget string,
	method string,
	errorDetails ...string,
) (int, HTTPError) {
	now := time.Now()
	utcNow := now.UTC()
	iso8601Format := "2006-01-02T15:04:05.000Z" // ISO 8601形式
	isoUtcTime := utcNow.Format(iso8601Format)
	detailMessage := "id: " + operatorID + ", timeStamp: " + isoUtcTime + ", dataTarget: " + dataTarget + ", method: " + method

	for _, errorDetail := range errorDetails {
		errorMsg += fmt.Sprintf(", %s", errorDetail)
	}

	switch httpStatusCode {
	case 400:
		errorModel := HTTPError{
			Code:    formatErrorCode("BadRequest", source),
			Message: errorMsg,
			Detail:  detailMessage,
		}
		return 400, errorModel
	case 401:
		errorModel := HTTPError{
			Code:    formatErrorCode("Unauthorized", source),
			Message: errorMsg,
			Detail:  detailMessage,
		}
		return 401, errorModel
	case 403:
		errorModel := HTTPError{
			Code:    formatErrorCode("AccessDenied", source),
			Message: errorMsg,
			Detail:  detailMessage,
		}
		return 403, errorModel
	case 404:
		errorModel := HTTPError{
			Code:    formatErrorCode("NotFound", source),
			Message: errorMsg,
			Detail:  detailMessage,
		}
		return 404, errorModel
	case 500:
		errorModel := HTTPError{
			Code:    formatErrorCode("InternalServerError", source),
			Message: errorMsg,
			Detail:  detailMessage,
		}
		return 500, errorModel
	case 503:
		errorModel := HTTPError{
			Code:    formatErrorCode("ServiceUnavailable", source),
			Message: errorMsg,
			Detail:  detailMessage,
		}
		return 503, errorModel
	default:
		errorModel := HTTPError{
			Code:    formatErrorCode("InternalServerError", source),
			Message: errorMsg,
			Detail:  detailMessage,
		}
		return 500, errorModel
	}
}

// formatErrorCode
// Summary: This is the function to format error code.
// input: code(string) error code
// input: source(HTTPErrorSource) source of error
// output: (string) formatted error code
func formatErrorCode(code string, source HTTPErrorSource) string {
	return fmt.Sprintf("[%s] %s", source, code)
}

// CustomError
// Summary: This is structure which defines CustomError.
type CustomError struct {
	Code          CustomErrorCode
	Message       string
	MessageDetail *string
	Source        HTTPErrorSource
}

// NewCustomError
// Summary: This is the function to create new CustomError.
// input: code(CustomErrorCode) error code
// input: message(string) error message
// input: messageDetail(*string) error message detail
// input: source(HTTPErrorSource) source of error
// output: (*CustomError) CustomError object
func NewCustomError(code CustomErrorCode, message string, messageDetail *string, source HTTPErrorSource) *CustomError {
	return &CustomError{
		Code:          code,
		Message:       message,
		MessageDetail: messageDetail,
		Source:        source,
	}
}

// Error
// Summary: This is the function to get error message.
// output: (string) error message
func (e CustomError) Error() string {
	return e.Message
}

// IsWarn
// Summary: This is the function to check if the error is a warning.
// output: (bool) true: warning, false: not warning
func (e CustomError) IsWarn() bool {
	return e.Code >= 400 && e.Code < 500
}

// CustomErrorCode
// Summary: This is enum which defines CustomErrorCode.
type CustomErrorCode int

var (
	CustomErrorCode400 CustomErrorCode = http.StatusBadRequest
	CustomErrorCode401 CustomErrorCode = http.StatusUnauthorized
	CustomErrorCode403 CustomErrorCode = http.StatusForbidden
	CustomErrorCode404 CustomErrorCode = http.StatusNotFound
	CustomErrorCode500 CustomErrorCode = http.StatusInternalServerError
	CustomErrorCode503 CustomErrorCode = http.StatusServiceUnavailable
)

// FormatBindErrMsg
// Summary: This is the function to format bind error message.
// input: err(error) error object
// output: (string) formatted error message
func FormatBindErrMsg(err error) string {
	if strings.Contains(err.Error(), "Syntax error") {
		return extractSyntaxErrorMessage(err.Error())
	}
	return extractTypeErrorMessage(err.Error())
}

// extractSyntaxErrorMessage
// Summary: This is the function to extract syntax error message.
// input: errString(string) error string
// output: (string) extracted error message
func extractSyntaxErrorMessage(errString string) string {
	const prefix = "error="
	start := strings.Index(errString, prefix)
	if start == -1 {
		return ""
	}

	message := errString[start+len(prefix):]
	end := strings.Index(message, ", internal=")
	if end != -1 {
		message = message[:end]
	}

	return message
}

// extractTypeErrorMessage
// Summary: This is the function to extract type error message.
// input: errString(string) error string
// output: (string) extracted error message
func extractTypeErrorMessage(errString string) string {
	messageStart := strings.Index(errString, "message=")
	if messageStart == -1 {
		return ""
	}
	message := errString[messageStart+len("message="):]

	fieldStart := strings.Index(message, ", field=")
	if fieldStart == -1 {
		return ""
	}
	field := message[fieldStart+len(", field="):]

	errorDescription := message[:fieldStart]

	nextComma := strings.Index(field, ",")
	if nextComma != -1 {
		field = field[:nextComma]
	}

	formattedMessage := fmt.Sprintf("%s: %s.", field, errorDescription)
	return formattedMessage
}

// UnexpectedQueryParameter
// Summary: This is the function to format unexpected query parameter error message.
// input: param(string) parameter name
// output: (string) formatted error message
func UnexpectedQueryParameter(param string) string {
	return fmt.Sprintf("%v: Unexpected query parameter", param)
}

// InvalidUUIDError
// Summary: This is the function to format invalid UUID error message.
// input: name(string) target name
// output: (string) formatted error message
func InvalidUUIDError(name string) string {
	return fmt.Sprintf("%s: invalid UUID.", name)
}

// OnlyOneCanBeSpecified
// Summary: This is the function to format only one can be specified error message.
// input: name1(string) name1
// input: name2(string) name2
// output: (string) formatted error message
func OnlyOneCanBeSpecified(name1 string, name2 string) string {
	return fmt.Sprintf("only one of %v and %v can be set.", name1, name2)
}

// FieldIsImutable
// Summary: This is the function to format field is immutable message.
// input: name(string) name
// output: (string) formatted error message
func FieldIsImutable(name string) string {
	return fmt.Sprintf("%s: field is immutable.", name)
}

// LimitLessThanError
// Summary: This is the function to format limit less than error message.
// input: min(int) min
// input: limit(int) limit
// output: (string) formatted error message
func LimitLessThanError(min int, limit int) string {
	return fmt.Sprintf("limit less than %v error. get value: %v", min, limit)
}

// LimitUpperError
// Summary: This is the function to format limit upper error message.
// input: limit(int) limit
// output: (string) formatted error message
func LimitUpperError(limit int) string {
	return fmt.Sprintf("limit upper limit error. get value: %v.", limit)
}

// UnexpectedEnumError
// Summary: This is the function to format unexpected enum error message.
// input: name(string) name
// input: value(string) value
// output: (string) formatted error message
func UnexpectedEnumError(name string, value string) string {
	return fmt.Sprintf("unexpected %v. get value: %v", name, value)
}

// NotFoundError
// Summary: This is the function to format not found error message.
// input: value(string) value
// output: (string) formatted error message
func NotFoundError(value string) string {
	return fmt.Sprintf("%v is not specified", value)
}

// ValidateStructureError
// Summary: This is the function to format validate structure error message.
// input: name(string) structure name
// input: err(error) error object
// output: (string) formatted error message
func ValidateStructureError(name string, err error) string {
	return fmt.Sprintf("%v: (%v)", name, err)
}

// DeleteTableError
// Summary: This is the function to format delete table error message.
// input: name(string) table name
// input: err(error) error object
// output: (string) formatted error message
func DeleteTableError(name string, err error) string {
	return fmt.Sprintf("failed to physically delete record from table %v : %v", name, err)
}

// UnmatchValuesError
// Summary: This is the function to format Uumatch values error message.
// input: value1(string) target value
// input: value2(string) target value
// output: (string) formatted error message
func UnmatchValuesError(value1 string, value2 string) string {
	return fmt.Sprintf("The %v in the token does not match the %v in the request body", value1, value2)
}

// DuplicateOperatorError
// Summary: This is function which get message of duplicate operator error.
// input: pgError(PgError) PgError object
// input: globalOperatorID(*string) Pointer of globalOperatorID
// output: (string) formatted error message
func DuplicateOperatorError(pgError *pgconn.PgError, globalOperatorID *string) string {
	msg := pgError.Message
	var dupKey string
	if strings.Contains(msg, "unique_global_operator_id") {
		dupKey = "globalOperatorId"
	}
	var dupValue string
	switch dupKey {
	case "globalOperatorId":
		if globalOperatorID != nil {
			dupValue = *globalOperatorID
		}
	}

	return fmt.Sprintf("%s: %s is already exists.", dupKey, dupValue)
}

// DuplicatePlantError
// Summary: This is function which get message of duplicate plant error.
// input: pgError(PgError) PgError object
// input: openPlantID(string) Value of openPlantID
// input: globalPlantID(*string) Pointer of globalPlantID
// output: (string) formatted error message
func DuplicatePlantError(pgError *pgconn.PgError, openPlantID string, globalPlantID *string) string {
	msg := pgError.Message

	var dupKey string
	if strings.Contains(msg, "unique_open_plant_id_operator_id") {
		dupKey = "openPlantId"
	} else if strings.Contains(msg, "unique_global_plant_id_operator_id") {
		dupKey = "globalPlantId"
	}

	var dupValue string
	switch dupKey {
	case "openPlantId":
		dupValue = openPlantID
	case "globalPlantId":
		if globalPlantID != nil {
			dupValue = *globalPlantID
		}
	}

	return fmt.Sprintf("%s: %s is already exists.", dupKey, dupValue)
}
