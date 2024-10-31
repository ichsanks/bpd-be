package producer

import "gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/event/model"

// Producer represents an event producer interface.
type Producer interface {
	Publish(request model.PublishRequest)
}
