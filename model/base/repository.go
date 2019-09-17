package base

import "github.com/jinzhu/gorm"

// 基础 repository 接口
type Repository interface {

	// 新增
	Insert(m interface{}) error

	// 完全覆盖
	Save(m interface{}) error

	// 更新局部内容
	Update(m interface{}, fields map[string]interface{}) error

	// 删除
	Delete(m interface{}) error

	// 根据 id 查询
	FindOne(id string) interface{}

	// 根据条件 查询单条记录
	FindSingle(condition string, params ...interface{}) interface{}

	// 根据条件查询多个结果
	FindMore(condition string, params ...interface{}) interface{}

	/** 分页查询 */
	FindPage(page, size int, sort, selection string, andCons, orCons map[string]interface{}) (pageBean *PageBean)
}

type RepoWriter struct {
	SQL *gorm.DB
}

func (w *RepoWriter) Insert(m interface{}) error {
	err := w.SQL.Create(m).Error
	return err
}

func (w *RepoWriter) Update(m interface{}) error {
	err := w.SQL.Save(m).Error
	return err
}

func (w *RepoWriter) Delete(m interface{}) error {
	err := w.SQL.Delete(m).Error
	return err
}
