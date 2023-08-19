package jobs

import (
	"smatflow/platform-installer/pkg/pubsub/redis_events"
	"smatflow/platform-installer/pkg/resources/db"
	"smatflow/platform-installer/pkg/resources/websocket"
)

func redis_pub_listeners(Ref string) func() {
	// Redis Event Logs Listeners
	close1 := redis_events.ResourceProviningLogsEvents(
		Ref,
		[]redis_events.SubscriberFunc{
			db.Jobs.Job_ListenResourceProviningLogs,
			websocket.EmitEncodedEvent,
		},
	)

	// Redis Event Status Listeners
	close2 := redis_events.ResourceProviningStatusEvents(
		Ref,
		[]redis_events.SubscriberFunc{
			db.Jobs.Job_ListenResourceProviningStatus,
			websocket.EmitEncodedEvent,
		},
	)

	// Redis Event Credentials Listeners
	close3 := redis_events.ResourceProviningCredentialsEvents(
		Ref,
		[]redis_events.SubscriberFunc{
			db.ResourceState.ResourceState_ListenResourceProviningCredentials,
			websocket.EmitEncodedEvent,
		},
	)

	return func() {
		close1()
		close3()
		close2()
	}
}
