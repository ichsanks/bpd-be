package notification

import (
	"fmt"

	"github.com/NaySoftware/go-fcm"
)

func PushNotification(title string, body string, data map[string]interface{}, token []string) (err error) {
	var NP fcm.NotificationPayload
	NP.Title = title
	NP.Body = body

	c := fcm.NewFcmClient("AAAAlI1S77U:APA91bEzWnxGfL_UjjJMJBthW_YJ4lsaH6deq8lj7AF3ikOEjm91NaicXZs9dc4U9BFFf5EY02up_9Hi9VRcQu8ikriQvIEsiMmyJ9enw7ibuxFRsiaPK-icoC0npQGTyqthoNNYvLEJ")
	c.NewFcmRegIdsMsg(token, data)
	//c.AppendDevices(xds)
	c.SetNotificationPayload(&NP)
	status, err := c.Send()
	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}
	return
}
