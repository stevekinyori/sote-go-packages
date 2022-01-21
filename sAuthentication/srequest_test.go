package sAuthentication

import (
	"reflect"
	"testing"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
)

var stagingExpToken = "eyJraWQiOiJlOCt4TW4rOGYrZmlIXC9OZDNDZGNxOVRvU3FPKzdZYldcL1wvSUxCYVJyTElNPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJjYzJiNGYwYS0xYmE0LTQzNTEtYmRmMS0wYmM3NTRhN2NlNjMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xX2E1MjFiZTA5LTQxMDQtNDc1MC1iZTQwLTQ2NTExYzczYzA2MCIsImNvZ25pdG86Z3JvdXBzIjpbIjEwMDM2Il0sImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS13ZXN0LTEuYW1hem9uYXdzLmNvbVwvZXUtd2VzdC0xX2ZwWkN5cWxRTiIsImNsaWVudF9pZCI6IjNlMzN0NGVjb2Vpam5ibTNscGVoZmNuaWcwIiwiZXZlbnRfaWQiOiIwOGVkNjhmZC1hZDk4LTRiMjEtYjNhZC0wN2U0MTk0YzFkOTkiLCJ0b2tlbl91c2UiOiJhY2Nlc3MiLCJzY29wZSI6ImF3cy5jb2duaXRvLnNpZ25pbi51c2VyLmFkbWluIiwiYXV0aF90aW1lIjoxNjIzODcwNTIzLCJleHAiOjE2MjM4NzA4MjMsImlhdCI6MTYyMzg3MDUyNCwianRpIjoiMjFmNWViZmItNDIxMC00MWU1LTllNTAtMWI2ZmNlOWExZDA0IiwidXNlcm5hbWUiOiJjYzJiNGYwYS0xYmE0LTQzNTEtYmRmMS0wYmM3NTRhN2NlNjMifQ.jngSSScEF3Nm1-XghY6-TGljngjKHk8LaXvBAofjLJB1NzYVh3CEk8UC4JdPswjhqc0_xdEclI6XRxHyM_44uxhFzFyzpHU39x1ly5ROfls_rSvucxgVbHGoaOkmtenOOlBwFarszgzcCdngkY_rujr-r_YhIrAYH-Y1JnpxZ_idPBzgfQ9W1O0UPmFOZFxyDPsC2kAGU0Zl_cBcPrqXMULywx6nc8Vxo3sqDTy1UFE80Ysyz8d_SmeCDihL-YXa2Lm4p_wGnk54rUJ5VG_eZJ6n5FEopdhXoPqY0lvkuWz01NbwtunPKENGTA-qi-o__UBHfpYcOA7Nx5HkI48j9Q"

func TestValidateBody(t *testing.T) {
	type args struct {
		data         []byte
		tEnvironment string
	}
	tests := []struct {
		name    string
		args    args
		want    *RequestHeaderSchema
		wantErr bool
	}{
		{
			name: "negative case: access missing file",
			args: args{
				data: []byte(`{
					"aws-user-name": "soteuser",
					"organizations-id": 10003,
					"device-id": 123456789
				}`),
				tEnvironment: sConfigParams.STAGING,
			},
			wantErr: true,
		},
		{
			name: "negative case: missing aws-user-name",
			args: args{
				data: []byte(`{
					"aws-user-name": "",
					"organizations-id": 10003,
					"device-id": 123456789
				}`),
				tEnvironment: sConfigParams.STAGING,
			},
			wantErr: true,
		},
		{
			name: "negative case: missing organizations-id",
			args: args{
				data: []byte(`{
					"aws-user-name": "soteuser",
					"device-id": 123456789
				}`),
				tEnvironment: sConfigParams.STAGING,
			},
			wantErr: true,
		},
		{
			name: "negative case: invalid kid",
			args: args{
				data: []byte(`{
					"json-web-token": "` + stagingExpToken + `",
					"aws-user-name": "soteuser",
					"organizations-id": 10003
				}`),
				tEnvironment: sConfigParams.STAGING,
			},
			wantErr: true,
		},
		{
			name: "negative case: invalid json web token",
			args: args{
				data: []byte(`{
					"json-web-token": "eyJraWQiOvxxx",
					"aws-user-name": "soteuser",
					"organizations-id": 10003
				}`),
				tEnvironment: sConfigParams.STAGING,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ValidateBody(tt.args.data, tt.args.tEnvironment)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
