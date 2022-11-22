package debugger

import (
	"encoding/json"
	"errors"
	"fmt"
)

const MessageTypeEventDebugger MessageType = 1

type MessageType int

func (mt MessageType) IsEventDebugger() bool {
	return mt == MessageTypeEventDebugger
}

// TODO: Promote as WS, not only debugger

// Message ...
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (m Message) EventDebugger() (EventDebugger, error) {
	if !m.Type.IsEventDebugger() {
		return EventDebugger{}, errors.New("not a event debugger")
	}

	var evt EventDebugger
	if err := json.Unmarshal(m.Payload, &evt); err != nil {
		return EventDebugger{}, fmt.Errorf("unknown event format: %s", err)
	}

	return evt, nil
}

const (
	EventDebuggerTypeRequestReceived EventDebuggerType = iota + 1
	EventDebuggerTypeRequestContinued
	EventDebuggerTypeImposterReceived
	EventDebuggerTypeImposterContinued
	EventDebuggerTypeResponseReceived
	EventDebuggerTypeResponseContinued
)

type EventDebuggerType int

func (t EventDebuggerType) IsRequestReceived() bool {
	return t == EventDebuggerTypeRequestReceived
}
func (t EventDebuggerType) IsRequestContinued() bool {
	return t == EventDebuggerTypeRequestContinued
}
func (t EventDebuggerType) IsImposterFound() bool {
	return t == EventDebuggerTypeImposterReceived
}
func (t EventDebuggerType) IsImposterProcessed() bool {
	return t == EventDebuggerTypeImposterContinued
}
func (t EventDebuggerType) IsResponsePrepared() bool {
	return t == EventDebuggerTypeResponseReceived
}
func (t EventDebuggerType) IsResponseContinued() bool {
	return t == EventDebuggerTypeResponseContinued
}

// EventDebugger ...
type EventDebugger struct {
	Type     EventDebuggerType `json:"type"`
	Imposter []byte            `json:"imposter"`
	Request  []byte            `json:"request"`
	Response []byte            `json:"response"`
}

type Wait interface {
	Wait() EventDebugger
}

type debuggerWait struct {
	ch chan EventDebugger
}

func (w *debuggerWait) Wait() EventDebugger {
	return <-w.ch
}
