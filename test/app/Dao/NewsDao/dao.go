package NewsDao

import (
	"github.com/dunpju/higo-gin/higo"
	"github.com/dunpju/higo-gin/higo/exceptions/DaoException"
	"github.com/dunpju/higo-gin/higo/sql"
	"github.com/dunpju/higo-gin/test/app/Entity/NewsEntity"
	"github.com/dunpju/higo-gin/test/app/Models/NewsModel"
	"github.com/dunpju/higo-gin/test/app/errcode"
	"github.com/dunpju/higo-utils/utils"
	"strings"
)

type Dao struct {
	model *NewsModel.Impl
}

func New() *Dao {
	return &Dao{model: NewsModel.New()}
}

func (this *Dao) Orm() *higo.Orm {
	return this.model.Orm
}

func (this *Dao) Model() *NewsModel.Impl {
	return NewsModel.New()
}

func (this *Dao) Models() []*NewsModel.Impl {
	return make([]*NewsModel.Impl, 0)
}

func (this *Dao) SetData(entity *NewsEntity.Impl) {
	if entity.IsEdit() { //编辑
		if entity.PriEmpty() {
			DaoException.Throw("NewsId"+errcode.PrimaryIdError.Message(), int(errcode.PrimaryIdError))
		}
		if !this.GetByNewsId(entity.NewsId).Exist() {
			DaoException.Throw(errcode.NotExistError.Message(), int(errcode.NotExistError))
		}
		builder := this.model.Update(this.model.TableName()).Where(NewsModel.NewsId, entity.NewsId)
		if NewsEntity.FlagDelete == entity.Flag() {

		} else {

		}
	} else { //新增
		this.model.Insert(this.model.TableName()).
			Set(NewsModel.NewsId, entity.NewsId).        //主键
			Set(NewsModel.Title, entity.Title).          //标题
			Set(NewsModel.Clicknum, entity.Clicknum).    //点击量
			Set(NewsModel.CreateTime, entity.CreateTime) //创建时间
	}
	this.model.Build()
}

//添加
func (this *Dao) Add() int {
	higo.Result(this.model.Mapper(this.model.GetBuilder()).InsertGetId().Error).Unwrap()
	return int(this.model.LastInsertId())
}

//更新
func (this *Dao) Update() bool {
	higo.Result(this.model.Mapper(this.model.GetBuilder()).Exec().Error).Unwrap()
	return true
}

//id查询
func (this *Dao) GetByNewsId(newsId int, fields ...string) *NewsModel.Impl {
	if len(fields) == 0 {
		fields = append(fields, "*")
	}
	model := this.Model()
	model.Mapper(sql.Select(fields...).
		From(this.model.TableName()).
		Where("`news_id` = ?", newsId).
		ToSql()).Query().Scan(&model)
	return model
}

//id集查询
func (this *Dao) GetByNewsIds(newsIds []interface{}, fields ...string) []*NewsModel.Impl {
	if len(fields) == 0 {
		fields = append(fields, "*")
	}
	models := this.Models()
	this.Model().Mapper(sql.Select(fields...).
		From(this.model.TableName()).
		Where("`news_id` IN (?)", strings.Join(utils.Convert.Slice(newsIds), ",")).
		ToSql()).Query().Scan(&models)
	return models
}

//硬删除
func (this *Dao) DeleteByAdminId(newsId int) {
	higo.Result(this.model.Mapper(sql.Delete(this.model.TableName()).
		DeleteBuilder().
		Where("news_id = ?", newsId).
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
