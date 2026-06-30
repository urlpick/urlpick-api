package config

import (
	"testing"
)

var requiredKeys = []string{
	"APP_PORT", "MYSQL_HOST", "MYSQL_PORT", "MYSQL_DATABASE",
	"MYSQL_USER", "MYSQL_PASSWORD", "BASE_URL", "TURNSTILE_SECRET_KEY",
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		set     map[string]string
		wantErr bool
	}{
		{
			name: "all present",
			set: map[string]string{
				"APP_PORT": "8080", "MYSQL_HOST": "db", "MYSQL_PORT": "3306",
				"MYSQL_DATABASE": "urlpick", "MYSQL_USER": "u", "MYSQL_PASSWORD": "p",
				"BASE_URL": "https://x.io", "TURNSTILE_SECRET_KEY": "secret",
			},
			wantErr: false,
		},
		{
			name:    "all missing",
			set:     map[string]string{},
			wantErr: true,
		},
		{
			name: "one missing",
			set: map[string]string{
				"APP_PORT": "8080", "MYSQL_HOST": "db", "MYSQL_PORT": "3306",
				"MYSQL_DATABASE": "urlpick", "MYSQL_USER": "u", "MYSQL_PASSWORD": "p",
				"BASE_URL": "https://x.io",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, k := range requiredKeys {
				t.Setenv(k, "")
			}
			for k, v := range tc.set {
				t.Setenv(k, v)
			}

			cfg, err := Parse()
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if cfg.AppPort != tc.set["APP_PORT"] {
				t.Fatalf("expected AppPort %q, got %q", tc.set["APP_PORT"], cfg.AppPort)
			}
		})
	}
}
