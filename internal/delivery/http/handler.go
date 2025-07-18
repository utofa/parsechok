package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	_ "whatsapp-parser/docs"
	"whatsapp-parser/internal/domain"
	"whatsapp-parser/internal/delivery/http/middleware"
)

// Handler структура для HTTP обработчиков
type Handler struct {
	sessionUseCase domain.SessionUseCase
}

// @title WhatsApp Parser API
// @version 1.0
// @description API для управления WhatsApp Web через Selenium
// @host localhost:8081
// @BasePath /
// @schemes http
func NewHandler(sessionUseCase domain.SessionUseCase) *Handler {
	return &Handler{
		sessionUseCase: sessionUseCase,
	}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	// Apply CORS middleware
	r.Use(middleware.CORS)

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// API endpoints
	r.HandleFunc("/session", h.CreateSession).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/session/{id}", h.RestoreSession).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/session/{id}/message", h.SendMessage).Methods(http.MethodPost, http.MethodOptions)
}

type SendMessageRequest struct {
	PhoneNumber string `json:"phone_number" example:"1234567890"`
	Message     string `json:"message" example:"Hello, World!"`
}

// CreateSession godoc
// @Summary Создать новую сессию WhatsApp
// @Description Создает новую сессию и возвращает QR код для авторизации
// @Tags session
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /session [post]
func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	session, qrCode, err := h.sessionUseCase.CreateSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"session_id": session.ID,
		"qr_code":    qrCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RestoreSession godoc
// @Summary Восстановить существующую сессию
// @Description Восстанавливает сохраненную сессию WhatsApp по ID
// @Tags session
// @Accept json
// @Produce json
// @Param id path string true "ID сессии"
// @Success 200 {string} string "OK"
// @Router /session/{id} [post]
func (h *Handler) RestoreSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	if err := h.sessionUseCase.RestoreSession(sessionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SendMessage godoc
// @Summary Отправить сообщение
// @Description Отправляет сообщение через WhatsApp используя указанную сессию
// @Tags message
// @Accept json
// @Produce json
// @Param id path string true "ID сессии"
// @Param message body SendMessageRequest true "Данные сообщения"
// @Success 200 {string} string "OK"
// @Router /session/{id}/message [post]
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.sessionUseCase.SendMessage(sessionID, req.PhoneNumber, req.Message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
} 