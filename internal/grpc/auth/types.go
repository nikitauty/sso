package auth

type LoginReq struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
	AppId    int32  `validate:"required"`
}

type RegisterReq struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type IsAdminReq struct {
	UserID int64
}
