package handler

import (
	"fmt"
	"log"
	"net/http"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/logic"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/validate"
	"github.com/labstack/echo/v4"
)

type PatientHandler struct {
	logic logic.PatientLogic
}

func NewPatientHandler(logic logic.PatientLogic) *PatientHandler {
	return &PatientHandler{logic: logic}
}

func (h *PatientHandler) GetPatientByID(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("patient-handler: patient fetching with ID: %d", ID)

	patient, err := h.logic.GetPatientByID(uint(ID))
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}

	return response.WriteSuccess(c, response.SuccessPatientFound, http.StatusOK, patient)
}

func (h *PatientHandler) GetPatientByDNI(c echo.Context) error {
	DNI := c.QueryParam("dni")
	if DNI == "" {
		return response.WriteError(c, response.ErrorDoctorDNIRequired.Error(), http.StatusBadRequest)
	}

	log.Printf("patient-handler: request received to fetch patient with DNI: %s", DNI)

	patient, err := h.logic.GetPatientByDNI(DNI)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}

	return response.WriteSuccess(c, response.SuccessDoctorFound, http.StatusOK, patient)
}

func (h *PatientHandler) GetAllPatients(c echo.Context) error {
	log.Println("patient-handler: request received in GetAllPatients")

	patients, err := h.logic.GetAllPatients()
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}

	if len(patients) == 0 {
		return response.WriteError(c, response.SuccessPatiensListEmpty, http.StatusNoContent)
	}

	return response.WriteSuccess(c, response.SuccessPatientsFound, http.StatusOK, patients)
}

func (h *PatientHandler) CreatePatient(c echo.Context) error {
	log.Println("patient-handler: request received in CreatePatient")

	patient := model.Patient{}
	if err := c.Bind(&patient); err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			return response.WriteError(c, fmt.Sprintf("Error al procesar los datos: %s", he.Message), http.StatusBadRequest)
		}

		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	if err := c.Validate(&patient); err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	if err := h.logic.CreatePatient(&patient); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessPatientCreated, http.StatusCreated, nil)
}

func (h *PatientHandler) UpdatePatient(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("patient-handler: request received in UpdatePatient with ID: %d", ID)

	patient := model.Patient{}
	if err := c.Bind(&patient); err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			return response.WriteError(c, fmt.Sprintf("Error al procesar los datos: %s", he.Message), http.StatusBadRequest)
		}

		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	if err := c.Validate(&patient); err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	if err := h.logic.UpdatePatient(uint(ID), &patient); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessPatientUpdated, http.StatusOK, nil)
}

func (h *PatientHandler) DeletePatient(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("patient-handler: request received in DeletePatient with ID: %d", ID)

	if err := h.logic.DeletePatient(uint(ID)); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessPatientDeleted, http.StatusOK, nil)
}
