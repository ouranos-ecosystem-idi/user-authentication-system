package usecase

// IResetUsecase
// Summary: This is the interface which defines the reset usecase.
//
//go:generate mockery --name IResetUsecase --output ../test/mock --case underscore
type IResetUsecase interface {
	Reset(apikey string) error
}
