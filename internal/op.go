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

const OP_MONGO = "mongo"
const OP_NIX_SHELL = "nixShellCommand"
const OP_SQL = "sql"
const OP_TEMPLATE = "template"
const OP_WAIT = "wait"
