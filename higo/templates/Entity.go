package templates

type Entity struct {
	PackageName   string
	StructName    string
	PrimaryId     string
	StructFields  []StructField
	HasCreateTime bool
	HasUpdateTime bool
	OutStruct     string
	OutDir        string
	File          string
}

type StructField struct {
	FieldName         string
	FieldType         string
	TableFieldName    string
	TableFieldComment string
}
