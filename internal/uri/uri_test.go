package uri

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Parsed
		wantErr bool
	}{
		{
			name:  "valid URI",
			input: "kc://github/token",
			want:  Parsed{Service: "github", Key: "token"},
		},
		{
			name:  "valid URI with dots and dashes",
			input: "kc://my-service/api.key",
			want:  Parsed{Service: "my-service", Key: "api.key"},
		},
		{
			name:  "valid URI with underscores and numbers",
			input: "kc://aws_prod/access_key_1",
			want:  Parsed{Service: "aws_prod", Key: "access_key_1"},
		},
		{
			name:    "missing scheme",
			input:   "github/token",
			wantErr: true,
		},
		{
			name:    "wrong scheme",
			input:   "http://github/token",
			wantErr: true,
		},
		{
			name:    "missing key",
			input:   "kc://github",
			wantErr: true,
		},
		{
			name:    "missing key with trailing slash",
			input:   "kc://github/",
			wantErr: true,
		},
		{
			name:    "empty service",
			input:   "kc:///token",
			wantErr: true,
		},
		{
			name:    "nested key with slash",
			input:   "kc://github/org/token",
			wantErr: true,
		},
		{
			name:    "invalid characters in service",
			input:   "kc://my service/token",
			wantErr: true,
		},
		{
			name:    "invalid characters in key",
			input:   "kc://github/my token",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "just scheme",
			input:   "kc://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got != tt.want {
				t.Errorf("Parse(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsKCURI(t *testing.T) {
	if !IsKCURI("kc://github/token") {
		t.Error("expected true for kc://github/token")
	}
	if IsKCURI("http://example.com") {
		t.Error("expected false for http://example.com")
	}
	if IsKCURI("plain_value") {
		t.Error("expected false for plain_value")
	}
}

func TestFormat(t *testing.T) {
	got := Format("github", "token")
	want := "kc://github/token"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}
