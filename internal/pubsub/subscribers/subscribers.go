package subscribers

import (
	"smatflow/platform-installer/internal/pubsub"
	proxyhost "smatflow/platform-installer/internal/resources/proxy_host"
	"smatflow/platform-installer/internal/resources/utilities"
)

func EventSubscribers() {
	// Notifier
	pubsub.BusEvent.Subscribe(pubsub.RESOURCES_NOTIFIER_EVENT, utilities.SendNotification)

	// Remove Proxy Host
	pubsub.BusEvent.Subscribe(pubsub.RESOURCES_CLEANUP_EVENT, proxyhost.DeleteProxyHost)
}
