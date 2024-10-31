package shipment

// ErrorReponse is a format error from transaction domain
type ErrorReponse struct {
	Error string `json:"error"`
}

type RequestRegisterTracking struct {
	CourierCode    string `json:"courierCode"`
	ExternalID     string `json:"externalId"`
	PickupProvider string `json:"pickupProvider"`
	TrackingID     string `json:"trackingId"`
}
