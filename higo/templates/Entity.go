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

const (
	EntityStructName = "Impl"
	EntityDirSuffix  = "Entity"
	EntityFileName   = "entity"
)

func NewEntity(modelTool ModelTool, model Model) *Entity {
	packageName := model.HumpUnpreTableName + EntityDirSuffix
	return &Entity{
		PackageName:   packageName,
		StructName:    EntityStructName,
		PrimaryId:     model.PrimaryId,
		StructFields:  model.StructFields,
		HasCreateTime: model.HasCreateTime,
		HasUpdateTime: model.HasUpdateTime,
		OutDir:        modelTool.OutEntityDir,
		FileName:      EntityFileName,
	}
}
