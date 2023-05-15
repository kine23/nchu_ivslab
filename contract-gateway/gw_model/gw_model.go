package model

// Project項目列表
type Project struct {
	Table               string `json:"table" form:"table"`  //數據庫標記
	Manufacturer        string `json:"Manufacturer"`        //製造商
	ManufactureLocation string `json:"ManufactureLocation"` //製造地點
	PartName            string `json:"PartName"`            //零件名稱
	PartNumber          string `json:"PartNumber"`          //零件批號
	SerialNumber        string `json:"SerialNumber"`        //產品序號
	ManufactureDate     string `json:"ManufactureDate"`     //製造日期
	Item                string `json:"Item"`                //項目
	ID                  string `json:"ID"`                  //項目唯一ID
	Category            string `json:"Category"`            //所屬類別
	Describes           string `json:"Describes"`           //描述
	Developer           string `json:"Developer"`           //開發者
	Organization        string `json:"Organization"`        //組織
}

func (o *Project) Index() string {
	o.Table = "project"
	return o.ID
}

func (o *Project) IndexKey() string {
	return "table~ID~manufacturer~manufacturelocation~partname~partnumber~serialnumber~manufacturedate~organization"
}

func (o *Project) IndexAttr() []string {
	return []string{o.Table, o.ID, o.Manufacturer, o.ManufactureLocation, o.PartName, o.PartNumber, o.SerialNumber, o.ManufactureDate, o.Organization}
}

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
	return "table~username~name"
}

func (o *User) IndexAttr() []string {
	return []string{o.Table, o.Username, o.Name}
}

