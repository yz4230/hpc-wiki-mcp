package cmd

import "testing"

func TestChooseServeMode(t *testing.T) {
	tests := []struct {
		name    string
		stdin   bool
		http    bool
		want    serveMode
		wantErr bool
	}{
		{name: "default HTTP", want: serveModeHTTP},
		{name: "explicit HTTP", http: true, want: serveModeHTTP},
		{name: "stdio", stdin: true, want: serveModeStdin},
		{name: "conflict", stdin: true, http: true, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := chooseServeMode(test.stdin, test.http)
			if (err != nil) != test.wantErr || got != test.want {
				t.Fatalf("chooseServeMode() = %q, %v; want %q, error=%v", got, err, test.want, test.wantErr)
			}
		})
	}
}

func TestResolveHTTPPort(t *testing.T) {
	t.Setenv(portEnvName, "18080")
	got, err := resolveHTTPPort(defaultHTTPPort)
	if err != nil || got != 18080 {
		t.Fatalf("resolveHTTPPort() = %d, %v", got, err)
	}

	t.Setenv(portEnvName, "invalid")
	if _, err := resolveHTTPPort(defaultHTTPPort); err == nil {
		t.Fatal("resolveHTTPPort() accepted an invalid PORT")
	}
}
