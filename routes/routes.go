package routes

import (
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/db"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/handler"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/logic"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/labstack/echo/v4"
)

const (
	idPath   = "/:id"
	voidPath = ""
)

func InitEnpoints(e *echo.Echo) {
	api := e.Group("/api/v1")
	setUpService(api)
	setUpPackage(api)
}

func setUpService(api *echo.Group) {
	serviceRepository := repository.NewRepository[model.Service](db.GDB)
	serviceLogic := logic.NewServiceLogic(serviceRepository)
	serviceHandler := handler.NewServiceHandler(serviceLogic)

	service := api.Group("/service")

	service.GET(idPath, serviceHandler.GetServiceByID)
	service.GET(voidPath, serviceHandler.GetAllServices)
	service.POST(voidPath, serviceHandler.CreateService)
	service.PUT(idPath, serviceHandler.UpdateService)
	service.DELETE(idPath, serviceHandler.DeleteService)
}

func setUpPackage(api *echo.Group) {
	packageRepository := repository.NewRepository[model.Package](db.GDB)
	packageLogic := logic.NewPackageLogic(packageRepository)
	packageHandler := handler.NewPackageHandler(packageLogic)

	packageServices := api.Group("/package")

	packageServices.GET(idPath, packageHandler.GetPackageByID)
	packageServices.GET(voidPath, packageHandler.GetAllPackages)
	packageServices.POST(voidPath, packageHandler.CreatePackage)
	packageServices.PUT(idPath, packageHandler.UpdatePackage)
	packageServices.DELETE(idPath, packageHandler.DeletePackage)
}
