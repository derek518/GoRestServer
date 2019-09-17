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

var roleService = &RoleService{}

type RoleService struct {
	model *model_auth.RoleModel
}

func RoleServiceInstance(model *model_auth.RoleModel) *RoleService {
	if roleService.model != model {
		roleService.model = model
	}
	return roleService
}

func RoleCacheKey(key string) string {
	return "role_" + key
}

func (s *RoleService) SaveOrUpdate(item *model_auth.Role, fields map[string]interface{}) error {
	if item == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	cache.Delete(RoleCacheKey("all"))

	// 判断 新增还是更新
	if strings.TrimSpace(item.ID) == "" {
		return s.model.Insert(item)
	} else {
		// Clear cache
		cache.Delete(RoleCacheKey(item.ID))
		return s.model.Update(item, fields)
	}
}

func (s *RoleService) GetByID(id string) *model_auth.Role {
	if strings.TrimSpace(id) == "" {
		return nil
	}

	// Get from cache
	var cacheItem model_auth.Role
	if err := cache.Get(RoleCacheKey(id), &cacheItem); err == nil {
		return &cacheItem
	}

	item := s.model.FindOne(id).(*model_auth.Role)

	// Cache
	if item != nil {
		cache.Add(RoleCacheKey(id), item)
	}

	return item
}

func (s *RoleService) GetByName(name string) *model_auth.Role {
	item := s.model.FindSingle("role_name = ?", name).(*model_auth.Role)
	return item
}

func (s *RoleService) DeleteByID(id string) error {
	// ClearCache
	cache.Delete(RoleCacheKey(id))
	cache.Delete(RoleCacheKey("all"))

	return s.model.Delete(&model_auth.Role{Model: base.Model{ID: id}})
}

func (s *RoleService) GetAll() []*model_auth.Role {
	// Get cache
	cacheItems := make([]*model_auth.Role, 0)
	if err := cache.Get(RoleCacheKey("all"), &cacheItems); err == nil {
		return cacheItems
	}

	items := s.model.FindMore("1=1").([]*model_auth.Role)

	// Cache
	if items != nil {
		cache.Add(RoleCacheKey("all"), items)
	}

	return items
}

func (s *RoleService) GetPage(query *base.QueryCondition) *base.PageBean {
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
