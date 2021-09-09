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



func NewEntity() *Entity {
	return &Entity{}
}
