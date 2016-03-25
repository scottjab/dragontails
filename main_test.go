package main

import "testing"

func TestParseEvent(t *testing.T) {
	passive := "[1457038519] PASSIVE SERVICE CHECK: s0137;fs_/boot;0;OK - 30.4% used (0.07 of 0.2 GB), (levels at 92.00/95.00%), trend: 0.00B / 24 hours"
	event, err := parseEvent(passive)
	if err != nil {
		t.Error(
			"Got Error", err,
			"Expected Parsed Event",
		)
	}
	if event.Host != "s0137" {
		t.Error(
			"Host parse error, Got", event.Host,
			"Expected s0137",
		)
	}
	if event.State != "ok" {
		t.Error(
			"State parse error, Got", event.State,
			"expected ok",
		)
	}
}
