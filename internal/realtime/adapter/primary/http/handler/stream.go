package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// Stream serves the authenticated realtime SSE stream for the current user.
func (h *Handler) Stream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set(headerContentType, contentTypeEventStream)
	w.Header().Set(headerCacheControl, cacheControlNoCache)
	w.Header().Set(headerConnection, connectionKeepAlive)
	w.Header().Set(headerAccelBuffering, accelBufferingNo)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	events, cleanup := h.Service.Subscribe(ctx, userID)
	defer cleanup()

	h.Logger.InfowCtx(ctx, logRealtimeConnected, commonkeys.UserID, strconv.FormatUint(userID, 10))

	if err := writeSSE(w, sseEventConnected, map[string]any{
		"type":           "connected",
		"userId":         strconv.FormatUint(userID, 10),
		"connectedAtUTC": time.Now().UTC().Format(time.RFC3339Nano),
	}); err == nil {
		flusher.Flush()
	}

	heartbeat := time.NewTicker(h.heartbeatInterval())
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done():
			h.Logger.InfowCtx(context.Background(), logRealtimeDisconnected, commonkeys.UserID, strconv.FormatUint(userID, 10))
			return
		case event, ok := <-events:
			if !ok {
				return
			}
			if err := writeSSE(w, sseEventRecordProjectionChanged, event); err != nil {
				return
			}
			flusher.Flush()
		case <-heartbeat.C:
			if _, err := fmt.Fprint(w, sseCommentHeartbeat); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func writeSSE(w http.ResponseWriter, eventName string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "%s%s\n", sseEventPrefix, eventName); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%s%s\n\n", sseDataPrefix, body); err != nil {
		return err
	}
	return nil
}
