package iputil

import (
	"fmt"
	"strconv"
	"strings"
)

func IPToID(ip string) (uint32, error) {
	ip = strings.TrimSpace(ip)

	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid IP format: expected 4 octets, got %d", len(parts))
	}

	var octets [4]uint32

	for i, part := range parts {
		num, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid octet %d: %w", i+1, err)
		}

		if num > 255 {
			return 0, fmt.Errorf("octet %d out of range: %d (must be 0-255)", i+1, num)
		}

		octets[i] = uint32(num)
	}
	ipID := (16777216 * octets[0]) + (65536 * octets[1]) + (256 * octets[2]) + octets[3]

	return ipID, nil
}
