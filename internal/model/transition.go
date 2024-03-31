package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Transition struct {
	id, workflowID          uuid.UUID
	fromNode, toNode, label string
	createdAt               time.Time
	walltime                time.Duration
	previous, next          *Transition
}

func (t *Transition) ID() uuid.UUID {
	return t.id
}

func (t *Transition) WorkflowID() uuid.UUID {
	return t.workflowID
}

func (t *Transition) FromTo() (string, string) {
	return t.fromNode, t.toNode
}

func (t *Transition) Label() string {
	return t.label
}

func (t *Transition) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Transition) Walltime() time.Duration {
	return t.walltime
}

func (t *Transition) Next() *Transition {
	return t.next
}

type tnHistoryBuilder struct {
	startTime  time.Time
	head, tail *Transition
	idx        map[uuid.UUID]*Transition
}

type RawTransition struct {
	ID, WID         uuid.UUID
	Prev, Next      uuid.NullUUID
	From, To, Label string
	Created         time.Time
}

func NewTransitionHistory(startTime time.Time, rts []RawTransition) (*Transition, error) {
	if len(rts) == 0 {
		return nil, fmt.Errorf("empty transition list")
	}
	b := &tnHistoryBuilder{
		startTime: startTime,
		idx:       make(map[uuid.UUID]*Transition),
	}
	for _, rt := range rts {
		if err := b.add(rt); err != nil {
			return nil, fmt.Errorf("failed to add transition %v: %w", rt, err)
		}
	}
	head, err := b.build()
	if err != nil {
		return nil, fmt.Errorf("failed to build history: %w", err)
	}
	return head, nil
}

func (th *tnHistoryBuilder) add(rt RawTransition) error {
	if _, found := th.idx[rt.ID]; found {
		return fmt.Errorf("transition %s already exists", rt.ID)
	}
	t := &Transition{
		id:         rt.ID,
		workflowID: rt.WID,
		fromNode:   rt.From,
		toNode:     rt.To,
		label:      rt.Label,
		createdAt:  rt.Created,
	}
	th.idx[rt.ID] = t

	if rt.Prev.Valid {
		if pt, found := th.idx[rt.Prev.UUID]; found {
			pt.next = t
			t.previous = pt
		}
	} else {
		if th.head == nil {
			th.head = t
		} else {
			return fmt.Errorf("head already set")
		}
	}
	if rt.Next.Valid {
		if nt, found := th.idx[rt.Next.UUID]; found {
			nt.previous = t
			t.next = nt
		}
	} else {
		if th.tail == nil {
			th.tail = t
		} else {
			return fmt.Errorf("tail already set")
		}
	}
	return nil
}

func (th *tnHistoryBuilder) build() (*Transition, error) {
	if th.head == nil {
		return nil, fmt.Errorf("head not set")
	}
	t := th.head
	prevTs := th.startTime
	cnt := len(th.idx)
	wid := t.workflowID
	for t != nil && cnt > 0 {
		t.walltime = t.createdAt.Sub(prevTs)
		prevTs = t.createdAt
		if t.workflowID != wid {
			return nil, fmt.Errorf("workflow ID mismatch: %s != %s", t.workflowID, wid)
		}
		t = t.next
		cnt--
	}
	if cnt != 0 {
		return nil, fmt.Errorf("corrupted history on forward pass: %d unreachable transitions detected", cnt)
	}
	if t != nil {
		return nil, fmt.Errorf("tail not reached on forward pass, possible loop detected; stopped at %v", t)
	}

	t = th.tail
	cnt = len(th.idx)
	for t != nil && cnt > 0 {
		t = t.previous
		cnt--
	}

	if cnt != 0 {
		return nil, fmt.Errorf("corrupted history on backward pass: %d unreachable transitions detected", cnt)
	}
	if t != nil {
		return nil, fmt.Errorf("head not reached on backward pass, possible loop detected; stopped at %v", t)
	}
	return th.head, nil
}
