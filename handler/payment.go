package handler

import (
	"log"
	"net/http"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/logic"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
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
		return response.WriteError(c, "Invalid request payload", http.StatusBadRequest)
	}

	paymentResponse, err := h.logic.PaymentRegister(&payment)
	if err != nil {
		return response.WriteError(c, "Error processing payment", http.StatusInternalServerError)
	}

	return response.WriteSuccessPayment(c, response.SuccessPaymentRegister, http.StatusCreated, *paymentResponse)

}
