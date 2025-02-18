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
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("service-handler: medical service fetching with ID: %d", ID)

	service, err := h.logic.GetServiceByID(uint(ID))
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusNotFound,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessServiceFound,
		Status:  http.StatusOK,
		Data:    service,
	})
}

func (h *ServiceHandler) GetAllServices(c echo.Context) error {
	log.Println("service-handler: request received in GetAllServices")

	services, err := h.logic.GetAllServices()
	if len(services) == 0 {
		return response.WriteSuccess(&response.WriteResponse{
			C:       c,
			Message: response.SuccessServicesListEmpty,
			Status:  http.StatusOK,
			Data:    []model.Service{},
		})
	}

	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusNotFound,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessServicesFound,
		Status:  http.StatusOK,
		Data:    services,
	})
}

func (h *ServiceHandler) CreateService(c echo.Context) error {
	log.Println("service-handler: request received in CreateService")

	service := model.Service{}
	if err := c.Bind(&service); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := h.logic.CreateService(&service); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessServiceCreated,
		Status:  http.StatusCreated,
		Data:    nil,
	})
}

func (h *ServiceHandler) UpdateService(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Println("service-handler: request received in UpdateService")

	service := model.Service{}
	if err := c.Bind(&service); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := h.logic.UpdateService(uint(ID), &service); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessServiceUpdated,
		Status:  http.StatusOK,
		Data:    nil,
	})
}

func (h *ServiceHandler) DeleteService(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Println("service-handler: request received in DeleteService")

	if err := h.logic.DeleteService(ID); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessServiceDeleted,
		Status:  http.StatusOK,
		Data:    nil,
	})
}
