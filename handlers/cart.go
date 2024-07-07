package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"order_service/db"
	"order_service/logger"
	"order_service/web/utils"
)

func NewCart(w http.ResponseWriter, r *http.Request) {

	order_id, err := db.GetOrderTypeRepo().NewCart(12) // 12 userId
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
	}
	utils.SendData(w, fmt.Sprintf("Order ID: %d", order_id))
}

type Productinfo struct {
	Vendor_id           int     `json:"vendor_id"`
	Product_name        string  `json:"product_name"`
	Product_price       float32 `json:"product_price"`
	Product_quantity    int     `json:"product_quantity"`
	Product_description string  `json:"product_description"`
}

type Apires struct {
	Data    Productinfo `json:"data"`
	Message string      `json:"message"`
	Status  bool        `json:"status"`
}

func InsertProduct(w http.ResponseWriter, r *http.Request) {
	var item Productinfo
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		slog.Error("Failed to get user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": item,
		}))
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	// Define the base URL for the service
	baseURL := "http://192.168.48.31:3000/users/addproduct"

	method := r.Method

	query_params := map[string]interface{}{}

	body_param := map[string]interface{}{
		"Vendor_id":           item.Vendor_id,
		"Product_name":        item.Product_name,
		"Product_price":       item.Product_price,
		"Product_quantity":    item.Product_quantity,
		"Product_description": item.Product_description,
	}

	body, err := RestApiCall(baseURL, method, body_param, query_params)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
	}

	// Unmarshal the response body into the ApiResponse struct
	var apiResponse Apires
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
		return
	}

	utils.SendData(w, apiResponse.Data)
}
