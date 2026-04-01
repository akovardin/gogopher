package main

import (
	"fmt"
	"testing"

	"gopher/gopher"

	"github.com/playdate-go/pdgo"
)

type mockClient struct {
	requestAccessFunc func(pdgo.AccessCallback) pdgo.AccessReply
	getDirectoryFunc  func(string, func([]gopher.Item), func(error))
}

func (m *mockClient) RequestAccess(callback pdgo.AccessCallback) pdgo.AccessReply {
	return m.requestAccessFunc(callback)
}

func (m *mockClient) GetDirectory(selector string, onItems func([]gopher.Item), onError func(error)) {
	m.getDirectoryFunc(selector, onItems, onError)
}

type mockView struct {
	dataFunc func(Data)
}

func (m *mockView) Data(d Data) {
	m.dataFunc(d)
}

type mockSystem struct {
	logFunc func(msg string)
}

func (m *mockSystem) LogToConsole(msg string) {
	m.logFunc(msg)
}

func TestNewController(t *testing.T) {
	view := &mockView{dataFunc: func(Data) {}}
	client := &mockClient{
		requestAccessFunc: func(pdgo.AccessCallback) pdgo.AccessReply { return pdgo.AccessDeny },
	}
	system := &mockSystem{logFunc: func(msg string) {}}

	c := NewController(view, client, system)
	if c == nil {
		t.Fatal("NewController returned nil")
	}
}

func TestPressA_AccessDenied(t *testing.T) {
	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	client := &mockClient{
		requestAccessFunc: func(callback pdgo.AccessCallback) pdgo.AccessReply {
			return pdgo.AccessDeny
		},
	}
	view := &mockView{dataFunc: func(d Data) {
		t.Error("Data should not be called when access is denied")
	}}

	c := NewController(view, client, system)
	c.PressA()

	if len(logged) == 0 {
		t.Error("expected LogToConsole to be called")
	}
}

func TestPressA_AccessAllowed_ReturnsItems(t *testing.T) {
	items := []gopher.Item{
		{Type: '0', Name: "File1", Selector: "/file1", Host: "host", Port: 70},
		{Type: '1', Name: "Dir1", Selector: "/dir1", Host: "host", Port: 70},
	}

	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	client := &mockClient{
		requestAccessFunc: func(callback pdgo.AccessCallback) pdgo.AccessReply {
			return pdgo.AccessAllow
		},
		getDirectoryFunc: func(selector string, onItems func([]gopher.Item), onError func(error)) {
			if selector != "" {
				t.Errorf("expected empty selector, got %q", selector)
			}
			onItems(items)
		},
	}

	var received Data
	view := &mockView{dataFunc: func(d Data) {
		received = d
	}}

	c := NewController(view, client, system)
	c.PressA()

	if len(received.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(received.Items))
	}
	if received.Items[0].Name != "File1" {
		t.Errorf("expected File1, got %s", received.Items[0].Name)
	}
}

func TestPressA_GetDirectoryError(t *testing.T) {
	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	client := &mockClient{
		requestAccessFunc: func(callback pdgo.AccessCallback) pdgo.AccessReply {
			return pdgo.AccessAllow
		},
		getDirectoryFunc: func(selector string, onItems func([]gopher.Item), onError func(error)) {
			onError(fmt.Errorf("connection failed"))
		},
	}

	view := &mockView{dataFunc: func(d Data) {
		t.Error("Data should not be called on error")
	}}

	c := NewController(view, client, system)
	c.PressA()

	found := false
	for _, msg := range logged {
		if msg == "Error: connection failed" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error log, got %v", logged)
	}
}

func TestPressA_AccessCallbackCalled(t *testing.T) {
	var callbackCalled bool
	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	client := &mockClient{
		requestAccessFunc: func(callback pdgo.AccessCallback) pdgo.AccessReply {
			callbackCalled = true
			callback(true)
			return pdgo.AccessAllow
		},
		getDirectoryFunc: func(selector string, onItems func([]gopher.Item), onError func(error)) {
			onItems(nil)
		},
	}

	view := &mockView{dataFunc: func(d Data) {}}

	c := NewController(view, client, system)
	c.PressA()

	if !callbackCalled {
		t.Error("expected access callback to be called")
	}
}

func TestPressA_AccessAsk_NoImmediateReply(t *testing.T) {
	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	client := &mockClient{
		requestAccessFunc: func(callback pdgo.AccessCallback) pdgo.AccessReply {
			return pdgo.AccessAsk
		},
	}

	view := &mockView{dataFunc: func(d Data) {
		t.Error("Data should not be called when access is not yet allowed")
	}}

	c := NewController(view, client, system)
	c.PressA()

	for _, msg := range logged {
		if msg == fmt.Sprintf("TCP access reply immediate: %d", int(pdgo.AccessAsk)) {
			t.Error("immediate reply log should not appear for AccessAsk")
		}
	}
}
