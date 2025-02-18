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

type PackageHandler struct {
	logicPkg  logic.PackageLogic
	logicServ logic.ServiceLogic
}

func NewPackageHandler(logicPkg logic.PackageLogic, logicServ logic.ServiceLogic) *PackageHandler {
	return &PackageHandler{logicPkg: logicPkg, logicServ: logicServ}
}

func (h *PackageHandler) GetPackageByID(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("package-handler: package fetching with ID: %d", ID)

	packageService, err := h.logicPkg.GetPackageByID(uint(ID))
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
		Message: response.SuccessPackageFound,
		Status:  http.StatusOK,
		Data:    packageService,
	})
}

func (h *PackageHandler) GetAllPackages(c echo.Context) error {
	log.Println("package-handler: request received in GetAllPackages")

	packageServices, err := h.logicPkg.GetAllPackages()
	if len(packageServices) == 0 {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: response.SuccessPatiensListEmpty,
			Status:  http.StatusNoContent,
			Data:    []model.Package{},
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
		Message: response.SuccessPackagesFound,
		Status:  http.StatusOK,
		Data:    packageServices,
	})
}

func (h *PackageHandler) CreatePackage(c echo.Context) error {
	log.Println("package-handler: request received in CreatePackage")

	var pkgRequest model.CreatePackageRequest
	if err := c.Bind(&pkgRequest); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if len(pkgRequest.ServiceIDs) == 0 {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: "Service IDs cannot be empty",
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := h.logicPkg.CreatePackage(&pkgRequest); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessPackageCreated,
		Status:  http.StatusCreated,
		Data:    nil,
	})
}

func (h *PackageHandler) UpdatePackage(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("package-handler: request received in UpdatePackage with ID: %d", ID)

	packageServices := model.CreatePackageRequest{}
	if err := c.Bind(&packageServices); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	if err := h.logicPkg.UpdatePackage(uint(ID), &packageServices); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessPackageUpdated,
		Status:  http.StatusOK,
		Data:    nil,
	})
}

func (h *PackageHandler) DeletePackage(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusBadRequest,
			Data:    nil,
		})
	}

	log.Printf("package-handler: request received in DeletePackage with ID: %d", ID)

	if err := h.logicPkg.DeletePackage(uint(ID)); err != nil {
		return response.WriteError(&response.WriteResponse{
			C:       c,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
			Data:    nil,
		})
	}

	return response.WriteSuccess(&response.WriteResponse{
		C:       c,
		Message: response.SuccessPackageDeleted,
		Status:  http.StatusOK,
		Data:    nil,
	})
}
