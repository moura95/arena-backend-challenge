package iputil

import (
	"testing"
)

func TestIPToID(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		want    uint32
		wantErr bool
	}{
		{
			name:    "valid IP - 8.8.8.8",
			ip:      "8.8.8.8",
			want:    134744072,
			wantErr: false,
		},
		{
			name:    "valid IP - 192.168.1.1",
			ip:      "192.168.1.1",
			want:    3232235777,
			wantErr: false,
		},
		{
			name:    "valid IP - 0.0.0.0",
			ip:      "0.0.0.0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "valid IP - 255.255.255.255",
			ip:      "255.255.255.255",
			want:    4294967295,
			wantErr: false,
		},
		{
			name:    "valid IP - 10.0.0.1",
			ip:      "10.0.0.1",
			want:    167772161,
			wantErr: false,
		},
		{
			name:    "valid IP with spaces",
			ip:      "  8.8.8.8  ",
			want:    134744072,
			wantErr: false,
		},
		{
			name:    "invalid IP - empty string",
			ip:      "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid IP - missing octets",
			ip:      "192.168.1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid IP - too many octets",
			ip:      "192.168.1.1.1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid IP - non-numeric octet",
			ip:      "192.168.a.1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid IP - octet out of range",
			ip:      "192.168.256.1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid IP - negative octet",
			ip:      "192.168.-1.1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid IP - float octet",
			ip:      "192.168.1.1.5",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IPToID(tt.ip)

			if (err != nil) != tt.wantErr {
				t.Errorf("IPToID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("IPToID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIPToID(b *testing.B) {
	testIP := "192.168.1.1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IPToID(testIP)
	}
}
