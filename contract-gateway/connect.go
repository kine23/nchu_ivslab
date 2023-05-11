package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/kine23/nchu_ivslab/contract-gateway"
)

func main() {
    fmt.Println(brandmspID)       
    fmt.Println(brandcryptoPath)  
    fmt.Println(brandcertPath)    
    fmt.Println(brandkeyPath)     
    fmt.Println(brandtlsCertPath) 
    fmt.Println(brandpeerEndpoint)
    fmt.Println(brandgatewayPeer) 

    fmt.Println(securitymspID)       
    fmt.Println(securitycryptoPath)  
    fmt.Println(securitycertPath)    
    fmt.Println(securitykeyPath)     
    fmt.Println(securitytlsCertPath) 
    fmt.Println(securitypeerEndpoint)
    fmt.Println(securitygatewayPeer) 	

    fmt.Println(networkmspID)       
    fmt.Println(networkcryptoPath)  
    fmt.Println(networkcertPath)    
    fmt.Println(networkkeyPath)     
    fmt.Println(networktlsCertPath) 
    fmt.Println(networkpeerEndpoint)
    fmt.Println(networkgatewayPeer) 

    fmt.Println(cmosmspID)       
    fmt.Println(cmoscryptoPath)  
    fmt.Println(cmoscertPath)    
    fmt.Println(cmoskeyPath)     
    fmt.Println(cmosytlsCertPath)
    fmt.Println(cmospeerEndpoint)
    fmt.Println(cmosgatewayPeer) 
	
    fmt.Println(videocodecmspID)       
    fmt.Println(videocodeccryptoPath)  
    fmt.Println(videocodeccertPath)    
    fmt.Println(videocodeckeyPath)     
    fmt.Println(videocodectlsCertPath) 
    fmt.Println(videocodecpeerEndpoint)
    fmt.Println(videocodecgatewayPeer) 	
}

// 建立指向聯盟網路的gRPC連接
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

// 根據用戶指定的X.509證書為這個連接創造一個客戶端標識。
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

// 載入client證書
func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// 使用私鑰生成數位簽章
func newSign() identity.Sign {
	files, err := ioutil.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := ioutil.ReadFile(path.Join(keyPath, files[0].Name()))

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