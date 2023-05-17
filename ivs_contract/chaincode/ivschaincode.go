package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"reflect"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const index = "madein~serialnumber"
const (RoleAdmin = "admin")

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// HistoryQueryResult structure used for returning result of history query.
type HistoryQueryResult struct {
	Record    *Asset    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// PaginatedQueryResult structure used for returning paginated query results and metadata.
type PaginatedQueryResult struct {
	Records             []*Asset `json:"records"`
	FetchedRecordsCount int32    `json:"fetchedRecordsCount"`
	Bookmark            string   `json:"bookmark"`
}

// Part represents a product in the ledger.
type Part struct {
	DocType             string `json:"docType"`             // DocType is used to distinguish the various types of objects in state database
	PID					string `json:"PID"`					// 零件唯ID
	Manufacturer        string `json:"Manufacturer"`        // 製造商
	ManufactureLocation string `json:"ManufactureLocation"` // 製造地點
	PartName            string `json:"PartName"`            // 零件名稱
	PartNumber          string `json:"PartNumber"`          // 零件批號
	Organization        string `json:"Organization"`        // 組織
	ManufactureDate     string `json:"ManufactureDate"`     // 零件製造日期
	TransferDate        string `json:"TransferDate"`        // 零件交易日期
}

// Asset represents a product in the ledger.
type Asset struct {
	DocType             string `json:"docType"`             // DocType is used to distinguish the various types of objects in state database
	ID                	string `json:"ID"`                	// 項目唯一ID
	MadeBy        		string `json:"MadeBy"`        		// 品牌商
	MadeIn 				string `json:"MadeIn"` 				// 組裝地點
	SerialNumber        string `json:"SerialNumber"`        // 產品序號
	SecurityChip        Part   `json:"SecurityChip"`        // 安全晶片組織
	NetworkChip         Part   `json:"NetworkChip"`         // 網路晶片組織
	CMOSChip            Part   `json:"CMOSChip"`            // CMOS晶片組織
	VideoCodecChip		Part   `json:"VideoCodecChip"`      // VideoCodec晶片組織
	ProductionDate      string `json:"ProductionDate"`      // 產品生產日期
}

// InitLedger adds a base set of parts to the ledger.
func (t *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	parts := []Part{
		{PID: "IVSLAB-S23FA0001", Manufacturer: "Security.Co", ManufactureLocation: "Taiwan", PartName: "SecurityChip-v1", PartNumber: "SPN3R1C00AA1", ManufactureDate: "2023-05-15", Organization: "Security-Org"},
		{PID: "IVSLAB-N23FA0001", Manufacturer: "Network.Co", ManufactureLocation: "Taiwan", PartName: "NetworkChip-v1", PartNumber: "NPN3R1C00AA1", ManufactureDate: "2023-05-15", Organization: "Network-Org"},
		{PID: "IVSLAB-C23FA0001", Manufacturer: "CMOS.Co", ManufactureLocation: "USA", PartName: "CMOSChip-v1", PartNumber: "CPN3R1C00AA1", ManufactureDate: "2023-05-15", Organization: "CMOS-Org"},
		{PID: "IVSLAB-V23FA0001", Manufacturer: "VideoCodec.Co", ManufactureLocation: "USA", PartName: "VideoCodecChip-v1", PartNumber: "VPN3R1C00AA1", ManufactureDate: "2023-05-15", Organization: "VideoCodec-Org"},
	}

	for _, part := range parts {
		err := t.CreatePart(ctx, part.PID, part.Manufacturer, part.ManufactureLocation, part.PartName, part.PartNumber, part.ManufactureDate, part.Organization)
		if err != nil {
			return err
		}
	}

	// Transfer the parts to Brand-Org
	for _, part := range parts {
		_, err := t.TransferPart(ctx, part.PID, "2023-05-15", "Brand-Org")
		if err != nil {
			return err
		}
	}

	// Create an asset with the transferred parts
	asset := Asset{
		ID:                  	"IVSLAB-PVC23FG0001", 
		MadeBy:        		"Brand.Co", 
		MadeIn: 		"Taiwan", 
		SerialNumber:           "IVSPN902300AACDC01", 
		SecurityChip:        	parts[0], 
		NetworkChip:         	parts[1], 
		CMOSChip:            	parts[2], 
		VideoCodecChip:         parts[3], 
		ProductionDate:        	"2023-05-15",
	}

	err := t.CreateAsset(ctx, asset.ID, asset.MadeBy, asset.MadeIn, asset.SerialNumber, asset.SecurityChip, asset.NetworkChip, asset.CMOSChip, asset.VideoCodecChip, asset.ProductionDate)
	if err != nil {
		return err
	}

	return nil
}

// CheckRole checks if the user has the required role.
func CheckRole(ctx contractapi.TransactionContextInterface, requiredRole string) error {
	attr, ok, err := ctx.GetClientIdentity().GetAttributeValue("role")
	if err != nil {
		return fmt.Errorf("failed to get the 'role' attribute: %v", err)
	}
	if !ok || attr != requiredRole {
		return fmt.Errorf("unauthorized user role")
	}
	return nil
}

// CheckRoleAndRetrieveData checks the role and retrieves the data.
func CheckRoleAndRetrieveData(ctx contractapi.TransactionContextInterface, role string, ID string) ([]byte, error) {
	err := CheckRole(ctx, role)
	if err != nil {
		return nil, err
	}

	dataBytes, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if dataBytes == nil {
		return nil, fmt.Errorf("the data %s does not exist", ID)
	}
	return dataBytes, nil
}

// CreateCompositeKeyAndPutState creates a composite key and puts it into the state.
func CreateCompositeKeyAndPutState(ctx contractapi.TransactionContextInterface, objectType string, attributes []string, ID string, dataBytes []byte) error {
	err := ctx.GetStub().PutState(ID, dataBytes)
	if err != nil {
		return err
	}
	compositeKey, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return err
	}
	value := []byte{0x00}
	return ctx.GetStub().PutState(compositeKey, value)
}

// DeleteStateAndCompositeKey deletes the state and the composite key.
func DeleteStateAndCompositeKey(ctx contractapi.TransactionContextInterface, objectType string, attributes []string, ID string) error {
	err := ctx.GetStub().DelState(ID)
	if err != nil {
		return fmt.Errorf("failed to delete data %s: %v", ID, err)
	}
	compositeKey, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return err
	}
	// Delete index entry
	return ctx.GetStub().DelState(compositeKey)
}

// CheckExists returns true when part with given ID exists in world state.
func (t *SmartContract) checkExistence(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	err := CheckRole(ctx, RoleAdmin)
	if err != nil {
		return false, err
	}
	bytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state. %v", err)
	}

	return bytes != nil, nil
}

// PartExists returns true when part with given ID exists in world state.
func (t *SmartContract) PartExists(ctx contractapi.TransactionContextInterface, partID string) (bool, error) {
	return t.checkExistence(ctx, partID)
}

// AssetExists returns true when asset with given ID exists in world state
func (t *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	return t.checkExistence(ctx, assetID)
}

// CreateItem initializes a new item in the ledger.
func (t *SmartContract) createItem(ctx contractapi.TransactionContextInterface, id string, item interface{}) error {
	err := CheckRole(ctx, RoleAdmin)
	if err != nil {
		return err
	}

	exists, err := t.checkExistence(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the item %s already exists", id)
	}

	dataBytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal data %v: %v", item, err)
	}

	var attributes []string
	switch v := item.(type) {
	case *Part:
		attributes = []string{v.Manufacturer, v.PID}
	case *Asset:
		attributes = []string{v.MadeBy, v.ID}
	default:
		return fmt.Errorf("unknown data type: %v", reflect.TypeOf(item))
	}

	return CreateCompositeKeyAndPutState(ctx, index, attributes, id, dataBytes)
}

// CreatePart initializes a new part in the ledger.
func (t *SmartContract) CreatePart(ctx contractapi.TransactionContextInterface, partID, manufacturer string, manufacturelocation string, partname string, partnumber string, manufacturedate string, organization string) error {
	part := &Part{
		DocType:             "part",
		PID:                 partID,
		Manufacturer:        manufacturer,
		ManufactureLocation: manufacturelocation,
		PartName:            partname,
		PartNumber:          partnumber,
		Organization:        organization,
		ManufactureDate:     manufacturedate,
	}

	return t.createItem(ctx, partID, part)
}
// CreateAsset initializes a new part in the ledger.
func (t *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, madeBy string, madeIn string, serialNumber string, securityChip Part, networkChip Part, cmosChip Part, videoCodecChip Part, productionDate string) error {
	// Ensure all parts belong to 'Brand-Org'
	parts := []Part{securityChip, networkChip, cmosChip, videoCodecChip}
	for _, part := range parts {
		if part.Organization != "Brand-Org" {
			return fmt.Errorf("part %s does not belong to Brand-Org", part.PID)
		}
	}

	// Create the asset
	asset := Asset{
		ID:              id,
		MadeBy:          madeBy,
		MadeIn:          madeIn,
		SerialNumber:    serialNumber,
		SecurityChip:    securityChip,
		NetworkChip:     networkChip,
		CMOSChip:        cmosChip,
		VideoCodecChip:  videoCodecChip,
		ProductionDate:  productionDate,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	// Use CreateCompositeKeyAndPutState to put the asset into the state
	attributes := []string{madeBy, id}
	return CreateCompositeKeyAndPutState(ctx, "asset", attributes, id, assetJSON)
}

// TransferPart updates the Organization field of Part with given partID in world state, and returns the old Organization.
func (t *SmartContract) TransferPart(ctx contractapi.TransactionContextInterface, partID string, assetTransferDate string, newOrganization string) (string, error) {
	partBytes, err := CheckRoleAndRetrieveData(ctx, RoleAdmin, partID)
	if err != nil {
		return "", err
	}

	var part Part
	err = json.Unmarshal(partBytes, &part)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal part %s: %v", string(partBytes), err)
	}

	oldOrganization := part.Organization
	part.Organization = newOrganization
	part.TransferDate = assetTransferDate

	partBytes, err = json.Marshal(part)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(partID, partBytes)
	if err != nil {
		return "", fmt.Errorf("failed to update part %s: %v", partID, err)
	}

	return oldOrganization, nil
}

// Read retrieves a asset/part from the ledger.
func read(ctx contractapi.TransactionContextInterface, id string, objectType reflect.Type) (interface{}, error) {
	dataBytes, err := CheckRoleAndRetrieveData(ctx, RoleAdmin, id)
	if err != nil {
		return nil, err
	}

	data := reflect.New(objectType).Interface()
	err = json.Unmarshal(dataBytes, data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data %s: %v", string(dataBytes), err)
	}

	return data, nil
}

// ReadPart retrieves a Part from the ledger.
func (t *SmartContract) ReadPart(ctx contractapi.TransactionContextInterface, partID string) (*Part, error) {
	data, err := read(ctx, partID, reflect.TypeOf(Part{}))
	if err != nil {
		return nil, err
	}

	return data.(*Part), nil
}

// ReadAsset retrieves an asset from the ledger.
func (t *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetID string) (*Asset, error) {
	data, err := read(ctx, assetID, reflect.TypeOf(Asset{}))
	if err != nil {
		return nil, err
	}

	return data.(*Asset), nil
}

// Delete retrieves a asset/part from the ledger.
func delete(ctx contractapi.TransactionContextInterface, id string, objectType reflect.Type) error {
	data, err := read(ctx, id, objectType)
	if err != nil {
		return err
	}

	var attributes []string
	switch v := data.(type) {
	case *Part:
		attributes = []string{v.Organization, v.PID}
	case *Asset:
		attributes = []string{v.MadeBy, v.ID}
	default:
		return fmt.Errorf("unknown data type: %v", objectType)
	}

	return DeleteStateAndCompositeKey(ctx, index, attributes, id)
}

// DeletePart removes a part key-value pair from the ledger.
func (t *SmartContract) DeletePart(ctx contractapi.TransactionContextInterface, partID string) error {
	return delete(ctx, partID, reflect.TypeOf(Part{}))
}

// DeleteAsset removes an asset key-value pair from the ledger
func (t *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	return delete(ctx, assetID, reflect.TypeOf(Asset{}))
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface, dataType reflect.Type) (interface{}, error) {
	var data []interface{}
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		item := reflect.New(dataType).Interface()
		err = json.Unmarshal(queryResult.Value, item)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	// return an empty slice instead of nil if there are no data
	if len(data) == 0 {
		return reflect.MakeSlice(reflect.SliceOf(dataType), 0, 0).Interface(), nil
	}

	return data, nil
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string, dataType reflect.Type) (interface{}, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator, dataType)
}

// GetAll asset/part from the ledger.
func getAll(ctx contractapi.TransactionContextInterface, objectType reflect.Type) (interface{}, error) {
	err := CheckRole(ctx, RoleAdmin)
	if err != nil {
		return nil, err
	}
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	data, err := constructQueryResponseFromIterator(resultsIterator, objectType)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetAllParts returns all parts found in world state.
func (t *SmartContract) GetAllParts(ctx contractapi.TransactionContextInterface) ([]*Part, error) {
	data, err := getAll(ctx, reflect.TypeOf(Part{}))
	if err != nil {
		return nil, err
	}

	return data.([]*Part), nil
}

// GetAllAssets returns all assets found in world state.
func (t *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	data, err := getAll(ctx, reflect.TypeOf(Asset{}))
	if err != nil {
		return nil, err
	}

	return data.([]*Asset), nil
}

// Get asset/part By Range from the ledger.
func getByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string, objectType reflect.Type) (interface{}, error) {
	err := CheckRole(ctx, RoleAdmin)
	if err != nil {
		return nil, err
	}
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	items, err := constructQueryResponseFromIterator(resultsIterator, objectType)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetPartsByRange returns all parts in the given range.
func (t *SmartContract) GetPartsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Part, error) {
	parts, err := getByRange(ctx, startKey, endKey, reflect.TypeOf(Part{}))
	if err != nil {
		return nil, err
	}

	return parts.([]*Part), nil
}

// GetAssetsByRange returns all assets in the given range.
func (t *SmartContract) GetAssetsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Asset, error) {
	assets, err := getByRange(ctx, startKey, endKey, reflect.TypeOf(Asset{}))
	if err != nil {
		return nil, err
	}

	return assets.([]*Asset), nil
}

// Query asset/part By Owner from the ledger.
func queryByOwner(ctx contractapi.TransactionContextInterface, ownerKey, ownerValue string, objectType reflect.Type) (interface{}, error) {
	err := CheckRole(ctx, RoleAdmin)
	if err != nil {
		return nil, err
	}
	queryString := fmt.Sprintf(`{"selector":{"docType":"%s","%s":"%s"}}`, objectType.Name(), ownerKey, ownerValue)
	items, err := getQueryResultForQueryString(ctx, queryString, objectType)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// QueryPartsByOwner returns all parts owned by the given organization.
func (t *SmartContract) QueryPartsByOwner(ctx contractapi.TransactionContextInterface, organization string) ([]*Part, error) {
	parts, err := queryByOwner(ctx, "organization", organization, reflect.TypeOf(Part{}))
	if err != nil {
		return nil, err
	}

	return parts.([]*Part), nil
}

// QueryAssetsByOwner returns all assets made by the given manufacturer.
func (t *SmartContract) QueryAssetsByOwner(ctx contractapi.TransactionContextInterface, madeby string) ([]*Asset, error) {
	assets, err := queryByOwner(ctx, "madeby", madeby, reflect.TypeOf(Asset{}))
	if err != nil {
		return nil, err
	}

	return assets.([]*Asset), nil
}

// Query asset/part All from the ledger.
func queryItems(ctx contractapi.TransactionContextInterface, queryString string, objectType reflect.Type) (interface{}, error) {
	err := CheckRole(ctx, RoleAdmin)
	if err != nil {
		return nil, err
	}
	items, err := getQueryResultForQueryString(ctx, queryString, objectType)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// QueryParts returns all parts that satisfy the provided query string.
func (t *SmartContract) QueryParts(ctx contractapi.TransactionContextInterface, queryString string) ([]*Part, error) {
	parts, err := queryItems(ctx, queryString, reflect.TypeOf(Part{}))
	if err != nil {
		return nil, err
	}

	return parts.([]*Part), nil
}

// QueryAssets returns all assets that satisfy the provided query string.
func (t *SmartContract) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	assets, err := queryItems(ctx, queryString, reflect.TypeOf(Asset{}))
	if err != nil {
		return nil, err
	}

	return assets.([]*Asset), nil
}

func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string, objectType reflect.Type) (*PaginatedQueryResult, error) {
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	data, err := constructQueryResponseFromIterator(resultsIterator, objectType)
	if err != nil {
		return nil, err
	}

	// Convert the data to a slice of *Asset
	items, ok := data.([]*Asset)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to []*Asset")
	}

	return &PaginatedQueryResult{
		Records:             items,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

func queryItemsWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int, bookmark string, objectType reflect.Type) (*PaginatedQueryResult, error) {
	paginatedQueryResult, err := getQueryResultForQueryStringWithPagination(ctx, queryString, int32(pageSize), bookmark, objectType)
	if err != nil {
		return nil, err
	}

	return paginatedQueryResult, nil
}

func (t *SmartContract) GetAssetsByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {
	queryString := fmt.Sprintf(`{"selector":{"SerialNumber":{"$gte":"%s","$lte":"%s"}}}`, startKey, endKey)
	return queryItemsWithPagination(ctx, queryString, pageSize, bookmark, reflect.TypeOf(Asset{}))
}

func (t *SmartContract) QueryAssetsWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {
	return queryItemsWithPagination(ctx, queryString, pageSize, bookmark, reflect.TypeOf(Asset{}))
}

// GetAssetHistory returns the chain of custody for an asset since issuance.
func (t *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetSerialNumber string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: SerialNumber %v", assetSerialNumber)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetSerialNumber)
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
				SerialNumber: assetSerialNumber,
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
