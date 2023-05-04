package Dao

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-gin/higo/errcode"
	"github.com/dunpju/higo-gin/higo/exceptions/DaoException"
	"github.com/dunpju/higo-gin/higo/sql"
	"github.com/dunpju/higo-gin/test/app/Entity/UserEntity"
	"github.com/dunpju/higo-gin/test/app/Models/UserModel"
)

type UserDao struct {
	model  *UserModel.Impl
	entity *UserEntity.Impl
}

func NewUserDao() *UserDao {
	return &UserDao{model: UserModel.New()}
}

func (this *UserDao) Orm() *higo.Orm {
	return this.model.Orm
}

func (this *UserDao) Model() *UserModel.Impl {
	return UserModel.New()
}

func (this *UserDao) Models() []*UserModel.Impl {
	return make([]*UserModel.Impl, 0)
}

func (this *UserDao) SetData(entity *UserEntity.Impl) {
	this.entity = entity
	if !entity.PriEmpty() || entity.IsEdit() { //编辑
		if !this.GetById(entity.Id).Exist() {
			DaoException.Throw(errcode.NotExistError.Message(), int(errcode.NotExistError))
		}
		_ = this.model.Update(this.model.TableName()).Where("`"+UserModel.Id+"`", entity.Id)
		if UserEntity.FlagDelete == entity.Flag() {

		} else {

		}
	} else { //新增
		this.model.Insert(this.model.TableName()).
			Set("`"+UserModel.Id+"`", entity.Id).       //
			Set("`"+UserModel.Uname+"`", entity.Uname). //
			Set("`"+UserModel.UTel+"`", entity.UTel).   //
			Set("`"+UserModel.Score+"`", entity.Score)  //
	}
	this.model.Build()
}

//添加
func (this *UserDao) Add() int64 {
	higo.Result(this.model.Mapper(this.model.GetBuilder()).InsertGetId().Error).Unwrap()
	return this.model.LastInsertId()
}

//更新
func (this *UserDao) Update() bool {
	if this.entity.PriEmpty() {
		DaoException.Throw("Id"+errcode.PrimaryIdError.Message(), int(errcode.PrimaryIdError))
	}
	higo.Result(this.model.Mapper(this.model.GetBuilder()).Exec().Error).Unwrap()
	return true
}

//id查询
func (this *UserDao) GetById(id int, fields ...string) *UserModel.Impl {
	if len(fields) == 0 {
		fields = append(fields, "*")
	}
	model := this.Model()
	model.Mapper(sql.Select(fields...).
		From(this.model.TableName()).
		Where("`"+UserModel.Id+"` = ?", id).
		ToSql()).Query().Scan(&model)
	model.CheckError()
	return model
}

//id集查询
func (this *UserDao) GetByIds(ids []int, fields ...string) []*UserModel.Impl {
	if len(fields) == 0 {
		fields = append(fields, "*")
	}
	models := this.Models()
	this.model.Mapper(sql.Select(fields...).
		From(this.model.TableName()).
		Where("`"+UserModel.Id+"` IN (?)", ids).
		ToSql()).Query().Scan(&models)
	this.model.CheckError()
	return models
}

//硬删除
func (this *UserDao) DeleteById(id int) {
	higo.Result(this.model.Mapper(sql.Delete(this.model.TableName()).
		DeleteBuilder().
		Where("`"+UserModel.Id+"` = ?", id).
		ToSql()).Exec().Error).Unwrap()
}

//列表
func (this *UserDao) List(perPage, page uint64, where map[string]interface{}, fields ...string) *higo.Pager {
	models := this.Models()
	pager := higo.NewPager(perPage, page)
	query := this.model.Table(this.model.TableName())
	query.Paginate(pager).Find(&models)
	query.CheckError()
	pager.Items = models
	return pager
}
