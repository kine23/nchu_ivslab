package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
// Project項目列表
type Asset struct {
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

// User用戶列表
type User struct {
	Table    string `json:"table" form:"table"`       //數據標記
	Username string `json:"username" form:"username"` //用戶帳號
	Name     string `json:"name" form:"name"`         //姓名
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitAssets(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "IVSLAB23FA01", Manufacturer: "Security.co", ManufactureLocation: "Taiwan", PartName: "SecurityChip-v1", PartNumber: "SPN300AA", SerialNumber: "SSN30A10AA", Organization: "Security-Org", ManufactureDate: "2023-05-15"},
		{ID: "IVSLAB23FA02", Manufacturer: "Network.co", ManufactureLocation: "Taiwan", PartName: "NetworkChip-v1", PartNumber: "NPN300AA", SerialNumber: "NSN30A10AA", Organization: "Network-Org", ManufactureDate: "2023-05-15"},
		{ID: "IVSLAB23FA03", Manufacturer: "CMOS.co", ManufactureLocation: "USA", PartName: "CMOSChip-v1", PartNumber: "CPN300AA", SerialNumber: "CSN30A10AA", Organization: "CMOS-Org", ManufactureDate: "2023-05-15"},
		{ID: "IVSLAB23FA04", Manufacturer: "VideoCodec.co", ManufactureLocation: "USA", PartName: "VideoCodecChip-v1", PartNumber: "VPN300AA", SerialNumber: "VSN30A10AA", Organization: "VideoCodec-Org", ManufactureDate: "2023-05-15"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = s.Insert(ctx, string(assetJSON)
		if err != nil {
			return err
		}
	}

	return nil
}

// InitLedger adds a base set of User to the ledger
func (s *SmartContract) InitUsers(ctx contractapi.TransactionContextInterface) error {
	users := []User{
		{Username: "SFChen", Name: "SFChen"},
	}

	for _, user := range users {
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(user.Username, userJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}
// 初始化智慧合約數據
func (s *IVSContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	err := s.InitAssets(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize projects: %v", err)
	}
	err = s.InitUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize users: %v", err)
	}
	return nil

}
// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, manufacturer string, manufactureLocation string, partname string, partnumber string, serialnumber string, organization string, manufacturedate string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:             		id,
		Manufacturer:          	manufacturer,
		ManufactureLocation:          manufactureLocation,
		PartName:          		partname,
		PartNumber: 		partnumber,
		SerialNumber:		serialnumber,
		Organization:		organization,
		ManufactureDate:		manufacturedate,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, manufacturer string, manufactureLocation string, partname string, partnumber string, serialnumber string, organization string, manufacturedate string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:             		id,
		Manufacturer:          	manufacturer,
		ManufactureLocation:          manufactureLocation,
		PartName:          		partname,
		PartNumber: 		partnumber,
		SerialNumber:		serialnumber,
		Organization:		organization,
		ManufactureDate:		manufacturedate,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOrganization string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldOrganization := asset.Organization
	asset.Organization = newOrganization

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return "", err
	}

	return oldOrganization, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
// GetAllUsers returns all users found in world state
func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var user []*User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		assets = append(user, &user)
	}

	return user, nil
}
