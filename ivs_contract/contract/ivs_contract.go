package contract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kine23/nchu_ivslab/ivs_contract/model"
	"github.com/kine23/nchu_ivslab/ivs_contract/tools"

)

type IVSContract struct {
	contractapi.Contract
}

// 判斷零件是否存在
func (o *IVSContract) Exists(ctx contractapi.TransactionContextInterface, index string) (bool, error) {
	resByte, err := ctx.GetStub().GetState(index)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return resByte != nil, nil
}

// 寫入新零件
func (o *IVSContract) Insert(ctx contractapi.TransactionContextInterface, pJSON string) error {
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

// 轉移零件
func (s *IVSContract) TransferProject(ctx contractapi.TransactionContextInterface, projectID string, newOwnerOrganization string) error {
	// 獲取項目
	projectJSON := fmt.Sprintf(`{"ID":"%s"}`, projectID)
	project, err := s.SelectByIndex(ctx, projectJSON)
	if err != nil {
		return fmt.Errorf("failed to get project: %v", err)
	}

	project.Organization = newOwnerOrganization

	projectJSON, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("failed to marshal project: %v", err)
	}
	err = ctx.GetStub().PutState(projectID, projectJSON)
	if err != nil {
		return fmt.Errorf("failed to put project to world state: %v", err)
	}
	
	fmt.Println("project transferred successfully")
	return nil
}

// 更新零件訊息
func (o *IVSContract) Update(ctx contractapi.TransactionContextInterface, pJSON string) error {
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
func (o *IVSContract) Delete(ctx contractapi.TransactionContextInterface, pJSON string) error {
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
func (o *IVSContract) SelectByIndex(ctx contractapi.TransactionContextInterface, pJSON string) (*model.Project, error) {
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
func (o *IVSContract) SelectAll(ctx contractapi.TransactionContextInterface) ([]*model.Project, error) {
	queryString := fmt.Sprintf(`{"selector":{"table":"project"}}`)
	fmt.Println("select string: ", queryString)
	return tools.SelectByQueryString[model.Project](ctx, queryString)
}

// 依索引查詢數據
func (o *IVSContract) SelectBySome(ctx contractapi.TransactionContextInterface, key, value string) ([]*model.Project, error) {
	queryString := fmt.Sprintf(`{"selector":{"%s":"%s", "table":"project"}}`, key, value)
	return tools.SelectByQueryString[model.Project](ctx, queryString)
}

// 多頁查詢所有數據
func (o *IVSContract) SelectAllWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"table":"project"}}`)
	fmt.Println("select string: ", queryString, "pageSize: ", pageSize, "bookmark", bookmark)
	res, err := tools.SelectByQueryStringWithPagination[model.Project](ctx, queryString, pageSize, bookmark)
	resb, _ := json.Marshal(res)
	fmt.Printf("select result: %v", res)
	return string(resb), err
}

// 按關鍵字多頁查詢

func (o *IVSContract) SelectBySomeWithPagination(ctx contractapi.TransactionContextInterface, key, value string, pageSize int32, bookmark string) (string, error) {
	queryString := fmt.Sprintf(`{"selector":{"%s":"%s","table":"project"}}`, key, value)
	fmt.Println("select string: ", queryString, "pageSize: ", pageSize, "bookmark", bookmark)
	res, err := tools.SelectByQueryStringWithPagination[model.Project](ctx, queryString, pageSize, bookmark)
	resb, _ := json.Marshal(res)
	fmt.Printf("select result: %v", res)
	return string(resb), err
}

// 按索引查詢數據歷史
func (o *IVSContract) SelectHistoryByIndex(ctx contractapi.TransactionContextInterface, pJSON string) (string, error) {
	var tx model.Project
	json.Unmarshal([]byte(pJSON), &tx)
	fmt.Println("select by tx: ", tx)
	res, err := tools.SelectHistoryByIndex[model.Project](ctx, tx.Index())
	resb, _ := json.Marshal(res)
	fmt.Printf("select result: %v", res)
	return string(resb), err
}

// 註冊新帳號
func (s *IVSContract) CreateUser(ctx contractapi.TransactionContextInterface, username string, name string) error {
	exists, err := s.UserExists(ctx, username)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the user %s already exists", username)
	}

	user := model.User{
		Username: username,
		Name:     name,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(username, userJSON)
}

// 讀取指定帳號訊息
func (s *IVSContract) ReadUser(ctx contractapi.TransactionContextInterface, username string) (*model.User, error) {
	userJSON, err := ctx.GetStub().GetState(username)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userJSON == nil {
		return nil, fmt.Errorf("the user %s does not exist", username)
	}

	var user model.User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 更新帳號訊息

func (s *IVSContract) UpdateUser(ctx contractapi.TransactionContextInterface, username string, name string) error {
	exists, err := s.UserExists(ctx, username)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the user %s does not exist", username)
	}

	user := model.User{
		Username: username,
		Name:     name,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(username, userJSON)
}

// 刪除指定ID帳號
func (s *IVSContract) DeleteUser(ctx contractapi.TransactionContextInterface, username string) error {
	exists, err := s.UserExists(ctx, username)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the user %s does not exist", username)
	}
	return ctx.GetStub().DelState(username)
}

// 判斷帳號是否存在
func (s *IVSContract) UserExists(ctx contractapi.TransactionContextInterface, username string) (bool, error) {
	userJSON, err := ctx.GetStub().GetState(username)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return userJSON != nil, nil
}

// 讀取所有帳號訊息
func (s *IVSContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*model.User, error) {
	// GetStateByRange 查詢參數兩個空字元就是查詢所有
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var users []*model.User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user model.User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// 初始化專案數據
func (s *IVSContract) InitProjects(ctx contractapi.TransactionContextInterface) error {
	projects := []model.Project{
		{
			ID:           "IVSLAB23FA05A1ADC01",
			Item:         "智慧影像監控產品追溯系統",
			Developer:    "SFChen",
			Organization: "Lab-IVSOrgs",
			Category:     "Blockchain",
			Describes:    "本研究旨在實現基於Hyperledger Fabric的區塊鏈溯源平台，通過對產品供應鏈進行可靠和透明的追踪，實現企業ESG目標。",
		},
	}
	for _, tx := range projects {
		txJSON, err := json.Marshal(tx)
		if err != nil {
			return err
		}
		err = s.Insert(ctx, string(txJSON))
		if err != nil {
			return err
		}
	}
	return nil
}

// 初始化使用者數據
func (s *IVSContract) InitUsers(ctx contractapi.TransactionContextInterface) error {
	users := []model.User{
		{
			Username: "SFChen",
			Name:     "SFChen",
		},
	}
	for _, user := range users {
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(user.Username, userJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}
	return nil
}

// 初始化智慧合約數據
func (s *IVSContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	err := s.InitProjects(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize projects: %v", err)
	}

	err = s.InitUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize users: %v", err)
	}

	return nil
}
