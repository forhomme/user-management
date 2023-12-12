package domain

type UsersDatabase struct {
	RoleId   int    `json:"role_id"`
	RoleName string `json:"role_name"`
	UserId   string `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	Email        string `json:"email" validate:"required,email"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse = LoginResponse

type SignUpRequest struct {
	RoleId   int    `json:"role_id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignUpResponse = LoginResponse

type ChangePasswordRequest struct {
	UserId      string `json:"user_id"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type MenuRequest struct {
	UserId string `json:"user_id" validate:"required,uuid"`
}

type Menu struct {
	RoleView     bool   `json:"role_view"`
	RoleAdd      bool   `json:"role_add"`
	RoleEdit     bool   `json:"role_edit"`
	RoleDelete   bool   `json:"role_delete"`
	MenuId       int    `json:"menu_id"`
	MenuParentId int    `json:"menu_parent_id"`
	MenuName     string `json:"menu_name"`
	MenuPath     string `json:"menu_path"`
	MenuIcon     string `json:"menu_icon"`
}

type MenuResponse struct {
	RoleId   int     `json:"role_id"`
	UserId   string  `json:"user_id"`
	RoleName string  `json:"role_name"`
	Menus    []*Menu `json:"menus"`
}
