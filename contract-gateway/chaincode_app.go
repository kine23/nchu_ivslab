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
	certPath= cryptoPath + "/registers_users/admin1/msp/signcerts/cert.pem"	// client數位簽章
	keyPath= cryptoPath + "/registers_users/admin1/msp/keystore/"			// client私鑰路徑
	tlsCertPath= cryptoPath + "/alliance/tls-ca-cert.pem"					// client tls證書
	peerEndpoint= "peer1.brand.ivsorg.net:7151"								// peer節點地址
	gatewayPeer= "peer1.brand.ivsorg.net"									// peer節點名稱
)

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

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
	getAllAssets(contract)
	getAllParts(contract)
	createAsset(contract)
	readAssetByID(contract)
	readPartByID(contract)
	getAssetsByRange(contract)
	getPartsByRange(contract)
	queryPartsByOwner(contract)
	queryAssets(contract)
	queryAssetsWithPagination(contract)
	getAssetHistory(contract)
	exampleErrorHandling(contract)
	
	createPart(contract, PartArgs{
    		PID: "IVSLAB-S23FA0002",
    		Manufacturer: "Security.Co",
    		ManufactureLocation: "Taiwan",
    		PartName: "SecurityChip-v1",
    		PartNumber: "SPN3R1C00AA2",
    		Organization: "Security-Org",
    		ManufactureDate: "2023-05-17",
		})

	createPart(contract, PartArgs{
    		PID: "IVSLAB-N23FA0002",
    		Manufacturer: "Network.Co",
    		ManufactureLocation: "Taiwan",
    		PartName: "NetworkChip-v1",
    		PartNumber: "NPN3R1C00AA2",
    		Organization: "Network-Org",
    		ManufactureDate: "2023-05-17",
		})

	createPart(contract, PartArgs{
    		PID: "IVSLAB-C23FA0002",
    		Manufacturer: "CMOS.Co",
    		ManufactureLocation: "USA",
    		PartName: "CMOSChip-v1",
    		PartNumber: "CPN3R1C00AA2",
    		Organization: "CMOS-Org",
    		ManufactureDate: "2023-05-17",
		})

	createPart(contract, PartArgs{
    		PID: "IVSLAB-V23FA0002",
    		Manufacturer: "VideoCodec.Co",
   		ManufactureLocation: "USA",
    		PartName: "VideoCodecChip-v1",
    		PartNumber: "VPN3R1C00AA2",
    		Organization: "VideoCodec-Org",
    		ManufactureDate: "2023-05-17",
		})
	
	transferPartAsync(contract, TransferPartArgs{
    		PID: "IVSLAB-S23FA0002",
    		Organization: "Brand-Org",
    		ManufactureDate: "2023-05-17",
		})

	transferPartAsync(contract, TransferPartArgs{
    		PID: "IVSLAB-N23FA0002",
    		Organization: "Brand-Org",
    		ManufactureDate: "2023-05-17",
		})

	transferPartAsync(contract, TransferPartArgs{
    		PID: "IVSLAB-C23FA0002",
    		Organization: "Brand-Org",
    		ManufactureDate: "2023-05-17",
		})

	transferPartAsync(contract, TransferPartArgs{
    		PID: "IVSLAB-V23FA0002",
    		Organization: "Brand-Org",
    		ManufactureDate: "2023-05-17",
		})	
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
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func getAllParts(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAllParts, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllParts")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Evaluate a transaction to query ledger state.
func getAllAssets(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}
type PartArgs struct {
    PID                  string
    Manufacturer         string
    ManufactureLocation  string
    PartName             string
    PartNumber           string
    Organization         string
    ManufactureDate      string
}

func createPart(contract *client.Contract, args PartArgs) {
    fmt.Printf("\n--> Submit Transaction: CreatePart, creates new part with PID, Manufacturer, ManufactureLocation, PartName, PartNumber, ManufactureDate, Organization \n")
    _, err := contract.SubmitTransaction("CreatePart", args.PID, args.Manufacturer, args.ManufactureLocation, args.PartName, args.PartNumber, args.Organization, args.ManufactureDate)
    if err != nil {
        panic(fmt.Errorf("failed to submit transaction: %w", err))
    }
    fmt.Printf("*** Transaction committed successfully\n")
}

type TransferPartArgs struct {
    PID                  string
    Organization         string
    ManufactureDate      string
}

func transferPartAsync(contract *client.Contract, args TransferPartArgs) {
	fmt.Printf("\n--> Async Submit Transaction: TransferPart, updates existing part Organization and TransferDate")
	submitResult, commit, err := contract.SubmitAsync("TransferPart", client.WithArguments(args.PID, args.ManufactureDate, args.Organization))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	fmt.Printf("\n*** Successfully submitted transaction to transfer ownership from %s to %d. \n", string(submitResult))
	fmt.Println("*** Waiting for transaction commit.")
	if commitStatus, err := commit.Status(); err != nil {
		panic(fmt.Errorf("failed to get commit status: %w", err))
	} else if !commitStatus.Successful {
		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
	}
	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func createAsset(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: CreateAsset, creates new asset with ID, MadeBy, MadeIn, SerialNumber, SecurityChip, NetworkChip, CMOSChip, VideoCodecChip, ProductionDate \n")

	_, err := contract.SubmitTransaction("CreateAsset", "IVSLAB-PVC23FG0002", "Bard.Co", "Taiwan", "IVSPN902300AACDC02", "IVSLAB-S23FA0002", "IVSLAB-N23FA0002", "IVSLAB-C23FA0002", "IVSLAB-V23FA0002", "2023-05-18")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction by assetID to query ledger state.
func readPartByID(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: ReadPart, function returns asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadPart", "IVSLAB-S23FA0002")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func readAssetByID(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsser, function returns asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", "IVSLAB-PVC23FG0001")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func getPartsByRange(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetPartsByRange, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetPartsByRange", "IVSLAB-V23FA0001", "IVSLAB-V23FA0004")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func getAssetsByRange(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAssetsByRange, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAssetsByRange", "IVSLAB-PVC23FG0001", "IVSLAB-PVC23FG0003")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func queryPartsByOwner(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: QueryPartsByOwner, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("QueryPartsByOwner", "Network-Org")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	// Add a check here for empty result
	if len(evaluateResult) == 0 {
		fmt.Println("*** No assets found for the specified organization")
		return
	}

	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func queryAssets(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: QueryAssets, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("QueryAssets", "IVSLAB-PVC23FG0001")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	// Add a check here for empty result
	if len(evaluateResult) == 0 {
		fmt.Println("*** No assets found for the specified organization")
		return
	}

	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func queryAssetsWithPagination(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: QueryAssetsWithPagination, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("QueryAssetsWithPagination", `{"selector":{"docType":"asset","madeby":"Brand.Co"}, "use_index":["_design/indexMadeByDoc", "indexMadeBy"]}`, "1", "")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	// Check if result is null
	if evaluateResult == nil {
		fmt.Println("*** No assets found for the specified madeby")
		return
	}

	result := formatJSON(evaluateResult)
	fmt.Printf("*** Result:%s\n", result)
}

func getAssetHistory(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAssetHistory, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAssetHistory", "IVSLAB-PVC23FG0001")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}
// Submit transaction, passing in the wrong number of arguments ,expected to throw an error containing details of any error responses from the smart contract.
func exampleErrorHandling(contract *client.Contract) {
	fmt.Println("\n--> Submit Transaction: UpdateAsset IVSLAB-N23FA03, IVSLAB-N23FA01 does not exist and should return an error")

	_, err := contract.SubmitTransaction("UpdateAsset", "IVSLAB-N23FA03", "Network.co", "Taiwan", "NetworkChip-v1", "NPN303AA", "SNN30A13AA", "Network-Org", "2023-05-15")
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
