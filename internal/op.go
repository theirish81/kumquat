package internal

import (
	"context"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// IOperation is the interface for operations
type IOperation interface {
	Run(ctx context.Context) error
	GetResult() any
}

const OpMongo = "mongo"
const OpNixShell = "nixShellCommand"
const OpSql = "sql"
const OpTemplate = "template"
const OpWait = "wait"
const OpFile = "file"
