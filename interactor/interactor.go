package interactor

import (
	"authenticator-backend/config"
	firebase_client "authenticator-backend/infrastructure/firebase"
	"authenticator-backend/infrastructure/firebase/repository"
	"authenticator-backend/infrastructure/persistence/datastore"
	"authenticator-backend/presentation/http/echo/handler"
	"authenticator-backend/presentation/http/echo/middleware"
	"authenticator-backend/usecase"

	firebase "firebase.google.com/go/v4"
	"gorm.io/gorm"
)

// Interactor
// Summary: This is interface which defines the methods that the interactor struct should implement.
type Interactor interface {
	NewAppHandler() handler.AppHandler
	NewAuthMiddleware() middleware.AuthMiddleware
}

// interactor
// Summary: This is structure which defines the fields that the interactor struct should have.
type interactor struct {
	cfg            *config.Config
	db             *gorm.DB
	firebaseConfig *firebase.Config
}

// NewInteractor
// Summary: This is function to create a new interactor struct.
// input: cfg(*config.Config) configuration
// input: db(*gorm.DB) database connection
// input: fc(*firebase.Config) firebase configuration
// output: Interactor
func NewInteractor(
	cfg *config.Config,
	db *gorm.DB,
	fc *firebase.Config,
) Interactor {
	return &interactor{
		cfg,
		db,
		fc,
	}
}

// appHandler
// Summary: This is structure which defines the fields that the appHandler struct should have.
type appHandler struct {
	handler.AuthHandler
	handler.OuranosHandler
}

// NewAppHandler
// Summary: This is function to create a new appHandler struct.
// output: handler.AppHandler
func (i *interactor) NewAppHandler() handler.AppHandler {
	firebaseCli, _ := firebase_client.NewClient(i.firebaseConfig.ProjectID, i.cfg.FirebaseAuthEmulatorHost)

	ouranosRepository := datastore.NewOuranosRepository(i.db)
	authRepository := datastore.NewAuthRepository(i.db)
	firebaseRepository := repository.NewFirebase(firebaseCli, i.cfg.IDPSignInURL, i.cfg.IDPAPIKey, i.cfg.SecureTokenAPIKey, i.cfg.SecureTokenAPI)

	authUsecase := usecase.NewAuthUsecase(firebaseRepository)
	verifyUsecase := usecase.NewVerifyUsecase(firebaseRepository, authRepository)
	operatorUsecase := usecase.NewOperatorUsecase(ouranosRepository)
	plantUsecase := usecase.NewPlantUsecase(ouranosRepository)
	resetUsecase := usecase.NewResetUsecase(ouranosRepository, authRepository)

	operatorHandler := handler.NewOperatorHandler(operatorUsecase)
	plantHandler := handler.NewPlantHandler(plantUsecase)
	resetHandler := handler.NewResetHandler(resetUsecase)

	authHandler := handler.NewAuthHandler(
		authUsecase,
		verifyUsecase,
	)
	ouranosHandler := handler.NewOuranosHandler(
		operatorHandler,
		plantHandler,
		resetHandler,
	)

	appHandler := &appHandler{
		AuthHandler:    authHandler,
		OuranosHandler: ouranosHandler,
	}
	return appHandler
}

// NewAuthMiddleware
// Summary: This is function to create a new authMiddleware struct.
// output: middleware.AuthMiddleware
func (i *interactor) NewAuthMiddleware() middleware.AuthMiddleware {
	firebaseCli, _ := firebase_client.NewClient(i.firebaseConfig.ProjectID, i.cfg.FirebaseAuthEmulatorHost)

	authRepository := datastore.NewAuthRepository(i.db)
	firebaseRepository := repository.NewFirebase(firebaseCli, i.cfg.IDPSignInURL, i.cfg.IDPAPIKey, i.cfg.SecureTokenAPIKey, i.cfg.SecureTokenAPI)

	verifyUsecase := usecase.NewVerifyUsecase(firebaseRepository, authRepository)

	return middleware.NewAuthMiddleware(verifyUsecase)
}
