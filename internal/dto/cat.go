package dto

import (
	"github.com/sprint-id/catsocialx/internal/entity"
)

type (
	ReqAddCat struct {
		Name        string   `json:"name" validate:"required,min=1,max=30"`
		Race        string   `json:"race" validate:"required,oneof=Persian Maine Coon Siamese Ragdoll Bengal Sphynx British Shorthair Abyssinian Scottish Fold Birman"`
		Sex         string   `json:"sex" validate:"required,oneof=male female"`
		AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
		Description string   `json:"description" validate:"required,min=1,max=200"`
		ImageUrls   []string `json:"imageUrls" validate:"required"`
	}
	ParamGetCat struct {
		ID               string `json:"id"`
		Limit            int    `json:"limit"`
		Offset           int    `json:"offset"`
		Race             string `json:"race"`
		Sex              string `json:"sex"`
		IsAlreadyMatched bool   `json:"isAlreadyMatched"`
		AgeInMonth       int    `json:"ageInMonth"`
		Owned            bool   `json:"owned"`
		Search           string `json:"search"`
	}
	ResGetCat struct {
		ID               string   `json:"id"`
		Name             string   `json:"name"`
		Race             string   `json:"race"`
		Sex              string   `json:"sex"`
		AgeInMonth       int      `json:"ageInMonth"`
		ImageUrls        []string `json:"imageUrls"`
		IsAlreadyMatched bool     `json:"isAlreadyMatched"`
		CreatedAt        string   `json:"createdAt"`
	}
)

func (d *ReqAddCat) ToCatEntity(userId string) entity.Cat {
	return entity.Cat{
		Name:        d.Name,
		Race:        d.Race,
		Sex:         d.Sex,
		AgeInMonth:  d.AgeInMonth,
		Description: d.Description,
		ImageUrls:   d.ImageUrls,
		UserID:      userId,
	}
}
