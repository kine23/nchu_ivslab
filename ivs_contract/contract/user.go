package contract

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kine23/nchu_ivslab/ivs_contract/model"
)

type UserContract struct {
	contractapi.Contract
}

// 註冊新帳號
func (s *UserContract) CreateUser(ctx contractapi.TransactionContextInterface, username string, name string) error {
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
func (s *UserContract) ReadUser(ctx contractapi.TransactionContextInterface, username string) (*model.User, error) {
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
func (s *UserContract) UpdateUser(ctx contractapi.TransactionContextInterface, username string, name string) error {
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
func (s *UserContract) DeleteUser(ctx contractapi.TransactionContextInterface, username string) error {
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
func (s *UserContract) UserExists(ctx contractapi.TransactionContextInterface, username string) (bool, error) {
	userJSON, err := ctx.GetStub().GetState(username)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return userJSON != nil, nil
}

// 讀取所有帳號訊息
func (s *UserContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*model.User, error) {
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

func (o *UserContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	txs := []model.User{
		{
			Username: "SF.Chen",
			Name:     "SF.Chen",
		},
	}
	for _, tx := range txs {
		txJSON, err := json.Marshal(tx)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(tx.Username, txJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}
