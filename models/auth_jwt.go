package models

type UserAuth struct {
	ID uint64            `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	Phone string `json:"phone"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
type Todo struct {
	UserID string `json:"user_id"`
	Title string `json:"title"`
}
type AccessDetails struct {
	AccessUuid string
	UserId   string
}
