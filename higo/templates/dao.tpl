package {{.PackageName}}

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/higo/errcodg"
	"github.com/dengpju/higo-gin/higo/exceptions/DaoException"
	"github.com/dengpju/higo-gin/higo/sql"
	"github.com/dengpju/higo-utils/utils"
	"strings"
	{{- range $impo := .Imports}}
    {{$impo}}
    {{- end}}
)

type {{.StructName}} struct {
	model *{{.ModelPackageName}}.{{.ModelName}}
}

func New() *{{.StructName}} {
	return &{{.StructName}}{model: {{.ModelPackageName}}.New()}
}

func (this *{{.StructName}}) Orm() *higo.Orm {
	return this.model.Orm
}

func (this *{{.StructName}}) Model() *{{.ModelPackageName}}.{{.ModelName}} {
	return {{.ModelPackageName}}.New()
}

func (this *{{.StructName}}) Models() []*{{.ModelPackageName}}.{{.ModelName}} {
	return make([]*{{.ModelPackageName}}.{{.ModelName}}, 0)
}

func (this *{{.StructName}}) SetData(entity *{{.EntityPackageName}}.{{.EntityName}}) {
	if entity.IsEdit() { //编辑
		if entity.PriEmpty() {
			DaoException.Throw("{{.PrimaryId}}"+errcodg.PrimaryIdError.Message(), int(errcodg.PrimaryIdError))
		}
		if !this.GetBy{{.PrimaryId}}(entity.{{.PrimaryId}}).Exist() {
			DaoException.Throw(errcodg.NotExistError.Message(), int(errcodg.NotExistError))
		}

		{{- if .HasUpdateTime}}
		builder := this.model.Update(this.model.TableName()).Where({{.ModelPackageName}}.{{.PrimaryId}}, entity.{{.PrimaryId}})
		{{- else}}
		_ = this.model.Update(this.model.TableName()).Where({{.ModelPackageName}}.{{.PrimaryId}}, entity.{{.PrimaryId}})
		{{- end}}
		if {{.EntityPackageName}}.FlagDelete == entity.Flag() {

		} else {

		}
		{{- if .HasUpdateTime}}
		builder.Set({{.ModelPackageName}}.{{.EntityUpdateTimeField}}, entity.{{.EntityUpdateTimeField}})
		{{- end}}
	} else { //新增
		this.model.Insert(this.model.TableName()).
		{{- range $v := .ModelFields}}
		{{- if ne $v.FieldName $.EntityDeleteTimeField}}
		    {{- if ne $v.FieldName $.ModelEndField}}
            Set({{$.ModelPackageName}}.{{$v.FieldName}}, entity.{{$v.FieldName}}).  //{{$v.TableFieldComment}}
            {{- else}}
            Set({{$.ModelPackageName}}.{{$v.FieldName}}, entity.{{$v.FieldName}})  //{{$v.TableFieldComment}}
            {{- end}}
		{{- end}}
        {{- end}}
	}
	this.model.Build()
}

//添加
func (this *{{.StructName}}) Add() {{.PrimaryIdType}} {
	higo.Result(this.model.Mapper(this.model.GetBuilder()).InsertGetId().Error).Unwrap()
	return this.model.LastInsertId()
}

//更新
func (this *{{.StructName}}) Update() bool {
	higo.Result(this.model.Mapper(this.model.GetBuilder()).Exec().Error).Unwrap()
	return true
}

//id查询
func (this *{{.StructName}}) GetBy{{.PrimaryId}}({{.SmallHumpPrimaryId}} {{.PrimaryIdType}}, fields ...string) *{{.ModelPackageName}}.{{.ModelName}} {
	if len(fields) == 0 {
		fields = append(fields, "*")
	}
	model := this.Model()
	model.Mapper(sql.Select(fields...).
		From(this.model.TableName()).
		Where({{.ModelPackageName}}.{{.PrimaryId}}+" = ?", {{.SmallHumpPrimaryId}}).
		{{- if .HasDeleteTime}}
        Where("isnull(`"+{{.ModelPackageName}}.{{.EntityDeleteTimeField}}+"`)").
        {{- end}}
		ToSql()).Query().Scan(&model)
	model.CheckError()
	return model
}

//id集查询
func (this *Dao) GetBy{{.PrimaryId}}s({{.SmallHumpPrimaryId}}s []{{.PrimaryIdType}}, fields ...string) []*{{.ModelPackageName}}.{{.ModelName}} {
	if len(fields) == 0 {
		fields = append(fields, "*")
	}
	models := this.Models()
	this.model.Mapper(sql.Select(fields...).
		From(this.model.TableName()).
		Where({{.ModelPackageName}}.{{.PrimaryId}}+" IN (?)", {{.SmallHumpPrimaryId}}s).
		{{- if .HasDeleteTime}}
		Where("isnull(`"+{{.ModelPackageName}}.{{.EntityDeleteTimeField}}+"`)").
		{{- end}}
		ToSql()).Query().Scan(&models)
	this.model.CheckError()
	return models
}

//硬删除
func (this *Dao) DeleteBy{{.PrimaryId}}({{.SmallHumpPrimaryId}} {{.PrimaryIdType}}) {
	higo.Result(this.model.Mapper(sql.Delete(this.model.TableName()).
		DeleteBuilder().
		Where({{.ModelPackageName}}.{{.PrimaryId}}+" = ?", {{.SmallHumpPrimaryId}}).
		ToSql()).Exec().Error).Unwrap()
}

//列表
func (this *Dao) List(perPage, page uint64, where map[string]interface{}, fields ...string) *higo.Pager {
	models := this.Models()
	pager := higo.NewPager(perPage, page)
	query := this.model.Table(this.model.TableName())
	query.Paginate(pager).Find(&models)
	query.CheckError()
	pager.Items = models
	return pager
}
