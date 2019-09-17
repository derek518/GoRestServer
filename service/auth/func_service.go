package service_auth

import (
	"errors"
	"github.com/fatih/structs"
	"github.com/rs/zerolog/log"
	"strings"
	"GoRestServer/helper"
	"GoRestServer/model/auth"
	"GoRestServer/model/base"
	"GoRestServer/pkg/cache"
)

var functionService = &FunctionService{}

type FunctionService struct {
	model *model_auth.FunctionModel
}

func FunctionServiceInstance(model *model_auth.FunctionModel) *FunctionService {
	if functionService.model != model {
		functionService.model = model
	}
	return functionService
}

func FunctionCacheKey(key string) string {
	return "function_" + key
}

func (s *FunctionService) SaveOrUpdate(item *model_auth.Function, fields map[string]interface{}) error {
	if item == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	// Clear cache
	cache.Delete(FunctionCacheKey("all"))

	// 判断 新增还是更新
	if strings.TrimSpace(item.ID) == "" {
		return s.model.Insert(item)
	} else {
		// Clear cache
		cache.Delete(FunctionCacheKey(item.ID))
		return s.model.Update(item, fields)
	}
}

func (s *FunctionService) GetByID(id string) *model_auth.Function {
	if strings.TrimSpace(id) == "" {
		return nil
	}

	// Get from cache
	var cacheItem model_auth.Function
	if err := cache.Get(FunctionCacheKey(id), &cacheItem); err == nil {
		return &cacheItem
	}

	item := s.model.FindOne(id).(*model_auth.Function)

	// Cache
	if item != nil {
		cache.Add(FunctionCacheKey(id), item)
	}

	return item
}

func (s *FunctionService) GetByName(name string) *model_auth.Function {
	item := s.model.FindSingle("fun_name = ?", name).(*model_auth.Function)
	return item
}

func (s *FunctionService) DeleteByID(id string) error {
	// ClearCache
	cache.Delete(FunctionCacheKey(id))
	cache.Delete(FunctionCacheKey("all"))

	return s.model.Delete(&model_auth.Function{Model: base.Model{ID: id}})
}

func (s *FunctionService) GetAll() []*model_auth.Function {
	// Get cache
	cacheItems := make([]*model_auth.Function, 0)
	if err := cache.Get(FunctionCacheKey("all"), &cacheItems); err == nil {
		return cacheItems
	}

	items := s.model.FindMore("1=1").([]*model_auth.Function)

	// Cache
	if items != nil {
		cache.Add(FunctionCacheKey("all"), items)
	}

	return items
}

func (s *FunctionService) GetPage(query *base.QueryCondition) *base.PageBean {
	log.Debug().Fields(structs.Map(query)).Caller().Msg("query condition")

	var andCons, orCons map[string]interface{}
	if query.AndCons != nil {
		andCons = helper.BuildSqlWhere(query.AndCons.(map[string]interface{}))
	}
	if query.OrCons != nil {
		orCons = helper.BuildSqlWhere(query.OrCons.(map[string]interface{}))
	}

	items := s.model.FindPage(query.Page, query.Size, query.Sort, query.Selection, andCons, orCons)
	return items
}
