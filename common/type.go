package common

type SearchRestaurantRequest struct {
	RestaurantName string `form:"restaurantName"`
	ZipCode        int    `form:"zipCode"`
	OrderBy        string `form:"orderBy"`
	Ascend         bool   `form:"ascend"`
}

type SearchRestaurantResponseItem struct {
	RestaurantID   int    `json:"restaurantID"`
	RestaurantName string `json:"restaurantName"`
	ZipCode        int    `json:"zipCode"`
	RestaurantAddr string `json:"restaurantAddr"`
}

type SearchRestaurantResponse []SearchRestaurantResponseItem

type searchDishRequest struct {
	RestaurantID int `json:"restaurantID"`
}

type searchDishResponseItem struct {
	DishID   int     `json:"dishID"`
	DishName string  `json:"dishName"`
	Price    float64 `json:"price"`
}

type searchDishResponse []searchDishResponseItem

type placeOrderRequest struct {
	DishIDList []int `json:"dishIDList"`
	UserID     int   `json:"userID"`
}

type placeOrderResponseItem struct {
	OrderID int  `json:"orderID"`
	Success bool `json:"success"`
}

type placeOrderResponse []placeOrderResponseItem

type deleteOrderRequest struct {
	OrderID int `json:"orderID"`
}

type deleteOrderResponse struct {
	Success bool `json:"success"`
}

type addCommentRequest struct {
	OrderID int    `json:"orderID"`
	Rating  int    `json:"rating"`
	Content string `json:"content"`
}

type addCommentResponse struct {
	CommentID int  `json:"commentID"`
	Success   bool `json:"success"`
}

type updateCommentRequest struct {
	CommentID int    `json:"commentID"`
	Rating    int    `json:"rating"`
	Content   string `json:"content"`
}

type updateCommentResponse struct {
	Success bool `json:"success"`
}

type deleteCommentRequest struct {
	CommentID int `json:"commentID"`
}

type deleteCommentResponse struct {
	Success bool `json:"success"`
}

type historyOrderRequest struct {
	UserID   int     `json:"userID"`
	MaxPrice float64 `json:"maxPrice"`
	MinPrice float64 `json:"minPrice"`
	MaxTime  uint64  `json:"maxTime"`
	MinTime  uint64  `json:"minTime"`
}

type historyOrderResponseItem struct {
	OrderID   int     `json:"orderID"`
	OrderTime uint64  `json:"orderTime"`
	DishPrice float64 `json:"dishPrice"`
	DishID    int     `json:"dishID"`
	RiderID   int     `json:"riderID"`
}

type historyOrderResponse []historyOrderResponseItem
