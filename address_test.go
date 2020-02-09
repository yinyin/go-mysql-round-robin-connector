package mysqlroundrobinconnector

import (
	"testing"
)

func TestParseAddress01a(t *testing.T) {
	r := parseAddress("loc-1")
	if r.locationName != "loc-1" {
		t.Errorf("unexpect location name: %s", r.locationName)
	}
	if r.orderedCount != -1 {
		t.Errorf("unexpect order count: %d", r.orderedCount)
	}
}

func TestParseAddress01b(t *testing.T) {
	r := parseAddress("loc-1/-")
	if r.locationName != "loc-1" {
		t.Errorf("unexpect location name: %s", r.locationName)
	}
	if r.orderedCount != 0 {
		t.Errorf("unexpect order count: %d", r.orderedCount)
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
}

func TestParseAddress03a(t *testing.T) {
	r := parseAddress("loc-3/5/7")
	if r.locationName != "loc-3" {
		t.Errorf("unexpect location name: %s", r.locationName)
	}
	if r.orderedCount != 5 {
		t.Errorf("unexpect ordered count: %d", r.orderedCount)
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
}

func TestParseAddress03c(t *testing.T) {
	r := parseAddress("loc-3/-/7")
	if r.locationName != "loc-3" {
		t.Errorf("unexpect location name: %s", r.locationName)
	}
	if r.orderedCount != 0 {
		t.Errorf("unexpect ordered count: %d", r.orderedCount)
	}
}

func TestCacheParsedAddress(t *testing.T) {
	r := parseAddress("loc-c/9/11")
	if (r.locationName != "loc-c") || (r.orderedCount != 9) {
		t.Errorf("unexpect result: (name: %s, ordered: %d)", r.locationName, r.orderedCount)
	}
	if r2 := checkAddrCache("loc-c/9/11"); (r2 == nil) || (r.locationName != "loc-c") || (r.orderedCount != 9) {
		t.Errorf("unexpect caching result: %#v", r2)
	}
}
