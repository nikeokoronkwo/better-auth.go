package main

// TODO: Mock via github.com/pashagolub/pgxmock
import (
	"database/sql"
	"embed"
)

//go:embed sql
var migrationScripts embed.FS

// Applies the SQL code in the `sql/` folder
func ApplySQL(db *sql.DB, user, session string) {
	
}

type Model struct {
	// the name of the model
	ModelName string
	// a map of fields for the given model
	Fields map[string] Field
}

type Field struct {
	// the type of the field
	// can be any common type used in SQL
	Type FieldType

	// whether the given field is unique
	Unique bool

	// whether the given field is a primary key
	Primary bool

	// whether the given field is not null
	NotNull bool

	// default value, if any
	Default any

	// default expression to use
	DefaultExpression string

	// Provides the model and key of the reference
	// if referencing "users.id", you pass (<user model>, "id")
	References func() (Model, string)
}

type FieldType string
const (
	String = "TEXT"
	Timestamp = "TIMESTAMP"
	Timestamptz = "TIMESTAMPTZ"
	Boolean = "BOOLEAN"
)