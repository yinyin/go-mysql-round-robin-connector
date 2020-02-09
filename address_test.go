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

func TestParseAddress02(t *testing.T) {
	r := parseAddress("loc-2/3")
	if r.locationName != "loc-2" {
		t.Errorf("unexpect location name: %s", r.locationName)
	}
	if r.orderedCount != 3 {
		t.Errorf("unexpect ordered count: %d", r.orderedCount)
	}
	if r.shuffleCount != -1 {
		t.Errorf("unexpect shuffle count: %d", r.shuffleCount)
	}
}

func TestParseAddress03a(t *testing.T) {
	r := parseAddress("loc-3/5/7")
	if r.locationName != "loc-3" {
		t.Errorf("unexpect location name: %s", r.locationName)
	}
	if r.orderedCount != 5 {
		t.Errorf("unexpect ordered count: %d", r.orderedCount)
	}
	if r.shuffleCount != 7 {
		t.Errorf("unexpect shuffle count: %d", r.shuffleCount)
	}
}

func TestParseAddress03b(t *testing.T) {
	r := parseAddress("loc-3//7")
	if r.locationName != "loc-3" {
		t.Errorf("unexpect location name: %s", r.locationName)
	}
	if r.orderedCount != -1 {
		t.Errorf("unexpect ordered count: %d", r.orderedCount)
	}
	if r.shuffleCount != 7 {
		t.Errorf("unexpect shuffle count: %d", r.shuffleCount)
	}
}

func TestParseAddress03c(t *testing.T) {
	r := parseAddress("loc-3/-/7")
	if r.locationName != "loc-3" {
		t.Errorf("unexpect location name: %s", r.locationName)
	}
	if r.orderedCount != -1 {
		t.Errorf("unexpect ordered count: %d", r.orderedCount)
	}
	if r.shuffleCount != 7 {
		t.Errorf("unexpect shuffle count: %d", r.shuffleCount)
	}
}

func TestCacheParsedAddress(t *testing.T) {
	r := parseAddress("loc-c/9/11")
	if (r.locationName != "loc-c") || (r.orderedCount != 9) || (r.shuffleCount != 11) {
		t.Errorf("unexpect result: (name: %s, ordered: %d, shuffle: %d)", r.locationName, r.orderedCount, r.shuffleCount)
	}
	if r2 := checkAddrCache("loc-c/9/11"); (r2 == nil) || (r.locationName != "loc-c") || (r.orderedCount != 9) || (r.shuffleCount != 11	) {
		t.Errorf("unexpect caching result: %#v", r2)
	}
}
