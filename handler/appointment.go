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

type AppointmentHandler struct {
	logic logic.AppointmentLogic
}

func NewAppointmentHandler(logic logic.AppointmentLogic) *AppointmentHandler {
	return &AppointmentHandler{logic: logic}
}

func (h *AppointmentHandler) GetAppointmentByID(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorInvalidID.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("handler: appointment fetching with ID: %d", ID)

	appointment, err := h.logic.GetAppointmentByID(ID)
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
		Message: response.SuccessAppointmentsFound,
		Status:  http.StatusOK,
		Data:    appointment,
	})
}

func (h *AppointmentHandler) GetAllAppointments(c echo.Context) error {
	log.Println("handler: request received in GetAllAppointments")

	appointments, err := h.logic.GetAllAppointments()
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
		Message: response.SuccessAppointmentFound,
		Status:  http.StatusOK,
		Data:    appointments,
	})
}

func (h *AppointmentHandler) CreateAppointment(c echo.Context) error {
	log.Println("handler: request received in CreateAppointment")

	appointment := model.Appointment{}
	if err := c.Bind(&appointment); err != nil {
		log.Printf("handler: error binding request: %v", err)
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorBadRequest.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := c.Validate(&appointment); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorBadRequest.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	finalPricePkg, err := h.logic.CreateAppointment(&appointment)
	if err != nil {
		log.Printf("handler: error in business logic: %v", err)
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorToCreatedAppointment.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccessAppointmentDesc(&response.WriteResponse{
		C:       c,
		Message: response.SuccessAppointmentCreated,
		Status:  http.StatusCreated,
	},
		finalPricePkg,
		appointment.Patient.Insurance,
	)
}

func (h *AppointmentHandler) UpdateAppointment(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorInvalidID.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("handler: request received in UpdateAppointment with ID: %d", ID)

	appointment := model.Appointment{}
	if err := c.Bind(&appointment); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorBadRequest.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := c.Validate(&appointment); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorBadRequest.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := h.logic.UpdateAppointment(ID, &appointment); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorToUpdatedAppointment.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessAppointmentUpdated,
		Status:  http.StatusOK,
		Data:    nil,
	})
}

func (h *AppointmentHandler) DeleteAppointment(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorInvalidID.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("handler: request received in DeleteAppointment with ID: %d", ID)

	if err := h.logic.DeleteAppointment(ID); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorToDeletedDoctor.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessAppointmentDeleted,
		Status:  http.StatusOK,
		Data:    nil,
	})
}
