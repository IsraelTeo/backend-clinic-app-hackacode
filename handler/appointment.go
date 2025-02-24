package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/appointment"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/validate"
	"github.com/labstack/echo/v4"
)

type AppointmentHandler struct {
	logicAppointment appointment.AppointmentLogic
}

func NewAppointmentHandler(logicAppointment appointment.AppointmentLogic) *AppointmentHandler {
	return &AppointmentHandler{logicAppointment: logicAppointment}
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

	log.Printf("appointment-handler: appointment fetching with ID: %d", ID)

	appointment, err := h.logicAppointment.GetAppointmentByID(ID)
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
	log.Println("appointment-handler: request received in GetAllAppointments")

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	appointments, err := h.logicAppointment.GetAllAppointments(limit, offset)
	if len(appointments) == 0 {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.SuccessPatiensListEmpty,
			Status:  http.StatusOK,
			Data:    []model.Appointment{},
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
		Message: response.SuccessAppointmentFound,
		Status:  http.StatusOK,
		Data:    appointments,
	})
}

func (h *AppointmentHandler) CreateAppointment(c echo.Context) error {
	log.Println("appointment-handler: request received in CreateAppointment")

	appointment := model.Appointment{}

	err := c.Bind(&appointment)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorBadRequest.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	finalPrice, err := h.logicAppointment.CreateAppointment(&appointment)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(), Status: http.StatusInternalServerError,
			Data: nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessAppointmentCreated,
		Status:  http.StatusCreated,
		Data:    finalPrice,
	})
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

	log.Printf("appointment-handler: request received in UpdateAppointment with ID: %d", ID)

	appointment := model.Appointment{}

	err = c.Bind(&appointment)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorBadRequest.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	finalPrice, err := h.logicAppointment.UpdateAppointment(ID, &appointment)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessAppointmentUpdated,
		Status:  http.StatusOK,
		Data:    finalPrice,
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

	log.Printf("appointment-handler: request received in DeleteAppointment with ID: %d", ID)

	err = h.logicAppointment.DeleteAppointment(ID)
	if err != nil {
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
