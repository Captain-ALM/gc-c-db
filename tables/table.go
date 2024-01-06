package tables

type Table interface {
	TableName() string
	GetConstraints() []Constraint
	GetChildTables() []Table
	GetID() interface{}
	GetIDObject() Table
	GetNullableColumns() []string
}
