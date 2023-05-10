package ivscontract

import (
	"encoding/json"
	"fmt"
	
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kine23/nchu_ivslab/ivs_contract/model"
)

type UserContract struct {
	contractapi.Contract
}

// CreateUser 註冊新帳號
func (u *UserContract) CreateUser(ctx contractapi.TransactionContextInterface, userJSON string) (string, error) {
	var user model.User
	err := json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(user.GetKey(), []byte(userJSON))
	if err != nil {
		return "", fmt.Errorf("failed to put user to world state. %v", err)
	}

	return user.GetKey(), nil
}

// ReadUser 讀取指定帳號訊息
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

// UpdateUser 更新帳號訊息
func (u *UserContract) UpdateUser(ctx contractapi.TransactionContextInterface, userJSON string) (string, error) {
	var user model.User
	err := json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(user.GetKey(), []byte(userJSON))
	if err != nil {
		return "", err
	}

	return user.GetKey(), nil
}

//  DeleteUser 刪除帳號
func (u *UserContract) DeleteUser(ctx contractapi.TransactionContextInterface, userID string) error {
	err := ctx.GetStub().DelState(userID)
	return err
}

// UserExists 判斷帳號是否存在
func (s *UserContract) UserExists(ctx contractapi.TransactionContextInterface, username string) (bool, error) {
	userJSON, err := ctx.GetStub().GetState(username)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return userJSON != nil, nil
}

// GetUser 讀取帳號訊息
func (u *UserContract) GetUser(ctx contractapi.TransactionContextInterface, userID string) (*model.User, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil, err
	}
	if userBytes == nil {
		return nil, fmt.Errorf("用戶 %s 不存在", userID)
	}

	var user model.User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// InitLedger 初始化智能合約數據，只在智能合約實例化時使用
func (o *UserContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	txs := []model.User{
		{
			Username: "ivspoc",
			Name:     "PoC",
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
