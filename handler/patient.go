package handler

import (
	"log"
	"net/http"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/logic"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/validate"
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
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("patient-handler: patient fetching with ID: %d", ID)

	patient, err := h.logic.GetPatientByID(uint(ID))
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
		Message: response.SuccessPatientFound,
		Status:  http.StatusOK,
		Data:    patient,
	})
}

func (h *PatientHandler) GetPatientByDNI(c echo.Context) error {
	DNI := c.QueryParam("dni")
	if DNI == "" {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorDoctorDNIRequired.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("patient-handler: request received to fetch patient with DNI: %s", DNI)

	patient, err := h.logic.GetPatientByDNI(DNI)
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
		Message: response.SuccessDoctorFound,
		Status:  http.StatusOK,
		Data:    patient,
	})
}

func (h *PatientHandler) GetAllPatients(c echo.Context) error {
	log.Println("patient-handler: request received in GetAllPatients")

	patients, err := h.logic.GetAllPatients()
	if len(patients) == 0 {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.SuccessPatiensListEmpty,
			Status:  http.StatusOK,
			Data:    []model.Patient{},
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
		Message: response.SuccessPatientsFound,
		Status:  http.StatusOK,
		Data:    patients,
	})

}

func (h *PatientHandler) CreatePatient(c echo.Context) error {
	log.Println("patient-handler: request received in CreatePatient")

	patient := model.Patient{}
	if err := c.Bind(&patient); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := c.Validate(&patient); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := h.logic.CreatePatient(&patient); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessPatientCreated,
		Status:  http.StatusCreated,
		Data:    nil,
	})
}

func (h *PatientHandler) UpdatePatient(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("patient-handler: request received in UpdatePatient with ID: %d", ID)

	patient := model.Patient{}
	if err := c.Bind(&patient); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := c.Validate(&patient); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := h.logic.UpdatePatient(uint(ID), &patient); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessPatientUpdated,
		Status:  http.StatusOK,
		Data:    nil,
	})
}

func (h *PatientHandler) DeletePatient(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("patient-handler: request received in DeletePatient with ID: %d", ID)

	if err := h.logic.DeletePatient(uint(ID)); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessPatientDeleted,
		Status:  http.StatusOK,
		Data:    nil,
	})
}
