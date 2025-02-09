package routes

import (
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/auth"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/db"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/handler"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/logic"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/labstack/echo/v4"
)

const (
	idPath    = "/:id"
	voidPath  = ""
	loginPath = "/login"
)

func InitEnpoints(e *echo.Echo) {
	api := e.Group("/api/v1")
	setUpService(api)
	setUpPackage(api)
	setUpDoctor(api)
	setUpPatient(api)
	setUpAuth(api)
	setUpAppointment(api)
	setUpPayment(api)
}

func setUpAuth(api *echo.Group) {
	authRepository := repository.NewUserRepository(db.GDB)
	authLogic := auth.NewLoginService(authRepository)

	auth := api.Group("/auth")

	auth.POST(loginPath, authLogic.Login)
}

func setUpService(api *echo.Group) {
	serviceRepository := repository.NewRepository[model.Service](db.GDB)
	serviceLogic := logic.NewServiceLogic(serviceRepository)
	serviceHandler := handler.NewServiceHandler(serviceLogic)

	service := api.Group("/services")

	service.GET(idPath, serviceHandler.GetServiceByID)
	service.GET(voidPath, serviceHandler.GetAllServices)
	service.POST(voidPath, serviceHandler.CreateService)
	service.PUT(idPath, serviceHandler.UpdateService)
	service.DELETE(idPath, serviceHandler.DeleteService)
}

func setUpPackage(api *echo.Group) {
	packageRepository := repository.NewRepository[model.Package](db.GDB)
	packageRepositoryMain := repository.NewRepository[model.Service](db.GDB)
	packageLogic := logic.NewPackageLogic(packageRepository, packageRepositoryMain)
	serviceLogic := logic.NewServiceLogic(packageRepositoryMain)
	packageHandler := handler.NewPackageHandler(packageLogic, serviceLogic)

	packageServices := api.Group("/packages")

	packageServices.GET(idPath, packageHandler.GetPackageByID)
	packageServices.GET(voidPath, packageHandler.GetAllPackages)
	packageServices.POST(voidPath, packageHandler.CreatePackage)
	packageServices.PUT(idPath, packageHandler.UpdatePackage)
	packageServices.DELETE(idPath, packageHandler.DeletePackage)
}

func setUpDoctor(api *echo.Group) {
	doctorRepository := repository.NewRepository[model.Doctor](db.GDB)
	doctorLogic := logic.NewDoctorLogic(doctorRepository)
	doctorHandler := handler.NewDoctorHandler(doctorLogic)

	doctor := api.Group("/doctors")

	doctor.GET(idPath, doctorHandler.GetDoctorByID)
	doctor.GET(voidPath, doctorHandler.GetAllDoctors)
	doctor.POST(voidPath, doctorHandler.CreateDoctor)
	doctor.PUT(idPath, doctorHandler.UpdateDoctor)
	doctor.DELETE(idPath, doctorHandler.DeleteDoctor)
}

func setUpPatient(api *echo.Group) {
	patientRepository := repository.NewRepository[model.Patient](db.GDB)
	patientLogic := logic.NewPatientLogic(patientRepository)
	patientHandler := handler.NewPatientHandler(patientLogic)

	patient := api.Group("/patients")

	patient.GET(idPath, patientHandler.GetPatientByID)
	patient.GET(voidPath, patientHandler.GetAllPatients)
	patient.POST(voidPath, patientHandler.CreatePatient)
	patient.PUT(idPath, patientHandler.UpdatePatient)
	patient.DELETE(idPath, patientHandler.DeletePatient)
}

func setUpAppointment(api *echo.Group) {
	appointmentRepository := repository.NewRepository[model.Appointment](db.GDB)
	appointmentRepositoryMain := repository.NewAppointmentRepository(db.GDB)
	doctorRepository := repository.NewRepository[model.Doctor](db.GDB)
	patientRepository := repository.NewRepository[model.Patient](db.GDB)
	patientLogic := logic.NewPatientLogic(patientRepository)
	packageRepository := repository.NewRepository[model.Package](db.GDB)
	serviceRepository := repository.NewRepository[model.Service](db.GDB)
	patientRepositoryMain := repository.NewPatientRepository(db.GDB)
	packageRepositoryMain := repository.NewPackageRepository(db.GDB)
	appointmentLogic := logic.NewAppointmentLogic(
		appointmentRepository,
		doctorRepository,
		patientRepository,
		appointmentRepositoryMain,
		packageRepository,
		serviceRepository,
		patientLogic,
		patientRepositoryMain,
		packageRepositoryMain,
	)
	appointmentHandler := handler.NewAppointmentHandler(appointmentLogic)

	appointment := api.Group("/appointments")

	appointment.GET(idPath, appointmentHandler.GetAppointmentByID)
	appointment.GET(voidPath, appointmentHandler.GetAllAppointments)
	appointment.POST(voidPath, appointmentHandler.CreateAppointment)
	appointment.PUT(idPath, appointmentHandler.UpdateAppointment)
	appointment.DELETE(idPath, appointmentHandler.DeleteAppointment)
}

func setUpPayment(api *echo.Group) {
	paymentRepositoryMain := repository.NewAppointmentRepository(db.GDB)
	paymentLogic := logic.NewPaymentLogic(paymentRepositoryMain)
	paymentHandler := handler.NewPaymentHandler(paymentLogic)

	payment := api.Group("/payment/register")

	payment.POST(voidPath, paymentHandler.PaymentRegister)
}
