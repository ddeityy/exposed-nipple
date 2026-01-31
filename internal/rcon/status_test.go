package rcon

import "testing"

const testString = `
hostname: Spire Server
version : 10336138/24 10336138 secure
udp/ip  : 127.0.0.1:69696  (local: 172.17.0.2:27015)  (public IP from Steam: 127.0.0.1)
steamid : [A:1:922269719:48625] (90280836254447639)
account : not logged in  (No account specified)
map     : cp_badlands at: 0 x, 0 y, 0 z
tags    : cp,nocrits
sourcetv:  127.0.0.1:69696, delay 90.0s  (local: 172.17.0.2:27020)
players : 1 humans, 1 bots (25 max)
edicts  : 416 used of 2048 max
# userid name                uniqueid            connected ping loss state  adr
#      2 "Spire Server TV"   BOT                                     active
#      4 "Deity. 2026 edition" [U:1:115754284]   00:33      165    0 active
`

func TestParseStatus(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected Status
	}{
		{
			name:  "default",
			input: testString,
			expected: Status{
				Hostname: "Spire Server",
				SDR:      "127.0.0.1:69696",
				Map:      "cp_badlands",
				Players:  1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseStatus(tc.input)
			if got != tc.expected {
				t.Errorf("got %v, expected %v", got, tc.expected)
			}
		})
	}
}
