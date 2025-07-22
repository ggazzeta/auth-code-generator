package service

import (
	"fmt"
	"hash/fnv"
	"main/internal/store"
	"main/pkg/models"
	"math/rand"
	"time"
)

const validitySeconds = 30

type CodeService interface {
	Generate(userID, userEmail string) (models.StoredCode, error)
	Verify(userID, userCode string) (bool, string)
}

type codeService struct {
	repo store.CodeRepository // Depends on the interface
}

func NewCodeService(r store.CodeRepository) CodeService {
	return &codeService{repo: r}
}

func (s *codeService) Generate(userID, userEmail string) (models.StoredCode, error) {
	nowUTC := time.Now().UTC()
	second := nowUTC.Second()
	var windowStart time.Time
	if second < 30 {
		windowStart = nowUTC.Truncate(time.Minute)
	} else {
		windowStart = nowUTC.Truncate(time.Minute).Add(30 * time.Second)
	}

	code, err := s.generateDeterministicCode(userID, windowStart)
	if err != nil {
		return models.StoredCode{}, err
	}

	storedCode := models.StoredCode{
		Code:      code,
		UserID:    userID,
		UserEmail: userEmail,
		CreatedAt: windowStart,
	}

	err = s.repo.Save(storedCode)
	return storedCode, err
}

func (s *codeService) Verify(userID, userCode string) (bool, string) {
	storedCode, found, err := s.repo.Get(userID)
	if err != nil || !found {
		return false, "No code found for user. Please request a new one."
	}

	expiresAt := storedCode.CreatedAt.Add(validitySeconds * time.Second)
	if time.Now().UTC().After(expiresAt) {
		return false, "Code has expired."
	}

	if storedCode.Code != userCode {
		return false, "Invalid code."
	}

	return true, "Code verified successfully."
}

func (s *codeService) generateDeterministicCode(userID string, t time.Time) (string, error) {
	seedSource := fmt.Sprintf("%s-%d", userID, t.Unix())
	h := fnv.New64a()
	h.Write([]byte(seedSource))
	seed := h.Sum64()
	r := rand.New(rand.NewSource(int64(seed)))
	code := r.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code), nil
}
