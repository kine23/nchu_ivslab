package main

const (
	brandmspID         = "brandMSP"															// 所屬組織的MSPID
	brandcryptoPath    = "/root/MyLab_IVS/organizations/brand.ivsorg.net"					// 中間變量
	brandcertPath      = cryptoPath + "/registers_users/user1/msp/signcerts/cert.pem"		// client數位簽章
	brandkeyPath       = cryptoPath + "/registers_users/user1/msp/keystore/"				// client私鑰路徑
	brandtlsCertPath   = cryptoPath + "/alliance/tls-ca-cert.pem"							// client tls證書
	brandpeerEndpoint  = "peer1.brand.ivsorg.net:7151"										// peer節點地址
	brandgatewayPeer   = "peer1.brand.ivsorg.net"											// peer節點名稱
	
	securitymspID         = "securityMSP"													
	securitycryptoPath    = "/root/MyLab_IVS/organizations/security.ivsorg.net"				
	securitycertPath      = cryptoPath + "/registers_users/user1/msp/signcerts/cert.pem"	
	securitykeyPath       = cryptoPath + "/registers_users/user1/msp/keystore/"				
	securitytlsCertPath   = cryptoPath + "/alliance/tls-ca-cert.pem"						
	securitypeerEndpoint  = "peer1.security.ivsorg.net:7251"								
	securitygatewayPeer   = "peer1.security.ivsorg.net"										
	
	networkmspID         = "networkMSP"														
	networkcryptoPath    = "/root/MyLab_IVS/organizations/network.ivsorg.net"				
	networkcertPath      = cryptoPath + "/registers_users/user1/msp/signcerts/cert.pem"		
	networkkeyPath       = cryptoPath + "/registers_users/user1/msp/keystore/"				
	networktlsCertPath   = cryptoPath + "/alliance/tls-ca-cert.pem"							
	networkpeerEndpoint  = "peer1.network.ivsorg.net:7351"									
	networkgatewayPeer   = "peer1.network.ivsorg.net"										
	
	cmosmspID         = "cmosMSP"															
	cmoscryptoPath    = "/root/MyLab_IVS/organizations/cmos.ivsorg.net"					
	cmoscertPath      = cryptoPath + "/registers_users/user1/msp/signcerts/cert.pem"	
	cmoskeyPath       = cryptoPath + "/registers_users/user1/msp/keystore/"				
	cmosytlsCertPath   = cryptoPath + "/alliance/tls-ca-cert.pem"						
	cmospeerEndpoint  = "peer1.cmos.ivsorg.net:7451"									
	cmosgatewayPeer   = "peer1.cmos.ivsorg.net"											
	
	videocodecmspID         = "videocodecMSP"													
	videocodeccryptoPath    = "/root/MyLab_IVS/organizations/videocodec.ivsorg.net"				
	videocodeccertPath      = cryptoPath + "/registers_users/user1/msp/signcerts/cert.pem"	
	videocodeckeyPath       = cryptoPath + "/registers_users/user1/msp/keystore/"				
	videocodectlsCertPath   = cryptoPath + "/alliance/tls-ca-cert.pem"						
	videocodecpeerEndpoint  = "peer1.videocodec.ivsorg.net:7551"								
	videocodecgatewayPeer   = "peer1.videocodec.ivsorg.net"										
)
