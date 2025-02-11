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

	service.GET(idPath, auth.ValidateJWT(serviceHandler.GetServiceByID))
	service.GET(voidPath, auth.ValidateJWT(serviceHandler.GetAllServices))
	service.POST(voidPath, auth.ValidateJWT(serviceHandler.CreateService))
	service.PUT(idPath, auth.ValidateJWT(serviceHandler.UpdateService))
	service.DELETE(idPath, auth.ValidateJWT(serviceHandler.DeleteService))
}

func setUpPackage(api *echo.Group) {
	packageRepository := repository.NewRepository[model.Package](db.GDB)
	packageRepositoryMain := repository.NewRepository[model.Service](db.GDB)
	packageLogic := logic.NewPackageLogic(packageRepository, packageRepositoryMain)
	serviceLogic := logic.NewServiceLogic(packageRepositoryMain)
	packageHandler := handler.NewPackageHandler(packageLogic, serviceLogic)

	packageServices := api.Group("/packages")

	packageServices.GET(idPath, auth.ValidateJWT(packageHandler.GetPackageByID))
	packageServices.GET(voidPath, auth.ValidateJWT(packageHandler.GetAllPackages))
	packageServices.POST(voidPath, auth.ValidateJWT(packageHandler.CreatePackage))
	packageServices.PUT(idPath, auth.ValidateJWT(packageHandler.UpdatePackage))
	packageServices.DELETE(idPath, auth.ValidateJWT(packageHandler.DeletePackage))
}

func setUpDoctor(api *echo.Group) {
	doctorRepository := repository.NewRepository[model.Doctor](db.GDB)
	doctorRepositoryMain := repository.NewDoctorRepository(db.GDB)
	doctorLogic := logic.NewDoctorLogic(doctorRepository, doctorRepositoryMain)
	doctorHandler := handler.NewDoctorHandler(doctorLogic)

	doctor := api.Group("/doctors")

	doctor.GET(idPath, auth.ValidateJWT(doctorHandler.GetDoctorByID))
	doctor.GET(voidPath, auth.ValidateJWT(doctorHandler.GetAllDoctors))
	doctor.POST(voidPath, auth.ValidateJWT(doctorHandler.CreateDoctor))
	doctor.PUT(idPath, auth.ValidateJWT(doctorHandler.UpdateDoctor))
	doctor.DELETE(idPath, auth.ValidateJWT(doctorHandler.DeleteDoctor))
}

func setUpPatient(api *echo.Group) {
	patientRepository := repository.NewRepository[model.Patient](db.GDB)
	patientRepositoryMain := repository.NewPatientRepository(db.GDB)
	appointmentRepositoryMain := repository.NewAppointmentRepository(db.GDB)
	patientLogic := logic.NewPatientLogic(patientRepository, patientRepositoryMain, appointmentRepositoryMain)
	patientHandler := handler.NewPatientHandler(patientLogic)

	patient := api.Group("/patients")

	patient.GET(idPath, auth.ValidateJWT(patientHandler.GetPatientByID))
	patient.GET(voidPath, auth.ValidateJWT(patientHandler.GetPatientByDNI))
	patient.GET("dni", auth.ValidateJWT(patientHandler.GetAllPatients))
	patient.POST(voidPath, auth.ValidateJWT(patientHandler.CreatePatient))
	patient.PUT(idPath, auth.ValidateJWT(patientHandler.UpdatePatient))
	patient.DELETE(idPath, auth.ValidateJWT(patientHandler.DeletePatient))
}

func setUpAppointment(api *echo.Group) {
	appointmentRepository := repository.NewRepository[model.Appointment](db.GDB)
	appointmentRepositoryMain := repository.NewAppointmentRepository(db.GDB)
	doctorRepository := repository.NewRepository[model.Doctor](db.GDB)
	patientRepository := repository.NewRepository[model.Patient](db.GDB)
	patientRepositoryMain := repository.NewPatientRepository(db.GDB)
	patientLogic := logic.NewPatientLogic(patientRepository, patientRepositoryMain, appointmentRepositoryMain)
	packageRepository := repository.NewRepository[model.Package](db.GDB)
	serviceRepository := repository.NewRepository[model.Service](db.GDB)
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

	appointment.GET(idPath, auth.ValidateJWT(appointmentHandler.GetAppointmentByID))
	appointment.GET(voidPath, auth.ValidateJWT(appointmentHandler.GetAllAppointments))
	appointment.POST(voidPath, auth.ValidateJWT(appointmentHandler.CreateAppointment))
	appointment.PUT(idPath, auth.ValidateJWT(appointmentHandler.UpdateAppointment))
	appointment.DELETE(idPath, auth.ValidateJWT(appointmentHandler.DeleteAppointment))
}

func setUpPayment(api *echo.Group) {
	paymentRepositoryMain := repository.NewAppointmentRepository(db.GDB)
	paymentLogic := logic.NewPaymentLogic(paymentRepositoryMain)
	paymentHandler := handler.NewPaymentHandler(paymentLogic)

	payment := api.Group("/payment/register")

	payment.POST(voidPath, auth.ValidateJWT(paymentHandler.PaymentRegister))
}
