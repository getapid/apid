package check

// import (
// 	"github.com/iv-p/apid/svc/cli/http"
// )

// // Checker is the service that checks all the checks
// type Checker struct {
// 	client *http.Client
// }

// // Check is the internal representation of a check
// type Check struct {
// 	Name  string
// 	Steps []Step
// }

// type Step struct {
// 	Name     string
// 	Request  *Request
// 	Response *Response
// }

// type Request struct {
// 	Type     string
// 	Endpoint string
// 	Headers  Headers
// 	Body     string
// }

// type Headers map[string]string

// type Response struct {
// 	Headers Headers
// }

// type Result struct {
// 	Name      string
// 	Response  http.Response
// 	Timestamp int64
// }

// // NewChecker creates and returns a new Checker instance
// func NewChecker(client *http.Client) *Checker {
// 	return &Checker{client}
// }

// func Check() Result {

// }
