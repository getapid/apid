package runner

import (
	"fmt"
	"os"

	"github.com/getapid/apid/log"
	"github.com/getapid/apid/spec"
	"github.com/getapid/apid/step"
	"github.com/getapid/apid/variables"
	"github.com/getapid/apid/writer"
)

type job struct {
	name string
	spec spec.Spec
}

type ParallelSpecRunner struct {
	runner      step.Runner
	writer      writer.Writer
	parallelism int
}

func NewParallelSpecRunner(parallelism int, runner step.Runner, writer writer.Writer) ParallelSpecRunner {
	if parallelism < 1 {
		fmt.Printf("parallelism level is negative or zero ( %d ), exiting.\n", parallelism)
		os.Exit(1)
	}
	return ParallelSpecRunner{
		runner:      runner,
		writer:      writer,
		parallelism: parallelism,
	}
}

func (r *ParallelSpecRunner) Run(specs []spec.Spec) bool {
	var jobs []job
	for _, spec := range specs {
		jobs = append(jobs, job{name: spec.Name, spec: spec})
	}
	jobsChan := make(chan job)
	resultsChan := make(chan writer.Result)

	workers := r.parallelism
	if workers > len(jobs) {
		workers = len(jobs)
	}
	for w := 0; w < workers; w++ {
		go r.worker(jobsChan, resultsChan)
	}

	go func() {
		for _, job := range jobs {
			jobsChan <- job
		}
		close(jobsChan)
	}()

	success := true
	for a := 0; a < len(jobs); a++ {
		res := <-resultsChan
		r.writer.Write(res)
		for _, step := range res.Steps {
			success = success && step.Pass
		}
	}
	close(resultsChan)

	return success
}

func (r *ParallelSpecRunner) worker(jobs <-chan job, results chan<- writer.Result) {
	for j := range jobs {
		results <- r.doJob(j)
	}
}

func (r *ParallelSpecRunner) doJob(j job) writer.Result {
	log.L.Infof("running spec %s", j.name)
	vars := variables.New()

	var err error
	var results []step.Result
	var res step.Result

	for _, step := range j.spec.Steps {
		log.L.Infof("running step %s", step.Name)

		var exported variables.Variables
		res, exported, err = r.runner.Run(step, vars)
		if err != nil {
			res.FailedChecks = append(res.FailedChecks, fmt.Sprintf("error while making request %v", err))
		}
		results = append(results, res)

		if err != nil || !res.Pass {
			break
		}
		vars = vars.Merge(exported)
	}
	return writer.Result{Name: j.spec.Name, Steps: results}

}
