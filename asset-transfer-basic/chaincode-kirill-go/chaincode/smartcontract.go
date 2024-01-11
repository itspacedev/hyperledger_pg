package chaincode

import (
	"encoding/json"
	"fmt"
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
const index_product_map = "map~companyid~productid"

type Product struct {
	ID    int    `json:"ID"`
	Title string `json:"Title"`
}
type ProductMap struct {
	MapID     string  `json:"MapID"`
	CompanyID int     `json:"CompanyID"`
	ProductID int     `json:"ProductID"`
	Price     float64 `json:"Price"`
	Quantity  int     `json:"Quantity"`
}
type ProductInfo struct {
	ID          int     `json:"ID"`
	Title       string  `json:"Title"`
	CompanyID   int     `json:"CompanyID"`
	CompanyName string  `json:"CompanyName"`
	Price       float64 `json:"Price"`
	Quantity    int     `json:"Quantity"`
}

// InitLedger adds a base set of assets to the ledger
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
		ID:    productID,
		Title: title,
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

	mapID := strings.Join([]string{"map", strconv.Itoa(companyID), strconv.Itoa(product.ID)}, "_")
	productMap := &ProductMap{
		MapID:     mapID,
		CompanyID: companyID,
		ProductID: productID,
		Price:     price,
		Quantity:  quantity,
	}
	productMapJson, err := json.Marshal(productMap)
	if err != nil {
		return err
	}
	productMapStateKey := strings.Join([]string{"map", strconv.Itoa(companyID), strconv.Itoa(product.ID)}, "_")
	err = ctx.GetStub().PutState(productMapStateKey, productMapJson)
	if err != nil {
		return err
	}
	productMapIndexKey, err := ctx.GetStub().CreateCompositeKey(index_product_map, []string{"map", strconv.Itoa(companyID), strconv.Itoa(product.ID)})
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(productMapIndexKey, emptyValue)
	if err != nil {
		return err
	}
	return nil
}

// GetProducts Get information of all products
func (s *SmartContract) GetProducts(ctx contractapi.TransactionContextInterface) ([]*ProductInfo, error) {
	// Execute a key range query on all keys starting with 'map'
	resultIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(index_product_map, []string{"map"})
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	var productsInfo []*ProductInfo

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
			key_companyID := compositeKeyParts[1]
			key_productID := compositeKeyParts[2]

			companyStateKey := strings.Join([]string{"company", key_companyID}, "_")
			companyJson, err := ctx.GetStub().GetState(companyStateKey)
			if err != nil {
				return nil, err
			}
			var company Company
			err = json.Unmarshal(companyJson, &company)
			if err != nil {
				return nil, err
			}

			productStateKey := strings.Join([]string{"product", key_productID}, "_")
			productJson, err := ctx.GetStub().GetState(productStateKey)
			if err != nil {
				return nil, err
			}
			var product Product
			err = json.Unmarshal(productJson, &product)
			if err != nil {
				return nil, err
			}

			mapStateKey := strings.Join([]string{"map", key_companyID, key_productID}, "_")
			mapJson, err := ctx.GetStub().GetState(mapStateKey)
			if err != nil {
				return nil, err
			}
			var productMap ProductMap
			err = json.Unmarshal(mapJson, &productMap)
			if err != nil {
				return nil, err
			}

			// Create an info object
			productInfo := &ProductInfo{
				ID:          product.ID,
				Title:       product.Title,
				CompanyID:   company.ID,
				CompanyName: company.Name,
				Price:       productMap.Price,
				Quantity:    productMap.Quantity,
			}
			productsInfo = append(productsInfo, productInfo)
		}
	}
	return productsInfo, nil
}

// GetCompanyProducts Get information of all products of a particular company
func (s *SmartContract) GetCompanyProducts(ctx contractapi.TransactionContextInterface, companyID int) ([]*ProductInfo, error) {
	// Execute a key range query on all keys starting with 'map~companyID'
	resultIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(index_product_map, []string{"map", strconv.Itoa(companyID)})
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	var productsInfo []*ProductInfo

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
			key_companyID := compositeKeyParts[1]
			key_productID := compositeKeyParts[2]

			companyStateKey := strings.Join([]string{"company", key_companyID}, "_")
			companyJson, err := ctx.GetStub().GetState(companyStateKey)
			if err != nil {
				return nil, err
			}
			var company Company
			err = json.Unmarshal(companyJson, &company)
			if err != nil {
				return nil, err
			}

			productStateKey := strings.Join([]string{"product", key_productID}, "_")
			productJson, err := ctx.GetStub().GetState(productStateKey)
			if err != nil {
				return nil, err
			}
			var product Product
			err = json.Unmarshal(productJson, &product)
			if err != nil {
				return nil, err
			}

			mapStateKey := strings.Join([]string{"map", key_companyID, key_productID}, "_")
			mapJson, err := ctx.GetStub().GetState(mapStateKey)
			if err != nil {
				return nil, err
			}
			var productMap ProductMap
			err = json.Unmarshal(mapJson, &productMap)
			if err != nil {
				return nil, err
			}

			// Create an info object
			productInfo := &ProductInfo{
				ID:          product.ID,
				Title:       product.Title,
				CompanyID:   company.ID,
				CompanyName: company.Name,
				Price:       productMap.Price,
				Quantity:    productMap.Quantity,
			}
			productsInfo = append(productsInfo, productInfo)
		}
	}
	return productsInfo, nil
}

// BuyProduct Add a new product
func (s *SmartContract) BuyProduct(ctx contractapi.TransactionContextInterface, buyerCompanyID int, sellerCompanyID int, productID int, quantity int) error {

	// Полачаем компанию кто будет покупать - Покупатель
	buyerStateKey := strings.Join([]string{"company", strconv.Itoa(buyerCompanyID)}, "_")
	buyerJson, err := ctx.GetStub().GetState(buyerStateKey)
	if err != nil {
		return err
	}
	var buyerCompany Company
	err = json.Unmarshal(buyerJson, &buyerCompany)
	if err != nil {
		return err
	}

	// Полачаем компанию кто будет продавать - Продавец
	sellerStateKey := strings.Join([]string{"company", strconv.Itoa(sellerCompanyID)}, "_")
	sellerJson, err := ctx.GetStub().GetState(sellerStateKey)
	if err != nil {
		return err
	}
	var sellerCompany Company
	err = json.Unmarshal(sellerJson, &sellerCompany)
	if err != nil {
		return err
	}

	// Получаем товар
	productStateKey := strings.Join([]string{"product", strconv.Itoa(productID)}, "_")
	productJson, err := ctx.GetStub().GetState(productStateKey)
	if err != nil {
		return err
	}
	var product Product
	err = json.Unmarshal(productJson, &product)
	if err != nil {
		return err
	}

	// Находим связку Продавец-Товар
	// Execute a key range query on all keys starting with 'map~sellerCompanyID'
	sellerMapStateKey := strings.Join([]string{"map", strconv.Itoa(sellerCompanyID), strconv.Itoa(productID)}, "_")
	sellerMapJson, err := ctx.GetStub().GetState(sellerMapStateKey)
	if err != nil {
		return err
	}
	if sellerMapJson == nil {
		return fmt.Errorf("Seller [%s] does not have this product", sellerCompany.Name)
	}
	var sellerProductMap ProductMap
	err = json.Unmarshal(sellerMapJson, &sellerProductMap)
	if err != nil {
		return err
	}

	// Проверка 1
	// Проверяем что продавец имеет в наличие количество товара которое хотим купить
	if sellerProductMap.Quantity < quantity {
		return fmt.Errorf("Seller [%s] does not have enough items to sell. In stock: %d, Needed: %d", sellerCompany.Name, sellerProductMap.Quantity, quantity)
	}

	// Проверка 2
	// Проверяем что у покупателя достаточно денег для покупки
	totalAmount := sellerProductMap.Price * float64(quantity)
	if buyerCompany.Balance < totalAmount {
		return fmt.Errorf("Buyer [%s] does not have enough balance to buy. Available balance: %d, Needed: %d", buyerCompany.Name, buyerCompany.Balance, totalAmount)
	}

	// Находим связку Покупатель-Товар
	// Execute a key range query on all keys starting with 'map~sellerCompanyID'
	buyerMapStateKey := strings.Join([]string{"map", strconv.Itoa(buyerCompanyID), strconv.Itoa(productID)}, "_")
	buyerMapJson, err := ctx.GetStub().GetState(buyerMapStateKey)
	if err != nil {
		return err
	}
	var buyerProductMap ProductMap
	if buyerMapJson == nil {
		// Create an empty map for buyer
		buyerProductMap := &ProductMap{
			MapID:     buyerMapStateKey,
			CompanyID: buyerCompanyID,
			ProductID: productID,
			Price:     sellerProductMap.Price, // Ставим цену продавца потому что другой пока нет
			Quantity:  0,
		}
		buyerProductMapJson, err := json.Marshal(buyerProductMap)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(buyerMapStateKey, buyerProductMapJson)
		if err != nil {
			return err
		}
		buyerProductMapIndexKey, err := ctx.GetStub().CreateCompositeKey(index_product_map, []string{"map", strconv.Itoa(buyerCompanyID), strconv.Itoa(productID)})
		if err != nil {
			return err
		}
		emptyValue := []byte{0x00}
		err = ctx.GetStub().PutState(buyerProductMapIndexKey, emptyValue)
		if err != nil {
			return err
		}
	} else {
		// Update existing map
		err = json.Unmarshal(buyerMapJson, &buyerProductMap)
		if err != nil {
			return err
		}
	}

	// Обновляем баланс Продавца
	sellerCompany.Balance = sellerCompany.Balance + totalAmount
	sellerCompanyJson, err := json.Marshal(sellerCompany)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(sellerStateKey, sellerCompanyJson)
	if err != nil {
		return err
	}

	// Обновляем баланс Покупателя
	buyerCompany.Balance = buyerCompany.Balance - totalAmount
	buyerCompanyJson, err := json.Marshal(buyerCompany)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(buyerStateKey, buyerCompanyJson)
	if err != nil {
		return err
	}

	// Обновляем связку Продавец-Товар
	sellerProductMap.Quantity = sellerProductMap.Quantity - quantity
	sellerProductMapJson, err := json.Marshal(sellerProductMap)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(sellerMapStateKey, sellerProductMapJson)
	if err != nil {
		return err
	}

	// Обновляем связку Покупатель-Товар
	buyerProductMap.Quantity = buyerProductMap.Quantity + quantity
	buyerProductMapJson, err := json.Marshal(buyerProductMap)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(buyerMapStateKey, buyerProductMapJson)
	if err != nil {
		return err
	}

	return nil
}
