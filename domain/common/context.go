package common

import (
	"authenticator-backend/extension/logger"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// QueryParamPtr
// Summary: This is function which returns the query parameter value as a pointer
// input: c(echo.Context): echo context
// input: paramName(string): query parameter name
// output: (*string) pointer to the query parameter value
func QueryParamPtr(c echo.Context, paramName string) *string {
	val := c.QueryParam(paramName)
	if val == "" {
		return nil
	}
	return &val
}

// QueryParamUUID
// Summary: This is function which returns the query parameter value as a UUID slice
// input: c(echo.Context): echo context
// input: paramName(string): query parameter name
// output: ([]uuid.UUID) UUID value
// output: (error) error object
func QueryParamUUIDs(c echo.Context, paramName string) ([]uuid.UUID, error) {
	val := c.QueryParam(paramName)
	if val == "" {
		return nil, nil
	}

	var parsedUUIDs []uuid.UUID
	for _, v := range strings.Split(val, ",") {
		if len(v) != UUIDLen {
			logger.Set(c).Warnf(UnexpectedQueryParameter(paramName))

			return nil, fmt.Errorf(UnexpectedQueryParameter(paramName))
		}

		parsedUUID, err := uuid.Parse(v)
		if err != nil {
			logger.Set(c).Warnf(err.Error())

			return nil, err
		}
		parsedUUIDs = append(parsedUUIDs, parsedUUID)
	}

	return parsedUUIDs, nil

}

// QueryParamExists
// Summary: This is function which checks whether the query parameter exists
// input: c(echo.Context): echo context
// input: paramName(string): query parameter name
// output: (bool) true if the query parameter exists, false otherwise
func QueryParamExists(c echo.Context, paramName string) bool {
	queryParams := c.QueryParams()
	_, ok := queryParams[paramName]
	return ok
}
