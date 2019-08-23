package http

type Response struct {
	StatusCode int
	Timings    Timings
	Headers    Headers
	Body       string
}
