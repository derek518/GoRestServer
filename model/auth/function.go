package model_auth

import (
	"github.com/jinzhu/gorm"
	"GoRestServer/helper"
	"GoRestServer/model/base"
)

// 功能菜单结构体
type Function struct {
	base.Model

	/** 功能名称 */
	Name string `gorm:"type:varchar(64);unique_index;not null" json:"name"`

	/** 访问路径 */
	Url string `gorm:"type:varchar(64);index" json:"url"`

	/** 功能分组 */
	Group string `gorm:"type:varchar(32);not null" json:"group"`

	/** 是否生成菜单  */
	IsMenu bool `json:"is_menu"`

	/** 图标 */
	Icon string `gorm:"type:varchar(128)" json:"icon"`

	/** 序号 */
	Seq int `json:"seq"`

	/** 父功能 id */
	PId *string `json:"p_id" `

	/** 父功能 */
	ParentFunction *Function `gorm:"foreignkey:PId;save_associations:false:"`

	/** 子功能 */
	ChildFunctions []*Function `gorm:"foreignkey:ID"`

	/** 对应的角色 */
	Roles []*Role `gorm:"many2many:role_functions;" json:"-"`
}

type FunctionModel struct {
	SQL *gorm.DB
}

var funcModel = &FunctionModel{}

// 实例化存储对象
func FunctionModelInstance(db *gorm.DB) *FunctionModel {
	if funcModel.SQL != db {
		funcModel.SQL = db
	}
	return funcModel
}

func (fm *FunctionModel) Insert(m interface{}) error {
	err := fm.SQL.Create(m).Error
	return err
}

func (fm *FunctionModel) Save(m interface{}) error {
	err := fm.SQL.Save(m).Error
	return err
}

func (fm *FunctionModel) Update(m interface{}, fields map[string]interface{}) error {
	err := fm.SQL.Model(m).Updates(fields).Error
	return err
}

func (fm *FunctionModel) Delete(m interface{}) error {
	err := fm.SQL.Delete(m).Error
	return err
}

func (fm *FunctionModel) FindOne(id string) interface{} {
	var item Function
	fm.SQL.Where("id = ?", id).First(&item)
	return &item
}

func (fm *FunctionModel) FindSingle(condition string, params ...interface{}) interface{} {
	var item Function
	fm.SQL.Model(&Function{}).Where(condition, params).First(&item)
	return &item
}

func (fm *FunctionModel) FindMore(condition string, params ...interface{}) interface{} {
	rows := make([]*Function, 0)
	fm.SQL.Where(condition, params).Find(&rows)
	return rows
}

func (fm *FunctionModel) FindPage(page, size int, sort, selection string, andCons, orCons map[string]interface{}) (pageBean *base.PageBean) {
	total := 0
	rows := make([]*Function, 0)

	fm.SQL = helper.BuildWhereIntoDB(fm.SQL, andCons, orCons)
	if len(sort) == 0 {
		sort = "created_at desc"
	}
	if len(selection) > 0 {
		fm.SQL = fm.SQL.Select(selection)
	}
	fm.SQL.Limit(size).Offset((page - 1) * size).Order(sort).Find(&rows).Count(&total)
	return &base.PageBean{Page: page, Size: size, Total: total, Rows: rows}
}
