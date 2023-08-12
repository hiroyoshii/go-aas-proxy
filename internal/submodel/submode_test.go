package submodel

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var (
	sy submodelYaml
)

func TestGetSubmodel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("select(.+)").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "manufacturername"}).
				AddRow(1, "Manufacturername1"))

	s, err := newSubmodelWithDb(sy, map[string]*sql.DB{"submodel1": db})
	assert.Nil(t, err)
	b, err := s.Get("454576463545648365874", "https://www.hsu-hh.de/aut/aas/identification", "Identification")
	assert.Nil(t, err)

	var res map[string]interface{}
	err = json.Unmarshal(b, &res)
	assert.Nil(t, err)
	assert.Equal(t, "Identification", res["idShort"])
	assert.Equal(t, 27, len(res["submodelElements"].([]interface{})))
	assert.Equal(t, "Manufacturername1", (res["submodelElements"].([]interface{}))[0].(map[string]interface{})["value"])

}

func TestMain(m *testing.M) {
	dir, _ := os.Getwd()
	os.Setenv("SUBMODEL_CONFIG_PATH", dir+"/submodel_config.yaml")
	ym, err := os.ReadFile(dir + "/submodel_config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(ym), &sy)
	if err != nil {
		panic(err)
	}
	sy.SubmodelConfigs[0].QueryTemplates[0].Path = dir + "/" + sy.SubmodelConfigs[0].QueryTemplates[0].Path
	sy.SubmodelConfigs[0].ResponseTemplatePath = dir + "/" + sy.SubmodelConfigs[0].ResponseTemplatePath

	os.Exit(m.Run())
}
