package subscribers

import (
	"smatflow/platform-installer/pkg/pubsub"
	proxyhost "smatflow/platform-installer/pkg/resources/proxy_host"
	"smatflow/platform-installer/pkg/resources/utilities"
)

func EventSubscribers() {
	// Notifier
	pubsub.BusEvent.Subscribe(pubsub.RESOURCES_NOTIFIER_EVENT, utilities.SendNotification)

	// Remove Proxy Host
	pubsub.BusEvent.Subscribe(pubsub.RESOURCES_CLEANUP_EVENT, proxyhost.DeleteProxyHost)
}
