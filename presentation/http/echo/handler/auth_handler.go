package handler

import (
	"authenticator-backend/usecase"

	"github.com/labstack/echo/v4"
)

type (
	AuthHandler interface {
		Login(c echo.Context) error
		Refresh(c echo.Context) error
		ChangePassword(c echo.Context) error
		TokenIntrospection(c echo.Context) error
		ApiKey(c echo.Context) error
	}

	authHandler struct {
		AuthUsecase   usecase.IAuthUsecase
		VerifyUsecase usecase.IVerifyUsecase
	}
)

func NewAuthHandler(
	authUsecase usecase.IAuthUsecase,
	verifyUsecase usecase.IVerifyUsecase,
) AuthHandler {
	return &authHandler{
		AuthUsecase:   authUsecase,
		VerifyUsecase: verifyUsecase,
	}
}
