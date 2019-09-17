package model_auth

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
	"GoRestServer/helper"
	"GoRestServer/model/base"
	"GoRestServer/pkg/config"
)

// 用户登录结构体
type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 用户结构体
type User struct {
	base.Model

	/** 用户名 */
	Username string `gorm:"type:varchar(32);unique_index;not null" json:"username"`

	/** 密码  */
	Password string `gorm:"type:varchar(64);not null" json:"password,omitempty"`

	/** 电话 */
	Phone string `gorm:"type:varchar(11)" json:"phone" `

	/** 电话 */
	Description string `gorm:"type:varchar(64)" json:"description" `

	/** 标志 1 表示这个账号是由管理方为商户添加的账号 */
	IsBuiltin bool `json:"is_builtin" `

	/** 登陆次数 */
	LogonCount int `json:"logon_count" `

	/** 状态  0 正常  */
	Status int `json:"status" `

	/** 最后一次登陆时间 */
	LoginTime *time.Time `gorm:"default:null" json:"login_time"`

	/** 用户对应的角色 */
	Role Role `gorm:"foreignkey:RoleId;save_associations:false:" json:"role,omitempty"`

	/** 外键 */
	RoleId string `gorm:"type:varchar(36)" json:"role_id" `
}

// 校验表单中提交的参数是否合法
func (user *User) Validator() error {

	if ok, err := helper.MatchLetterNumMinAndMax(user.Username, 4, 16, "用户名"); !ok {
		return err
	}
	if ok, err := helper.MatchStrongPassword(user.Password, 6, 13); !ok && strings.TrimSpace(user.ID) == "" {
		return err
	}
	if ok, err := helper.IsPhone(user.Phone); !ok && strings.TrimSpace(user.Phone) != "" {
		return err
	}

	return nil
}

type UserModel struct {
	base.RepoWriter
}

var userModel = &UserModel{}

// 实例化存储对象
func UserModelInstance(db *gorm.DB) *UserModel {
	if userModel.SQL != db {
		userModel.SQL = db
	}
	return userModel
}

func (um *UserModel) Insert(m interface{}) error {
	err := um.SQL.Create(m).Error
	return err
}

func (um *UserModel) Save(m interface{}) error {
	err := um.SQL.Save(m).Error
	return err
}

func (um *UserModel) Update(m interface{}, fields map[string]interface{}) error {
	err := um.SQL.Model(m).Omit("username", "password", "is_builtin", "logon_count").Updates(fields).Error
	return err
}

func (um *UserModel) Delete(m interface{}) error {
	err := um.SQL.Delete(m).Error
	return err
}

func (um *UserModel) FindOne(id string) interface{} {
	var item User
	um.SQL.Preload("Role").Where("id = ?", id).First(&item)
	return &item
}

func (um *UserModel) FindSingle(condition string, params ...interface{}) interface{} {
	var item User
	um.SQL.Preload("Role").Where(condition, params).First(&item)
	return &item
}

func (um *UserModel) FindMore(condition string, params ...interface{}) interface{} {
	rows := make([]*User, 0)
	um.SQL.Preload("Role").Where(condition, params).Find(&rows)
	return rows
}

func (um *UserModel) FindPage(page, size int, sort, selection string, andCons, orCons map[string]interface{}) (pageBean *base.PageBean) {
	total := 0
	rows := make([]*User, 0)

	// 外联表名称加上前缀
	helper.AddTablePrefixForWhere(andCons, config.Database.TablePrefix, "role.")
	helper.AddTablePrefixForWhere(orCons, config.Database.TablePrefix, "role.")
	if config.Database.TablePrefix != "" && len(sort) > 0 {
		sort = strings.ReplaceAll(sort, "role.", config.Database.TablePrefix+"role.")
	}

	um.SQL = helper.BuildWhereIntoDB(um.SQL, andCons, orCons)
	if len(sort) == 0 {
		sort = "created_at desc"
	}
	if len(selection) > 0 {
		um.SQL = um.SQL.Select(selection)
	}

	// 通过JOIN解决外键表排序和检索功能
	joinStr := fmt.Sprintf("JOIN %srole ON %srole.id = %suser.role_id", config.Database.TablePrefix, config.Database.TablePrefix, config.Database.TablePrefix)
	um.SQL.Joins(joinStr).Preload("Role").Limit(size).Offset((page - 1) * size).Order(sort).Find(&rows).Count(&total)
	return &base.PageBean{Page: page, Size: size, Total: total, Rows: rows}
}

func (um *UserModel) SetLogin(m interface{}) error {
	err := um.SQL.Model(m).Update(map[string]interface{}{"logon_count": gorm.Expr("logon_count + ?", 1), "login_time": time.Now()}).Error
	return err
}
