package model

type Asset interface {
	GetKey() string
}

// User用戶列表
type User struct {
	Table    string `json:"table" form:"table"`       //數據標記
	Username string `json:"username" form:"username"` //用戶帳號
	Name     string `json:"name" form:"name"`         //姓名
}

func (u *User) GetKey() string {
	return u.Username
}