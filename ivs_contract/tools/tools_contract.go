package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/kine23/nchu_ivslab/ivs_contract/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// buildQueryString 用於構建用於查詢的JSON字符串
func buildQueryString(table, key, value string) string {
	var buffer bytes.Buffer
	buffer.WriteString(`{"selector":{"table":"`)
	buffer.WriteString(table)
	buffer.WriteString(`"`)

	if key != "" && value != "" {
		buffer.WriteString(`,"`)
		buffer.WriteString(key)
		buffer.WriteString(`":"`)
		buffer.WriteString(value)
		buffer.WriteString(`"`)
	}

	buffer.WriteString("}}")
	return buffer.String()
}

// SelectByQueryString 通用查詢功能
func SelectByQueryString[T interface{}](ctx contractapi.TransactionContextInterface, queryString string) ([]*T, error) {
    resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var results []*T
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }
        var item T
        err = json.Unmarshal(queryResponse.Value, &item)
        if err != nil {
            return nil, err
        }
        results = append(results, &item)
    }

    return results, nil
}

// SelectByQueryStringWithPagination 通用分頁查詢功能
func SelectByQueryStringWithPagination[T interface{}](ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*model.PaginatedQueryResult[T], error) {
    resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var results []*T
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }
        var item T
        err = json.Unmarshal(queryResponse.Value, &item)
        if err != nil {
            return nil, err
        }
        results = append(results, &item)
    }

	type PaginatedQueryResult[T any] struct {
		Records             []*T             `json:"records"`
		FetchedRecordsCount int32            `json:"fetched_records_count"`
		Bookmark            string           `json:"bookmark"`
	}
    return &paginatedQueryResult, nil
}

// SelectHistoryByIndex 通用歷史查詢功能
func SelectHistoryByIndex[T interface{}](ctx contractapi.TransactionContextInterface, index string) ([]*model.HistoryQueryResult[T], error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(index)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []*model.HistoryQueryResult[T]
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var item T
		err = json.Unmarshal(queryResponse.Value, &item)
		if err != nil {
			return nil, err
		}
	type HistoryQueryResult[T any] struct {
		Record    T                   `json:"record"`
		TxId      string              `json:"tx_id"`
		Timestamp *timestamppb.Timestamp `json:"timestamp"`
		IsDelete  bool                `json:"is_delete"`
	}
		results = append(results, &historyItem)
	}
	return results, nil
}

// SelectByIndexAndPagination 通用按索引分頁查詢功能
func SelectByIndexAndPagination[T interface{}](ctx contractapi.TransactionContextInterface, table, key, value string, pageSize int32, bookmark string) (*model.PaginatedQueryResult[T], error) {
	queryString := buildQueryString(table, key, value)
	return SelectByQueryStringWithPagination[T](ctx, queryString, pageSize, bookmark)
}

// 交易創建後的所有變化
func SelectHistoryByIndex[T interface{}](ctx contractapi.TransactionContextInterface, index string) ([]model.HistoryQueryResult[T], error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(index)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []model.HistoryQueryResult[T]
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var tx T
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &tx)
			if err != nil {
				return nil, err
			}
		}
		timestamp := time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)) // 將時間戳轉換為 time.Time 類型
		record := model.HistoryQueryResult[T]{
			TxId:      response.TxId,
			Record:    tx,
			IsDelete:  response.IsDelete,
			Timestamp: timestamp,
		}
		records = append(records, record)
	}
	return records, nil
}


