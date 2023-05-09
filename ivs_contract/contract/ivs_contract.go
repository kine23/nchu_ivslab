package contract

import (
	"encoding/json"
	"fmt"
	
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kine23/nchu_ivslab/ivs_contract/model"
	"github.com/kine23/nchu_ivslab/ivs_contract/tools"
)

type ProjectContract struct {
	contractapi.Contract
}

// 判斷零件是否存在
func (o *ProjectContract) Exists(ctx contractapi.TransactionContextInterface, index string) (bool, error) {
	resByte, err := ctx.GetStub().GetState(index)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return resByte != nil, nil
}

// 寫入新零件
func (o *ProjectContract) Insert(ctx contractapi.TransactionContextInterface, pJSON string) error {
	var tx model.Project
	json.Unmarshal([]byte(pJSON), &tx)
	exists, err := o.Exists(ctx, tx.Index())
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the data %s already exists", tx.Index())
	}
	txb, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(tx.Index(), txb)
	indexKey, err := ctx.GetStub().CreateCompositeKey(tx.IndexKey(), tx.IndexAttr())
	if err != nil {
		return err
	}
	value := []byte{0x00}
	fmt.Println("create success: ", tx)
	return ctx.GetStub().PutState(indexKey, value)
}

// 更新零件訊息
func (o *ProjectContract) Update(ctx contractapi.TransactionContextInterface, pJSON string) error {
	var tx model.Project
	json.Unmarshal([]byte(pJSON), &tx)

	otx, err := o.SelectByIndex(ctx, pJSON)
	if err != nil {
		return err
	}
	if otx == nil {
		return fmt.Errorf("the tx %s does not exist", tx.Index())
	}

	// 刪除舊索引
	indexKey, err := ctx.GetStub().CreateCompositeKey(otx.IndexKey(), otx.IndexAttr())
	if err != nil {
		return err
	}
	ctx.GetStub().DelState(indexKey)

	txb, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(tx.Index(), txb)

	if indexKey, err = ctx.GetStub().CreateCompositeKey(tx.IndexKey(), tx.IndexAttr()); err != nil {
		return err
	}
	value := []byte{0x00}
	return ctx.GetStub().PutState(indexKey, value)
}

// 刪除零件
func (o *ProjectContract) Delete(ctx contractapi.TransactionContextInterface, pJSON string) error {
	var tx model.Project
	json.Unmarshal([]byte(pJSON), &tx)

	anstx, err := o.SelectByIndex(ctx, pJSON)
	if err != nil {
		return err
	}
	if anstx == nil {
		return fmt.Errorf("the tx %s does not exist", tx.Index())
	}
	err = ctx.GetStub().DelState(anstx.Index())
	if err != nil {
		return fmt.Errorf("failed to delete transaction %s: %v", anstx.Index(), err)
	}

	indexKey, err := ctx.GetStub().CreateCompositeKey(tx.IndexKey(), tx.IndexAttr())
	if err != nil {
		return err
	}

	// Delete index entry
	return ctx.GetStub().DelState(indexKey)
}

// 查詢指定零件紀錄
func (o *ProjectContract) SelectByIndex(ctx contractapi.TransactionContextInterface, pJSON string) (*model.Project, error) {
	tx := model.Project{}
	json.Unmarshal([]byte(pJSON), &tx)
	queryString := fmt.Sprintf(`{"selector":{"ID":"%s", "table":"project"}}`, tx.ID)
	fmt.Println("select string: ", queryString)
	res, err := tools.SelectByQueryString[model.Project](ctx, queryString)
	if len(res) == 0 {
		return nil, err
	}
	return res[0], err
}

// 查詢所有紀錄
func (o *ProjectContract) SelectAll(ctx contractapi.TransactionContextInterface) ([]*model.Project, error) {
	queryString := fmt.Sprintf(`{"selector":{"table":"project"}}`)
	fmt.Println("select string: ", queryString)
	return tools.SelectByQueryString[model.Project](ctx, queryString)
}

// 依索引查詢數據
func (o *ProjectContract) SelectBySome(ctx contractapi.TransactionContextInterface, key, value string) ([]*model.Project, error) {
	queryString := fmt.Sprintf(`{"selector":{"%s":"%s", "table":"project"}}`, key, value)
	return tools.SelectByQueryString[model.Project](ctx, queryString)
}

// 多頁查詢所有數據
func (o *ProjectContract) SelectAllWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"table":"project"}}`)
	fmt.Println("select string: ", queryString, "pageSize: ", pageSize, "bookmark", bookmark)
	res, err := tools.SelectByQueryStringWithPagination[model.Project](ctx, queryString, pageSize, bookmark)
	resb, _ := json.Marshal(res)
	fmt.Printf("select result: %v", res)
	return string(resb), err
}

// 按關鍵字多頁查詢
func (o *ProjectContract) SelectBySomeWithPagination(ctx contractapi.TransactionContextInterface, key, value string, pageSize int32, bookmark string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"%s":"%s","table":"project"}}`, key, value)
	fmt.Println("select string: ", queryString, "pageSize: ", pageSize, "bookmark", bookmark)
	res, err := tools.SelectByQueryStringWithPagination[model.Project](ctx, queryString, pageSize, bookmark)
	resb, _ := json.Marshal(res)
	fmt.Printf("select result: %v", res)
	return string(resb), err
}

// 按索引查詢數據歷史
func (o *ProjectContract) SelectHistoryByIndex(ctx contractapi.TransactionContextInterface, pJSON string) (string, error) {
	var tx model.Project
	json.Unmarshal([]byte(pJSON), &tx)
	fmt.Println("select by tx: ", tx)
	res, err := tools.SelectHistoryByIndex[model.Project](ctx, tx.Index())
	resb, _ := json.Marshal(res)
	fmt.Printf("select result: %v", res)
	return string(resb), err
}

// 初始化智慧合約數據
func (s *ProjectContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	projects := []model.Project{
		{ID: "IVSLAB23FA05A1ADC01",
			Name:         "智慧影像監控產品追溯系統",
			Developer:    "PoC",
			Organization: "IVS-Orgs",
			Category:     "blockchain",
			Describes:    "本研究旨在實現基於Hyperledger Fabric的區塊鏈溯源平台，通過對產品供應鏈進行可靠和透明的追踪，實現企業ESG目標。",
		},
	}
	for _, tx := range projects {
		txJsonByte, err := json.Marshal(tx)
		if err != nil {
			return err
		}
		err = s.Insert(ctx, string(txJsonByte))
		if err != nil {
			return err
		}
	}
	return nil
}
