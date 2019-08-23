package check

// type Submitter struct {
// 	api     string
// 	results chan *Result
// }

// func NewSubmitter(results chan *Result) *Submitter {
// 	return &Submitter{
// 		results: results,
// 	}
// }

// func (s *Submitter) Start() {
// 	for {
// 		select {
// 		case r := <-s.results:
// 			r.Timings.Print()
// 		}
// 	}
// }

// func (s *Submitter) submitResult(result *Result) {
// }
