package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const index = "manufacturer~manufacturelocation~serialnumber"

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
// Project項目列表
type Asset struct {
	DocType        	string `json:"docType"` 	      //docType is used to distinguish the various types of objects in state database
	ID                  string `json:"ID"`                  //項目唯一ID
	Manufacturer        string `json:"Manufacturer"`        //製造商
	ManufactureLocation string `json:"ManufactureLocation"` //製造地點
	PartName            string `json:"PartName"`            //零件名稱
	PartNumber          string `json:"PartNumber"`          //零件批號
	SerialNumber        string `json:"SerialNumber"`        //產品序號
	Organization        string `json:"Organization"`        //組織
	ManufactureDate     string `json:"ManufactureDate"`     //製造日期
	TransferDate	string `json:"TransferDate"`        //交易日期
//	Category            string `json:"Category"`            //所屬類別
//	Describes           string `json:"Describes"`           //描述
//	Developer           string `json:"Developer"`           //開發者
}

// User用戶列表
//type User struct {
//	Username string `json:"username" form:"username"` //用戶帳號
//	Name     string `json:"name" form:"name"`         //姓名
//}

// InitLedger adds a base set of assets to the ledger
func (t *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "IVSLAB-S23FA01", Manufacturer: "Security.co", ManufactureLocation: "Taiwan", PartName: "SecurityChip-v1", PartNumber: "SPN300AA", SerialNumber: "SSN30A10AA", Organization: "Security-Org", ManufactureDate: "2023-05-15"},
		{ID: "IVSLAB-N23FA01", Manufacturer: "Network.co", ManufactureLocation: "Taiwan", PartName: "NetworkChip-v1", PartNumber: "NPN300AA", SerialNumber: "NSN30A10AA", Organization: "Network-Org", ManufactureDate: "2023-05-15"},
		{ID: "IVSLAB-C23FA01", Manufacturer: "CMOS.co", ManufactureLocation: "USA", PartName: "CMOSChip-v1", PartNumber: "CPN300AA", SerialNumber: "CSN30A10AA", Organization: "CMOS-Org", ManufactureDate: "2023-05-15"},
		{ID: "IVSLAB-V23FA01", Manufacturer: "VideoCodec.co", ManufactureLocation: "USA", PartName: "VideoCodecChip-v1", PartNumber: "VPN300AA", SerialNumber: "VSN30A10AA", Organization: "VideoCodec-Org", ManufactureDate: "2023-05-15"},
	}

	for _, asset := range assets {
		err := t.CreateAsset(ctx, asset.ID, asset.Manufacturer, asset.ManufactureLocation, asset.PartName, asset.PartNumber, asset.SerialNumber, asset.Organization, asset.ManufactureDate)
		if err != nil {
			return err
		}
	}

	return nil
}

// InitLedger adds a base set of User to the ledger
//func (t *SmartContract) InitUsers(ctx contractapi.TransactionContextInterface) error {
//	users := []User{
//		{Username: "SFChen", Name: "SFChen"},
//	}
//
//	for _, user := range users {
//		userJSON, err := json.Marshal(user)
//		if err != nil {
//			return err
//		}
//
//		err = ctx.GetStub().PutState(user.Username, userJSON)
//		if err != nil {
//			return fmt.Errorf("failed to put to world state. %v", err)
//		}
//	}
//
//	return nil
//}
// 初始化智慧合約數據
//func (t *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
//	err := s.InitAssets(ctx)
//	if err != nil {
//		return fmt.Errorf("failed to initialize projects: %v", err)
//	}
//
//	err = s.InitUsers(ctx)
//	if err != nil {
//		return fmt.Errorf("failed to initialize users: %v", err)
//	}
//
//	return nil
//}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *Asset    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// PaginatedQueryResult structure used for returning paginated query results and metadata
type PaginatedQueryResult struct {
	Records             []*Asset `json:"records"`
	FetchedRecordsCount int32    `json:"fetchedRecordsCount"`
	Bookmark            string   `json:"bookmark"`
}

// CreateAsset initializes a new asset in the ledger
func (t *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetID string, manufacturer string, manufacturelocation string, partname string, partnumber string, serialnumber string, organization string, manufacturedate string) error {
	exists, err := t.AssetExists(ctx, assetID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", assetID)
	}

	asset := &Asset{
		DocType:        		"asset",
		ID:             		assetID,
		Manufacturer:          	manufacturer,
		ManufactureLocation:          manufacturelocation,
		PartName:          		partname,
		PartNumber: 		partnumber,
		SerialNumber:		serialnumber,
		Organization:		organization,
		ManufactureDate:		manufacturedate,
	}
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(assetID, assetBytes)
	if err != nil {
		return err
	}
	ivsIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.Organization, asset.ID})
	if err != nil {
		return err
	}
	value := []byte{0x00}
	return ctx.GetStub().PutState(ivsIndexKey, value)	
}

// ReadAsset retrieves an asset from the ledger
func (t *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetID string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset %s: %v", assetID, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", assetID)
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (t *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, assetID string, manufacturer string, manufacturelocation string, partname string, partnumber string, serialnumber string, organization string, manufacturedate string) error {
	exists, err := t.AssetExists(ctx, assetID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", assetID)
	}

	// overwriting original asset with new asset
	asset := &Asset{
		DocType:        		"asset",
		ID:             		assetID,
		Manufacturer:          	manufacturer,
		ManufactureLocation:          manufacturelocation,
		PartName:          		partname,
		PartNumber: 		partnumber,
		SerialNumber:		serialnumber,
		Organization:		organization,
		ManufactureDate:		manufacturedate,
	}
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(assetID, assetBytes)
	if err != nil {
		return err
	}
	ivsIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.Organization, asset.ID})
	if err != nil {
		return err
	}
	value := []byte{0x00}
	return ctx.GetStub().PutState(ivsIndexKey, value)	
}

// DeleteAsset removes an asset key-value pair from the ledger
func (t *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	asset, err := t.ReadAsset(ctx, assetID)
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(assetID)
	if err != nil {
		return fmt.Errorf("failed to delete asset %s: %v", assetID, err)
	}

	ivsIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.Organization, asset.ID})
	if err != nil {
		return err
	}

	// Delete index entry
	return ctx.GetStub().DelState(ivsIndexKey)
}

// TransferAsset updates the Organization and TransferDate field of asset with given id in world state, and returns the old Organization.
func (t *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, assetID, assetTransferDate string, newOrganization string) (string, error) {
	asset, err := t.ReadAsset(ctx, assetID)
	if err != nil {
		return "", fmt.Errorf("failed to read asset: %v", err)
	}

	oldOrganization := asset.Organization
	asset.Organization = newOrganization
	asset.TransferDate = assetTransferDate
	
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("failed to marshal asset: %v", err)
	}

	err = ctx.GetStub().PutState(assetID, assetBytes)
	if err != nil {
		return "", fmt.Errorf("failed to update asset: %v", err)
	}

	return oldOrganization, nil
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Asset, error) {
	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Asset
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	// return an empty slice instead of nil if there are no assets
	if len(assets) == 0 {
		return []*Asset{}, nil
	}

	return assets, nil
}

// GetAllAssets returns all assets found in world state
func (t *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SmartContract) GetAssetsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SmartContract) QueryAssetsByOrganization(ctx contractapi.TransactionContextInterface, organization string) ([]*Asset, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"asset","organization":"%s"}}`, organization)
	return getQueryResultForQueryString(ctx, queryString)
}

func (t *SmartContract) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SmartContract) GetAssetsByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetStateByRangeWithPagination(startKey, endKey, int32(pageSize), bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	assets, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             assets,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

func (t *SmartContract) QueryAssetsWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {

	return getQueryResultForQueryStringWithPagination(ctx, queryString, int32(pageSize), bookmark)
}

func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	assets, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             assets,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// GetAssetHistory returns the chain of custody for an asset since issuance.
func (t *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", assetID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Asset{
				ID: assetID,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &asset,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

// AssetExists returns true when asset with given ID exists in world state
func (t *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return false, fmt.Errorf("failed to read asset %s from world state. %v", assetID, err)
	}

	return assetBytes != nil, nil
}

// GetAllUsers returns all users found in world state
//func (t *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
//	// range query with empty string for startKey and endKey does an
//	// open-ended query of all assets in the chaincode namespace.
//	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
//	if err != nil {
//		return nil, err
//	}
//	defer resultsIterator.Close()
//
//	var users []*User
//	for resultsIterator.HasNext() {
//		queryResponse, err := resultsIterator.Next()
//		if err != nil {
//			return nil, err
//		}
//
//		var user User
//		err = json.Unmarshal(queryResponse.Value, &user)
//		if err != nil {
//			return nil, err
//		}
//		users = append(users, &user)
//	}
//
//	return users, nil
//}
