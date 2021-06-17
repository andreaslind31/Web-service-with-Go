package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init(){
	fmt.Println("loading products..")
	prodMap, err:= loadProductMap()
	productMap.m = prodMap
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded...\n ", len(productMap.m))
}

func loadProductMap()  (map[int]Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}
	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}
	return prodMap, nil
}

func getProduct(productId int) *Product {
	productMap.RLock() //ReadLock
	defer productMap.RUnlock() //ReadUnlock
	if product, ok := productMap.m[productId]; ok {
		return &product
	}
	return nil
}

func getProductList() []Product {
	productMap.RLock()
	products := make([]Product, 0, len(productMap.m))
	for _, value := range productMap.m {
		products = append(products, value)
	}
	productMap.RUnlock()
	return products
}

func getProductIds() []int {
	productMap.RLock()
	productIds := []int{}
	for key := range productMap.m {
		productIds = append(productIds, key)
	}
	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds	
}

func getNextProductId()  int {
	productIds := getProductIds()
	return productIds[len(productIds)-1] + 1
}

func removeProduct(productID int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, productID)
}

func addOrUpdateProduct(product Product) (int, error) {
	//if the product id is set, update, otherwise add
	addOrUpdateId := -1
	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		//if exists, replace it, otherwise return error
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
		}
		addOrUpdateId = product.ProductID
	} else {
		addOrUpdateId = getNextProductId()
		product.ProductID = addOrUpdateId
	}
	productMap.Lock()
	productMap.m[addOrUpdateId] = product
	productMap.Unlock()
	return addOrUpdateId, nil
}