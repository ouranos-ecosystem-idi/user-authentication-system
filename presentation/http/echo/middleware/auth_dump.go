package middleware

import (
	"encoding/json"
	"path"
	"time"

	"authenticator-backend/domain/common"
	"authenticator-backend/extension/logger"
	"authenticator-backend/usecase/input"
	"authenticator-backend/usecase/output"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

const (
	systemAuthResourceToken    = "token"
	systemAuthResourceAPIKey   = "apiKey"
	authResourceLogin          = "login"
	authResourceRefresh        = "refresh"
	authResourceChangePassword = "change"

	eventToken          = "operatorToken"
	eventAPIKey         = "apiKey"
	eventLogin          = "operatorLogin"
	eventRefresh        = "operatorRefreshToken"
	eventChangePassword = "operatorChangePassword"
)

// AuthDump
// Summary: This is the function which dumps the authentication information.
// output: (echo.MiddlewareFunc) echo middleware function
func AuthDump() echo.MiddlewareFunc {
	return echo_middleware.BodyDumpWithConfig(echo_middleware.BodyDumpConfig{
		Skipper: authDumpSkipper,
		Handler: authDumpHandler,
	})
}

// authDumpHandler
// Summary: This is the function which handles the authentication dump.
// input: c(echo.Context): echo context
// input: reqBody([]byte): request body
// input: resBody([]byte): response body
func authDumpHandler(c echo.Context, reqBody, resBody []byte) {
	resource := path.Base(c.Request().URL.Path)

	switch resource {
	case systemAuthResourceToken:
		tokenDumpHandler(c, reqBody, resBody)
	case systemAuthResourceAPIKey:
		apiKeyDumpHandler(c, reqBody, resBody)
	case authResourceLogin:
		loginDumpHandler(c, reqBody, resBody)
	case authResourceRefresh:
		refreshDumpHandler(c, reqBody, resBody)
	case authResourceChangePassword:
		changeDumpHandler(c, reqBody, resBody)
	}
}

// tokenDumpHandler
// Summary: This is the function which dumps the token authentication information.
// input: c(echo.Context): echo context
// input: reqBody([]byte): request body
// input: resBody([]byte): response body
func tokenDumpHandler(c echo.Context, reqBody, resBody []byte) {
	var req input.VerifyTokenParam
	if err := json.Unmarshal(reqBody, &req); err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}
	var res output.VerifyTokenResponse
	if err := json.Unmarshal(resBody, &res); err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}
	result := res.OperatorID != nil
	authDump(c, req, res, eventToken, result)
}

// apiKeyDumpHandler
// Summary: This is the function which dumps the API key authentication information.
// input: c(echo.Context): echo context
// input: reqBody([]byte): request body
// input: resBody([]byte): response body
func apiKeyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	var req input.VerifyAPIKeyParam
	if err := json.Unmarshal(reqBody, &req); err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}
	var res output.VerifyApiKeyResponse
	if err := json.Unmarshal(resBody, &res); err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}
	result := res.IsAPIKeyValid && res.IsIPAddressValid
	authDump(c, req, res, eventAPIKey, result)
}

// loginDumpHandler
// Summary: This is the function which dumps the login authentication information.
// input: c(echo.Context): echo context
// input: reqBody([]byte): request body
// input: resBody([]byte): response body
func loginDumpHandler(c echo.Context, reqBody, resBody []byte) {
	var req input.LoginParam
	if err := json.Unmarshal(reqBody, &req); err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}
	req.Mask()

	var res output.LoginResponse
	if err := json.Unmarshal(resBody, &res); err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}
	res.Mask()

	result := c.Response().Status == 201
	authDump(c, req, res, eventLogin, result)
}

// refreshDumpHandler
// Summary: This is the function which dumps the refresh authentication information.
// input: c(echo.Context): echo context
// input: reqBody([]byte): request body
// input: resBody([]byte): response body
func refreshDumpHandler(c echo.Context, reqBody, resBody []byte) {
	var req input.RefreshParam
	if err := json.Unmarshal(reqBody, &req); err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}
	req.Mask()

	var res output.RefreshResponse
	if err := json.Unmarshal(resBody, &res); err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}

	result := c.Response().Status == 201
	authDump(c, req, res, eventRefresh, result)
}

// changeDumpHandler
// Summary: This is the function which dumps the change password authentication information.
// input: c(echo.Context): echo context
// input: reqBody([]byte): request body
// input: resBody([]byte): response body
func changeDumpHandler(c echo.Context, reqBody, resBody []byte) {
	var req input.ChangePasswordParam
	if err := json.Unmarshal(reqBody, &req); err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}
	req.Mask()

	var res common.EmptyBody
	if err := json.Unmarshal(resBody, &res); err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}

	result := c.Response().Status == 201
	authDump(c, req, res, eventChangePassword, result)
}

// authDumpInfo
// Summary: This is the structure which defines the authentication dump information.
type authDumpInfo struct {
	Event            string      `json:"event"`
	IsRequestResult  bool        `json:"isRequestResult"`
	RequestBody      interface{} `json:"requestBody"`
	ResponseBody     interface{} `json:"responseBody"`
	TimeStamp        time.Time   `json:"timeStamp"`
	RequestIpAddress string      `json:"requestIpAddress"`
	RequestApiKey    string      `json:"requestApiKey"`
}

// authDump
// Summary: This is the function which dumps the authentication information.
// input: c(echo.Context): echo context
// input: reqBody(interface{}): request body
// input: resBody(interface{}): response body
// input: event(string): event type to dump
// input: isRequestResult(bool): is request successful
func authDump(c echo.Context, reqBody, resBody interface{}, event string, isRequestResult bool) {

	tempReqBody := reqBody
	tempResBody := resBody
	if zap.S().Level() != zap.DebugLevel {
		tempReqBody = "******"
		tempResBody = "******"
	}

	dump := authDumpInfo{
		Event:            event,
		IsRequestResult:  isRequestResult,
		RequestBody:      tempReqBody,
		ResponseBody:     tempResBody,
		TimeStamp:        time.Now(),
		RequestIpAddress: c.RealIP(),
		RequestApiKey:    c.Request().Header.Get("apiKey"),
	}

	b, err := json.Marshal(dump)
	if err != nil {
		logger.Set(c).Warnf(err.Error())

		return
	}

	if zap.S().Level() == zap.DebugLevel {
		logger.Set(c).Debugf(string(b))
	} else {
		logger.Set(c).Infof(string(b))
	}
}

// authDumpSkipper
// Summary: This is the function which skips the authentication dump.
// input: c(echo.Context): echo context
// output: (bool) fixed to false
func authDumpSkipper(c echo.Context) bool {
	return false
}
