package endpoint

import "testing"

func Test_apidEndpointProvider_GetForRegion(t *testing.T) {
	provider := NewAPIDEndpointProvider()

	type args struct {
		region string
	}
	tests := []struct {
		name    string
		p       EndpointProvider
		args    args
		want    string
		wantErr bool
	}{
		{
			"legit",
			provider,
			args{
				"us-east",
			},
			"https://use.api.getapid.com/executor",
			false,
		},
		{
			"missing",
			provider,
			args{
				"us-wrong-1",
			},
			"",
			true,
		},
		{
			"empty",
			provider,
			args{
				"",
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := apidEndpointProvider{}
			got, err := p.GetForRegion(tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("apidEndpointProvider.GetForRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("apidEndpointProvider.GetForRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}
