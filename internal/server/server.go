package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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
		slog.Error(fmt.Sprintf("failed to retrieved aas list: %v", err))
		return err
	}
	var res []*basyxAas.AssetAdministrationShell
	err = json.Unmarshal(b, &res)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to unmarshal: %v", err))
		return err
	}
	return ctx.JSON(http.StatusOK, res)
}

// Deletes a specific Asset Administration Shell at the Asset Administration Shell repository
// (DELETE /shells/{aasId})
func (s Server) DeleteAssetAdministrationShellById(ctx echo.Context, aasId string) error {
	err := s.aasCli.Delete(aasId)
	if err != nil {
		return err
	}
	b := true
	return ctx.JSON(http.StatusOK, &basyxAas.Result{Success: &b})
}

// Retrieves a specific Asset Administration Shell from the Asset Administration Shell repository
// (GET /shells/{aasId})
func (s Server) GetShellsAasId(ctx echo.Context, aasId string) error {
	ctx.Logger().Info(ctx.Path())
	b, err := s.aasCli.Get(aasId)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to retrieved aas: %v", err))
		return err
	}
	if len(b) == 0 {
		return ctx.JSON(http.StatusNotFound, "No Asset Administration Shell found")
	}
	var res *basyxAas.AssetAdministrationShell
	err = json.Unmarshal(b, &res)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to unmarshal: %v", err))
		return err
	}
	return ctx.JSON(http.StatusOK, res)
}

// Creates or updates a Asset Administration Shell at the Asset Administration Shell repository
// (PUT /shells/{aasId})
func (s Server) PutAssetAdministrationShell(ctx echo.Context, aasId string) error {
	req, err := ctx.Request().GetBody()
	if err != nil {
		return err
	}
	b, err := io.ReadAll(req)
	if err != nil {
		return err
	}
	var aas basyxAas.AssetAdministrationShell
	err = json.Unmarshal(b, &aas)
	if err != nil {
		return err
	}
	_, err = s.aasCli.CreateOrUpdate(aasId, b)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, aas)
}

// Retrieves a specific Asset Administration Shell from the Asset Administration Shell repository
// (GET /shells/{aasId}/aas)
func (s Server) GetAssetAdministrationShellById(ctx echo.Context, aasId string) error {
	return s.GetShellsAasId(ctx, aasId)
}

// Retrieves all Submodels from the  Asset Administration Shell
// (GET /shells/{aasId}/aas/submodels)
func (s Server) ShellRepoGetSubmodelsFromShell(ctx echo.Context, aasId string) error {
	id2Semantic, err := s.aasCli.GetSubmodelIds(aasId, "")
	if err != nil {
		slog.Error(fmt.Sprintf("failed to retrieved submodel semantic id: %v", err))
		return err
	}
	res := []map[string]interface{}{}
	for id, semantic := range id2Semantic {
		b, err := s.submodelCli.Get(aasId, semantic, id)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to retrieved submodel: %v", err))
			return err
		}
		if len(b) == 0 {
			return ctx.JSON(http.StatusNotFound, fmt.Sprintf("Submodel(idShort: %s) is not found", id))
		}
		var r map[string]interface{}
		err = json.Unmarshal(b, &r)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to unmarshal: %v", err))
			return err
		}
		res = append(res, r)
	}
	return ctx.JSON(http.StatusOK, res)
}

// Deletes a specific Submodel from the Asset Administration Shell
// (DELETE /shells/{aasId}/aas/submodels/{submodelIdShort})
func (s Server) ShellRepoDeleteSubmodelFromShellByIdShort(ctx echo.Context, aasId string, submodelIdShort string) error {
	err := s.aasCli.DeleteSubmodel(aasId, submodelIdShort)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to retrieved submodel semantic id: %v", err))
		return err
	}
	b := true
	return ctx.JSON(http.StatusOK, &basyxAas.Result{Success: &b})
}

// Retrieves the Submodel from the Asset Administration Shell
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort})
func (s Server) GetShellsAasIdAasSubmodelsSubmodelIdShort(ctx echo.Context, aasId string, submodelIdShort string) error {
	res, err := s.getSubmodel(aasId, submodelIdShort)
	if err != nil {
		return ctx.JSON(err.(*HttpError).Status, err.Error())
	}
	return ctx.JSON(http.StatusOK, res)
}

func (s Server) getSubmodel(aasId string, submodelIdShort string) (map[string]interface{}, error) {
	id2Semantic, err := s.aasCli.GetSubmodelIds(aasId, submodelIdShort)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to retrieved submodel semantic id: %v", err))
		return nil, &HttpError{Status: http.StatusInternalServerError, Err: err}
	}
	if len(id2Semantic) != 1 {
		return nil, &HttpError{Status: http.StatusNotFound, Err: fmt.Errorf("No related submodel found")}
	}
	var res map[string]interface{}
	for _, semantic := range id2Semantic {
		b, err := s.submodelCli.Get(aasId, semantic, submodelIdShort)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to retrieved submodel: %v", err))
			return nil, &HttpError{Status: http.StatusInternalServerError, Err: err}
		}
		if len(b) == 0 {
			return nil, &HttpError{Status: http.StatusNotFound, Err: fmt.Errorf("No Submodel found")}
		}
		err = json.Unmarshal(b, &res)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to unmarshal: %v", err))
			return nil, &HttpError{Status: http.StatusInternalServerError, Err: err}
		}
	}
	return res, err
}

// Creates or updates a Submodel to an existing Asset Administration Shell
// (PUT /shells/{aasId}/aas/submodels/{submodelIdShort})
func (s Server) ShellRepoPutSubmodelToShell(ctx echo.Context, aasId string, submodelIdShort string) error {
	req, err := ctx.Request().GetBody()
	if err != nil {
		return err
	}
	b, err := io.ReadAll(req)
	if err != nil {
		return err
	}
	var subm *basyxAas.Submodel
	err = json.Unmarshal(b, &subm)
	if err != nil {
		return err
	}

	// check existence of submodel
	semanticId := ""
	if len(subm.SemanticId.Keys) == 1 {
		semanticId = subm.SemanticId.Keys[0].Value
	}
	sb, err := s.submodelCli.Get(aasId, semanticId, submodelIdShort)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to retrieved submodel: %v", err))
		return err
	}
	if len(sb) == 0 {
		return ctx.JSON(http.StatusNotFound, "No Submodel found")
	}
	err = s.aasCli.CreateOrUpdateSubmodel(aasId, subm.Identification.Id, semanticId, submodelIdShort)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, subm)
}

// Retrieves the Submodel from the Asset Administration Shell
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel)
func (s Server) ShellRepoGetSubmodelFromShellByIdShort(ctx echo.Context, aasId string, submodelIdShort string) error {
	return s.GetShellsAasIdAasSubmodelsSubmodelIdShort(ctx, aasId, submodelIdShort)
}

// Retrieves the minimized version of a Submodel, i.e. only the values of SubmodelElements are serialized and returned
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/values)
func (s Server) ShellRepoGetSubmodelValues(ctx echo.Context, aasId string, submodelIdShort string) error {
	res, err := s.getSubmodel(aasId, submodelIdShort)
	if err != nil {
		return ctx.JSON(err.(*HttpError).Status, err.Error())
	}
	result := map[string]interface{}{}
	for _, ses := range res["submodelElements"].([]interface{}) {
		e := ses.(map[string]interface{})
		recursive(e["idShort"].(string), e["value"], result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func recursive(key string, value interface{}, result map[string]interface{}) {
	if arr, ok := value.([]interface{}); ok {
		r := map[string]interface{}{}
		for _, a := range arr {
			v := a.(map[string]interface{})
			recursive(v["idShort"].(string), v["value"], r)
		}
		result[key] = r
	} else {
		result[key] = value
	}
}

// Retrieves all Submodel-Elements from the Submodel
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements)
func (s Server) ShellRepoGetSubmodelElements(ctx echo.Context, aasId string, submodelIdShort string) error {
	res, err := s.getSubmodel(aasId, submodelIdShort)
	if err != nil {
		return ctx.JSON(err.(*HttpError).Status, err.Error())
	}
	return ctx.JSON(http.StatusOK, res["submodelElements"])
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
	return ctx.JSON(http.StatusNotFound, nil)
}

// Creates or updates a Submodel-Element at the Submodel
// (PUT /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements/{seIdShortPath})
func (s Server) ShellRepoPutSubmodelElement(ctx echo.Context, aasId string, submodelIdShort string, seIdShortPath string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// Retrieves the value of a specific Submodel-Element from the Submodel
// (GET /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements/{seIdShortPath}/value)
func (s Server) ShellRepoGetSubmodelElementValueByIdShort(ctx echo.Context, aasId string, submodelIdShort string, seIdShortPath string) error {
	return ctx.JSON(http.StatusNotFound, nil)
}

// Updates the Submodel-Element's value
// (PUT /shells/{aasId}/aas/submodels/{submodelIdShort}/submodel/submodelElements/{seIdShortPath}/value)
func (s Server) ShellRepoPutSubmodelElementValueByIdShort(ctx echo.Context, aasId string, submodelIdShort string, seIdShortPath string) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
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
