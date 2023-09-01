package submodel

import (
	"database/sql"
	"errors"
	"fmt"
	"hiroyoshii/go-aas-proxy/internal/sqlutility"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/Masterminds/sprig/v3"
	"github.com/caarlos0/env"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/sijms/go-ora/v2"
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
	SubmodelConfigPath string `env:"SUBMODEL_CONFIG_PATH,required"`
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
		DbType   string `yaml:"dbType"`
		DbConfig struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
			Sslmode  string `yaml:"sslmode"`
			Options  string `yaml:"options"`
		} `yaml:"dbConfig"`
	} `yaml:"databases"`
}

func NewSubmodel() (Submodel, error) {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	ym, err := os.ReadFile(cfg.SubmodelConfigPath)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	t := submodelYaml{}

	err = yaml.Unmarshal([]byte(ym), &t)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	dbs := map[string]*sql.DB{}
	for _, c := range t.DatabaseConfigs {
		d := c.DbConfig
		source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
			d.Host, d.Port, d.User, d.Password, d.Database)
		switch c.DbType {
		case "mysql":
			source = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", d.User, d.Password, d.Host, d.Port, d.Database, d.Options)
		case "oracle":
			source = fmt.Sprintf("oracle://%s:%s@%s:%d/%s?%s", d.User, d.Password, d.Host, d.Port, d.Database, d.Options)

		case "postgres":
			source = fmt.Sprintf("%s sslmode=%s", source, d.Sslmode)
		}
		db, err := sql.Open(c.DbType, source)
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

	respErr := sync.Map{}
	respMap := sync.Map{}
	wg := &sync.WaitGroup{}
	for k, v := range a.queryDbTplNameMap[semanticID] {
		for _, file := range v {
			wg.Add(1)
			go func(file *template.Template) {
				key := fileNameToCamel(file.Name())
				defer wg.Done()
				writer := new(strings.Builder)
				err := file.Execute(writer, map[string]interface{}{"AasID": aasId, "SubmodelIdShort": submodelIdShort})
				if err != nil {
					slog.Error(err.Error())
					respErr.Store(key, err)
					return
				}
				rows, err := a.dbs[k].Query(writer.String())
				if err != nil {
					slog.Error(err.Error())
					respErr.Store(key, err)
					return
				}
				defer rows.Close()
				m, mm, err := sqlutility.RowsToMap(rows)
				if err != nil {
					slog.Error(err.Error())
					respErr.Store(key, err)
					return
				}
				respMap.Store(key, map[string]interface{}{
					"Results": m,
					"Columns": mm,
				},
				)
			}(file)
		}
	}
	wg.Wait()
	errs := []error{}
	respErr.Range(func(key interface{}, value interface{}) bool {
		errs = append(errs, value.(error))
		return true
	})
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	m := map[string]interface{}{}
	respMap.Range(func(key interface{}, value interface{}) bool {
		m[key.(string)] = value
		return true
	})
	slog.Debug(fmt.Sprintf("response funcMap: %v", m))
	writer := new(strings.Builder)
	err := a.respTpl[semanticID].Execute(writer, m)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return []byte(writer.String()), nil
}

func fileNameToCamel(name string) string {
	n := strings.Split(name, ".")[0]
	return strcase.ToCamel(n)
}
