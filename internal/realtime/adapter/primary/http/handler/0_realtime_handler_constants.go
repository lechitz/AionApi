package handler

const (
	streamRoute = "/events/stream"

	headerContentType    = "Content-Type"
	headerCacheControl   = "Cache-Control"
	headerConnection     = "Connection"
	headerAccelBuffering = "X-Accel-Buffering"

	contentTypeEventStream = "text/event-stream"
	cacheControlNoCache    = "no-cache"
	connectionKeepAlive    = "keep-alive"
	accelBufferingNo       = "no"

	sseEventConnected               = "connected"
	sseEventRecordProjectionChanged = "record_projection_changed"
	sseDataPrefix                   = "data: "
	sseEventPrefix                  = "event: "
	sseCommentHeartbeat             = ": keepalive\n\n"

	logRealtimeConnected    = "realtime stream connected"
	logRealtimeDisconnected = "realtime stream disconnected"
)
