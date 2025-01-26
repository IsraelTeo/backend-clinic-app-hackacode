package handler

import (
	"log"
	"net/http"
	"strconv"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/logic"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/labstack/echo/v4"
)

type ServiceHandler struct {
	logic logic.ServiceLogic
}

func NewServiceHandler(logic logic.ServiceLogic) *ServiceHandler {
	return &ServiceHandler{logic: logic}
}

func (h *ServiceHandler) GetServiceByID(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ID <= 0 {
		return response.WriteError(c, string(response.ErrorInvalidId), http.StatusBadRequest)
	}

	log.Printf("handler: medical service fetching with ID: %d", ID)

	service, err := h.logic.GetServiceByID(uint(ID))
	if err != nil {
		return response.WriteError(c, string(response.ErrorServiceNotFound), http.StatusNotFound)
	}

	return response.WriteSuccess(c, string(response.SuccessServiceFound), http.StatusOK, service)
}
