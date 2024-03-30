package dto

import (
	"github.com/syarifid/bankx/internal/entity"
	"github.com/syarifid/bankx/pkg/auth"
)

type (
	ReqRegister struct {
		Email    string `json:"email" validate:"required,email"`
		Name     string `json:"name" validate:"required,min=5,max=50"`
		Password string `json:"password" validate:"required,min=5,max=15"`
	}
	ResRegister struct {
		Email       string `json:"email,omitempty"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
	ReqLogin struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=5,max=15"`
	}
	ResLogin struct {
		Email       string `json:"email,omitempty"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
	ReqUpdateAccount struct {
		ImageURL string `json:"imageUrl" validate:"required,url"`
		Name     string `json:"name" validate:"required,min=5,max=50"`
	}
)

func (d *ReqRegister) ToEntity(cryptCost int) entity.User {
	return entity.User{Name: d.Name, Password: auth.HashPassword(d.Password, cryptCost), Email: d.Email}
}
