package api

// TODO:
// All request structs for *GET* method should be useless...
// It's better to process all query strings by DefaultQuery() than Bind() in *GET* handlers
// See ./restaurantDishOrder.go/searchRestaurant()
// If Bind() is used, json should be changed to form in all request structs

type searchRestaurantRequest struct {
	RestaurantName string `form:"restaurantName"`
	ZipCode        int    `form:"zipCode"`
	OrderBy        string `form:"orderBy"`
	Ascend         bool   `form:"ascend"`
}

type searchRestaurantResponseItem struct {
	RestaurantID   int    `json:"restaurantID"`
	RestaurantName string `json:"restaurantName"`
	ZipCode        int    `json:"zipCode"`
	RestaurantAddr string `json:"restaurantAddr"`
}

type searchRestaurantResponse []searchRestaurantResponseItem

type searchDishRequest struct {
	RestaurantID int `json:"restaurantID"`
}

type searchDishResponseItem struct {
	DishID   int     `json:"dishID"`
	DishName string  `json:"dishName"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type searchDishResponse []searchDishResponseItem

type placeOrderRequest struct {
	DishIDList []int `json:"dishIDList"`
	UserID     int   `json:"userID"`
}

type placeOrderResponseItem struct {
	OrderID int `json:"orderID"`
	RiderID int `json:"riderID"`
}

type placeOrderResponse []placeOrderResponseItem

type deleteOrderRequest struct {
	OrderID int `json:"orderID"`
}

type deleteOrderResponse struct {
	OrderID int    `json:"orderID"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type getCommentRequest struct {
	OrderID int `json:"orderID"`
}

type getCommentResponse struct {
	Rating  int    `json:"rating"`
	Content string `json:"content"`
}

type updateCommentRequest struct {
	OrderID int    `json:"orderID"`
	Rating  int    `json:"rating"`
	Content string `json:"content"`
}

type updateCommentResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type deleteCommentRequest struct {
	CommentID int `json:"commentID"`
}

type deleteCommentResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type historyOrderRequest struct {
	UserID   int     `json:"userID"`
	MaxPrice float64 `json:"maxPrice"`
	MinPrice float64 `json:"minPrice"`
	MaxTime  uint64  `json:"maxTime"`
	MinTime  uint64  `json:"minTime"`
	OrderBy  string  `json:"orderBy"`
	Ascend   bool    `json:"ascend"`
}

type historyOrderResponseItem struct {
	OrderID   int     `json:"orderID"`
	OrderTime uint64  `json:"orderTime"`
	DishPrice float64 `json:"dishPrice"`
	DishID    int     `json:"dishID"`
	RiderID   int     `json:"riderID"`
}
type historyOrderResponse []historyOrderResponseItem

type advancedCustomersResponseItem struct {
	UserID         int    `json:"userID"`
	UserName       string `json:"userName"`
	NumberOfOrders int    `json:"numberOfOrders"`
}

type advancedCustomersResponse []advancedCustomersResponseItem

type advancedRestaurantsResponseItem struct {
	RestaurantID   int     `json:"restaurantID"`
	RestaurantName string  `json:"restaurantName"`
	AvgPrice       float64 `json:"avgPrice"`
}

type advancedRestaurantsResponse []advancedRestaurantsResponseItem
