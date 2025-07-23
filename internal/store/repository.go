package store

import "auth-code-generator/pkg/models"

type CodeRepository interface {
	Save(code models.StoredCode) error
	Get(userID string) (models.StoredCode, bool, error)
}
