package presenter

// LoginInput is presenter to handler login request
type LoginInput struct {
	UserAccess string `json:"user_access"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Name         string `json:"name"`
}

// AuthInfoResponse auth info response
type AuthInfoResponse struct {
	Email  string  `json:"email"`
	Name   string  `json:"name"`
	UserID uint64  `json:"-"`
	SiteID []int64 `json:"site_id"`
	Role   string  `json:"role"`
}
