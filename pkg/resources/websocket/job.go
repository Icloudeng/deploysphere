package websocket

import (
	"encoding/json"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/pubsub"
)

func EmitJobEvent(job *database.Job) {
	playload := map[string]interface{}{
		"job_id": job.ID,
		"status": job.Status,
	}

	if playload_json, err := json.Marshal(playload); err == nil {
		EmitDecodedEvent(pubsub.NetworkEventPayload{
			Type:      pubsub.REDIS_EVENT_TYPE_JOBS,
			Channel:   job.Ref + "-" + pubsub.REDIS_EVENT_TYPE_JOBS,
			Reference: job.Ref,
			Payload:   string(playload_json),
		})
	}
}
