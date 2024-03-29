package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestJobLabels(t *testing.T) {
	d, err := time.Parse(time.RFC3339, "2021-09-01T00:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}
	j := NewJob(JobID{}, []string{"a", "b"}, d)
	got := j.Labels().Slice()
	want := []string{"a", "b"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected difference: %v", diff)
	}

	if err = j.ClaimAt(d.Add(10 * time.Second)); err != nil {
		t.Fatalf("failed to claim job: %v", err)
	}
	if err = j.CloseAt(d.Add(20*time.Second), "a"); err != nil {
		t.Fatalf("failed to close job: %v", err)
	}
	if got, want := j.ResultLabel(), "a"; got != want {
		t.Errorf("unexpected result label: %v; want %v", got, want)
	}
}

func TestJobLabelsError(t *testing.T) {
	d, err := time.Parse(time.RFC3339, "2021-09-01T00:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}
	j := NewJob(JobID{}, []string{"a", "b"}, d)
	if err = j.ClaimAt(d.Add(10 * time.Second)); err != nil {
		t.Fatalf("failed to claim job: %v", err)
	}
	if err = j.CloseAt(d.Add(20*time.Second), "c"); err == nil {
		t.Fatalf("expected error on closing job with invalid result label")
	}
	if got, want := j.ResultLabel(), ""; got != want {
		t.Errorf("unexpected result label: %v; want %v", got, want)
	}
}
