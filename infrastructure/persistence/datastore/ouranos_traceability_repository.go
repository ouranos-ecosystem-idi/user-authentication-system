package datastore

import (
	"authenticator-backend/domain/repository"

	"gorm.io/gorm"
)

// ouranosRepository
// Summary: This is the structure which defines the repository for the Ouranos.
type ouranosRepository struct {
	db *gorm.DB
}

// NewOuranosRepository
// Summary: This is the function which creates the Ouranos repository.
// input: db(gorm.DB) database
// output: (repository.OuranosRepository) Ouranos repository
func NewOuranosRepository(
	db *gorm.DB,
) repository.OuranosRepository {
	return &ouranosRepository{db}
}
