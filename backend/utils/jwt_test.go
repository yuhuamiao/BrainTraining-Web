package utils

import "testing"

func TestValidateJWTConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		secret  string
		wantErr bool
	}{
		{name: "debug allows local secret", mode: "debug", wantErr: false},
		{name: "release rejects missing secret", mode: "release", wantErr: true},
		{name: "release rejects development secret", mode: "release", secret: developmentJWTSecret, wantErr: true},
		{name: "release rejects short secret", mode: "release", secret: "too-short", wantErr: true},
		{name: "release accepts strong secret", mode: "release", secret: "0123456789abcdef0123456789abcdef", wantErr: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("GIN_MODE", test.mode)
			t.Setenv("JWT_SECRET", test.secret)
			err := ValidateJWTConfiguration()
			if (err != nil) != test.wantErr {
				t.Fatalf("ValidateJWTConfiguration() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
