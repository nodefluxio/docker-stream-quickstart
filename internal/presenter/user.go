package presenter

import (
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

type UserRequest struct {
	ID         uint64  `json:"id"`
	Email      string  `json:"email" validate:"required"`
	Username   string  `json:"Username" validate:"required"`
	Password   string  `json:"password" validate:"required"`
	RePassword string  `json:"re_password" validate:"required"`
	Fullname   string  `json:"fullname" validate:"required"`
	Avatar     string  `json:"avatar"`
	Role       string  `json:"role" validate:"required"`
	SiteID     []int64 `json:"site_id" validate:"required"`
}

type UserChangePassRequest struct {
	ID         uint64 `json:"id"`
	Email      string `json:"email" validate:"required"`
	Username   string `json:"Username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	RePassword string `json:"re_password" validate:"required"`
}

type UserPaging struct {
	util.PaginationDetails
	Users []*entity.User `json:"users"`
}
