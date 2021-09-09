package AdminDao

import (
	"github.com/dengpju/higo-gin/higo"
	"github.com/dengpju/higo-gin/higo/errcodg"
	"github.com/dengpju/higo-gin/higo/exceptions/DaoException"
	"github.com/dengpju/higo-gin/higo/sql"
	"github.com/dengpju/higo-gin/test/app/Entity/AdminEntity"
	"github.com/dengpju/higo-gin/test/app/Models/AdminModel"
	"github.com/dengpju/higo-utils/utils"
	"strings"
)

type Dao struct {
	model *AdminModel.Impl
}

func New() *Dao {
	return &Dao{model: AdminModel.New()}
}

func (this *Dao) Orm() *higo.Orm {
	return this.model.Orm
}

func (this *Dao) Model() *AdminModel.Impl {
	return AdminModel.New()
}

func (this *Dao) Models() []*AdminModel.Impl {
	return make([]*AdminModel.Impl, 0)
}

func (this *Dao) SetData(entity *AdminEntity.Impl) {
	if entity.IsEdit() { //编辑
		if entity.PriEmpty() {
			DaoException.Throw("AdminId"+errcodg.PrimaryIdError.Message(), int(errcodg.PrimaryIdError))
		}
		if !this.GetByAdminId(entity.AdminId).Exist() {
			DaoException.Throw(errcodg.NotExistError.Message(), int(errcodg.NotExistError))
		}
		builder := this.model.Update(this.model.TableName()).Where("admin_id", entity.AdminId)
		if AdminEntity.FlagDelete == entity.Flag() {

		} else {

		}
		builder.Set("update_time", entity.UpdateTime)
	} else { //新增
		this.model.Insert(this.model.TableName()).
			Set("admin_name", entity.AdminName).
			Set("user_id", entity.UserId).
			Set("create_time", entity.CreateTime).
			Set("update_time", entity.UpdateTime)
	}
	this.model.Build()
}

//添加
func (this *Dao) Add() int64 {
	higo.Result(this.model.Mapper(this.model.GetBuilder()).InsertGetId().Error).Unwrap()
	return this.model.LastInsertId()
}

//更新
func (this *Dao) Update() bool {
	higo.Result(this.model.Mapper(this.model.GetBuilder()).Exec().Error).Unwrap()
	return true
}

//id查询
func (this *Dao) GetByAdminId(adminId int64, fields ...string) *AdminModel.Impl {
	if len(fields) == 0 {
		fields = append(fields, "*")
	}
	model := this.Model()
	model.Mapper(sql.Select(fields...).
		From(this.model.TableName()).
		Where("`admin_id` = ?", adminId).
		Where("isnull(`delete_time`)").
		ToSql()).Query().Scan(&model)
	return model
}

//id集查询
func (this *Dao) GetByAdminIds(adminIds []interface{}, fields ...string) []*AdminModel.Impl {
	if len(fields) == 0 {
		fields = append(fields, "*")
	}
	models := this.Models()
	this.Model().Mapper(sql.Select(fields...).
		From(this.model.TableName()).
		Where("`admin_id` IN (?)", strings.Join(utils.ConvStrSlice(adminIds), ",")).
		Where("isnull(`delete_time`)").
		ToSql()).Query().Scan(&models)
	return models
}

//硬删除
func (this *Dao) DeleteByAdminId(adminId int64) {
	higo.Result(this.model.Mapper(sql.Delete(this.model.TableName()).
		DeleteBuilder().
		Where("admin_id = ?", adminId).
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
