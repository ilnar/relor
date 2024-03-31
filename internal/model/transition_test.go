package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidTransitionList(t *testing.T) {
	wid := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	startTs := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	rts := []RawTransition{
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000002"),
			WID:     wid,
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000001"), Valid: true},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000003"), Valid: true},
			From:    "b",
			To:      "c",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 10, 0, time.UTC),
		},
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000001"),
			WID:     wid,
			Prev:    uuid.NullUUID{},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Valid: true},
			From:    "a",
			To:      "b",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 5, 0, time.UTC),
		},
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000003"),
			WID:     wid,
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Valid: true},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000004"), Valid: true},
			From:    "c",
			To:      "a",
			Label:   "reset",
			Created: time.Date(2024, 4, 1, 0, 0, 20, 0, time.UTC),
		},
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000004"),
			WID:     wid,
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000003"), Valid: true},
			Next:    uuid.NullUUID{},
			From:    "a",
			To:      "d",
			Label:   "exit",
			Created: time.Date(2024, 4, 1, 0, 0, 40, 0, time.UTC),
		},
	}
	th, err := NewTransitionHistory(startTs, rts)
	if err != nil {
		t.Fatalf("failed to create transition history: %v", err)
	}
	want := []struct {
		id, wid         uuid.UUID
		from, to, label string
		walltime        time.Duration
	}{
		{
			id:       uuid.MustParse("11111111-0000-0000-0000-000000000001"),
			wid:      wid,
			from:     "a",
			to:       "b",
			label:    "ok",
			walltime: 5 * time.Second,
		},
		{
			id:       uuid.MustParse("11111111-0000-0000-0000-000000000002"),
			wid:      wid,
			from:     "b",
			to:       "c",
			label:    "ok",
			walltime: 5 * time.Second,
		},
		{
			id:       uuid.MustParse("11111111-0000-0000-0000-000000000003"),
			wid:      wid,
			from:     "c",
			to:       "a",
			label:    "reset",
			walltime: 10 * time.Second,
		},
		{
			id:       uuid.MustParse("11111111-0000-0000-0000-000000000004"),
			wid:      wid,
			from:     "a",
			to:       "d",
			label:    "exit",
			walltime: 20 * time.Second,
		},
	}
	for i, w := range want {
		t.Run(w.id.String(), func(t *testing.T) {
			if th == nil {
				t.Fatal("unexpected nil transition history")
			}
			if th.ID() != w.id {
				t.Errorf("unexpected id: %s; want %s", th.ID(), w.id)
			}
			if th.WorkflowID() != w.wid {
				t.Errorf("unexpected workflow id: %s; want %s", th.WorkflowID(), w.wid)
			}
			from, to := th.FromTo()
			if from != w.from {
				t.Errorf("unexpected from node: %s; want %s", from, w.from)
			}
			if to != w.to {
				t.Errorf("unexpected to node: %s; want %s", to, w.to)
			}
			if th.Label() != w.label {
				t.Errorf("unexpected label: %s; want %s", th.Label(), w.label)
			}
			if th.Walltime() != w.walltime {
				t.Errorf("unexpected walltime: %v; want %v", th.Walltime(), w.walltime)
			}
			if i < len(want)-1 {
				th = th.Next()
			} else if th.Next() != nil {
				t.Errorf("unexpected next transition: %v; want nil", th.Next())
			}
		})
	}
}

func TestEmptyHistory(t *testing.T) {
	startTs := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	th, err := NewTransitionHistory(startTs, nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if th != nil {
		t.Errorf("unexpected transition history: %v; want nil", th)
	}
}

func TestMissingTail(t *testing.T) {
	wid := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	startTs := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	rts := []RawTransition{
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000002"),
			WID:     wid,
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000001"), Valid: true},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000003"), Valid: true},
			From:    "b",
			To:      "c",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 10, 0, time.UTC),
		},
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000001"),
			WID:     wid,
			Prev:    uuid.NullUUID{},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Valid: true},
			From:    "a",
			To:      "b",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 5, 0, time.UTC),
		},
	}
	th, err := NewTransitionHistory(startTs, rts)
	t.Log(err)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if th != nil {
		t.Errorf("unexpected transition history: %v; want nil", th)
	}
}

func TestUnreachableNodes(t *testing.T) {
	wid := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	startTs := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	rts := []RawTransition{
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000003"),
			WID:     wid,
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Valid: true},
			Next:    uuid.NullUUID{},
			From:    "b",
			To:      "c",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 10, 0, time.UTC),
		},
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000001"),
			WID:     wid,
			Prev:    uuid.NullUUID{},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Valid: true},
			From:    "a",
			To:      "b",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 5, 0, time.UTC),
		},
	}
	th, err := NewTransitionHistory(startTs, rts)
	t.Log(err)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if th != nil {
		t.Errorf("unexpected transition history: %v; want nil", th)
	}
}

func TestMissingHead(t *testing.T) {
	wid := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	startTs := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	rts := []RawTransition{
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000002"),
			WID:     wid,
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000001"), Valid: true},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000003"), Valid: true},
			From:    "b",
			To:      "c",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 10, 0, time.UTC),
		},
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000003"),
			WID:     wid,
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Valid: true},
			Next:    uuid.NullUUID{},
			From:    "a",
			To:      "b",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 5, 0, time.UTC),
		},
	}
	th, err := NewTransitionHistory(startTs, rts)
	t.Log(err)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if th != nil {
		t.Errorf("unexpected transition history: %v; want nil", th)
	}
}

func TestDuplicateTransition(t *testing.T) {
	wid := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	startTs := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	rts := []RawTransition{
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000002"),
			WID:     wid,
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000001"), Valid: true},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000003"), Valid: true},
			From:    "b",
			To:      "c",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 10, 0, time.UTC),
		},
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000002"),
			WID:     wid,
			Prev:    uuid.NullUUID{},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Valid: true},
			From:    "a",
			To:      "b",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 5, 0, time.UTC),
		},
	}
	th, err := NewTransitionHistory(startTs, rts)
	t.Log(err)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if th != nil {
		t.Errorf("unexpected transition history: %v; want nil", th)
	}
}

func TestDifferentWorkflowIDs(t *testing.T) {
	startTs := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	rts := []RawTransition{
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000002"),
			WID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000001"), Valid: true},
			Next:    uuid.NullUUID{},
			From:    "b",
			To:      "c",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 10, 0, time.UTC),
		},
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000001"),
			WID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Prev:    uuid.NullUUID{},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Valid: true},
			From:    "a",
			To:      "b",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 5, 0, time.UTC),
		},
	}
	th, err := NewTransitionHistory(startTs, rts)
	t.Log(err)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if th != nil {
		t.Errorf("unexpected transition history: %v; want nil", th)
	}
}

func TestLoop(t *testing.T) {
	wid := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	startTs := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	rts := []RawTransition{
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000002"),
			WID:     wid,
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000001"), Valid: true},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000003"), Valid: true},
			From:    "b",
			To:      "c",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 10, 0, time.UTC),
		},
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000001"),
			WID:     wid,
			Prev:    uuid.NullUUID{},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Valid: true},
			From:    "a",
			To:      "b",
			Label:   "ok",
			Created: time.Date(2024, 4, 1, 0, 0, 5, 0, time.UTC),
		},
		{
			ID:      uuid.MustParse("11111111-0000-0000-0000-000000000003"),
			WID:     wid,
			Prev:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Valid: true},
			Next:    uuid.NullUUID{UUID: uuid.MustParse("11111111-0000-0000-0000-000000000001"), Valid: true},
			From:    "c",
			To:      "a",
			Label:   "reset",
			Created: time.Date(2024, 4, 1, 0, 0, 20, 0, time.UTC),
		},
	}
	th, err := NewTransitionHistory(startTs, rts)
	t.Log(err)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if th != nil {
		t.Errorf("unexpected transition history: %v; want nil", th)
	}
}
