package templates

type Entity struct {
	PackageName   string
	StructName    string
	PrimaryId     string
	StructFields  []StructField
	HasCreateTime bool
	HasUpdateTime bool
	OutDir        string
	FileName      string
}

type StructField struct {
	FieldName         string
	FieldType         string
	TableFieldName    string
	TableFieldComment string
}

func NewEntity() *Entity {
	return &Entity{}
}
