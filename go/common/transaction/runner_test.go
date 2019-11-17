package transaction

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iv-p/apid/common/mock"
	"github.com/iv-p/apid/common/result"
	"github.com/stretchr/testify/suite"

	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/variables"
)

var (
	rootVars = variables.New()
	txVars   = map[string]interface{}{
		"url": "one",
	}
	stepVars = map[string]interface{}{
		"url": "two",
	}

	errStepErr = errors.New("error")

	okStep = step.Step{
		ID:        "ok-step",
		Variables: stepVars,
		Export:    step.Export{},
	}
	okStepResult = step.Result{
		Step: step.PreparedStep(okStep),
		Valid: step.ValidationResult{
			Errors: map[string]string{},
		},
		Exported: step.Exported{},
	}

	errStep = step.Step{
		ID:        "err-step",
		Variables: stepVars,
		Export:    step.Export{},
	}
	errStepResult = step.Result{
		Step: step.PreparedStep(errStep),
		Valid: step.ValidationResult{
			Errors: map[string]string{
				"error-one": "this is why",
			},
		},
		Exported: step.Exported{},
	}
)

type RunnerSuite struct {
	suite.Suite
}

func (s *RunnerSuite) TestTransactionRunner_Run() {
	type args struct {
		transactions []Transaction
		vars         variables.Variables
	}
	type want struct {
		Results []result.TransactionResult
	}
	tests := []struct {
		name string
		args args
		want want
		ok   bool
	}{
		{
			"empty",
			args{
				transactions: []Transaction{},
				vars:         rootVars,
			},
			want{
				Results: []result.TransactionResult{},
			},
			true,
		},
		{
			"single tx without steps",
			args{
				[]Transaction{
					{
						"test-id-1",
						txVars,
						[]step.Step{},
					},
				},
				rootVars,
			},
			want{
				Results: []result.TransactionResult{
					{Id: "test-id-1"},
				},
			},
			true,
		},
		{
			"multiple tx with steps",
			args{
				[]Transaction{
					{
						ID:        "test-id-2",
						Variables: txVars,
						Steps: []step.Step{
							okStep,
							okStep,
						},
					},
					{
						ID:        "test-id-3",
						Variables: txVars,
						Steps: []step.Step{
							okStep,
							okStep,
							okStep,
							okStep,
						},
					},
					{
						ID:        "test-id-4",
						Variables: txVars,
						Steps: []step.Step{
							okStep,
							okStep,
							okStep,
						},
					},
				},
				rootVars,
			},
			want{
				[]result.TransactionResult{
					{
						Id: "test-id-2",
						Steps: []step.Result{
							okStepResult,
							okStepResult,
						},
					},
					{
						Id: "test-id-3",
						Steps: []step.Result{
							okStepResult,
							okStepResult,
							okStepResult,
							okStepResult,
						},
					},
					{
						Id: "test-id-4",
						Steps: []step.Result{
							okStepResult,
							okStepResult,
							okStepResult,
						},
					},
				},
			},
			true,
		},
		{
			"multiple tx with steps and errs",
			args{
				[]Transaction{
					{
						ID:        "test-id-5",
						Variables: txVars,
						Steps: []step.Step{
							okStep,
							okStep,
						},
					},
					{
						ID:        "test-id-6",
						Variables: txVars,
						Steps: []step.Step{
							okStep,
							errStep,
							okStep,
							okStep,
						},
					},
					{
						ID:        "test-id-7",
						Variables: txVars,
						Steps: []step.Step{
							okStep,
							okStep,
							okStep,
						},
					},
				},
				rootVars,
			},
			want{
				[]result.TransactionResult{
					{
						Id: "test-id-5",
						Steps: []step.Result{
							okStepResult,
							okStepResult,
						},
					},
					{
						Id: "test-id-6",
						Steps: []step.Result{
							okStepResult,
							errStepResult,
						},
					},
					{
						Id: "test-id-7",
						Steps: []step.Result{
							okStepResult,
							okStepResult,
							okStepResult,
						},
					},
				},
			},
			false,
		},
	}
	mockCtrl := gomock.NewController(s.T())
	defer mockCtrl.Finish()

	for _, tt := range tests {
		stepRunner := mock.NewMockRunner(mockCtrl)
		var writerCalls []*gomock.Call
		for _, tx := range tt.args.transactions {
			exported := variables.New()
			for _, step := range tx.Steps {
				if step.ID == okStep.ID {
					writerCalls = append(writerCalls,
						stepRunner.EXPECT().
							Run(step, variables.New(
								variables.WithOther(rootVars),
								variables.WithVars(txVars),
								variables.WithOther(exported),
								variables.WithVars(step.Variables),
							)).
							Return(okStepResult, nil))
					exported = variables.New(
						variables.WithOther(exported),
						variables.WithRaw(
							map[string]interface{}{
								step.ID: okStepResult.Exported,
							},
						),
					)
				} else {
					writerCalls = append(writerCalls,
						stepRunner.EXPECT().
							Run(step, variables.New(
								variables.WithOther(rootVars),
								variables.WithVars(txVars),
								variables.WithOther(exported),
								variables.WithVars(step.Variables),
							)).
							Return(errStepResult, errStepErr))
					break
				}
			}
		}
		gomock.InOrder(writerCalls...)

		writer := mock.NewMockWriter(mockCtrl)
		writerCalls = []*gomock.Call{}
		for _, txResult := range tt.want.Results {
			writerCalls = append(writerCalls,
				writer.EXPECT().
					Write(txResult).
					Return())
		}

		writerCalls = append(writerCalls, writer.EXPECT().Close())
		gomock.InOrder(writerCalls...)

		r := &TransactionRunner{
			stepRunner: stepRunner,
			writer:     writer,
		}
		ok := r.Run(tt.args.transactions, tt.args.vars)
		s.Equal(tt.ok, ok, tt.name)
	}
}

func TestRunnerSuite(t *testing.T) {
	suite.Run(t, new(RunnerSuite))
}
