package handler

import (
	"log"
	"net/http"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/logic"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	logic logic.PaymentLogic
}

func NewPaymentHandler(logic logic.PaymentLogic) *PaymentHandler {
	return &PaymentHandler{logic: logic}
}

func (h *PaymentHandler) PaymentRegister(c echo.Context) error {
	log.Println("handler: request received in PaymentRegister")

	payment := model.Payment{}

	if err := c.Bind(&payment); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	paymentResponse, err := h.logic.PaymentRegister(&payment)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: "Error processing payment",
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessPaymentRegister,
		Status:  http.StatusCreated,
		Data:    paymentResponse,
	})
}
