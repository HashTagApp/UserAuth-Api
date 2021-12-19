package response

import (
	"context"
	"encoding/json"
	"net/http"
)

type PingResponse struct {
	status string `json:status`
}

func EncodePing(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	res := response.(string)
	jsonR := PingResponse{}
	jsonR.status = res
	w.Header().Set("content-Type", "application/json")
	if ctx.Value(`X-Token-Refresh`) == true {
		w.Header().Set("X-Token-Refresh", "true")
	} else {
		w.Header().Set("X-Token-Refresh", "false")
	}
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}
