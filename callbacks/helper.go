package callbacks

import (
	"sort"

	"github.com/joelazar/gorm"
	"github.com/joelazar/gorm/clause"
)

// ConvertMapToValuesForCreate convert map to values
func ConvertMapToValuesForCreate(stmt *gorm.Statement, mapValue map[string]interface{}) (values clause.Values) {
	values.Columns = make([]clause.Column, 0, len(mapValue))
	selectColumns, restricted := stmt.SelectAndOmitColumns(true, false)

	var keys = make([]string, 0, len(mapValue))
	for k := range mapValue {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		value := mapValue[k]
		if stmt.Schema != nil {
			if field := stmt.Schema.LookUpField(k); field != nil {
				k = field.DBName
			}
		}

		if v, ok := selectColumns[k]; (ok && v) || (!ok && !restricted) {
			values.Columns = append(values.Columns, clause.Column{Name: k})
			if len(values.Values) == 0 {
				values.Values = [][]interface{}{{}}
			}

			values.Values[0] = append(values.Values[0], value)
		}
	}
	return
}

// ConvertSliceOfMapToValuesForCreate convert slice of map to values
func ConvertSliceOfMapToValuesForCreate(stmt *gorm.Statement, mapValues []map[string]interface{}) (values clause.Values) {
	var (
		columns = make([]string, 0, len(mapValues))
	)

	// when the length of mapValues,return directly here
	// no need to call stmt.SelectAndOmitColumns method
	if len(mapValues) == 0 {
		stmt.AddError(gorm.ErrEmptySlice)
		return
	}

	var (
		result                    = make(map[string][]interface{}, len(mapValues))
		selectColumns, restricted = stmt.SelectAndOmitColumns(true, false)
	)

	for idx, mapValue := range mapValues {
		for k, v := range mapValue {
			if stmt.Schema != nil {
				if field := stmt.Schema.LookUpField(k); field != nil {
					k = field.DBName
				}
			}

			if _, ok := result[k]; !ok {
				if v, ok := selectColumns[k]; (ok && v) || (!ok && !restricted) {
					result[k] = make([]interface{}, len(mapValues))
					columns = append(columns, k)
				} else {
					continue
				}
			}

			result[k][idx] = v
		}
	}

	sort.Strings(columns)
	values.Values = make([][]interface{}, len(mapValues))
	values.Columns = make([]clause.Column, len(columns))
	for idx, column := range columns {
		values.Columns[idx] = clause.Column{Name: column}

		for i, v := range result[column] {
			if len(values.Values[i]) == 0 {
				values.Values[i] = make([]interface{}, len(columns))
			}

			values.Values[i][idx] = v
		}
	}
	return
}
