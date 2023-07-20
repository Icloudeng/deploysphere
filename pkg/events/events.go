package events

import (
	proxyhost "smatflow/platform-installer/pkg/resources/proxy_host"
	"smatflow/platform-installer/pkg/resources/utilities"

	"github.com/asaskevich/EventBus"
)

var BusEvent EventBus.Bus

// Events
const RESOURCES_CLEANUP_EVENT = "resources:cleanup"
const NOTIFIER_RESOURCES_EVENT = "resources:utilities:notifier"

func init() {
	BusEvent = EventBus.New()

	// Notifier
	BusEvent.Subscribe(NOTIFIER_RESOURCES_EVENT, utilities.SendNotification)

	// Remove Proxy Host
	BusEvent.Subscribe(RESOURCES_CLEANUP_EVENT, proxyhost.DeleteProxyHost)
}
