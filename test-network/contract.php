<?php

// To run this script, run in command line:
// php ./contract.php

class SmartContract {

    private $debug = false;
    private $delay = 3;

    /**
     * Prepare and execute BASH command
     * @param $organizationId
     * @param $queryType
     * @param $commandName
     * @param $commandParams
     * @return false|string|null
     */
    private function executeCommand(
        $organizationId,
        $queryType,
        $commandName,
        $commandParams = []
    ) {
        if ($queryType === 'QUERY') {
            // Read from blockchain
            $allParams = [$commandName];
            if (count($commandParams) > 0) {
                foreach ($commandParams as $commandParam) {
                    $allParams[] = $commandParam;
                }
            }
            $command = '\'{"Args":["'.implode('","', $allParams).'"]}\'';
        } elseif ($queryType === 'INVOKE') {
            // Change state in blockchain
            $command = '\'{"function":"'.$commandName.'", "Args":["'.implode('","', $commandParams).'"]}\'';
        } else {
            echo 'ERROR: Incorrect Query TYpe'.PHP_EOL;
            die;
        }
        $shellCommand = "./contract.sh $organizationId $queryType $command";

        if ($this->debug) {
             echo PHP_EOL;
             echo 'SHELL COMMAND: '.$shellCommand.PHP_EOL;
             echo PHP_EOL;
        }
        $output = shell_exec($shellCommand);

        // Надо подождать пару секунд чтобы изменения появились в блокчейне
        // Если сразу делать выборку, может вернуть пустой массив
        // Короче надо подождать скоко-то, я поставил 3 секунды
        if ($queryType === 'INVOKE') {
            //echo '===> WAIT '.$this->delay.' seconds...'.PHP_EOL;
            sleep($this->delay);
        }
        return $output;
    }

    /**
     * Display information about companies
     * @return void
     */
    public function displayCompanies() {
        $output = $this->executeCommand(1, 'QUERY', 'GetCompanies');
        $output = trim($output);
        $companies = json_decode($output, true);
        return $companies;

//        if ($this->debug) {
//             print_r($companies);
//        }
//
//        echo PHP_EOL;
//        echo '===== Companies ====='.PHP_EOL;
//        foreach ($companies as $company) {
//            echo $company['ID'].') '.$company['Name'].', Balance: '.$company['Balance'].' RUB'.PHP_EOL;
//        }
//        echo PHP_EOL;
    }

    /**
     * Add a new product
     * @param $productId
     * @param $companyId
     * @param $title
     * @param $price
     * @param $quantity
     * @return void
     */
    public function addProduct($productId, $companyId, $title, $price, $quantity) {
        $params = [$productId, $companyId, $title, $price, $quantity];
        $output = $this->executeCommand($companyId, 'INVOKE', 'AddProduct', $params);
    }

    /**
     * Display information about products
     * @return void
     */
    public function displayProducts() {
        $output = $this->executeCommand(1, 'QUERY', 'GetProducts');
        $output = trim($output);
        $products = json_decode($output, true);

        if ($this->debug) {
             print_r($products);
        }

//        echo PHP_EOL;
//        echo '===== Products ====='.PHP_EOL;

//        $outputIndex = 1;
//        foreach ($products as $product) {
//            echo $outputIndex.') '.$product['Title'].' ('.$product['Price'].' RUB - '.$product['Quantity'].' items) [Company #'.$product['CompanyID'].']'.PHP_EOL;
//            $outputIndex++;
//        }
//        echo PHP_EOL;
    }

    /**
     * Display information about products of a company
     * @return void
     */
    public function displayCompanyProducts($companyId) {
        $params = [$companyId];
        $output = $this->executeCommand($companyId, 'QUERY', 'GetCompanyProducts', $params);
        $output = trim($output);
        $products = json_decode($output, true);

        if ($this->debug) {
            print_r($products);
        }

//        echo PHP_EOL;
//        echo '===== Products ====='.PHP_EOL;
//
//        $outputIndex = 1;
//        foreach ($products as $product) {
//            echo $outputIndex.') '.$product['Title'].' ('.$product['Price'].' RUB - '.$product['Quantity'].' items) [Company #'.$product['CompanyID'].']'.PHP_EOL;
//            $outputIndex++;
//        }
//        echo PHP_EOL;
    }

    // BuyProduct(buyerCompanyID int, sellerCompanyID int, productID int, quantity int) error {
    public function buyProduct(
        $buyerCompanyId,
        $sellerCompanyId,
        $productId,
        $quantity
    ) {
        $params = [$buyerCompanyId, $sellerCompanyId, $productId, $quantity];
        $output = $this->executeCommand($buyerCompanyId, 'INVOKE', 'BuyProduct', $params);
    }
}

$smartContract = new SmartContract();

// Какой метод вызывать
$method_to_call = $_POST['method_to_call'];
$response = null;
if ($method_to_call === 'get_companies') {
    $response = $smartContract->displayCompanies();
}
if ($response !== null) {
    echo json_encode([
       'status' => 1,
       'data' => $response
    ]);
} else {
    echo json_encode([
        'status' => 0
    ]);
}

//// 2. Add products
//$smartContract->addProduct(1, 1, 'Sony Playstation 5 - Gaming Console', 300.50, 5);
//$smartContract->addProduct(2, 2, 'PS Game - The Last of Us 2', 42.22, 15);
//$smartContract->addProduct(2, 1, 'PS Game - The Last of Us 2', 35, 4);
//
//// 3. Display products
//$smartContract->displayProducts();
//
//// 4. Display products for company 1
//$smartContract->displayCompanyProducts(1);
//
//// 5. Display products for company 2
//$smartContract->displayCompanyProducts(2);
//
//// 6. Buy a product
//// Buyer - 1
//// Seller - 2
//// Product - 2
//// Quantity - 3
//$smartContract->buyProduct(1, 2, 2, 3);
//
//// 7. Display companies
//$smartContract->displayCompanies();
//
//// 8. Display products
//$smartContract->displayProducts();
//
//// 9. Display products for company 1
//$smartContract->displayCompanyProducts(1);
//
//// 10. Display products for company 2
//$smartContract->displayCompanyProducts(2);
