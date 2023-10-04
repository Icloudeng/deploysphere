package websocket

import (
	"encoding/base64"
	"encoding/json"

	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/http/ws"
	"github.com/icloudeng/platform-installer/internal/pubsub"
)

func EmitDecodedEvent(playload pubsub.NetworkEventPayload) {
	if json_data, err := json.Marshal(playload); err == nil {
		ws.Broadcast(json_data)
	}
}

func EmitEncodedEvent(playload pubsub.NetworkEventPayload) {
	data := playload
	decodedBytes, err := base64.StdEncoding.DecodeString(data.Payload)
	if err != nil {
		return
	}

	data.Payload = string(decodedBytes)

	if json_data, err := json.Marshal(data); err == nil {
		ws.Broadcast(json_data)
	}
}

func EmitJobEvent(job *entities.Job) {
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
