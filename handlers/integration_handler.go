package handlers

import (
	"net/http"

	"go-microservice-highload/services"
)

type IntegrationHandler struct {
	userSvc *services.UserService
	intSvc  *services.IntegrationService
}

func NewIntegrationHandler(userSvc *services.UserService, intSvc *services.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{
		userSvc: userSvc,
		intSvc:  intSvc,
	}
}

func (h *IntegrationHandler) Export(w http.ResponseWriter, r *http.Request) {
	users := h.userSvc.GetAll()
	err := h.intSvc.ExportUsers(r.Context(), users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Export successful"))
}
