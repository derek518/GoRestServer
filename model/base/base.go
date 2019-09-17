package base

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	"time"
)

// Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//    type User struct {
//      BaseModel
//    }
type Model struct {
	ID        string     `gorm:"type:varchar(36);primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (m *Model) BeforeCreate(scope *gorm.Scope) error {
	if m.ID == "" {
		id := xid.New()
		scope.Set("ID", &id)
		m.ID = id.String()
	}
	return nil
}
