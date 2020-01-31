package mysqlroundrobinconnector

import (
	"testing"
)

func TestParseAddress01(t *testing.T) {
	r := parseAddress("loc-1")
	if r.locationName != "loc-1" {
		t.Errorf("unexpect location name: %s", r.locationName)
	}
	if r.orderedCount != -1 {
		t.Errorf("unexpect ordered count: %d", r.orderedCount)
	}
	if r.shuffleCount != -1 {
		t.Errorf("unexpect shuffle count: %d", r.shuffleCount)
	}
}
