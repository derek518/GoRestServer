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

// 添加、修改角色信息
// @Summary 添加、修改角色信息
// @Tags RoleApi
// @Accept json
// @Produce json
// @Param id             body string false "角色记录id,新增时id为空"
// @Param name       	 body string true  "角色名称"
// @Param role_key     	 body string true  "角色类别标识"
// @Param description    body string false  "角色描述"
// @Param functions      body model_auth.Function false  "角色关联的功能"
// @Success 200 {object} base.JsonObject
// @Router /api/auth/save_role [post]
func SaveRole(context *gin.Context) {
	var role model_auth.Role
	fields := make(map[string]interface{})
	if err := context.ShouldBindBodyWith(&role, binding.JSON); err == nil {
		role.DeletedAt = nil

		// 获取需要更新的字段
		if body, ok := context.Get(gin.BodyBytesKey); ok {
			json.Unmarshal(body.([]byte), &fields)
		}
		roleService := service_auth.RoleServiceInstance(model_auth.RoleModelInstance(model.SQL))
		if err := roleService.SaveOrUpdate(&role, fields); err == nil {
			response.Succeed(context, helper.SaveStatusOK, nil)
		} else {
			response.FailedWithOK(context, helper.SaveStatusErr, err)
		}
	} else {
		response.Failed(context, http.StatusUnprocessableEntity, helper.SaveStatusErr, err)
	}
}

// 删除角色信息
// @Summary 删除角色信息
// @Tags RoleApi
// @Produce json
// @Param id query string true "角色记录id"
// @Success 200 {object} base.JsonObject
// @Router /api/auth/delete_role [post]
func DeleteRole(context *gin.Context) {
	id := context.Query("id")
	if strings.TrimSpace(id) == "" {
		response.FailedWithOK(context, helper.NoneParamErr, nil)
	}
	roleService := service_auth.RoleServiceInstance(model_auth.RoleModelInstance(model.SQL))
	if err := roleService.DeleteByID(id); err == nil {
		response.Succeed(context, helper.DeleteStatusOK, nil)
	} else {
		response.FailedWithOK(context, helper.DeleteStatusErr, err)
	}
}

// 获取单个角色信息
// @Summary 获取单个角色信息,传id按id查，传name按角色名称查
// @Tags RoleApi
// @Produce json
// @Param id        query string false "角色记录id"
// @Param name  	query string false "角色名"
// @Success 200 {object} model_auth.Role "base.JsonObject中content的内容"
// @Router /api/auth/get_role [get]
func GetRole(context *gin.Context) {
	id := context.Query("id")
	name := context.Query("name")
	if strings.TrimSpace(name+id) == "" {
		response.FailedWithOK(context, helper.NoneParamErr, nil)
	}
	roleService := service_auth.RoleServiceInstance(model_auth.RoleModelInstance(model.SQL))
	var role *model_auth.Role
	if id != "" {
		role = roleService.GetByID(id)
	} else {
		role = roleService.GetByName(name)
	}
	response.Succeed(context, 0, role)
}

// 获取所有角色信息
// @Summary 获取所有角色信息
// @Tags RoleApi
// @Produce json
// @Success 200 {object} base.JsonObject
// @Router /api/get_role_all [get]
func GetRoleAll(context *gin.Context) {
	roleService := service_auth.RoleServiceInstance(model_auth.RoleModelInstance(model.SQL))
	list := roleService.GetAll()
	response.Succeed(context, 0, list)
}

// 角色信息分页查询
// @Summary 角色信息分页查询；查询参数为QueryCondition。
// @Tags RoleApi
// @Accept json
// @Produce json
// @Param page body integer true "页码"
// @Param size body integer true "每页显示最大行"
// @Param sort body string false "排序"
// @Param selection body string false "字段删选"
// @Param and_cons body {object} false "查询条件(And)"
// @Param or_cons body {object} false "查询条件(Or)"
// @Success 200 {object} base.PageBean
// @Router /api/get_role_page [post]
func GetRolePage(context *gin.Context) {
	query := &base.QueryCondition{Page: 1, Size: 10}
	context.BindJSON(query)
	roleService := service_auth.RoleServiceInstance(model_auth.RoleModelInstance(model.SQL))
	pageBean := roleService.GetPage(query)
	response.Succeed(context, 0, pageBean)
}
