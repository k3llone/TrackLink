package shared

import "testing"

func TestIsUUID(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "canonical lowercase",
			value: "7607f3ca-90d7-4c47-b2f7-f968ad1f5f9a",
			want:  true,
		},
		{
			name:  "canonical uppercase",
			value: "7607F3CA-90D7-4C47-B2F7-F968AD1F5F9A",
			want:  true,
		},
		{
			name:  "empty",
			value: "",
			want:  false,
		},
		{
			name:  "malformed",
			value: "not-a-uuid",
			want:  false,
		},
		{
			name:  "non hex",
			value: "7607f3ca-90d7-4c47-b2f7-f968ad1f5f9z",
			want:  false,
		},
		{
			name:  "wrong segment lengths",
			value: "7607f3ca-90d7-4c47-b2f7f-968ad1f5f9a",
			want:  false,
		},
		{
			name:  "braced",
			value: "{7607f3ca-90d7-4c47-b2f7-f968ad1f5f9a}",
			want:  false,
		},
		{
			name:  "urn",
			value: "urn:uuid:7607f3ca-90d7-4c47-b2f7-f968ad1f5f9a",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUUID(tt.value); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
