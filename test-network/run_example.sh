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

# Only works on the original smart contract
function runExampleOriginal() {
    # Query the chaincode by Org1
    setEnvironementorg1
    echo "---------- ORG1 GetAllAssets ----------"
    peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'

    # Query the chaincode by Org2
    setEnvironementorg2
    echo "---------- ORG2 GetAllAssets ----------"
    peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'
}

function initFunction() {
    # Get all companies (Org1)
    setEnvironementorg1
    echo "---------- ORG1 GetAllCompaniesBalances ----------"
    peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllCompaniesBalances"]}'

    # Get all companies (Org2)
    setEnvironementorg2
    echo "---------- ORG2 GetAllCompaniesBalances ----------"
    peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllCompaniesBalances"]}'

    # Get all products (Org1)
    setEnvironementorg1
    echo "---------- ORG1 GetAllProducts ----------"
    peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllProducts"]}'

    # Get all products (Org2)
    setEnvironementorg2
    echo "---------- ORG2 GetAllProducts ----------"
    peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllProducts"]}'

    # Get system info (Org1)
    setEnvironementorg1
    echo "---------- ORG1 GetSystemInfo ----------"
    peer chaincode query -C mychannel -n basic -c '{"Args":["GetSystemInfo"]}'

    # Get system info (Org2)
    setEnvironementorg2
    echo "---------- ORG2 GetSystemInfo ----------"
    peer chaincode query -C mychannel -n basic -c '{"Args":["GetSystemInfo"]}'
}

function buyProduct() {
    # Buy a product
    # "Kirill Coffee Shop" buy "Coffee Machine"
    setEnvironementorg1
    echo "---------- ORG1 BuyProduct ----------"
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
      --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
      -C mychannel -n basic --peerAddresses localhost:7051 \
      --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
      --peerAddresses localhost:9051 \
      --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
      -c '{"function":"BuyProduct","Args":[]}'
}

function showTransactionLogs() {
    # Get transaction logs (Org1)
    setEnvironementorg1
    echo "---------- ORG1 GetTransactionLogs ----------"
    peer chaincode query -C mychannel -n basic -c '{"Args":["GetTransactionLogs"]}'
}

## Parse mode from command line arguments
if [[ $# -lt 1 ]] ; then
    echo "MODES: init, buy, logs"
    exit 0
else
  MODE=$1
  shift
fi

# Call a function based on mode
if [ "${MODE}" == "init" ]; then
    initFunction
elif [ "${MODE}" == "buy" ]; then
    buyProduct
elif [ "${MODE}" == "logs" ]; then
    showTransactionLogs
else
    echo "Available modes: init, buy, logs"
    exit 1
fi