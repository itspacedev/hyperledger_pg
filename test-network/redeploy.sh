#!/bin/bash

function setEnvironementorg1() {
  export PATH=${PWD}/../bin:$PATH
  export FABRIC_CFG_PATH=$PWD/../config/
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
  export CORE_PEER_ADDRESS=localhost:7051
}
function setEnvironementorg2() {
  export PATH=${PWD}/../bin:$PATH
  export FABRIC_CFG_PATH=$PWD/../config/
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org2MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
  export CORE_PEER_ADDRESS=localhost:9051
}

function runNetwork() {
    ./network.sh down
    ./network.sh up
    ./network.sh createChannel
}

function deployChaincode() {
    # Chaincode name
    CC_NAME=kirill
    CC_LABEL=${CC_NAME}_1.0
    echo "CC_NAME = ${CC_NAME}"
    echo "CC_LABEL = ${CC_LABEL}"

    # Compile chaincode package
    peer lifecycle chaincode package ${CC_NAME}.tar.gz --path ../asset-transfer-basic/chaincode-kirill-go/ --lang golang --label ${CC_LABEL}

    # Install chaincode on Org1
    setEnvironementorg1
    peer lifecycle chaincode install ${CC_NAME}.tar.gz
    # Install chaincode on Org2
    setEnvironementorg2
    peer lifecycle chaincode install ${CC_NAME}.tar.gz

    # Approve the chaincode by Org1
    setEnvironementorg1
    CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid ${CC_NAME}.tar.gz)
    echo "CC_PACKAGE_ID = ${CC_PACKAGE_ID}"
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --package-id ${CC_PACKAGE_ID} --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

    # Approve the chaincode by Org2
    setEnvironementorg2
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --package-id ${CC_PACKAGE_ID} --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

    # Commit the chaincode
    setEnvironementorg1
    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"


    # Check if chaincode is commited
    setEnvironementorg1
    peer lifecycle chaincode querycommitted --channelID mychannel --name basic

    # Initialize the chaincode
    # Here we call a smart contract's method "InitLedger" to generate initial state
    setEnvironementorg1
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'
}

runNetwork
deployChaincode