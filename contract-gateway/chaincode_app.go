/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	mspID="brandMSP"														//所屬組織的MSPID
	cryptoPath= "/root/MyLab_IVS/organizations/brand.ivsorg.net"			// 中間變量
	certPath= cryptoPath + "/registers_users/user1/msp/signcerts/cert.pem"	// client數位簽章
	keyPath= cryptoPath + "/registers_users/user1/msp/keystore/"			// client私鑰路徑
	tlsCertPath= cryptoPath + "/alliance/tls-ca-cert.pem"					// client tls證書
	peerEndpoint= "peer1.brand.ivsorg.net:7151"								// peer節點地址
	gatewayPeer= "peer1.brand.ivsorg.net"									// peer節點名稱
)

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)
var partId = fmt.Sprintf("part%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	chaincodeName := "ivs_basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "ivschannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	initLedger(contract)
//	transferPartAsync(contract)
//	transferPartsByOrganizationAsync(contract)
//	createPart(contract)
	getAllParts(contract)
//	createAsset(contract)
//	updateAsset(contract)
//	readAssetByID(contract)
//	readPartByID(contract)
//	queryAssets(contract)
//	queryAssetsBySerialNumber(contract)
	getAssetSerialNumberHistory(contract)
//	getAssetHistory(contract)	
//	exampleErrorHandling(contract)
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)
	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}
	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}
	return sign
}

// This type of transaction would typically only be run once by an application the first time it was started after its
// initial deployment. A new version of the chaincode deployed later would likely not need to run an "init" function.
func initLedger(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")
	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		fmt.Printf("failed to submit transaction: %s\n", err)
		return
	}
	fmt.Printf("*** Transaction committed successfully\n")
}
// Evaluate a transaction to query ledger state.
func getAllParts(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAllParts, function returns all the current parts on the ledger")
	evaluateResult, err := contract.EvaluateTransaction("GetAllParts")
	if err != nil {
		fmt.Printf("failed to submit transaction: %s\n", err)
		return
	}
	result := formatJSON(evaluateResult)
	fmt.Printf("*** Result:%s\n", result)
}

func createPart(contract *client.Contract) {
    partID := "IVSLAB-N23FA0004"
    fmt.Printf("\n--> Submit Transaction: CreatePart, creates new part with PID, Manufacturer, ManufactureLocation, PartName, PartNumber, Organization\n")
    _, err := contract.SubmitTransaction("CreatePart", partID, "Network.Co", "Taiwan", "NetworkChip-v1", "NPN3R1C00AA4", "Network-Org")
    if err != nil {
        fmt.Printf("failed to submit transaction: %s\n", err)
        return
    }
    fmt.Printf("*** Transaction committed Part %s created successfully\n", partID)
}

//func createPart(contract *client.Contract) {
//	fmt.Printf("\n--> Submit Transaction: CreatePart, creates new part with PID, Manufacturer, ManufactureLocation, PartName, PartNumber, Organization\n")
//	_, err := contract.SubmitTransaction("CreatePart", "IVSLAB-N23FA0004", "Network.Co", "Taiwan", "NetworkChip-v1", "NPN3R1C00AA4", "Network-Org")
//	if err != nil {
//		fmt.Printf("failed to submit transaction: %s\n", err)
//		return
//	}
//	fmt.Printf("*** Transaction committed Part part%d created successfully\n")
//}

// Submit transaction asynchronously, blocking until the transaction has been sent to the orderer, and allowing
// this thread to process the chaincode response (e.g. update a UI) without waiting for the commit notification
func transferPartAsync(contract *client.Contract) {
    partID := "IVSLAB-N23FA0001"
    fmt.Printf("\n--> Async Submit Transaction: TransferPart, changes existing part Organization and TransferDate")
    submitResult, commit, err := contract.SubmitAsync("TransferPart", client.WithArguments(partID, "Brand-Org"))
    if err != nil {
        panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
    }

    fmt.Printf("\n*** Successfully submitted transaction to transfer %s ownership from %s to Brand.Co. \n", partID, string(submitResult))
    fmt.Println("*** Waiting for transaction commit.")

    if commitStatus, err := commit.Status(); err != nil {
        panic(fmt.Errorf("failed to get commit status: %w", err))
    } else if !commitStatus.Successful {
        panic(fmt.Errorf("transaction %s failed to commit with status: Brand.Co", commitStatus.TransactionID, int32(commitStatus.Code)))
    }
    fmt.Printf("*** Transaction committed successfully\n")
}

//func transferPartAsync(contract *client.Contract) {
//	fmt.Printf("\n--> Async Submit Transaction: TransferPart, updates existing part Organization and TransferDate")
//	submitResult, commit, err := contract.SubmitAsync("TransferPart", client.WithArguments("IVSLAB-N23FA0001", "Brand-Org"))
//	if err != nil {
//		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
//	}
//	fmt.Printf("\n*** Successfully submitted transaction to transfer ownership from %s to Brand.Co. \n", string(submitResult))
//	fmt.Println("*** Waiting for transaction commit.")
//	if commitStatus, err := commit.Status(); err != nil {
//		panic(fmt.Errorf("failed to get commit status: %w", err))
//	} else if !commitStatus.Successful {
//		panic(fmt.Errorf("transaction %s failed to commit with status: Brand.Co", commitStatus.TransactionID, int32(commitStatus.Code)))
//	}
//	fmt.Printf("*** Transaction committed successfully\n")
//}

func transferPartsByOrganizationAsync(contract *client.Contract) {
    organization := "CMOS-Org"
    newOrganization := "Brand-Org"

	fmt.Printf("\n--> Async Submit Transaction: TransferPartsByOrganization, transfers all parts from %s to %s\n", organization, newOrganization)
    // Call TransferPartsByOrganization function asynchronously
    _, commit, err := contract.SubmitAsync("TransferPartsByOrganization", client.WithArguments(organization, newOrganization))
    if err != nil {
        panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
    }

    fmt.Printf("\n*** Successfully submitted transaction to transfer all parts from %s to %s. \n", organization, newOrganization)
    fmt.Println("*** Waiting for transaction commit.")
    commitStatus, err := commit.Status()
    if err != nil {
        panic(fmt.Errorf("failed to get commit status: %w", err))
    } else if !commitStatus.Successful {
        panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
    }

    fmt.Printf("*** Transaction committed successfully\n")
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func createAsset(contract *client.Contract) {
	assetID := "IVSLAB-PVC23FG0002"
	fmt.Printf("\n--> Submit Transaction: CreateAsset, creates new asset with ID, MadeBy, MadeIn, SerialNumber, SecurityChip, NetworkChip, CMOSChip, VideoCodecChip\n")
	_, err := contract.SubmitTransaction("CreateAsset", assetID, "Brand.Co", "Taiwan", "IVSPN902300AACDC02", "IVSLAB-S23FA0002", "IVSLAB-N23FA0002", "IVSLAB-C23FA0002", "IVSLAB-V23FA0002")
	if err != nil {
		fmt.Printf("failed to submit transaction: %s\n", err)
		return
	}
	fmt.Printf("*** Transaction committed Asset %s created successfully\n", assetID)
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func updateAsset(contract *client.Contract) {
	assetID := "IVSLAB-PVC23FG0002"
	fmt.Printf("\n--> Submit Transaction: UpdateAsset, update asset with ID, MadeBy, MadeIn, SerialNumber, SecurityChip, NetworkChip, CMOSChip, VideoCodecChip\n")
	_, err := contract.SubmitTransaction("UpdateAsset", assetID, "Brand.Co", "Taiwan", "IVSPN902300AACDC01", "IVSLAB-S23FA0003", "IVSLAB-N23FA0001", "IVSLAB-C23FA0001", "IVSLAB-V23FA0001")
	if err != nil {
		fmt.Printf("failed to submit transaction: %s\n", err)
		return
	}
	fmt.Printf("*** Transaction committed Asset %s updated successfully\n", assetID)
}

// Evaluate a transaction by partID to query ledger state.
func readPartByID(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: ReadPart, function returns part attributes\n")
	evaluateResult, err := contract.EvaluateTransaction("ReadPart", "IVSLAB-S23FA0002")
	if err != nil {
		fmt.Printf("failed to submit transaction: %s\n", err)
		return
	}
	result := formatJSON(evaluateResult)
	fmt.Printf("*** Result:%s\n", result)
}

// Evaluate a transaction by assetID to query ledger state.
func readAssetByID(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")
	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", "IVSLAB-PVC23FG0001")
	if err != nil {
		fmt.Printf("failed to submit transaction: %s\n", err)
		return
	}
	result := formatJSON(evaluateResult)
	fmt.Printf("*** Result:%s\n", result)
}

func queryAssetsBySerialNumber(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: QueryAssetsBySerialNumber, function returns the current assets By SerialNumber on the ledger")
	evaluateResult, err := contract.EvaluateTransaction("QueryAssetsBySerialNumber", "IVSPN902300AACDC01", "IVSPN902300AACDC02")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	// Add a check here for empty result
	if len(evaluateResult) == 0 {
		fmt.Println("*** No assets found for the specified Brand.Co")
		return
	}
	result := formatJSON(evaluateResult)
	fmt.Printf("*** Result:%s\n", result)
}

func queryAssets(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: QueryAssets, function returns the current assets made by Brand.Co on the ledger")
	queryString := "{\"selector\":{\"MadeBy\":\"Brand.Co\"}}"
	evaluateResult, err := contract.EvaluateTransaction("QueryAssets", queryString)
	if err != nil {
		fmt.Printf("failed to submit transaction: %s\n", err)
		return
	}
	result := formatJSON(evaluateResult)
	fmt.Printf("*** Result:%s\n", result)
}

func getAssetHistory(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAssetHistory, function returns all the current assets on the ledger")
	evaluateResult, err := contract.EvaluateTransaction("GetAssetHistory", "IVSLAB-PVC23FG0001")
	if err != nil {
		fmt.Printf("failed to submit transaction: %s\n", err)
		return
	}
	result := formatJSON(evaluateResult)
	fmt.Printf("*** Result:%s\n", result)
}

// Submit transaction, passing in the wrong number of arguments ,expected to throw an error containing details of any error responses from the smart contract.
func exampleErrorHandling(contract *client.Contract) {
	fmt.Println("\n--> Submit Transaction: UpdateAsset IVSLAB-N23FA04, IVSLAB-N23FA04 does not exist and should return an error")
	_, err := contract.SubmitTransaction("UpdateAsset", "IVSLAB-N23FA04", "Network.co", "Taiwan", "NetworkChip-v1", "NPN304AA", "SNN30A14AA", "Network-Org", "2023-05-15")
	if err == nil {
		panic("******** FAILED to return an error")
	}
	fmt.Println("*** Successfully caught the error:")
	switch err := err.(type) {
	case *client.EndorseError:
		fmt.Printf("Endorse error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
	case *client.SubmitError:
		fmt.Printf("Submit error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
	case *client.CommitStatusError:
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Timeout waiting for transaction %s commit status: %s", err.TransactionID, err)
		} else {
			fmt.Printf("Error obtaining commit status for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
		}
	case *client.CommitError:
		fmt.Printf("Transaction %s failed to commit with status %d: %s\n", err.TransactionID, int32(err.Code), err)
	default:
		panic(fmt.Errorf("unexpected error type %T: %w", err, err))
	}

	// Any error that originates from a peer or orderer node external to the gateway will have its details
	// embedded within the gRPC status error. The following code shows how to extract that.
	statusErr := status.Convert(err)
	details := statusErr.Details()
	if len(details) > 0 {
		fmt.Println("Error Details:")
		for _, detail := range details {
			switch detail := detail.(type) {
			case *gateway.ErrorDetail:
				fmt.Printf("- address: %s, mspId: %s, message: %s\n", detail.Address, detail.MspId, detail.Message)
			}
		}
	}
}

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
