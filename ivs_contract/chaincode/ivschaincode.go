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

const index = "madein~serialnumber"

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset Project項目列表.
type Asset struct {
	DocType             string `json:"docType"`             	// DocType is used to distinguish the various types of objects in state database
	ID                	string `json:"ID"`                		// 項目唯一ID
	MadeBy        		string `json:"MadeBy"`        			// 品牌商
	MadeIn 				string `json:"MadeIn"` 					// 組裝地點
	SerialNumber        string `json:"SerialNumber"`        	// 產品序號
	SecurityChip        Part   `json:"SecurityChip"`        	// 安全晶片組織
	NetworkChip         Part   `json:"NetworkChip"`         	// 網路晶片組織
	CMOSChip            Part   `json:"CMOSChip"`            	// CMOS晶片組織
	VideoCodecChip		Part   `json:"VideoCodecChip"`      	// VideoCodec晶片組織
	ProductionDate      string `json:"ProductionDate"`      	// 產品生產日期
	Updated				string `json:"Updated"`      			// 產品更新日期
}

// Part Project項目列表.
type Part struct {
	DocType             string `json:"docType"`             	// DocType is used to distinguish the various types of objects in state database
	PID					string `json:"PID"`						// 零件唯ID
	Manufacturer        string `json:"Manufacturer"`        	// 製造商
	ManufactureLocation string `json:"ManufactureLocation"` 	// 製造地點
	PartName            string `json:"PartName"`            	// 零件名稱
	PartNumber          string `json:"PartNumber"`          	// 零件批號
	Organization        string `json:"Organization"`        	// 組織
	ManufactureDate     string `json:"ManufactureDate"`     	// 零件製造日期
	TransferDate        string `json:"TransferDate"`        	// 零件交易日期
}

// InitLedger adds a base set of assets to the ledger
func (t *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	parts := []Part{
		{PID: "IVSLAB-S23FA0001", Manufacturer: "Security.Co", ManufactureLocation: "Taiwan", PartName: "SecurityChip-v1", PartNumber: "SPN3R1C00AA1", Organization: "Security-Org"},
		{PID: "IVSLAB-N23FA0001", Manufacturer: "Network.Co", ManufactureLocation: "Taiwan", PartName: "NetworkChip-v1", PartNumber: "NPN3R1C00AA1", Organization: "Network-Org"},
		{PID: "IVSLAB-C23FA0001", Manufacturer: "CMOS.Co", ManufactureLocation: "USA", PartName: "CMOSChip-v1", PartNumber: "CPN3R1C00AA1", Organization: "CMOS-Org"},
		{PID: "IVSLAB-V23FA0001", Manufacturer: "VideoCodec.Co", ManufactureLocation: "USA", PartName: "VideoCodecChip-v1", PartNumber: "VPN3R1C00AA1", Organization: "VideoCodec-Org"},
		{PID: "IVSLAB-S23FA0002", Manufacturer: "Security.Co", ManufactureLocation: "Taiwan", PartName: "SecurityChip-v1", PartNumber: "SPN3R1C00AA2", Organization: "Security-Org"},
		{PID: "IVSLAB-N23FA0002", Manufacturer: "Network.Co", ManufactureLocation: "Taiwan", PartName: "NetworkChip-v1", PartNumber: "NPN3R1C00AA2", Organization: "Network-Org"},
		{PID: "IVSLAB-C23FA0002", Manufacturer: "CMOS.Co", ManufactureLocation: "USA", PartName: "CMOSChip-v1", PartNumber: "CPN3R1C00AA2", Organization: "CMOS-Org"},
		{PID: "IVSLAB-V23FA0002", Manufacturer: "VideoCodec.Co", ManufactureLocation: "USA", PartName: "VideoCodecChip-v1", PartNumber: "VPN3R1C00AA2", Organization: "VideoCodec-Org"},
		{PID: "IVSLAB-S23FA0003", Manufacturer: "Security.Co", ManufactureLocation: "Taiwan", PartName: "SecurityChip-v1", PartNumber: "SPN3R1C00AA3", Organization: "Security-Org"},
		{PID: "IVSLAB-N23FA0003", Manufacturer: "Network.Co", ManufactureLocation: "Taiwan", PartName: "NetworkChip-v1", PartNumber: "NPN3R1C00AA3", Organization: "Network-Org"},
		{PID: "IVSLAB-C23FA0003", Manufacturer: "CMOS.Co", ManufactureLocation: "USA", PartName: "CMOSChip-v1", PartNumber: "CPN3R1C00AA3", Organization: "CMOS-Org"},
		{PID: "IVSLAB-V23FA0003", Manufacturer: "VideoCodec.Co", ManufactureLocation: "USA", PartName: "VideoCodecChip-v1", PartNumber: "VPN3R1C00AA3", Organization: "VideoCodec-Org"},		
	}

	for _, part := range parts {
		err := t.CreatePart(ctx, part.PID, part.Manufacturer, part.ManufactureLocation, part.PartName, part.PartNumber, part.Organization)
		if err != nil {
			return err
		}
	}

	return nil
}

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

// CreatePart initializes a new part in the ledger
func (t *SmartContract) CreatePart(ctx contractapi.TransactionContextInterface, partID, manufacturer string, manufacturelocation string, partname string, partnumber string, organization string) error {
	exists, err := t.PartExists(ctx, partID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the part %s already exists", partID)
	}

	part := &Part{
		DocType:             "part",
		PID:                 partID,
		Manufacturer:        manufacturer,
		ManufactureLocation: manufacturelocation,
		PartName:            partname,
		PartNumber:          partnumber,
		Organization:        organization,
		ManufactureDate:     time.Now().Format("2006-01-02"),
	}
	partBytes, err := json.Marshal(part)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(partID, partBytes)
	if err != nil {
		return err
	}
	ivsIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{part.Organization, part.PID})
	if err != nil {
		return err
	}
	value := []byte{0x00}
	return ctx.GetStub().PutState(ivsIndexKey, value)	
}

// GetPart retrieves a part from the ledger by its ID
func (t *SmartContract) GetPart(ctx contractapi.TransactionContextInterface, partID string) (*Part, error) {
	partBytes, err := ctx.GetStub().GetState(partID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if partBytes == nil {
		return nil, fmt.Errorf("the part %s does not exist", partID)
	}

	part := new(Part)
	err = json.Unmarshal(partBytes, part)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal part: %v", err)
	}

	return part, nil
}

// CreateAsset initializes a new asset in the ledger
func (t *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetID string, madeby string, madein string, serialnumber string, securitychipID string, networkchipID string, cmoschipID string, videocodecchipID string) error {
	exists, err := t.AssetExists(ctx, assetID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", assetID)
	}
	var securitychip *Part
	var networkchip *Part
	var cmoschip *Part
	var videocodecchip *Part
	// Get the Part instances from the ledger state
	securitychip, err = t.GetPart(ctx, securitychipID)
	if err != nil {
		return err
	}
	networkchip, err = t.GetPart(ctx, networkchipID)
	if err != nil {
		return err
	}
	cmoschip, err = t.GetPart(ctx, cmoschipID)
	if err != nil {
		return err
	}
	videocodecchip, err = t.GetPart(ctx, videocodecchipID)
	if err != nil {
		return err
	}	
	// Ensure all parts belong to 'Brand-Org'
	parts := []*Part{securitychip, networkchip, cmoschip, videocodecchip}
	for _, part := range parts {
		if part.Organization != "Brand-Org" {
			return fmt.Errorf("part %s does not belong to Brand-Org", part.PID)
		}
	}
	asset := Asset{
		DocType:        "asset",
		ID:             assetID,
		MadeBy:         madeby,
		MadeIn:         madein,
		SerialNumber:   serialnumber,
		SecurityChip:   *securitychip,
		NetworkChip:    *networkchip,
		CMOSChip:       *cmoschip,
		VideoCodecChip: *videocodecchip,
		ProductionDate: time.Now().Format("2006-01-02"),
	}
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(assetID, assetBytes)
	if err != nil {
		return err
	}
	ivsIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.MadeBy, asset.ID})
	if err != nil {
		return err
	}
	value := []byte{0x00}
	return ctx.GetStub().PutState(ivsIndexKey, value)	
}

// CreateAsset initializes a new asset in the ledger
//func (t *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetID string, madeby string, madein string, serialnumber string, securitychip Part, networkchip Part, cmoschip Part, videocodecchip Part) error {
//	exists, err := t.AssetExists(ctx, assetID)
//	if err != nil {
//		return err
//	}
//	if exists {
//		return fmt.Errorf("the asset %s already exists", assetID)
//	}
//	asset := Asset{
//		DocType:        "asset",
//		ID:              assetID,
//		MadeBy:          madeby,
//		MadeIn:          madein,
//		SerialNumber:    serialnumber,
//		SecurityChip:    securitychip,
//		NetworkChip:     networkchip,
//		CMOSChip:        cmoschip,
//		VideoCodecChip:  videocodecchip,
//		ProductionDate:  time.Now().Format("2006-01-02"),
//	}
//	assetBytes, err := json.Marshal(asset)
//	if err != nil {
//		return err
//	}
//
//	err = ctx.GetStub().PutState(assetID, assetBytes)
//	if err != nil {
//		return err
//	}
//	ivsIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.MadeBy, asset.ID})
//	if err != nil {
//		return err
//	}
//	value := []byte{0x00}
//	return ctx.GetStub().PutState(ivsIndexKey, value)	
//}

// ReadPart retrieves an part from the ledger
func (t *SmartContract) ReadPart(ctx contractapi.TransactionContextInterface, partID string) (*Part, error) {
	partBytes, err := ctx.GetStub().GetState(partID)
	if err != nil {
		return nil, fmt.Errorf("failed to get part %s: %v", partID, err)
	}
	if partBytes == nil {
		return nil, fmt.Errorf("part %s does not exist", partID)
	}

	var part Part
	err = json.Unmarshal(partBytes, &part)
	if err != nil {
		return nil, err
	}

	return &part, nil
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
func (t *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, assetID string, madeby string, madein string, serialnumber string, securitychip Part, networkchip Part, cmoschip Part, videocodecchip Part) error {
	exists, err := t.AssetExists(ctx, assetID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", assetID)
	}
	var securitychip *Part
	var networkchip *Part
	var cmoschip *Part
	var videocodecchip *Part
	// Get the Part instances from the ledger state
	securitychip, err = t.GetPart(ctx, securitychipID)
	if err != nil {
		return err
	}
	networkchip, err = t.GetPart(ctx, networkchipID)
	if err != nil {
		return err
	}
	cmoschip, err = t.GetPart(ctx, cmoschipID)
	if err != nil {
		return err
	}
	videocodecchip, err = t.GetPart(ctx, videocodecchipID)
	if err != nil {
		return err
	}
	// Ensure all parts belong to 'Brand-Org'
	parts := []*Part{securitychip, networkchip, cmoschip, videocodecchip}
	for _, part := range parts {
		if part.Organization != "Brand-Org" {
			return fmt.Errorf("part %s does not belong to Brand-Org", part.PID)
		}
	}
	// overwriting original asset with new asset
	asset := &Asset{
		DocType:        "asset",
		ID:              assetID,
		MadeBy:          madeby,
		MadeIn:          madein,
		SerialNumber:    serialnumber,
		SecurityChip:    *securitychip,
		NetworkChip:     *networkchip,
		CMOSChip:        *cmoschip,
		VideoCodecChip:  *videoCodecchip,
		ProductionDate:  time.Now().Format("2006-01-02"),
		Updated:  		 time.Now().Format("2006-01-02"),
	}
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(assetID, assetBytes)
	if err != nil {
		return err
	}
	ivsIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.MadeBy, asset.ID})
	if err != nil {
		return err
	}
	value := []byte{0x00}
	return ctx.GetStub().PutState(ivsIndexKey, value)	
}

// DeletePart removes an part key-value pair from the ledger
func (t *SmartContract) DeletePart(ctx contractapi.TransactionContextInterface, partID string) error {
	part, err := t.ReadPart(ctx, partID)
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(partID)
	if err != nil {
		return fmt.Errorf("failed to delete part %s: %v", partID, err)
	}

	ivsIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{part.Manufacturer, part.PID})
	if err != nil {
		return err
	}

	// Delete index entry
	return ctx.GetStub().DelState(ivsIndexKey)
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

	ivsIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.MadeBy, asset.ID})
	if err != nil {
		return err
	}

	// Delete index entry
	return ctx.GetStub().DelState(ivsIndexKey)
}

// TransferPart updates the Organization and TransferDate field of part with given id in world state, and returns the old Organization.
func (t *SmartContract) TransferPart(ctx contractapi.TransactionContextInterface, partID string, newOrganization string) (string, error) {
	part, err := t.ReadPart(ctx, partID)
	if err != nil {
		return "", fmt.Errorf("failed to read part: %v", err)
	}

	oldOrganization := part.Organization
	part.Organization = newOrganization

	// Set the transfer date to the current system date
	part.TransferDate = time.Now().Format("2006-01-02")

	partBytes, err := json.Marshal(part)
	if err != nil {
		return "", fmt.Errorf("failed to marshal part: %v", err)
	}

	err = ctx.GetStub().PutState(partID, partBytes)
	if err != nil {
		return "", fmt.Errorf("failed to write part: %v", err)
	}

	return oldOrganization, nil
}

// constructQueryResponseFromIteratorPart constructs a slice of parts from the resultsIterator
func constructQueryResponseFromIteratorPart(resultsIterator shim.StateQueryIteratorInterface) ([]*Part, error) {
	var parts []*Part
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var part Part
		err = json.Unmarshal(queryResult.Value, &part)
		if err != nil {
			return nil, err
		}
		parts = append(parts, &part)
	}

	// return an empty slice instead of nil if there are no parts
	if len(parts) == 0 {
		return []*Part{}, nil
	}

	return parts, nil
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

// GetAllParts returns all parts found in world state
func (t *SmartContract) GetAllParts(ctx contractapi.TransactionContextInterface) ([]*Part, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIteratorPart(resultsIterator)
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

func (t *SmartContract) GetPartsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Part, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIteratorPart(resultsIterator)
}

func (t *SmartContract) GetAssetsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SmartContract) QueryAssetsByOrganization(ctx contractapi.TransactionContextInterface, madeby string) ([]*Asset, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"asset","madeby":"%s"}}`, madeby)
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

// PartExists returns true when part with given ID exists in world state
func (t *SmartContract) PartExists(ctx contractapi.TransactionContextInterface, partID string) (bool, error) {
	partBytes, err := ctx.GetStub().GetState(partID)
	if err != nil {
		return false, fmt.Errorf("failed to read part %s from world state. %v", partID, err)
	}

	return partBytes != nil, nil
}

// AssetExists returns true when asset with given ID exists in world state
func (t *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return false, fmt.Errorf("failed to read asset %s from world state. %v", assetID, err)
	}

	return assetBytes != nil, nil
}
