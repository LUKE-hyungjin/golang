package validator

// RegisterRequest validates user registration
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,strong_password"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	FirstName       string `json:"first_name" binding:"max=50"`
	LastName        string `json:"last_name" binding:"max=50"`
}

// LoginRequest validates user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest validates refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest validates password change
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,strong_password"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// UpdateProfileRequest validates profile update
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"max=50"`
	LastName  string `json:"last_name" binding:"max=50"`
	Phone     string `json:"phone" binding:"omitempty,phone"`
}

// CreatePostRequest validates post creation
type CreatePostRequest struct {
	Title     string   `json:"title" binding:"required,min=3,max=200,no_sql_injection"`
	Content   string   `json:"content" binding:"required,min=10,max=10000,no_sql_injection"`
	Slug      string   `json:"slug" binding:"omitempty,slug,max=200"`
	Published bool     `json:"published"`
	Tags      []string `json:"tags" binding:"max=10,dive,min=1,max=30"`
}

// UpdatePostRequest validates post update
type UpdatePostRequest struct {
	Title     string   `json:"title" binding:"omitempty,min=3,max=200,no_sql_injection"`
	Content   string   `json:"content" binding:"omitempty,min=10,max=10000,no_sql_injection"`
	Slug      string   `json:"slug" binding:"omitempty,slug,max=200"`
	Published *bool    `json:"published"`
	Tags      []string `json:"tags" binding:"omitempty,max=10,dive,min=1,max=30"`
}

// PaginationRequest validates pagination parameters
type PaginationRequest struct {
	Page    int    `form:"page" binding:"min=1"`
	PerPage int    `form:"per_page" binding:"min=1,max=100"`
	Sort    string `form:"sort" binding:"omitempty,oneof=created_at updated_at title"`
	Order   string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// SearchRequest validates search parameters
type SearchRequest struct {
	Query string `form:"q" binding:"required,min=2,max=100,no_sql_injection"`
	Type  string `form:"type" binding:"omitempty,oneof=posts users tags"`
	PaginationRequest
}