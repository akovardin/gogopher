package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

func TestFail(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "simple error message",
			message: "File not found",
		},
		{
			name:    "unknown selector",
			message: "Unknown selector",
		},
		{
			name:    "empty message",
			message: "",
		},
		{
			name:    "message with spaces",
			message: "Something went wrong",
		},
		{
			name:    "message with special chars",
			message: "error: file.txt not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := net.Pipe()
			defer server.Close()
			defer client.Close()

			done := make(chan struct{})
			var got string

			go func() {
				defer close(done)
				reader := bufio.NewReader(client)
				line, err := reader.ReadString('\n')
				if err != nil {
					t.Errorf("failed to read line: %v", err)
					return
				}
				got = line

				term, err := reader.ReadString('\n')
				if err != nil {
					t.Errorf("failed to read terminator: %v", err)
					return
				}
				got += term
			}()

			fail(server, tt.message)

			<-done

			expected := "3" + tt.message + "\tfake\tfake\t0\r\n.\r\n"
			if got != expected {
				t.Errorf("fail() wrote %q, want %q", got, expected)
			}
		})
	}
}

func TestFail_ResponseFormat(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	done := make(chan struct{})
	var lines []string

	go func() {
		defer close(done)
		scanner := bufio.NewScanner(client)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	}()

	fail(server, "test error")
	server.Close()

	<-done

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}

	// First line: error response
	first := lines[0]
	if !strings.HasPrefix(first, "3") {
		t.Errorf("expected line to start with '3' (error type), got %q", first)
	}
	if !strings.Contains(first, "test error") {
		t.Errorf("expected line to contain 'test error', got %q", first)
	}
	if !strings.Contains(first, "\tfake\tfake\t0") {
		t.Errorf("expected line to contain '\\tfake\\tfake\\t0', got %q", first)
	}

	// Second line: terminator
	second := lines[1]
	if second != "." {
		t.Errorf("expected second line to be '.', got %q", second)
	}
}

func TestFail_WritesToConn(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	done := make(chan struct{})
	var buf []byte

	go func() {
		defer close(done)
		reader := bufio.NewReader(client)
		for {
			b, err := reader.ReadByte()
			if err != nil {
				return
			}
			buf = append(buf, b)
		}
	}()

	fail(server, "msg")
	server.Close()

	<-done

	expected := "3msg\tfake\tfake\t0\r\n.\r\n"
	if string(buf) != expected {
		t.Errorf("fail() wrote %q, want %q", string(buf), expected)
	}
}
