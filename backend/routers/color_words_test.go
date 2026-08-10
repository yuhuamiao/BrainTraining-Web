package routers

import (
	"net/http/httptest"
	"testing"
)

func TestValidWebSocketOrigin(t *testing.T) {
	tests := []struct {
		name          string
		host          string
		origin        string
		allowedOrigin string
		want          bool
	}{
		{name: "rejects missing origin", host: "training.example.com", want: false},
		{name: "allows same origin", host: "training.example.com", origin: "https://training.example.com", want: true},
		{name: "rejects foreign origin", host: "training.example.com", origin: "https://evil.example.com", want: false},
		{name: "allows configured origin", host: "api.example.com", origin: "https://app.example.com", allowedOrigin: "https://app.example.com", want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("CORS_ALLOWED_ORIGIN", test.allowedOrigin)
			request := httptest.NewRequest("GET", "http://"+test.host+"/api/v1/color_words/color", nil)
			request.Host = test.host
			if test.origin != "" {
				request.Header.Set("Origin", test.origin)
			}
			if got := validWebSocketOrigin(request); got != test.want {
				t.Fatalf("validWebSocketOrigin() = %v, want %v", got, test.want)
			}
		})
	}
}
