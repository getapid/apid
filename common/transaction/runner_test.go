package transaction

import (
	"errors"
	"testing"

	"github.com/getapid/apid-cli/common/mock"
	"github.com/getapid/apid-cli/common/result"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/getapid/apid-cli/common/step"
	"github.com/getapid/apid-cli/common/variables"
)

var (
	rootVars = variables.New()
	txVars   = newVarsWithVars(map[string]interface{}{
		"url": "one",
	})
	stepVars = newVarsWithVars(map[string]interface{}{
		"url": "two",
	})

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
						ID:        "test-id-1",
						Variables: txVars,
						Steps:     []step.Step{},
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
		var mockedCalls []*gomock.Call
		for _, tx := range tt.args.transactions {
			exported := variables.New()
			for _, step := range tx.Steps {
				if step.ID == okStep.ID {
					mockedCalls = append(mockedCalls,
						stepRunner.EXPECT().
							Run(step, rootVars.DeepCopy().
								Merge(txVars).
								Merge(exported).
								Merge(step.Variables),
							).
							Return(okStepResult, nil))
					exported = variables.New(
						variables.WithOther(exported),
						variables.WithRaw(
							map[string]interface{}{
								step.ID: okStepResult.Exported.Generic(),
							},
						),
					)
				} else {
					mockedCalls = append(mockedCalls,
						stepRunner.EXPECT().
							Run(step, rootVars.DeepCopy().
								Merge(txVars).
								Merge(exported).
								Merge(step.Variables),
							).
							Return(errStepResult, errStepErr))
					break
				}
			}
		}
		gomock.InOrder(mockedCalls...)

		writer := mock.NewMockWriter(mockCtrl)
		mockedCalls = []*gomock.Call{}
		for _, txResult := range tt.want.Results {
			mockedCalls = append(mockedCalls,
				writer.EXPECT().
					Write(txResult).
					Return())
		}

		mockedCalls = append(mockedCalls, writer.EXPECT().Close())
		gomock.InOrder(mockedCalls...)

		r := &TransactionRunner{
			stepRunner: stepRunner,
			writer:     writer,
		}
		ok := r.Run(tt.args.transactions, tt.args.vars)
		s.Equal(tt.ok, ok, tt.name)
	}
}

func (s *RunnerSuite) TestTransactionRunner_RunWithMatrix() {
	transaction := Transaction{
		ID:        "test-with-matrix",
		Variables: variables.Variables{},
		Steps:     []step.Step{okStep},
		Matrix: &Matrix{
			M: map[string][]interface{}{
				"var1": {1, 2},
				"var2": {"a", "b"},
			},
		},
	}
	matrixClone := Matrix{
		M: transaction.Matrix.M,
	}
	expectedVarSets := make([]variables.Variables, 0, 4)
	for matrixClone.HasNext() {
		expectedVarSets = append(expectedVarSets, matrixClone.NextSet())
	}

	mockCtrl := gomock.NewController(s.T())
	defer mockCtrl.Finish()
	mockStepRunner := mock.NewMockRunner(mockCtrl)
	mockWriter := mock.NewMockWriter(mockCtrl)

	for _, oneExpSet := range expectedVarSets {
		mockStepRunner.EXPECT().
			Run(okStep, stepVars.DeepCopy().Merge(oneExpSet)).
			Return(okStepResult, nil).
			Times(1)
		mockWriter.EXPECT().
			Write(gomock.Any()).
			Times(1)
	}
	mockWriter.EXPECT().Close()

	r := &TransactionRunner{
		stepRunner: mockStepRunner,
		writer:     mockWriter,
	}
	ok := r.Run([]Transaction{transaction}, rootVars)
	s.True(ok)
}

func newVarsWithVars(m map[string]interface{}) variables.Variables {
	return variables.New(variables.WithVars(m))
}

func TestRunnerSuite(t *testing.T) {
	suite.Run(t, new(RunnerSuite))
}
