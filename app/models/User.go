package models

import (
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Name     string `json:"name" xml:"name" form:"name" query:"name"`
	Password string `json:"-" xml:"-" form:"-" query:"-"`
	Email    string `json:"email" xml:"email" form:"email" query:"email"`
	RoleID   uint   `gorm:"column:role_id" json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type UserDto struct {
	Id       uint   `json:"id"`
	Name     string `json:"name" xml:"name" form:"name" query:"name"`
	Password string `json:"password" xml:"password" form:"-" query:"-"`
	Email    string `json:"email" xml:"email" form:"email" query:"email"`
	RoleID   uint   `json:"role_id"`
}
