package template

import (
	"reflect"
	"testing"

	"github.com/getapid/apid/variables"
)

var (
	nilData = variables.New()
	data    = map[string]interface{}{
		"array": []interface{}{"value"},
		"one":   "two",
		"nested": map[string]interface{}{
			"key":   "three",
			"array": []interface{}{"four"},
		},
		"int":   3,
		"float": 3.14,
	}
)

func TestGet(t *testing.T) {
	type args struct {
		template string
		data     variables.Variables
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"only text",
			args{
				"text",
				data,
			},
			"text",
			false,
		},
		{
			"simple root variable",
			args{
				"{{ one }}",
				data,
			},
			"two",
			false,
		},
		{
			"simple array variable",
			args{
				"{{ array.0 }}",
				data,
			},
			"value",
			false,
		},
		{
			"nested variable",
			args{
				"{{ nested.array.0 }}",
				data,
			},
			"four",
			false,
		},
		{
			"text variable",
			args{
				"pre {{ nested.key}} post",
				data,
			},
			"pre three post",
			false,
		},
		{
			"nil data only text",
			args{
				"text",
				nilData,
			},
			"text",
			false,
		},

		{
			"nil data variable text",
			args{
				"text {{ key }}",
				nilData,
			},
			"",
			true,
		},

		{
			"command text",
			args{
				"text {% echo '1' %}",
				nilData,
			},
			"text 1",
			false,
		},

		{
			"command with env text",
			args{
				"text {% echo $ONE %}",
				data,
			},
			"text two",
			false,
		},

		{
			"command with env branch text",
			args{
				"text {% echo $NESTED %}",
				data,
			},
			"text {\"array\":[\"four\"],\"key\":\"three\"}",
			false,
		},

		{
			"int in text",
			args{
				"text {{ int }}",
				data,
			},
			"text 3",
			false,
		},

		{
			"float in text",
			args{
				"text {{ float }}",
				data,
			},
			"text 3.14",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.args.template, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
