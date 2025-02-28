package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int       `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(255);not null"`
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"default:'male';not null;type:varchar(6)"`
	Role     int        `gorm:"default:1;type:int comment '1表示普通用户，2表示管理员用户'"`
}
