package internal

import (
	"context"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// SqlOp is the operation that runs SQL queries
type SqlOp struct {
	Driver  string
	URI     string
	Select  string
	Timeout time.Duration
	Result  []map[string]any
	scope   *Scope
}

// NewSqlOp constructor for SqlOp
func NewSqlOp(config map[string]any, scope *Scope) (*SqlOp, error) {
	config = SetDefault(config, "timeout", "10s")
	duration, err := time.ParseDuration(config["timeout"].(string))
	if err != nil {
		return nil, err
	}
	if err := PrototypeCheck(config, Proto{"driver": TYPE_STRING, "URI": TYPE_STRING, "select": TYPE_STRING}); err == nil {
		return &SqlOp{Driver: config["driver"].(string), URI: config["URI"].(string), Select: config["select"].(string),
			Timeout: duration, scope: scope}, nil
	} else {
		return nil, err
	}
}

// Run runs the query
func (o *SqlOp) Run(ctx context.Context) error {
	log.Debug().Str("select", o.Select).Msg("running SQL op")
	evalDriver, err := o.scope.Render(ctx, o.Driver)
	if err != nil {
		return err
	}
	evalURI, err := o.scope.Render(ctx, o.URI)
	if err != nil {
		return err
	}
	evalSelect, err := o.scope.Render(ctx, o.Select)
	if err != nil {
		return err
	}
	// Opening the connection to the SQL server
	if db, err := sqlx.Open(evalDriver, evalURI); err == nil {
		defer func() {
			_ = db.Close()
		}()
		res := make([]map[string]any, 0)
		// Performing the query
		if rows, err := db.Queryx(evalSelect); err == nil {
			defer func() {
				_ = rows.Close()
			}()
			// For each row...
			for rows.Next() {
				// ... we convert the row into a map, and append it to the response array
				rx := map[string]any{}

				if err := rows.MapScan(rx); err == nil {
					o.castMap(&rx, rows)
					res = append(res, rx)
				} else {
					// In case the map conversion fails
					return err
				}
			}
			o.Result = res
		} else {
			// In case the query failed
			return err
		}
	} else {
		// In case the connection failed
		return err
	}
	return nil
}

// castMap deals with the inadequacies of sqlx when it comes to type mapping
func (o *SqlOp) castMap(rx *map[string]any, rows *sqlx.Rows) {
	columnTypes, _ := rows.ColumnTypes()
	columnNames, _ := rows.Columns()
	for idx, columnName := range columnNames {
		cType := columnTypes[idx]
		rxData := *rx
		switch cType.ScanType().Name() {
		case "RawBytes", "NullTime":
			if rxData[columnName] != nil {
				if bytes, ok := rxData[columnName].([]byte); ok {
					data := string(bytes)
					rxData[columnName] = data
				}

			}
		case "float32", "float64", "NullFloat32", "NullFloat64":
			if rxData[columnName] != nil {
				if bytes, ok := rxData[columnName].([]byte); ok {
					data, _ := strconv.ParseFloat(string(bytes), 64)
					rxData[columnName] = data
				}
			}
		case "int", "int8", "int16", "int32", "int64", "NullInt64":
			if rxData[columnName] != nil {
				if bytes, ok := rxData[columnName].([]byte); ok {
					data, _ := strconv.Atoi(string(bytes))
					rxData[columnName] = data
				}
			}
		default:
			if rxData[columnName] != nil {
				switch cType.DatabaseTypeName() {
				case "TIMESTAMP", "DATETIME":
					if bytes, ok := rxData[columnName].([]byte); ok {
						data := string(bytes)
						rxData[columnName] = data
					}
				case "VARCHAR":
					if str, ok := rxData[columnName].(string); ok {
						rxData[columnName] = str
					}
				case "JSONB":
					if bytes, ok := rxData[columnName].([]byte); ok {
						js, _ := simplejson.NewJson(bytes)
						rxData[columnName] = js
					}
				}
			}
		}
	}
}

// GetResult will return the result of the query
func (o *SqlOp) GetResult() any {
	return o.Result
}
