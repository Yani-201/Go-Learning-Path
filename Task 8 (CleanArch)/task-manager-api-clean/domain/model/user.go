package model

type UserUpdate struct{
	Username       string    `json:"username"`
	Email string    `json:"email"`

}

type UserInfo struct {
	UserId string `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
}

type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdatePassword struct {
	Password string `json:"password"`
}
