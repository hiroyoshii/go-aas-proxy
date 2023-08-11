package submodel

import (
	"database/sql"
	"fmt"
	"hiroyoshii/go-aas-proxy/internal/sqlutility"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/Masterminds/sprig/v3"
	"github.com/caarlos0/env"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type Submodel interface {
	Get(aasId, semanticID, submodelIdShort string) ([]byte, error)
}

type submodel struct {
	dbs               map[string]*sql.DB
	respTpl           map[string]*template.Template
	queryDbTplNameMap map[string]map[string][]*template.Template
}

type config struct {
	SubmodelConfigPath string `env:"SUBMODEL_CONFIG_PATH" envDefault:"internal/submodel/submodel_config.yaml"`
}

type submodelYaml struct {
	SubmodelConfigs []struct {
		SemanticID     string `yaml:"semanticID"`
		QueryTemplates []struct {
			Path   string `yaml:"path"`
			DbName string `yaml:"dbName"`
		} `yaml:"queryTemplates"`
		ResponseTemplatePath string `yaml:"responseTemplatePath"`
	} `yaml:"submodels"`
	DatabaseConfigs []struct {
		Name     string `yaml:"name"`
		DbConfig struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
			Sslmode  string `yaml:"sslmode"`
		} `yaml:"dbConfig"`
	} `yaml:"databases"`
}

func NewSubmodel() (Submodel, error) {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	ym, err := os.ReadFile(cfg.SubmodelConfigPath)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}
	t := submodelYaml{}

	err = yaml.Unmarshal([]byte(ym), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	dbs := map[string]*sql.DB{}
	for _, c := range t.DatabaseConfigs {
		d := c.DbConfig
		source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			d.Host, d.Port, d.User, d.Password, d.Database, d.Sslmode)
		db, err := sql.Open("postgres", source)
		if err != nil {
			return nil, err
		}
		dbs[c.Name] = db
	}

	return newSubmodelWithDb(t, dbs)
}

func newSubmodelWithDb(t submodelYaml, dbs map[string]*sql.DB) (Submodel, error) {
	respTpl := map[string]*template.Template{}
	qfileMap := map[string]map[string][]*template.Template{}
	for _, s := range t.SubmodelConfigs {
		qtlps := map[string][]*template.Template{}
		for _, q := range s.QueryTemplates {
			tpl := template.Must(
				template.New(filepath.Base(q.Path)).Funcs(sprig.FuncMap()).ParseFiles(q.Path),
			)
			if qtlps[q.DbName] == nil {
				qtlps[q.DbName] = []*template.Template{}
			}
			qtlps[q.DbName] = append(qtlps[q.DbName], tpl)
		}
		qfileMap[s.SemanticID] = qtlps
		tpl := template.Must(
			template.New(filepath.Base(s.ResponseTemplatePath)).Funcs(sprig.FuncMap()).ParseFiles(s.ResponseTemplatePath),
		)
		respTpl[s.SemanticID] = tpl
	}
	return &submodel{
		dbs:               dbs,
		respTpl:           respTpl,
		queryDbTplNameMap: qfileMap,
	}, nil
}

func (a *submodel) Get(aasId, semanticID, submodelIdShort string) ([]byte, error) {
	if a.respTpl[semanticID] == nil {
		return []byte{}, nil
	}

	respMap := map[string]interface{}{}
	for k, v := range a.queryDbTplNameMap[semanticID] {
		for _, file := range v {
			writer := new(strings.Builder)
			err := file.Execute(writer, map[string]interface{}{"AasID": aasId, "SubmodelIdShort": submodelIdShort})
			if err != nil {
				log.Fatalln(err)
				return nil, err
			}
			rows, err := a.dbs[k].Query(writer.String())
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			m, mm, err := sqlutility.RowsToMap(rows)
			if err != nil {
				return nil, err
			}
			respMap[fileNameToCamel(file.Name())] = map[string]interface{}{
				"Results": m,
				"Columns": mm,
			}
		}
	}
	log.Printf("response funcMap: %v\n", respMap)
	writer := new(strings.Builder)
	err := a.respTpl[semanticID].Execute(writer, respMap)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return []byte(writer.String()), nil
}

func fileNameToCamel(name string) string {
	n := strings.Split(name, ".")[0]
	return strcase.ToCamel(n)
}
