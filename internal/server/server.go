package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	basyxAas "hiroyoshii/go-aas-proxy/gen/go"
	"hiroyoshii/go-aas-proxy/internal/aas"
	"hiroyoshii/go-aas-proxy/internal/submodel"

	"github.com/labstack/echo/v4"
)

type Server struct {
	aasCli      aas.Aas
	submodelCli submodel.Submodel
}

// Retrieves all Asset Administration Shells from the Asset Administration Shell repository
// (GET /shells)
func (s Server) GetAllAssetAdministrationShells(ctx echo.Context) error {
	b, err := s.aasCli.List()
	if err != nil {
		log.Printf("failed to retrieved aas list: %v\n", err)
		return err
	}
	var res []*basyxAas.AssetAdministrationShell
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Printf("failed to unmarshal: %v\n", err)
		return err
	}
	return ctx.JSON(http.StatusOK, res)
}

// Deletes a specific Asset Administration Shell at the Asset Administration Shell repository
// (DELETE /shells/{aasId})
func (s Server) DeleteAssetAdministrationShellById(ctx echo.Context, aasId string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Retrieves a specific Asset Administration Shell from the Asset Administration Shell repository
// (GET /shells/{aasId})
func (s Server) GetShellsAasId(ctx echo.Context, aasId string) error {
	ctx.Logger().Info(ctx.Path())
	b, err := s.aasCli.Get(aasId)
	if err != nil {
		log.Printf("failed to retrieved aas: %v\n", err)
		return err
	}
	if len(b) == 0 {
		return ctx.JSON(http.StatusNotFound, "No Asset Administration Shell found")
	}
	var res *basyxAas.AssetAdministrationShell
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Printf("failed to unmarshal: %v\n", err)
		return err
	}
	return ctx.JSON(http.StatusOK, res)
}

// Creates or updates a Asset Administration Shell at the Asset Administration Shell repository
// (PUT /shells/{aasId})
func (s Server) PutAssetAdministrationShell(ctx echo.Context, aasId string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Retrieves a specific Asset Administration Shell from the Asset Administration Shell repository
// (GET /shells/{aasId}/aas)
func (s Server) GetAssetAdministrationShellById(ctx echo.Context, aasId string) error {
	return s.GetShellsAasId(ctx, aasId)
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
	submodelSemanticId, err := s.aasCli.GetSubmodel(aasId, submodelIdShort)
	if err != nil {
		log.Printf("failed to retrieved submodel semantic id: %v\n", err)
		return err
	}
	b, err := s.submodelCli.Get(submodelSemanticId, submodelIdShort)
	if err != nil {
		log.Printf("failed to retrieved submodel: %v\n", err)
		return err
	}
	if len(b) == 0 {
		return ctx.JSON(http.StatusNotFound, "No Submodel found")
	}
	var res *basyxAas.Submodel
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Printf("failed to unmarshal: %v\n", err)
		return err
	}
	return ctx.JSON(http.StatusOK, res)
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
	aCli, err := aas.NewAas()
	if err != nil {
		return nil, err
	}
	smCli, err := submodel.NewSubmodel()
	server := Server{aasCli: aCli, submodelCli: smCli}

	// 自動生成されたハンドラ登録関数にServerInterfaceを満たすserverを渡す
	basyxAas.RegisterHandlers(instance, server)
	return instance, nil
}
