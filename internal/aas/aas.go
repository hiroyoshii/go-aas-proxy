package aas

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/caarlos0/env"
	_ "github.com/lib/pq"
)

var (
	submodelSql = `
select semantic_id from submodel_proxy
join aas on aas.global_id = submodel_proxy.aas_id
where aas.global_id = '%s' and submodel_proxy.short_id = '%s'
`
)

type Aas interface {
	List() ([]byte, error)
	Get(aasId string) ([]byte, error)
	GetSubmodel(aasId, submodelIdShort string) (string, error)
}

type aas struct {
	db  *sql.DB
	tpl *template.Template
}

type config struct {
	AasQuerySqlPath string `env:"AAS_QUERY_SQL_PATH" envDefault:"internal/aas/query_default.tpl.sql"`
	AasDbHost       string `env:"AAS_DB_HOST" envDefault:"127.0.0.1"`
	AasDbPort       int    `env:"AAS_DB_PORT" envDefault:"5432"`
	AasDbUser       string `env:"AAS_DB_USER" envDefault:"mebee"`
	AasDbPassword   string `env:"AAS_DB_PASSWORD" envDefault:"password"`
	AasDbDatabase   string `env:"AAS_DB_DATABASE" envDefault:"sample"`
	AasDbSslMode    string `env:"AAS_DB_SSL_MODE" envDefault:"disable"`
}

func NewAas() (Aas, error) {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	tpl := template.Must(template.ParseFiles(cfg.AasQuerySqlPath))
	tpl.Funcs(sprig.FuncMap())

	source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.AasDbHost, cfg.AasDbPort, cfg.AasDbUser, cfg.AasDbPassword, cfg.AasDbDatabase, cfg.AasDbSslMode)
	db, err := sql.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	return &aas{
		db:  db,
		tpl: tpl,
	}, nil
}

func (a *aas) List() ([]byte, error) {
	writer := new(strings.Builder)
	err := a.tpl.Execute(writer, nil)
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := a.db.Query(writer.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sjson := []string{}
	for rows.Next() {
		var b string
		err = rows.Scan(&b)
		if err != nil {
			return nil, err
		}
		sjson = append(sjson, b)
	}
	return []byte(fmt.Sprintf("[%s]", strings.Join(sjson, ","))), nil
}

// Get returns a single AAS. If empty result, return "" (empty string).
func (a *aas) Get(aasId string) ([]byte, error) {
	writer := new(strings.Builder)
	err := a.tpl.Execute(writer, map[string]interface{}{"AasID": aasId})
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := a.db.Query(writer.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sjson := ""
	for rows.Next() {
		err = rows.Scan(&sjson)
		if err != nil {
			return nil, err
		}
	}
	return []byte(sjson), nil
}

func (a *aas) GetSubmodel(aasId, submodelIdShort string) (string, error) {
	rows, err := a.db.Query(fmt.Sprintf(submodelSql, aasId, submodelIdShort))
	if err != nil {
		return "", err
	}
	defer rows.Close()
	semanticId := ""
	for rows.Next() {
		err = rows.Scan(&semanticId)
		if err != nil {
			return "", err
		}
	}
	return semanticId, nil
}
