package http

type Headers map[string]string

type Request struct {
	Method  string
	Url     string
	Headers Headers
	Body    string
}
