package api_auth

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strings"
	"GoRestServer/api/response"
	"GoRestServer/helper"
	"GoRestServer/model"
	"GoRestServer/model/auth"
	"GoRestServer/model/base"
	"GoRestServer/service/auth"
)

// 添加、修改用户信息
// @Summary 添加、修改用户信息
// @Tags UserApi
// @Accept json
// @Produce json
// @Param id             body string false "用户记录id,新增时id为空"
// @Param user_name      body string true  "用户名称"
// @Param password     	 body string true  "密码"
// @Param role_id      	 body string true  "用户关联的角色"
// @Param phone     	 body string true  "电话号码"
// @Param description    body string false  "用户描述"
// @Param is_builtin     body boolean false  "是否内置账号"
// @Param logon_count    body string false  "用户登录次数"
// @Param login_time     body string false  "用户最后一次登录时间"
// @Param logon_count    body string false  "用户登录次数"
// @Success 200 {object} base.JsonObject
// @Router /api/auth/save_user [post]
func SaveUser(context *gin.Context) {
	var user model_auth.User
	fields := make(map[string]interface{})
	if err := context.ShouldBindBodyWith(&user, binding.JSON); err == nil {
		user.DeletedAt = nil
		user.IsBuiltin = false
		user.LogonCount = 0
		user.LoginTime = nil

		// 获取需要更新的字段
		if body, ok := context.Get(gin.BodyBytesKey); ok {
			json.Unmarshal(body.([]byte), &fields)
		}
		userService := service_auth.UserServiceInstance(model_auth.UserModelInstance(model.SQL))
		if err := userService.SaveOrUpdate(&user, fields); err == nil {
			response.Succeed(context, helper.SaveStatusOK, nil)
		} else {
			response.FailedWithOK(context, helper.SaveStatusErr, err)
		}
	} else {
		response.Failed(context, http.StatusUnprocessableEntity, helper.SaveStatusErr, err)
	}
}

// 删除用户信息
// @Summary 删除用户信息
// @Tags UserApi
// @Produce json
// @Param id query string true "用户记录id"
// @Success 200 {object} base.JsonObject
// @Router /api/auth/delete_user [post]
func DeleteUser(context *gin.Context) {
	id := context.Query("id")
	if strings.TrimSpace(id) == "" {
		response.FailedWithOK(context, helper.NoneParamErr, nil)
	}
	userService := service_auth.UserServiceInstance(model_auth.UserModelInstance(model.SQL))
	if err := userService.DeleteByID(id); err == nil {
		response.Succeed(context, helper.DeleteStatusOK, nil)
	} else {
		response.FailedWithOK(context, helper.DeleteStatusErr, err)
	}
}

// 获取单个用户信息
// @Summary 获取单个用户信息,传id按id查，传name按用户名称查
// @Tags UserApi
// @Produce json
// @Param id        query string false "用户记录id"
// @Param user_name	query string false "用户名"
// @Success 200 {object} model_auth.User "base.JsonObject中content的内容"
// @Router /api/auth/get_user [get]
func GetUser(context *gin.Context) {
	id := context.Query("id")
	name := context.Query("user_name")
	if strings.TrimSpace(name+id) == "" {
		response.FailedWithOK(context, helper.NoneParamErr, nil)
	}
	userService := service_auth.UserServiceInstance(model_auth.UserModelInstance(model.SQL))
	var user *model_auth.User
	if id != "" {
		user = userService.GetByID(id)
	} else {
		user = userService.GetByName(name)
	}
	response.Succeed(context, 0, user)
}

// 获取所有用户信息
// @Summary 获取所有用户信息
// @Tags UserApi
// @Produce json
// @Success 200 {object} base.JsonObject
// @Router /api/get_user_all [get]
func GetUserAll(context *gin.Context) {
	userService := service_auth.UserServiceInstance(model_auth.UserModelInstance(model.SQL))
	list := userService.GetAll()
	response.Succeed(context, 0, list)
}

// 用户信息分页查询
// @Summary 用户信息分页查询；查询参数为QueryCondition。
// @Tags UserApi
// @Accept json
// @Produce json
// @Param page body integer true "页码"
// @Param size body integer true "每页显示最大行"
// @Param sort body string false "排序"
// @Param selection body string false "字段删选"
// @Param and_cons body {object} false "查询条件(And)"
// @Param or_cons body {object} false "查询条件(Or)"
// @Success 200 {object} base.PageBean
// @Router /api/get_user_page [post]
func GetUserPage(context *gin.Context) {
	query := &base.QueryCondition{Page: 1, Size: 10}
	context.BindJSON(query)
	userService := service_auth.UserServiceInstance(model_auth.UserModelInstance(model.SQL))
	pageBean := userService.GetPage(query)
	response.Succeed(context, 0, pageBean)
}
