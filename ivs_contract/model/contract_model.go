package model

type Project struct {
	Table			string `json:"table" form:"table"`  //數據庫標記
	Manufacturer         	string `json:"Manufacturer"`        //製造商
	ManufactureLocation  	string `json:"ManufactureLocation"` //製造地點
	PartName             	string `json:"PartName"`            //零件名稱
	BatchNumber          	string `json:"BatchNumber"`         //零件批號
	SerialNumber         	string `json:"SerialNumber"`        //產品序號
	ManufactureDate      	string `json:"ManufactureDate"`     //製造日期
	Name                 	string `json:"Name"`                //名稱
	ID                   	string `json:"ID"`                  //項目唯一ID
	Category             	string `json:"Category"`            //所屬類別
	Describes            	string `json:"Describes"`           //描述
	Developer            	string `json:"Developer"`           //開發者
	Organization         	string `json:"Organization"`        //組織
}

func (o *Project) Index() string {
	o.Table = "project"
	return o.ID
}

func (o *Project) IndexKey() string {
	return "table~ID~PartName~BatchNumber~SerialNumber"
}

func (o *Project) IndexAttr() []string {
	return []string{o.Table, o.ID, o.PartName, o.BatchNumber, o.SerialNumber}
}
