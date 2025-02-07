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
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("handler: package fetching with ID: %d", ID)

	packageService, err := h.logicPkg.GetPackageByID(uint(ID))
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}
	return response.WriteSuccess(c, response.SuccessPackageFound, http.StatusOK, packageService)
}

func (h *PackageHandler) GetAllPackages(c echo.Context) error {
	log.Println("handler: request received in GetAllPackages")

	packageServices, err := h.logicPkg.GetAllPackages()
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusNotFound)
	}

	return response.WriteSuccess(c, response.SuccessPackagesFound, http.StatusOK, packageServices)
}

func (h *PackageHandler) CreatePackage(c echo.Context) error {
	log.Println("handler: request received in CreatePackage")

	var pkgRequest model.CreatePackageRequest
	if err := c.Bind(&pkgRequest); err != nil {
		return response.WriteError(c, "Invalid request format", http.StatusBadRequest)
	}

	// Validar la entrada
	if len(pkgRequest.ServiceIDs) == 0 {
		return response.WriteError(c, "No service IDs provided", http.StatusBadRequest)
	}

	// Pasar el objeto al servicio de lógica
	if err := h.logicPkg.CreatePackage(&pkgRequest); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessPackageCreated, http.StatusCreated, nil)
}

func (h *PackageHandler) UpdatePackage(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("handler: request received in UpdatePackage with ID: %d", ID)

	packageServices := model.CreatePackageRequest{}
	if err := c.Bind(&packageServices); err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	if err := h.logicPkg.UpdatePackage(uint(ID), &packageServices); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessPackageUpdated, http.StatusOK, nil)
}

func (h *PackageHandler) DeletePackage(c echo.Context) error {
	ID, err := validate.ParseID(c)
	if err != nil {
		return response.WriteError(c, err.Error(), http.StatusBadRequest)
	}

	log.Printf("handler: request received in DeletePackage with ID: %d", ID)

	if err := h.logicPkg.DeletePackage(uint(ID)); err != nil {
		return response.WriteError(c, err.Error(), http.StatusInternalServerError)
	}

	return response.WriteSuccess(c, response.SuccessPackageDeleted, http.StatusOK, nil)
}
