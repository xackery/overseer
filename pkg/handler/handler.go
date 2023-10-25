package handler

var (
	windowCloseHandler []func(cancelled *bool, reason byte)
)

// WindowCloseSubscribe allows subscribing to close events
func WindowCloseSubscribe(fn func(cancelled *bool, reason byte)) {
	windowCloseHandler = append(windowCloseHandler, fn)
}

// WindowCloseInvoke invokes close events on the window
func WindowCloseInvoke(cancelled *bool, reason byte) {
	for _, fn := range windowCloseHandler {
		fn(cancelled, reason)
	}
}
