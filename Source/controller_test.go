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
	getFileFunc       func(string, func([]byte), func(error))
}

func (m *mockClient) RequestAccess(callback pdgo.AccessCallback) pdgo.AccessReply {
	return m.requestAccessFunc(callback)
}

func (m *mockClient) GetDirectory(selector string, onItems func([]gopher.Item), onError func(error)) {
	m.getDirectoryFunc(selector, onItems, onError)
}

func (m *mockClient) GetFile(selector string, onData func([]byte), onError func(error)) {
	m.getFileFunc(selector, onData, onError)
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

func TestPressB_AccessDenied(t *testing.T) {
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
	c.PressB()

	if len(logged) == 0 {
		t.Error("expected LogToConsole to be called")
	}
}

func TestPressB_AccessAllowed_ReturnsItems(t *testing.T) {
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
	c.PressB()

	if len(received.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(received.Items))
	}
	if received.Items[0].Name != "File1" {
		t.Errorf("expected File1, got %s", received.Items[0].Name)
	}
}

func TestPressB_GetDirectoryError(t *testing.T) {
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
	c.PressB()

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

func TestPressB_AccessCallbackCalled(t *testing.T) {
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
	c.PressB()

	if !callbackCalled {
		t.Error("expected access callback to be called")
	}
}

func TestPressB_AccessAsk_NoImmediateReply(t *testing.T) {
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
	c.PressB()

	for _, msg := range logged {
		if msg == fmt.Sprintf("TCP access reply immediate: %d", int(pdgo.AccessAsk)) {
			t.Error("immediate reply log should not appear for AccessAsk")
		}
	}
}

func TestPressA_DirectoryNavigation(t *testing.T) {
	items := []gopher.Item{
		{Type: '0', Name: "File1", Selector: "/file1", Host: "host", Port: 70},
	}

	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	cursor := gopher.Item{Type: '1', Name: "Dir1", Selector: "/dir1", Host: "host", Port: 70}

	client := &mockClient{
		getDirectoryFunc: func(selector string, onItems func([]gopher.Item), onError func(error)) {
			if selector != "/dir1" {
				t.Errorf("expected selector /dir1, got %q", selector)
			}
			onItems(items)
		},
	}

	var received Data
	view := &mockView{dataFunc: func(d Data) {
		received = d
	}}

	c := NewController(view, client, system)
	c.PressA(cursor)

	if len(received.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(received.Items))
	}
	if received.Items[0].Name != "File1" {
		t.Errorf("expected File1, got %s", received.Items[0].Name)
	}
}

func TestPressA_FileRetrieval(t *testing.T) {
	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	cursor := gopher.Item{Type: '0', Name: "File1", Selector: "/file1", Host: "host", Port: 70}
	fileData := []byte("hello gopher")

	client := &mockClient{
		getFileFunc: func(selector string, onData func([]byte), onError func(error)) {
			if selector != "/file1" {
				t.Errorf("expected selector /file1, got %q", selector)
			}
			onData(fileData)
		},
	}

	var received Data
	view := &mockView{dataFunc: func(d Data) {
		received = d
	}}

	c := NewController(view, client, system)
	c.PressA(cursor)

	if received.File != "hello gopher" {
		t.Errorf("expected 'hello gopher', got %q", received.File)
	}
}

func TestPressA_GetDirectoryError(t *testing.T) {
	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	cursor := gopher.Item{Type: '1', Name: "Dir1", Selector: "/dir1", Host: "host", Port: 70}

	client := &mockClient{
		getDirectoryFunc: func(selector string, onItems func([]gopher.Item), onError func(error)) {
			onError(fmt.Errorf("connection failed"))
		},
	}

	view := &mockView{dataFunc: func(d Data) {
		t.Error("Data should not be called on error")
	}}

	c := NewController(view, client, system)
	c.PressA(cursor)

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

func TestPressA_GetFileError(t *testing.T) {
	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	cursor := gopher.Item{Type: '0', Name: "File1", Selector: "/file1", Host: "host", Port: 70}

	client := &mockClient{
		getFileFunc: func(selector string, onData func([]byte), onError func(error)) {
			onError(fmt.Errorf("file not found"))
		},
	}

	view := &mockView{dataFunc: func(d Data) {
		t.Error("Data should not be called on error")
	}}

	c := NewController(view, client, system)
	c.PressA(cursor)

	found := false
	for _, msg := range logged {
		if msg == "Error: file not found" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error log, got %v", logged)
	}
}

func TestPressA_LogsSelector(t *testing.T) {
	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	cursor := gopher.Item{Type: '1', Selector: "/some/path"}

	client := &mockClient{
		getDirectoryFunc: func(selector string, onItems func([]gopher.Item), onError func(error)) {
			onItems(nil)
		},
	}

	view := &mockView{dataFunc: func(d Data) {}}

	c := NewController(view, client, system)
	c.PressA(cursor)

	found := false
	for _, msg := range logged {
		if msg == "Pressed A: /some/path" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'Pressed A: /some/path' log, got %v", logged)
	}
}

func TestPressA_UnknownTypeDoesNothing(t *testing.T) {
	var logged []string
	system := &mockSystem{logFunc: func(msg string) { logged = append(logged, msg) }}

	cursor := gopher.Item{Type: '7', Name: "Search", Selector: "/search", Host: "host", Port: 70}

	client := &mockClient{}

	view := &mockView{dataFunc: func(d Data) {
		t.Error("Data should not be called for unknown cursor type")
	}}

	c := NewController(view, client, system)
	c.PressA(cursor)
}
