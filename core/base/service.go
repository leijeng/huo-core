package base

import (
	"github.com/leijeng/huo-core/core"
	"github.com/leijeng/huo-core/core/cache"
	"gorm.io/gorm"
)

func NewService(dbname string) *BaseService {
	return &BaseService{
		DbName: dbname,
	}
}

type BaseService struct {
	DbName string
}

// 数据库
func (s *BaseService) DB() *gorm.DB {
	return core.Db(s.DbName)
}

// 缓存
func (s *BaseService) Cache() cache.ICache {
	return core.Cache
}

// 创建 结构体model
func (s *BaseService) Create(model any) error {
	return s.DB().Create(model).Error
}

// 更新整个模型 结构体model
func (s *BaseService) Save(model any) error {
	return s.DB().Save(model).Error
}

// 根据id获取模型
func (s *BaseService) Get(id any, model any) error {
	return s.DB().First(model, id).Error
}

// 条件删除，模型
func (s *BaseService) DelWhere(model any) error {
	return s.DB().Delete(model).Error
}

// 条件删除，模型 where 为map
func (s *BaseService) DelWhereMap(model any, where map[string]any) error {
	return s.DB().Model(model).Delete(where).Error
}

// 多个id删除
func (s *BaseService) DelIds(model any, ids any) error {
	return s.DB().Delete(model, ids).Error
}

// 分页获取
func (s *BaseService) Page(where any, data any, total *int64, limit, offset int) error {
	return s.DB().Where(where).Limit(limit).Offset(offset).
		Find(data).Limit(-1).Offset(-1).Count(total).Error
}

// 条件跟新
func (s *BaseService) UpdateWhere(model any, where any, updates map[string]any) error {
	return s.DB().Model(model).Where(where).Updates(updates).Error
}

// 模型更新
func (s *BaseService) UpdateWhereModel(where any, updates any) error {
	return s.DB().Where(where).Updates(updates).Error
}

// 根据模型id更新
func (s *BaseService) UpdateById(model any) error {
	return s.DB().Updates(model).Error
}

// 条件查询 model为数组
func (s *BaseService) GetByWhere(where any, model any) error {
	return s.DB().Where(where).Find(model).Error
}

// 条件查询 model为数组
func (s *BaseService) GetByMap(where map[string]any, model any) error {
	return s.DB().Where(where).Find(model).Error
}

func (s *BaseService) Query(where Query, models any) error {
	return s.DB().Scopes(s.MakeCondition(where)).Find(models).Error
}

// 分页获取
func (s *BaseService) QueryPage(where Query, models any, total *int64, limit, offset int) error {
	return s.DB().Scopes(s.MakeCondition(where)).Limit(limit).Offset(offset).
		Find(models).Limit(-1).Offset(-1).Count(total).Error
}

func (s *BaseService) Paginate(pageSize, pageIndex int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageIndex - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *BaseService) MakeCondition(q Query) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		condition := &GormCondition{
			GormPublic: GormPublic{},
			Join:       make([]*GormJoin, 0),
		}
		ResolveSearchQuery(core.Cfg.DBCfg.GetDriver(s.DbName), q, condition, q.TableName())
		for _, join := range condition.Join {
			if join == nil {
				continue
			}
			db = db.Joins(join.JoinOn)
			for k, v := range join.Where {
				db = db.Where(k, v...)
			}
			for k, v := range join.Or {
				db = db.Or(k, v...)
			}
			for _, o := range join.Order {
				db = db.Order(o)
			}
		}
		for k, v := range condition.Where {
			db = db.Where(k, v...)
		}
		for k, v := range condition.Or {
			db = db.Or(k, v...)
		}
		for _, o := range condition.Order {
			db = db.Order(o)
		}
		return db
	}
}
