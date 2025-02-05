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
	packageRepository3 := repository.NewRepository[model.Service](db.GDB)
	packageLogic := logic.NewPackageLogic(packageRepository, packageRepository3)
	serviceLogic := logic.NewServiceLogic(packageRepository3)
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
