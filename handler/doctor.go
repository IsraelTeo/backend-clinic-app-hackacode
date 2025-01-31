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
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("handler: doctor fetching with ID: %d", ID)

	doctor, err := h.logic.GetDoctorByID(ID)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}

	return response.WriteSuccess(c, response.SuccessDoctorFound, http.StatusOK, doctor)
}

func (h *DoctorHandler) GetAllDoctors(c echo.Context) error {
	log.Println("handler: request received in GetAllDoctors")

	doctors, err := h.logic.GetAllDoctors()
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}

	return response.WriteSuccess(c, response.SuccessDoctorsFound, http.StatusOK, doctors)
}

func (h *DoctorHandler) CreateDoctor(c echo.Context) error {
	log.Println("handler: request received in CreateDoctor")

	doctor := model.Doctor{}
	if err := c.Bind(&doctor); err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	if err := h.logic.CreateDoctor(&doctor); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessDoctorCreated, http.StatusCreated, nil)
}

func (h *DoctorHandler) UpdateDoctor(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("handler: request received in UpdateDoctor with ID: %d", ID)

	doctor := model.Doctor{}
	if err := c.Bind(&doctor); err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	if err := h.logic.UpdateDoctor(ID, &doctor); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessDoctorUpdated, http.StatusOK, nil)
}

func (h *DoctorHandler) DeleteDoctor(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("handler: request received in DeleteDoctor with ID: %d", ID)

	if err := h.logic.DeleteDoctor(ID); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessDoctorDeleted, http.StatusOK, nil)
}
