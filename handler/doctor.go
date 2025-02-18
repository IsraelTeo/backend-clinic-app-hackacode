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

type DoctorHandler struct {
	logic logic.DoctorLogic
}

func NewDoctorHandler(logic logic.DoctorLogic) *DoctorHandler {
	return &DoctorHandler{logic: logic}
}

func (h *DoctorHandler) GetDoctorByID(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("handler: doctor fetching with ID: %d", ID)

	doctor, err := h.logic.GetDoctorByID(ID)
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
		Data:    doctor,
	})
}

func (h *DoctorHandler) GetDoctorByDNI(c echo.Context) error {
	DNI := c.QueryParam("dni")
	if DNI == "" {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorPatientDNIRequired.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("handler: request received to fetch doctor with DNI: %s", DNI)

	doctor, err := h.logic.GetDoctorByDNI(DNI)
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
		Data:    doctor,
	})
}

func (h *DoctorHandler) GetAllDoctors(c echo.Context) error {
	log.Println("handler: request received in GetAllDoctors")

	doctors, err := h.logic.GetAllDoctors()
	if len(doctors) == 0 {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.SuccessPatiensListEmpty,
			Status:  http.StatusOK,
			Data:    []model.Doctor{},
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
		Message: response.SuccessDoctorsFound,
		Status:  http.StatusOK,
		Data:    doctors,
	})
}

func (h *DoctorHandler) CreateDoctor(c echo.Context) error {
	log.Println("handler: request received in CreateDoctor")

	doctor := model.Doctor{}
	if err := c.Bind(&doctor); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := c.Validate(&doctor); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorBadRequest.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := h.logic.CreateDoctor(&doctor); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessDoctorCreated,
		Status:  http.StatusCreated,
		Data:    nil,
	})
}

func (h *DoctorHandler) UpdateDoctor(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("handler: request received in UpdateDoctor with ID: %d", ID)

	doctor := model.Doctor{}
	if err := c.Bind(&doctor); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := c.Validate(&doctor); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.ErrorBadRequest.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := h.logic.UpdateDoctor(ID, &doctor); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessDoctorUpdated,
		Status:  http.StatusOK,
		Data:    nil,
	})
}

func (h *DoctorHandler) DeleteDoctor(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("handler: request received in DeleteDoctor with ID: %d", ID)

	if err := h.logic.DeleteDoctor(ID); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessDoctorDeleted,
		Status:  http.StatusOK,
		Data:    nil,
	})
}
