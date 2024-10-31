package notification

type RequestData struct {
	Notification Notification
	Priority     string `json:"priority"`
	Data         Data
}

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"` // Berisi Message
}
type Data struct {
	Title          string `json:"title"`
	Message        string `json:"message"`
	Sound          string `json:"sound"`
	ClickAction    string `json:"click_action"`
	Type           string `json:"type"`
	KodeNotifikasi string `json:"kode_notifikasi"`
}

// constructor function
func (std *RequestData) DefaultValue() {
	std.Priority = "high"
	std.Data.Sound = "default"
	std.Data.ClickAction = "FLUTTER_NOTIFICATION_CLICK"
}
