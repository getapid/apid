package interpolator

import (
	"encoding/json"
	"testing"
)

func TestSimpleStringInterpolator_Interpolate(t *testing.T) {
	type args struct {
		str       string
		variables []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"should resolve all", args{"{{ one }} {{ two }}", []byte(`{"one": "test", "two": "much"}`)}, "test much", false},
		{"should resolve all", args{"one", []byte(``)}, "one", false},
		{"should resolve just one", args{"{{ one }} {{ not_existent }}", []byte(`{"one": "test", "two": "much"}`)}, "test ", false},
		{"should resolve nested key", args{"{{ one.two }}", []byte(`{"one": {"two": "much"} }`)}, "much", false},
		{"should resolve array", args{"{{ one[0] }}", []byte(`{"one": [ "test" ] }`)}, "test", false},
		{"shouldn't resolve array", args{"{{ two[0] }}", []byte(`{"one": [ "test" ] }`)}, "", false},
		{"shouldn't resolve array", args{"{{ two.[0] }}", []byte(`{"one": [ "test" ] }`)}, "", false},
		{"should resolve array int", args{"{{ one[0] }}", []byte(`{"one": [ 1, 2.01, 3 ] }`)}, "1", false},
		{"should resolve array float", args{"{{ one[1] }}", []byte(`{"one": [ 1, 2.01, 3 ] }`)}, "2.01", false},
		{"shouldn't resolve array out of bounds", args{"{{ one[-1] }}", []byte(`{"one": [ 1, 2.01, 3 ] }`)}, "", false},
		{"shouldn't resolve array out of bounds", args{"{{ one[4] }}", []byte(`{"one": [ 1, 2.01, 3 ] }`)}, "", false},
		{"should resolve array string", args{"{{ one[2] }}", []byte(`{"one": [ 1, 2.01, "3" ] }`)}, "3", false},
		{"should resolve array and nested key", args{"{{ one[0].test }}", []byte(`{"one": [ { "test": "much" } ] }`)}, "much", false},
		{"should resolve array and nested key", args{"{{ one.[0].test }}", []byte(`{"one": [ { "test": "much" } ] }`)}, "much", false},
		{"should resolve root array", args{"{{ [0].test }}", []byte(`[ { "test": "much" } ]`)}, "much", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSimpleStringInterpolator()
			var data interface{}
			err := json.Unmarshal(tt.args.variables, &data)
			got, err := s.Interpolate(tt.args.str, data)
			if (err != nil) != tt.wantErr {
				t.Errorf("SimpleStringInterpolator.Interpolate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SimpleStringInterpolator.Interpolate() = %v, want %v", got, tt.want)
			}
		})
	}
}
