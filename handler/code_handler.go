package handler

import (
	"auth-code-generator/service"
	"encoding/json"
	"net/http"
	"time"
)

const validitySeconds = 60

type GenerateResponse struct {
	Code        string    `json:"code"`
	GeneratedAt time.Time `json:"generated_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}
type VerifyRequest struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	Code      string `json:"code"`
}
type VerifyResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

type CodeHandler struct {
	service service.CodeService
}

func NewCodeHandler(s service.CodeService) *CodeHandler {
	return &CodeHandler{service: s}
}

// GenerateCode godoc
// @Summary Generate a new 2FA code
// @Description Generates a new 6-digit code for a user, valid for a fixed 60-second UTC window.
// @ID generate-code
// @Produce  json
// @Param   userID    query   string  true  "User ID"
// @Param   userEmail query   string  true  "User Email"
// @Success 200 {object} GenerateResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /code [get]
func (h *CodeHandler) GenerateCode(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	userEmail := r.URL.Query().Get("userEmail")
	if userID == "" || userEmail == "" {
		http.Error(w, "userID and userEmail query parameters are required", http.StatusBadRequest)
		return
	}

	storedCode, err := h.service.Generate(userID, userEmail)
	if err != nil {
		http.Error(w, "Failed to generate code", http.StatusInternalServerError)
		return
	}

	response := GenerateResponse{
		Code:        storedCode.Code,
		GeneratedAt: storedCode.CreatedAt,
		ExpiresAt:   storedCode.CreatedAt.Add(validitySeconds * time.Second),
	}
	respondWithJSON(w, http.StatusOK, response)
}

// VerifyCode godoc
// @Summary Verify a 2FA code
// @Description Verifies a 2FA code submitted by a user.
// @ID verify-code
// @Accept  json
// @Produce  json
// @Param   verificationRequest body VerifyRequest true "Verification Request"
// @Success 200 {object} VerifyResponse
// @Failure 400 {string} string "Invalid request body"
// @Router /verify [post]
func (h *CodeHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valid, message := h.service.Verify(req.UserID, req.Code)
	respondWithJSON(w, http.StatusOK, VerifyResponse{Valid: valid, Message: message})
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
