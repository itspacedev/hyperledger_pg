package chaincode

import (
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
	"strings"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Company struct {
	ID      int     `json:"ID"`
	Name    string  `json:"Name"`
	Balance float64 `json:"Balance"`
}

// Ven. return fmt.Errorf("failed to put to world state. %v", err)

const index_company = "company~companyid"
const index_product = "product~productid"

type Product struct {
	ID        int     `json:"ID"`
	CompanyID int     `json:"CompanyID"`
	Title     string  `json:"Title"`
	Price     float64 `json:"Price"`
	Quantity  int     `json:"Quantity"`
}

// InitLedger adds a base set of assets to the ledger
// 1. Create 2 companies
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Create 2 companies
	companies := []Company{
		{ID: 1, Name: "Kirill Company LTD", Balance: 1300.0},
		{ID: 2, Name: "Veniamin Holding", Balance: 1200.31},
	}
	for _, company := range companies {
		companyJson, err := json.Marshal(company)
		if err != nil {
			return err
		}
		companyStateKey := strings.Join([]string{"company", strconv.Itoa(company.ID)}, "_")
		err = ctx.GetStub().PutState(companyStateKey, companyJson)
		if err != nil {
			return err
		}
		companyCompositeKey, err := ctx.GetStub().CreateCompositeKey(index_company, []string{"company", strconv.Itoa(company.ID)})
		if err != nil {
			return err
		}
		emptyValue := []byte{0x00}
		err = ctx.GetStub().PutState(companyCompositeKey, emptyValue)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCompanies Get information of all companies
func (s *SmartContract) GetCompanies(ctx contractapi.TransactionContextInterface) ([]*Company, error) {
	// Execute a key range query on all keys starting with 'company'
	resultIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(index_company, []string{"company"})
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	var companies []*Company

	for resultIterator.HasNext() {
		responseRange, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, err
		}
		if len(compositeKeyParts) > 1 {
			recordId := compositeKeyParts[1]
			companyStateKey := strings.Join([]string{"company", recordId}, "_")
			companyJson, err := ctx.GetStub().GetState(companyStateKey)
			if err != nil {
				return nil, err
			}
			var company Company
			err = json.Unmarshal(companyJson, &company)
			if err != nil {
				return nil, err
			}
			companies = append(companies, &company)
		}
	}
	return companies, nil
}

// AddProduct Add a new product
func (s *SmartContract) AddProduct(ctx contractapi.TransactionContextInterface, productID int, companyID int, title string, price float64, quantity int) error {
	product := &Product{
		ID:        productID,
		CompanyID: companyID,
		Title:     title,
		Price:     price,
		Quantity:  quantity,
	}
	productJson, err := json.Marshal(product)
	if err != nil {
		return err
	}
	productStateKey := strings.Join([]string{"product", strconv.Itoa(product.ID)}, "_")
	err = ctx.GetStub().PutState(productStateKey, productJson)
	if err != nil {
		return err
	}
	productIndexKey, err := ctx.GetStub().CreateCompositeKey(index_product, []string{"product", strconv.Itoa(product.ID)})
	if err != nil {
		return err
	}
	emptyValue := []byte{0x00}
	err = ctx.GetStub().PutState(productIndexKey, emptyValue)
	if err != nil {
		return err
	}
	return nil
}

// GetProducts Get information of all products
func (s *SmartContract) GetProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	// Execute a key range query on all keys starting with 'company'
	resultIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(index_product, []string{"product"})
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	var products []*Product

	for resultIterator.HasNext() {
		responseRange, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, err
		}
		if len(compositeKeyParts) > 1 {
			recordId := compositeKeyParts[1]
			productStateKey := strings.Join([]string{"product", recordId}, "_")
			productJson, err := ctx.GetStub().GetState(productStateKey)
			if err != nil {
				return nil, err
			}
			var product Product
			err = json.Unmarshal(productJson, &product)
			if err != nil {
				return nil, err
			}
			products = append(products, &product)
		}
	}
	return products, nil
}
