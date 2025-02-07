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
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("handler: appointment fetching with ID: %d", ID)

	appointment, err := h.logic.GetAppointmentByID(ID)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}

	return response.WriteSuccess(c, response.SuccessAppointmentsFound, http.StatusOK, appointment)
}

func (h *AppointmentHandler) GetAllAppointments(c echo.Context) error {
	log.Println("handler: request received in GetAllAppointments")

	appointments, err := h.logic.GetAllAppointments()
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}

	return response.WriteSuccess(c, response.SuccessAppointmentFound, http.StatusOK, appointments)
}
func (h *AppointmentHandler) CreateAppointment(c echo.Context) error {
	log.Println("handler: request received in CreateAppointment")

	// Vincular los datos de la solicitud a la estructura del modelo
	appointment := model.Appointment{}
	if err := c.Bind(&appointment); err != nil {
		log.Printf("handler: error binding request: %v", err)
		return response.WriteError(c, "Invalid request payload", http.StatusBadRequest)
	}

	// Llamar a la lógica para crear la cita
	originalPrice, packageDiscount, finalPrice, insuranceDiscount, err := h.logic.CreateAppointment(&appointment)
	if err != nil {
		log.Printf("handler: error in business logic: %v", err)
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	// Respuesta con los valores calculados correctamente
	return response.WriteSuccessAppointmentDesc(
		c,
		response.SuccessAppointmentCreated,
		http.StatusCreated,
		originalPrice,
		packageDiscount,
		insuranceDiscount,
		finalPrice,
		appointment.Patient.Insurance,
	)
}

func (h *AppointmentHandler) UpdateAppointment(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("handler: request received in UpdateAppointment with ID: %d", ID)

	appointment := model.Appointment{}
	if err := c.Bind(&appointment); err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	if err := h.logic.UpdateAppointment(ID, &appointment); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessAppointmentUpdated, http.StatusOK, nil)
}

func (h *AppointmentHandler) DeleteAppointment(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("handler: request received in DeleteAppointment with ID: %d", ID)

	if err := h.logic.DeleteAppointment(ID); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessAppointmentDeleted, http.StatusOK, nil)
}
