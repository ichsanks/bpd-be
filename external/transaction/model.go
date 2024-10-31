package transaction

// ErrorReponse is a format error from transaction domain
type ErrorReponse struct {
	Error string `json:"error"`
}

type CreateTransactionFormatResponse struct {
	Data struct {
		TotalItemInserted int            `json:"totalItemInserted"`
		ShipmentData      []ShipmentData `json:"shipmentData"`
	} `json:"data"`
}

type ShipmentData struct {
	CourierName string `json:"courierName"`
	TrackingID  string `json:"trackingId"`
	OrderID     string `json:"orderId"`
}
