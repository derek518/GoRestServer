package model_auth

import (
	"github.com/jinzhu/gorm"
	"github.com/thoas/go-funk"
	"GoRestServer/helper"
	"GoRestServer/model/base"
)

// 角色结构体
type Role struct {
	base.Model

	/** 角色名称 */
	Name string `gorm:"type:varchar(32);unique;not null" json:"name"`

	/** 角色类别标识 */
	RoleKey string `gorm:"type:varchar(32);unique;not null" json:"role_key"`

	/** 角色描述 */
	Description string `gorm:"type:varchar(128)" json:"description"`

	/** 角色关联的功能 */
	Functions []*Function `gorm:"many2many:role_functions;" json:"functions"`
}

type RoleModel struct {
	SQL *gorm.DB
}

var roleModel = &RoleModel{}

// 实例化存储对象
func RoleModelInstance(db *gorm.DB) *RoleModel {
	if roleModel.SQL != db {
		roleModel.SQL = db
	}
	return roleModel
}

// 因为在Insert/Update时,gorm会对Role.Functions进行自动关联修改，这里进行屏蔽处理
func (rm *RoleModel) FetchFunctions(functions []*Function) []*Function {
	list := funk.Filter(functions, func(fun *Function) bool {
		return fun.ID != ""
	}).([]*Function)

	if len(list) > 0 {
		funcIds := funk.Map(list, func(fun *Function) interface{} {
			return fun.ID
		}).([]interface{})
		list = FunctionModelInstance(rm.SQL).FindMore("id IN (?)", funcIds...).([]*Function)
	}

	return list
}

func (rm *RoleModel) Insert(m interface{}) error {
	item, ok := m.(*Role)
	if ok {
		item.Functions = rm.FetchFunctions(item.Functions)
	}
	err := rm.SQL.Create(m).Error
	return err
}

func (rm *RoleModel) Save(m interface{}) error {
	item, ok := m.(*Role)
	if ok {
		item.Functions = rm.FetchFunctions(item.Functions)
	}

	err := rm.SQL.Save(m).Error
	return err
}

func (rm *RoleModel) Update(m interface{}, fields map[string]interface{}) error {
	//functions, ok := fields["functions"]
	//if ok && len(functions.([]map[string]interface{})) > 0 {
	//
	//}
	err := rm.SQL.Model(m).Updates(fields).Error
	return err
}

func (rm *RoleModel) Delete(m interface{}) error {
	err := rm.SQL.Delete(m).Error
	return err
}

func (rm *RoleModel) FindOne(id string) interface{} {
	var item Role
	rm.SQL.Preload("Functions").Where("id = ?", id).First(&item)
	return &item
}

func (rm *RoleModel) FindSingle(condition string, params ...interface{}) interface{} {
	var item Role
	rm.SQL.Preload("Functions").Where(condition, params).First(&item)
	return &item
}

func (rm *RoleModel) FindMore(condition string, params ...interface{}) interface{} {
	rows := make([]*Role, 0)
	rm.SQL.Preload("Functions").Where(condition, params).Find(&rows)
	return rows
}

func (rm *RoleModel) FindPage(page, size int, sort, selection string, andCons, orCons map[string]interface{}) (pageBean *base.PageBean) {
	total := 0
	rows := make([]*Role, 0)

	rm.SQL = helper.BuildWhereIntoDB(rm.SQL, andCons, orCons)
	if len(sort) == 0 {
		sort = "created_at desc"
	}
	if len(selection) > 0 {
		rm.SQL = rm.SQL.Select(selection)
	}

	rm.SQL.Preload("Functions").Limit(size).Offset((page - 1) * size).Order(sort).Find(&rows).Count(&total)
	return &base.PageBean{Page: page, Size: size, Total: total, Rows: rows}
}
