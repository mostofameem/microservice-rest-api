package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	restapi "order_service/Rest-Api"
	"order_service/db"
	"order_service/logger"
	"order_service/web/utils"
	"strconv"
	"strings"
)

func UpdateOrderInfo(w http.ResponseWriter, r *http.Request) {
	id, err := ExtractID(r.URL.Path)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("error converting ID to integer"))
		return
	}
	db.OrderID = id

	err = db.GetOrderTypeRepo().UpdateOrderInfoDB()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}
	utils.SendData(w, "Updated Order Info Successful")

}

type ResponseProduct struct {
	/*
		"product_id": 1,
		"product_name": "Computer",
		"product_price": 45000,
		"product_quantity": 10,
		"vendor_id": 1
	*/
	Product_id       int    `json:"product_id"`
	Product_name     string `json:"product_name"`
	Product_price    int    `json:"product_price"`
	Product_quantity int    `json:"product_quantity"`
	Vendor_id        int    `json:"vendor_id"`
}

type UserRequirement struct {
	Product_name     string `json:"product_name"`
	Product_quantity int    `json:"product_quantity"`
	Vendor_id        int    `json:"vendor_id"`
}
type ApiResponse struct {
	Data    db.Product `json:"data"`
	Message string     `json:"message"`
	Status  bool       `json:"status"`
}

func AddToOrderList(w http.ResponseWriter, r *http.Request) {

	var item UserRequirement
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		slog.Error("Failed to get user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": item,
		}))
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	// Initialize query builder
	qb := restapi.NewQueryBuilder("http://192.168.48.31:3000/users/productinfo").
		QueryParam("product_name", item.Product_name).
		QueryParam("vendor_id", strconv.Itoa(item.Vendor_id))

	// Execute request
	respBody, err := qb.Execute(http.MethodGet)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var productInfo ResponseProduct
	err = json.Unmarshal(respBody, &productInfo)
	if err != nil {
		http.Error(w, "Failed to parse product info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendData(w, productInfo)

}

func GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	id, err := ExtractID(r.URL.Path)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("error converting ID to integer"))
		return
	}

	listch := make(chan []db.OrderList)
	totalch := make(chan int)

	go db.GetOrderTypeRepo().GetOrderList(id, listch)
	go db.GetOrderTypeRepo().GetSumOfTotalCost(id, totalch)

	list := <-listch
	total := <-totalch
	utils.SendBothData(w, fmt.Sprintf("Total =%d", total), list)
}

// ExtractID extracts the ID from the request URL and returns it as an integer.
func ExtractID(url string) (int, error) {
	parts := strings.Split(url, "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	return id, err
}
