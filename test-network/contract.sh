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

function executeQuery() {
    PARAM_ORGANIZATION_ID=$1
    PARAM_QUERY_TYPE=$2
    PARAM_QUERY_STRING=$3

    # echo "${PARAM_ORGANIZATION_ID} - ${PARAM_QUERY_TYPE} - ${PARAM_QUERY_STRING}"
    if [[ "${PARAM_ORGANIZATION_ID}" == "1" ]]; then
        setEnvironementorg1
    else
        setEnvironementorg2
    fi
    if [[ "${PARAM_QUERY_TYPE}" == "QUERY" ]]; then
        peer chaincode query -C mychannel -n basic -c "${PARAM_QUERY_STRING}"
    else
        peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
              --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
              -C mychannel -n basic --peerAddresses localhost:7051 \
              --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
              --peerAddresses localhost:9051 \
              --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
              -c "${PARAM_QUERY_STRING}"
    fi
}

## Parse mode from command line arguments
if [[ $# -lt 1 ]] ; then
    echo "ERROR: Invalid command"
    exit 0
else
    ORGANIZATION_ID=$1
    QUERY_TYPE=$2
    QUERY_STRING=$3
fi

executeQuery "${ORGANIZATION_ID}" "${QUERY_TYPE}" "${QUERY_STRING}"