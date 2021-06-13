package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Product struct {
	ProductID      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var productList []Product

func init() {
	productsJSON := `[
		{
			"productId": 1,
			"manufacturer": "Alibaba",
			"sku": "p5xh4asadv",
			"upc": "933340000",
			"pricePerUnit": "4.45",
			"quantityOnHand": 9703,
			"productName": "sticky note"
		},
		{
			"productId": 2,
			"manufacturer": "Tjotahejti",
			"sku": "p5xh4asoobn",
			"upc": "933340113",
			"pricePerUnit": "49.5",
			"quantityOnHand": 903,
			"productName": "The Witcher: 1"
		},
		{
			"productId": 3,
			"manufacturer": "Amazon",
			"sku": "p5xh4asggdf",
			"upc": "933340453",
			"pricePerUnit": "97",
			"quantityOnHand": 403,
			"productName": "Standard barbell"
		}
	]`
	err := json.Unmarshal([]byte(productsJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}
func getNextId() int {
	highestId := -1
	for _, product := range productList{
		if highestId < product.ProductID {
			highestId = product.ProductID
		}
	}
	return highestId +1
}
func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJSON, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJSON)
	case http.MethodPost:
		//add a new product to the list
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newProduct.ProductID = getNextId()
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusCreated)
		return
	}
}
func main() {
	http.HandleFunc("/products", productsHandler)
	http.ListenAndServe(":5000", nil)
}