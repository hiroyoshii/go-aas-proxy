package server

import (
	"context"
	"net/http"

	basyxAas "hiroyoshii/go-aas-proxy/gen/go"

	"github.com/labstack/echo/v4"
)

type Server struct{}

// Retrieves all Asset Administration Shells from the Asset Administration Shell repository
// (GET /shells)
func (s Server) GetAllAssetAdministrationShells(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, &basyxAas.AssetAdministrationShell{})
}

// Deletes a specific Asset Administration Shell at the Asset Administration Shell repository
// (DELETE /shells/{aasId})
func (s Server) DeleteAssetAdministrationShellById(ctx echo.Context, aasId string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Retrieves a specific Asset Administration Shell from the Asset Administration Shell repository
// (GET /shells/{aasId})
func (s Server) GetShellsAasId(ctx echo.Context, aasId string) error {
	return ctx.JSON(http.StatusOK, &basyxAas.AssetAdministrationShell{})
}

// Creates or updates a Asset Administration Shell at the Asset Administration Shell repository
// (PUT /shells/{aasId})
func (s Server) PutAssetAdministrationShell(ctx echo.Context, aasId string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Retrieves a specific Asset Administration Shell from the Asset Administration Shell repository
// (GET /shells/{aasId}/aas)
func (s Server) GetAssetAdministrationShellById(ctx echo.Context, aasId string) error {
	return ctx.JSON(http.StatusOK, &basyxAas.AssetAdministrationShell{})
}

// Retrieves all Submodels from the  Asset Administration Shell
// (GET /shells/{aasId}/aas/submodels)
func (s Server) ShellRepoGetSubmodelsFromShell(ctx echo.Context, aasId string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Deletes a specific Submodel from the Asset Administration Shell
// (DELETE /shells/{aasId}/aas/submodels/{submodelIdShort})
func (s Server) ShellRepoDeleteSubmodelFromShellByIdShort(ctx echo.Context, aasId string, submodelIdShort string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Retrieves the Submodel from the Asset Administration Shell
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort})
func (s Server) GetShellsAasIdAasSubmodelsSubmodelIdShort(ctx echo.Context, aasId string, submodelIdShort string) error {
	return ctx.JSON(http.StatusOK, &basyxAas.Submodel{})
}

// Creates or updates a Submodel to an existing Asset Administration Shell
// (PUT /shells/{aasId}/aas/submodels/{submodelIdShort})
func (s Server) ShellRepoPutSubmodelToShell(ctx echo.Context, aasId string, submodelIdShort string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Retrieves the Submodel from the Asset Administration Shell
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel)
func (s Server) ShellRepoGetSubmodelFromShellByIdShort(ctx echo.Context, aasId string, submodelIdShort string) error {
	return ctx.JSON(http.StatusOK, &basyxAas.Submodel{})
}

// Retrieves all Submodel-Elements from the Submodel
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements)
func (s Server) ShellRepoGetSubmodelElements(ctx echo.Context, aasId string, submodelIdShort string) error {
	return ctx.JSON(http.StatusOK, &basyxAas.SubmodelElement{})
}

// Retrieves the result of an asynchronously started operation
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements/{idShortPathToOperation}/invocationList/{requestId})
func (s Server) ShellRepoGetInvocationResultByIdShort(ctx echo.Context, aasId string, submodelIdShort string, idShortPathToOperation string, requestId string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Invokes a specific operation from the Submodel synchronously or asynchronously
// (POST /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements/{idShortPathToOperation}/invoke)
func (s Server) ShellRepoInvokeOperationByIdShort(ctx echo.Context, aasId string, submodelIdShort string, idShortPathToOperation string, params basyxAas.ShellRepoInvokeOperationByIdShortParams) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Deletes a specific Submodel-Element from the Submodel
// (DELETE /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements/{seIdShortPath})
func (s Server) ShellRepoDeleteSubmodelElementByIdShort(ctx echo.Context, aasId string, submodelIdShort string, seIdShortPath string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Retrieves a specific Submodel-Element from the Submodel
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements/{seIdShortPath})
func (s Server) ShellRepoGetSubmodelElementByIdShort(ctx echo.Context, aasId string, submodelIdShort string, seIdShortPath string) error {
	return ctx.JSON(http.StatusOK, &basyxAas.SubmodelElement{})
}

// Creates or updates a Submodel-Element at the Submodel
// (PUT /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements/{seIdShortPath})
func (s Server) ShellRepoPutSubmodelElement(ctx echo.Context, aasId string, submodelIdShort string, seIdShortPath string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Retrieves the value of a specific Submodel-Element from the Submodel
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements/{seIdShortPath}/value)
func (s Server) ShellRepoGetSubmodelElementValueByIdShort(ctx echo.Context, aasId string, submodelIdShort string, seIdShortPath string) error {
	return ctx.JSON(http.StatusOK, &basyxAas.SubmodelElement{})
}

// Updates the Submodel-Element's value
// (PUT /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements/{seIdShortPath}/value)
func (s Server) ShellRepoPutSubmodelElementValueByIdShort(ctx echo.Context, aasId string, submodelIdShort string, seIdShortPath string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Retrieves the minimized version of a Submodel, i.e. only the values of SubmodelElements are serialized and returned
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/values)
func (s Server) ShellRepoGetSubmodelValues(ctx echo.Context, aasId string, submodelIdShort string) error {
	return ctx.JSON(http.StatusOK, &basyxAas.Submodel{})
}

func NewServer(ctx context.Context) (*echo.Echo, error) {
	instance := echo.New()
	server := Server{}

	// 自動生成されたハンドラ登録関数にServerInterfaceを満たすserverを渡す
	basyxAas.RegisterHandlers(instance, server)
	return instance, nil
}
