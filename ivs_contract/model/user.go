package model

// User用戶列表
type User struct {
	Table    string `json:"table" form:"table"`       //數據標記
	Username string `json:"username" form:"username"` //用戶帳號
	Name     string `json:"name" form:"name"`         //姓名
}

func (o *User) Index() string {
	o.Table = "user"
	return o.Username
}

func (o *User) IndexKey() string {
	return "table~username"
}

func (o *User) IndexAttr() []string {
	return []string{o.Table, o.Username}
}
