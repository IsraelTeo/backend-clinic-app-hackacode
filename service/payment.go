package logic

import (
	"fmt"
	"log"
	"os"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type PaymentLogic interface {
	PaymentRegister(payment *model.Payment) (*model.PaymentResponse, error)
}

type paymentLogic struct {
	repositoryAppointmentMain repository.AppointmentRepository
}

func NewPaymentLogic(repositoryAppointmentMain repository.AppointmentRepository,
) PaymentLogic {
	return &paymentLogic{
		repositoryAppointmentMain: repositoryAppointmentMain,
	}
}

func (l paymentLogic) PaymentRegister(payment *model.Payment) (*model.PaymentResponse, error) {
	appointment, err := l.repositoryAppointmentMain.GetByID(payment.AppoimentID)
	if err != nil {
		log.Printf("appointment: Error fetching appointment with ID %d: %v", payment.AppoimentID, err)
		return nil, response.ErrorAppointmentNotFound
	}

	if !payment.Paid {
		return nil, response.ErrorPaidNotTrue
	}

	if payment.TotalAmount == 0 {
		return nil, response.ErrorTotalAmountEmpty
	}

	if payment.TotalAmount < appointment.TotalAmount {
		return nil, response.ErrorTotalAmountBadRequest
	}

	//métodos de pago
	validPaymentTypes := map[model.PaymentType]bool{
		"tarjeta":     true,
		"applicativo": true,
		"efectivo":    true,
	}

	// Variable para almacenar si el método de pago de pago existe
	isValid := false

	// Recorrer el mapa para verificar si el tipo de pago existe
	for key, value := range validPaymentTypes {
		if value == validPaymentTypes[key] {
			isValid = true
			break
		}
	}

	if !isValid {
		log.Println("payment: Error invalid payment type")
		return nil, response.ErrorInvalidPaymentType
	}

	err = l.repositoryAppointmentMain.UpdatePaid(payment.AppoimentID)
	if err != nil {
		log.Printf("payment: Error updating paid status for appointment ID %d: %v", payment.AppoimentID, err)
		return nil, response.ErrorToUpdatePaid
	}

	qrCodePath, err := GenerateQRCode(appointment, payment)
	if err != nil {
		log.Printf("payment: Error generating QR code for appointment ID %d: %v", payment.AppoimentID, err)
		return nil, response.ErrorGeneratingQRCode
	}

	pdfReceiptPath, err := GeneratePDFReceipt(appointment, payment, qrCodePath)
	if err != nil {
		log.Printf("payment: Error generating PDF receipt for appointment ID %d: %v", payment.AppoimentID, err)
		return nil, response.ErrorGeneratingPDF
	}

	paymentResponse := model.PaymentResponse{
		QRCode:     qrCodePath,
		PDFReceipt: pdfReceiptPath,
	}

	return &paymentResponse, nil
}

func GenerateQRCode(appointment *model.Appointment, payment *model.Payment) (string, error) {
	qrData := fmt.Sprintf(
		"Appointment ID: %d\nPatient: %s %s\nDate: %s\nStart Time: %s\nEnd Time: %s\nPaid: %v",
		appointment.ID, appointment.Patient.Name, appointment.Patient.LastName,
		appointment.Date, appointment.StartTime, appointment.EndTime, payment.Paid,
	)

	qrCodePath := fmt.Sprintf("qrcodes/appointment_%d.png", payment.AppoimentID)

	err := qrcode.WriteFile(qrData, qrcode.Medium, 256, qrCodePath)
	if err != nil {
		log.Printf("Error generating QR code: %v", err)
		return "", err
	}

	return qrCodePath, nil
}

func GeneratePDFReceipt(appointment *model.Appointment, payment *model.Payment, qrCodePath string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Recibo de Cita Médica")
	pdf.Ln(12)

	// Información de la cita
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("ID de Cita: %d", appointment.ID))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Paciente: %s %s", appointment.Patient.Name, appointment.Patient.LastName))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Fecha: %s", appointment.Date))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Hora de Inicio: %s", appointment.StartTime))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Hora de Fin: %s", appointment.EndTime))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Monto Total: %.2f", payment.TotalAmount))
	pdf.Ln(12)

	// Adjuntar QR
	qrImage, err := os.Open(qrCodePath)
	if err == nil {
		imgOpts := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: false}
		pdf.RegisterImageOptionsReader("qr", imgOpts, qrImage)
		pdf.Image("qr", 10, pdf.GetY(), 50, 50, false, "qr", 0, "")
		qrImage.Close()
	} else {
		log.Printf("Error attaching QR code to PDF: %v", err)
	}

	// Guardar PDF
	pdfPath := fmt.Sprintf("receipts/receipt_%d.pdf", payment.AppoimentID)
	err = pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return "", err
	}

	return pdfPath, nil
}
