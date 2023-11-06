package Dao

import (
	"github.com/dunpju/higo-gin/higo/exceptions/DaoException"
	"github.com/dunpju/higo-gin/test/app/Entity/UserEntity"
	"github.com/dunpju/higo-gin/test/app/Models/User"
	"github.com/dunpju/higo-orm/arm"
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

type UserDao struct {
	*arm.BaseDao
	model *User.Model
}

func NewUserDao() *UserDao {
	dao := &UserDao{model: User.New()}
	dao.BaseDao = arm.NewBaseDao(dao)
	return dao
}

func (this *UserDao) SetModel(model arm.IModel) {
	this.model = model.(*User.Model)
}

func (this *UserDao) GetModel() arm.IModel {
	return this.model
}

func (this *UserDao) Model() *User.Model {
	return User.New()
}

func (this *UserDao) Models() []*User.Model {
	return make([]*User.Model, 0)
}

func (this *UserDao) TX(tx *gorm.DB) *UserDao {
	this.model.TX(tx)
	return this
}

func (this *UserDao) SetData(entity *UserEntity.Entity) arm.IDao {
	return this.model.Builder(this, func() {
		if !entity.PrimaryEmpty() || entity.IsEdit() { //编辑
			if !this.GetById(entity.Id).Exist() {
				DaoException.Throw("不存在", 0)
			}
			this.model.Where(User.Id, "=", entity.Id)
			if entity.Equals(UserEntity.FlagDelete) {
				// todo::填充修改字段
			} else if entity.Equals(UserEntity.FlagUpdate) {
				// todo::填充修改字段
			}
		} else { //新增

		}
	})
}

// Add 添加
func (this *UserDao) Add() (gormDB *gorm.DB, lastInsertId int64) {
	gormDB, lastInsertId = this.model.Insert().LastInsertId()
	this.CheckError(gormDB)
	return
}

// Update 更新
func (this *UserDao) Update() *gorm.DB {
	gormDB, _ := this.model.Update().Exec()
	this.CheckError(gormDB)
	return gormDB
}

// GetBySchoolId schoolId查询
func (this *UserDao) GetById(id int) *User.Model {
	model := this.Model()
	gormDB := this.model.Select().Where(User.Id, "=", id).First(&model)
	this.CheckError(gormDB)
	return model
}

// GetBySchoolIds schoolId集查询
func (this *UserDao) GetBySchoolIds(schoolIds []int64) []*User.Model {
	models := this.Models()
	gormDB := this.model.Select().WhereIn(User.Id, schoolIds).Get(&models)
	this.CheckError(gormDB)
	return models
}

// DeleteBySchoolId 硬删除
func (this *UserDao) DeleteById(id int64) *gorm.DB {
	gormDB, _ := this.model.Delete().Where(User.Id, "=", id).Exec()
	this.CheckError(gormDB)
	return gormDB
}

// Paginate 列表
func (this *UserDao) Paginate(perPage, page uint64) him.Paginate {
	models := this.Models()
	gormDB, paginate := this.model.Select().Paginate(page, perPage, &models)
	this.CheckError(gormDB)
	return paginate
}
