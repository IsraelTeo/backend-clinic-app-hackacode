package handler

import (
	"log"
	"net/http"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/logic"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/validate"
	"github.com/labstack/echo/v4"
)

type ServiceHandler struct {
	logic logic.ServiceLogic
}

func NewServiceHandler(logic logic.ServiceLogic) *ServiceHandler {
	return &ServiceHandler{logic: logic}
}

func (h *ServiceHandler) GetServiceByID(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("handler: medical service fetching with ID: %d", ID)

	service, err := h.logic.GetServiceByID(uint(ID))
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}
	return response.WriteSuccess(c, response.SuccessServiceFound, http.StatusOK, service)
}

func (h *ServiceHandler) GetAllServices(c echo.Context) error {
	log.Println("handler: request received in GetAllServices")

	services, err := h.logic.GetAllServices()
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}

	return response.WriteSuccess(c, response.SuccessServicesFound, http.StatusOK, services)
}

func (h *ServiceHandler) CreateService(c echo.Context) error {
	log.Println("handler: request received in CreateService")

	service := model.Service{}
	if err := c.Bind(&service); err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	if err := h.logic.CreateService(&service); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessServiceCreated, http.StatusCreated, nil)
}

func (h *ServiceHandler) UpdateService(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Println("handler: request received in UpdateService")

	service := model.Service{}
	if err := c.Bind(&service); err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	if err := h.logic.UpdateService(uint(ID), &service); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessServiceUpdated, http.StatusOK, nil)
}

func (h *ServiceHandler) DeleteService(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Println("handler: request received in DeleteService")

	if err := h.logic.DeleteService(ID); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessServiceDeleted, http.StatusOK, nil)
}
