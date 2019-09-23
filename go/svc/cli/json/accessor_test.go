package json

import "testing"

func Test_getKeyAndIndex(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int
		want2 bool
	}{
		{"test", args{"afjnadfa[1]"}, "afjnadfa", 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := getKeyAndIndex(tt.args.key)
			if got != tt.want {
				t.Errorf("getKeyAndIndex() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getKeyAndIndex() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("getKeyAndIndex() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
