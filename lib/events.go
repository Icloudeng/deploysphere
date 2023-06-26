package lib

import (
	proxyhost "smatflow/platform-installer/lib/resources/proxy_host"

	"github.com/asaskevich/EventBus"
)

var BusEvent EventBus.Bus

// Events
const RESOURCES_CLEANUP_EVENT = "resources:cleanup"

func init() {
	BusEvent = EventBus.New()

	// Remove Proxy Host
	BusEvent.Subscribe(RESOURCES_CLEANUP_EVENT, proxyhost.DeleteProxyHost)
}
