package interfaces

import "dfs/internal/models"

type FileRepositoryInterface interface {
	GetByID(id string) (*models.FileMetadata, error)
	GetAll() (*[]*models.FileMetadata, error)
	Create(file *models.FileMetadata) (any, error)
}
