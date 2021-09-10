package {{.PackageName}}

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/higo/errcodg"
	"github.com/dengpju/higo-gin/higo/exceptions/DaoException"
	"github.com/dengpju/higo-gin/higo/sql"
	"github.com/dengpju/higo-gin/test/app/Entity/AdminEntity"
	"github.com/dengpju/higo-gin/test/app/Models/AdminModel"
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
	return make([]*{{.ModelPackageName}}.{{.ModelName}}``, 0)
}

func (this *{{.StructName}}) SetData(entity *{{.EntityPackageName}}.{{.EntityName}}) {
	if entity.IsEdit() { //编辑
		if entity.PriEmpty() {
			DaoException.Throw("{{.PrimaryId}}"+errcodg.PrimaryIdError.Message(), int(errcodg.PrimaryIdError))
		}
		if !this.GetBy{{.PrimaryId}}(entity.{{.PrimaryId}}).Exist() {
			DaoException.Throw(errcodg.NotExistError.Message(), int(errcodg.NotExistError))
		}
		builder := this.model.Update(this.model.TableName()).Where("{{.TablePrimaryId}}", entity.{{.PrimaryId}})
		if {{.EntityPackageName}}.FlagDelete == entity.Flag() {

		} else {

		}
		builder.Set("update_time", entity.UpdateTime)
	} else { //新增
		this.model.Insert(this.model.TableName()).
		    Set(string(AdminModel.AdminName), entity.AdminName).
			Set("admin_name", entity.AdminName).
			Set("user_id", entity.UserId).
			Set("create_time", entity.CreateTime).
			Set("update_time", entity.UpdateTime)
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
		Where("`{{.TablePrimaryId}}` = ?", {{.SmallHumpPrimaryId}}).
		{{- if .HasDeleteTime}}
        Where("isnull(`delete_time`)").
        {{- end}}
		ToSql()).Query().Scan(&model)
	return model
}

//id集查询
func (this *Dao) GetBy{{.PrimaryId}}s({{.SmallHumpPrimaryId}}s []interface{}, fields ...string) []*{{.ModelPackageName}}.{{.ModelName}} {
	if len(fields) == 0 {
		fields = append(fields, "*")
	}
	models := this.Models()
	this.Model().Mapper(sql.Select(fields...).
		From(this.model.TableName()).
		Where("`{{.TablePrimaryId}}` IN (?)", strings.Join(utils.ConvStrSlice({{.SmallHumpPrimaryId}}s), ",")).
		{{- if .HasDeleteTime}}
		Where("isnull(`delete_time`)").
		{{- end}}
		ToSql()).Query().Scan(&models)
	return models
}

//硬删除
func (this *Dao) DeleteByAdminId({{.SmallHumpPrimaryId}} {{.PrimaryIdType}}) {
	higo.Result(this.model.Mapper(sql.Delete(this.model.TableName()).
		DeleteBuilder().
		Where("{{.TablePrimaryId}} = ?", {{.SmallHumpPrimaryId}}).
		ToSql()).Exec().Error).Unwrap()
}

//列表
func (this *Dao) List(perPage, page uint64, where map[string]interface{}, fields ...string) *higo.Pager {
	models := this.Models()
	pager := higo.NewPager(perPage, page)
	query := this.Model().Table(this.model.TableName())
	query.Paginate(pager).Find(&models)
	pager.Items = models
	return pager
}
