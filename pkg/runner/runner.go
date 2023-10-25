package runner

// Runner wraps an executable and provides a way manage it's output
type Runner interface {
	Start()
	Stop() error
	PID() int
}
