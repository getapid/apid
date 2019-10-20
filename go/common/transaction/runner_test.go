package transaction

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	commonmock "github.com/iv-p/apid/common/mock"
	"github.com/iv-p/apid/common/result"
	climock "github.com/iv-p/apid/svc/cli/mock"

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
	}
	okStepResult = step.Result{
		Step: step.PreparedStep(okStep),
		Valid: step.ValidationResult{
			Errors: map[string]string{},
		},
	}

	errStep = step.Step{
		ID:        "err-step",
		Variables: stepVars,
	}
	errStepResult = step.Result{
		Step: step.PreparedStep(errStep),
		Valid: step.ValidationResult{
			Errors: map[string]string{
				"error-one": "this is why",
			},
		},
	}
)

func TestTransactionRunner_Run(t *testing.T) {
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
				[]Transaction{},
				rootVars,
			},
			want{
				[]result.TransactionResult{},
			},
			true,
		},
		{
			"single tx without steps",
			args{
				[]Transaction{
					{
						"test-id",
						txVars,
						[]step.Step{},
					},
				},
				rootVars,
			},
			want{
				[]result.TransactionResult{
					{},
				},
			},
			true,
		},
		{
			"multiple tx with steps",
			args{
				[]Transaction{
					{
						"test-id",
						txVars,
						[]step.Step{
							okStep,
							okStep,
						},
					},
					{
						"test-id",
						txVars,
						[]step.Step{
							okStep,
							okStep,
							okStep,
							okStep,
						},
					},
					{
						"test-id",
						txVars,
						[]step.Step{
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
						[]step.Result{
							okStepResult,
							okStepResult,
						},
					},
					{
						[]step.Result{
							okStepResult,
							okStepResult,
							okStepResult,
							okStepResult,
						},
					},
					{
						[]step.Result{
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
						"test-id",
						txVars,
						[]step.Step{
							okStep,
							okStep,
						},
					},
					{
						"test-id",
						txVars,
						[]step.Step{
							okStep,
							errStep,
							okStep,
							okStep,
						},
					},
					{
						"test-id",
						txVars,
						[]step.Step{
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
						[]step.Result{
							okStepResult,
							okStepResult,
						},
					},
					{
						[]step.Result{
							okStepResult,
							errStepResult,
						},
					},
					{
						[]step.Result{
							okStepResult,
							okStepResult,
						},
					},
				},
			},
			false,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stepRunner := commonmock.NewMockRunner(mockCtrl)
			var runs []*gomock.Call
			for _, tx := range tt.args.transactions {
				for _, step := range tx.Steps {
					if step.ID == okStep.ID {
						runs = append(runs,
							stepRunner.EXPECT().
								Run(step, rootVars.Merge(variables.New(variables.WithVars(step.Variables)))).
								Return(okStepResult))
					} else {
						runs = append(runs,
							stepRunner.EXPECT().
								Run(step, rootVars.Merge(variables.New(variables.WithVars(step.Variables)))).
								Return(errStepResult))
						break
					}
				}
			}
			gomock.InOrder(runs...)

			writer := climock.NewMockWriter(mockCtrl)
			runs = []*gomock.Call{}
			for _, txResult := range tt.want.Results {
				runs = append(runs,
					writer.EXPECT().
						Write(txResult).
						Return())
			}
			gomock.InOrder(runs...)

			r := &TransactionRunner{
				stepRunner: stepRunner,
				writer:     writer,
			}
			if ok := r.Run(tt.args.transactions, tt.args.vars); ok != tt.ok {
				t.Errorf("TransactionRunner.Run() error = %v, wantErr %v", ok, tt.ok)
			}
		})
	}
}
