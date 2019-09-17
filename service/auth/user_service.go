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

var userService = &UserService{}

type UserService struct {
	model *model_auth.UserModel
}

func UserServiceInstance(model *model_auth.UserModel) *UserService {
	if userService.model != model {
		userService.model = model
	}
	return userService
}

func UserCacheKey(key string) string {
	return "user_" + key
}

func (s *UserService) SaveOrUpdate(item *model_auth.User, fields map[string]interface{}) error {
	if item == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	// 判断 新增还是更新
	if strings.TrimSpace(item.ID) == "" {
		log.Debug().Msgf("plain password: %s", item.Password)
		item.Password = helper.SHA256(item.Password)
		log.Debug().Msgf("encrypted password: %s", item.Password)

		return s.model.Insert(item)
	} else {
		// Clear cache
		cache.Delete(UserCacheKey(item.ID))
		return s.model.Update(item, fields)
	}
}

func (s *UserService) GetByID(id string) *model_auth.User {
	if strings.TrimSpace(id) == "" {
		return nil
	}

	// Get from cache
	var cacheItem model_auth.User
	if err := cache.Get(UserCacheKey(id), &cacheItem); err == nil {
		return &cacheItem
	}

	item := s.model.FindOne(id).(*model_auth.User)

	// Cache
	if item != nil {
		cache.Add(UserCacheKey(id), item)
	}

	return item
}

func (s *UserService) GetByName(name string) *model_auth.User {
	item := s.model.FindSingle("username = ?", name).(*model_auth.User)
	return item
}

func (s *UserService) DeleteByID(id string) error {
	// ClearCache
	cache.Delete(UserCacheKey(id))

	return s.model.Delete(&model_auth.User{Model: base.Model{ID: id}})
}

func (s *UserService) GetAll() []*model_auth.User {
	items := s.model.FindMore("1=1").([]*model_auth.User)
	return items
}

func (s *UserService) GetPage(query *base.QueryCondition) *base.PageBean {
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

func (s *UserService) Login(login *model_auth.Login) (*model_auth.User, error) {
	if user := s.GetByName(login.Username); user != nil && user.Password == helper.SHA256(login.Password) {
		s.model.SetLogin(user)
		// Clear cache
		cache.Delete(UserCacheKey(user.ID))

		return user, nil
	}
	return nil, errors.New("invalid username or password")
}
