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
select short_id, semantic_id from submodel_proxy
join aas on aas.global_id = submodel_proxy.aas_id
`
)

type Aas interface {
	List() ([]byte, error)
	Get(aasId string) ([]byte, error)
	CreateOrUpdate(aasId string, jsonb []byte) ([]byte, error)
	Delete(aasId string) error
	GetSubmodelIds(aasId string, submodelIdShort string) (map[string]string, error)
	DeleteSubmodel(aasId, submodelIdShort string) error
	CreateOrUpdateSubmodel(aasId, globalId, semanticId, submodelIdShort string) error
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
	tpl := template.Must(
		template.New("base").Funcs(sprig.FuncMap()).ParseFiles(cfg.AasQuerySqlPath),
	)

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

func (a *aas) CreateOrUpdate(aasId string, jsonb []byte) ([]byte, error) {
	_, err := a.db.Exec(`
	insert into 
	  aas
		(global_id, content, updated_at, created_at) 
	values
		(?, ?, current_timestamp, current_timestamp)
	on conflict (global_id)
	do update set content=?, updated_at=current_timestamp`,
		aasId, jsonb, jsonb)
	if err != nil {
		return nil, err
	}
	return jsonb, nil
}

func (a *aas) Delete(aasId string) error {
	_, err := a.db.Exec("delete from aas where global_id = ?", aasId)
	if err != nil {
		return err
	}
	return nil
}

func (a *aas) GetSubmodel(aasId, submodelIdShort string) (string, error) {
	rows, err := a.db.Query(
		fmt.Sprintf(submodelSql+" where aas.global_id = '%s' and submodel_proxy.short_id = '%s'",
			"semantic_id", aasId, submodelIdShort),
	)
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
func (a *aas) DeleteSubmodel(aasId, submodelIdShort string) error {
	_, err := a.db.Exec("delete from submodel_proxy where asset_id = ? and short_id = ?", aasId, submodelIdShort)
	if err != nil {
		return err
	}
	return nil
}

func (a *aas) CreateOrUpdateSubmodel(aasId, globalId, semanticId, submodelIdShort string) error {
	_, err := a.db.Exec(`
	insert into 
	  submodel_proxy
		(short_id, semantic_id, global_id, aas_id, updated_at, created_at) 
	values
		(?, ?, ?, ?, current_timestamp, current_timestamp)
	on conflict (short_id, aas_id)
	do update set semanticId=?, updated_at=current_timestamp`,
		submodelIdShort, semanticId, globalId, aasId, semanticId)
	if err != nil {
		return err
	}
	return nil
}
func (a *aas) GetSubmodelIds(aasId string, submodelIdShort string) (map[string]string, error) {
	where := fmt.Sprintf(" where aas.global_id = '%s'", aasId)
	if submodelIdShort != "" {
		where += fmt.Sprintf(" and short_id = '%s'", submodelIdShort)
	}
	rows, err := a.db.Query(fmt.Sprintf(submodelSql + where))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	semanticIds := map[string]string{}
	for rows.Next() {
		var id, semanticId string
		err = rows.Scan(&id, &semanticId)
		if err != nil {
			return nil, err
		}
		semanticIds[id] = semanticId
	}
	return semanticIds, nil
}
