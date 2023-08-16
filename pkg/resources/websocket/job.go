package websocket

import (
	"encoding/base64"
	"encoding/json"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/events"
	"smatflow/platform-installer/pkg/events/redis_events"
)

func EmitJobEvent(job *database.Job) {
	playload := map[string]interface{}{
		"job_id": job.ID,
		"status": job.Status,
	}

	if playload_json, err := json.Marshal(playload); err == nil {
		EmitRedisEvent(redis_events.ResourceRedisEventPayload{
			Type:      events.REDIS_EVENT_TYPE_JOBS,
			Channel:   job.Ref + "-" + events.REDIS_EVENT_TYPE_JOBS,
			Reference: job.Ref,
			Payload:   base64.StdEncoding.EncodeToString(playload_json),
		})
	}
}
