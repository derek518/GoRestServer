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

// 添加、修改功能信息
// @Summary 添加、修改功能信息
// @Tags FunctionApi
// @Accept json
// @Produce json
// @Param id             body string false "功能记录id,新增时id为空"
// @Param name       	 body string true  "功能名称"
// @Param url       	 body string true  "功能URL"
// @Param action         body string true  "HTTP请求类型"
// @Param is_menu        body boolean true  "是否生成菜单"
// @Param seq            body integer true  "序号"
// @Param icon    		 body string true  "图标"
// @Param p_id           body string true  "父功能id"
// @Success 200 {object} base.JsonObject
// @Router /api/auth/save_function [post]
func SaveFunction(context *gin.Context) {
	var function model_auth.Function
	fields := make(map[string]interface{})
	if err := context.ShouldBindBodyWith(&function, binding.JSON); err == nil {
		function.DeletedAt = nil

		// 获取需要更新的字段
		if body, ok := context.Get(gin.BodyBytesKey); ok {
			json.Unmarshal(body.([]byte), &fields)
		}

		functionService := service_auth.FunctionServiceInstance(model_auth.FunctionModelInstance(model.SQL))
		if err := functionService.SaveOrUpdate(&function, fields); err == nil {
			response.Succeed(context, helper.SaveStatusOK, nil)
		} else {
			response.FailedWithOK(context, helper.SaveStatusErr, err)
		}
	} else {
		response.Failed(context, http.StatusUnprocessableEntity, helper.SaveStatusErr, err)
	}
}

// 删除功能信息
// @Summary 删除功能信息
// @Tags FunctionApi
// @Produce json
// @Param id query string true "功能记录id"
// @Success 200 {object} base.JsonObject
// @Router /api/auth/delete_function [post]
func DeleteFunction(context *gin.Context) {
	id := context.Query("id")
	if strings.TrimSpace(id) == "" {
		response.FailedWithOK(context, helper.NoneParamErr, nil)
	}
	functionService := service_auth.FunctionServiceInstance(model_auth.FunctionModelInstance(model.SQL))
	if err := functionService.DeleteByID(id); err == nil {
		response.Succeed(context, helper.DeleteStatusOK, nil)
	} else {
		response.FailedWithOK(context, helper.DeleteStatusErr, err)
	}
}

// 获取单个功能信息
// @Summary 获取单个功能信息,传id按id查，传name按功能名称查
// @Tags FunctionApi
// @Produce json
// @Param id        query string false "功能记录id"
// @Param name  	query string false "功能名"
// @Success 200 {object} base.JsonObject
// @Router /api/auth/get_function [get]
func GetFunction(context *gin.Context) {
	id := context.Query("id")
	name := context.Query("name")
	if strings.TrimSpace(name+id) == "" {
		response.FailedWithOK(context, helper.NoneParamErr, nil)
	}
	functionService := service_auth.FunctionServiceInstance(model_auth.FunctionModelInstance(model.SQL))
	var function *model_auth.Function
	if id != "" {
		function = functionService.GetByID(id)
	} else {
		function = functionService.GetByName(name)
	}
	response.Succeed(context, 0, function)
}

// 获取所有功能信息
// @Summary 获取所有功能信息
// @Tags FunctionApi
// @Produce json
// @Success 200 {object} base.JsonObject
// @Router /api/get_function_all [get]
func GetFunctionAll(context *gin.Context) {
	functionService := service_auth.FunctionServiceInstance(model_auth.FunctionModelInstance(model.SQL))
	list := functionService.GetAll()
	response.Succeed(context, 0, list)
}

// 功能信息分页查询
// @Summary 用户信息分页查询。查询参数为QueryCondition。
// @Tags FunctionApi
// @Accept json
// @Produce json
// @Param page body integer true "页码"
// @Param size body integer true "每页显示最大行"
// @Param sort body string false "排序"
// @Param selection body string false "字段删选"
// @Param and_cons body {object} false "查询条件(And)"
// @Param or_cons body {object} false "查询条件(Or)"
// @Success 200 {object} base.PageBean
// @Router /api/get_function_page [post]
func GetFunctionPage(context *gin.Context) {
	query := &base.QueryCondition{Page: 1, Size: 10}
	context.BindJSON(query)
	functionService := service_auth.FunctionServiceInstance(model_auth.FunctionModelInstance(model.SQL))
	pageBean := functionService.GetPage(query)
	response.Succeed(context, 0, pageBean)
}
