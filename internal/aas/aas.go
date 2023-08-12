package aas

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/caarlos0/env"
	_ "github.com/lib/pq"
)

type Aas interface {
	List() ([]byte, error)
	Get(aasId string) ([]byte, error)
	CreateOrUpdate(aasId string, jsonb []byte) ([]byte, error)
	Delete(aasId string) error
	GetSubmodelIds(aasId, submodelIdShort string) (map[string]string, error)
	DeleteSubmodel(aasId, submodelIdShort string) error
	CreateOrUpdateSubmodel(aasId, globalId, semanticId, submodelIdShort string) error
}

type aas struct {
	db      *sql.DB
	tpl     *template.Template
	sRefTpl *template.Template
}

type config struct {
	AasTablesCreated bool   `env:"AAS_TABLES_CREATED" envDefault:"false"`
	AasDbHost        string `env:"AAS_DB_HOST" envDefault:"127.0.0.1"`
	AasDbPort        int    `env:"AAS_DB_PORT" envDefault:"5432"`
	AasDbUser        string `env:"AAS_DB_USER" envDefault:"postgres"`
	AasDbPassword    string `env:"AAS_DB_PASSWORD" envDefault:"password"`
	AasDbDatabase    string `env:"AAS_DB_DATABASE" envDefault:"sample"`
	AasDbSslMode     string `env:"AAS_DB_SSL_MODE" envDefault:"disable"`
}

func NewAas() (Aas, error) {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		slog.Error(fmt.Sprintf("%+v\n", err))
		return nil, err
	}
	source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.AasDbHost, cfg.AasDbPort, cfg.AasDbUser, cfg.AasDbPassword, cfg.AasDbDatabase, cfg.AasDbSslMode)
	db, err := sql.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	if cfg.AasTablesCreated {
		_, err := db.Exec(aasInitSql)
		if err != nil {
			return nil, err
		}
	}

	return newAasWithDB(db)
}

func newAasWithDB(db *sql.DB) (Aas, error) {
	tpl := template.Must(
		template.New(filepath.Base("aasSql")).Funcs(sprig.FuncMap()).Parse(aasSql),
	)
	sRefTpl := template.Must(
		template.New(filepath.Base("submodelRefSql")).Funcs(sprig.FuncMap()).Parse(submodelRefSql),
	)

	return &aas{
		db:      db,
		tpl:     tpl,
		sRefTpl: sRefTpl,
	}, nil
}

func (a *aas) List() ([]byte, error) {
	writer := new(strings.Builder)
	err := a.tpl.Execute(writer, nil)
	if err != nil {
		slog.Error(err.Error())
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
		slog.Error(err.Error())
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
	writer := new(strings.Builder)
	err := a.sRefTpl.Execute(writer, map[string]interface{}{"AasID": aasId, "SubmodelIDShort": submodelIdShort})
	if err != nil {
		slog.Error(err.Error())
	}

	rows, err := a.db.Query(writer.String())

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
	writer := new(strings.Builder)
	err := a.sRefTpl.Execute(writer, map[string]interface{}{
		"AasID":           aasId,
		"SubmodelIDShort": submodelIdShort,
	})
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(writer.String())

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
