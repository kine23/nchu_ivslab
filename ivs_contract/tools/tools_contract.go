package tools

import (
	"bytes"
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kine23/nchu_ivslab/ivs_contract/model"
	"time"
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

	var results []T
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
		results = append(results, item)
	}

	paginatedQueryResult := &model.PaginatedQueryResult[T]{
		Records:             results,
		FetchedRecordsCount: metadata.FetchedRecordsCount,
		Bookmark:            metadata.Bookmark,
	}
	return paginatedQueryResult, nil
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

		// 在這裡添加類型斷言
		asset, ok := item.(model.Asset)
		if !ok {
			return nil, fmt.Errorf("無法將T轉換為model.Asset")
		}

		historyItem := &model.HistoryQueryResult[asset]{
			Record:    asset,
			TxId:      queryResponse.TxId,
			Timestamp: queryResponse.Timestamp.AsTime(),
			IsDelete:  queryResponse.IsDelete,
		}
		results = append(results, historyItem)
	}
	return results, nil
}

// SelectByIndexAndPagination 通用按索引分頁查詢功能
func SelectByIndexAndPagination[T interface{}](ctx contractapi.TransactionContextInterface, table, key, value string, pageSize int32, bookmark string) (*model.PaginatedQueryResult[T], error) {
	queryString := buildQueryString(table, key, value)
	return SelectByQueryStringWithPagination[T](ctx, queryString, pageSize, bookmark)
}
