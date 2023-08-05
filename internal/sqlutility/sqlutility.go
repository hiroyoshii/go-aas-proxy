package sqlutility

import (
	"database/sql"

	"github.com/iancoleman/strcase"
)

func RowsToMap(rows *sql.Rows) ([]map[string]interface{}, map[string]string, error) {
	var resultsMap []map[string]interface{}
	typeMap := map[string]string{}

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	coltypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}
	toCamelMap := map[string]string{}
	for i, col := range coltypes {
		typeMap[strcase.ToCamel(cols[i])] = col.DatabaseTypeName()
		toCamelMap[cols[i]] = strcase.ToCamel(cols[i])
	}

	for rows.Next() {
		var row = make([]interface{}, len(cols))
		var rowp = make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			rowp[i] = &row[i]
		}
		err = rows.Scan(rowp...)
		if err != nil {
			return nil, nil, err
		}

		rowMap := make(map[string]interface{})
		for i, col := range cols {
			rowMap[toCamelMap[col]] = row[i]
		}
		resultsMap = append(resultsMap, rowMap)
	}
	return resultsMap, typeMap, nil
}
