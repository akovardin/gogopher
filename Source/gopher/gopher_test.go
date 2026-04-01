package gopher

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Item
	}{
		{
			name:  "file type",
			input: "0About\t/about\tgopher.floodgap.com\t70",
			want:  Item{Type: '0', Name: "About", Selector: "/about", Host: "gopher.floodgap.com", Port: 70},
		},
		{
			name:  "directory type",
			input: "1Floodgap Home\t/h/\tgopher.floodgap.com\t70",
			want:  Item{Type: '1', Name: "Floodgap Home", Selector: "/h/", Host: "gopher.floodgap.com", Port: 70},
		},
		{
			name:  "info type",
			input: "iThis is an info line\tfake\terror.host\t1",
			want:  Item{Type: 'i', Name: "This is an info line", Selector: "fake", Host: "error.host", Port: 1},
		},
		{
			name:  "search type",
			input: "7Search\t/search\tgopher.example.com\t70",
			want:  Item{Type: '7', Name: "Search", Selector: "/search", Host: "gopher.example.com", Port: 70},
		},
		{
			name:  "error type",
			input: "3Error message\tfake\thost\t1",
			want:  Item{Type: '3', Name: "Error message", Selector: "fake", Host: "host", Port: 1},
		},
		{
			name:  "empty string",
			input: "",
			want:  Item{},
		},
		{
			name:  "no tabs",
			input: "0singlefield",
			want:  Item{Type: '0', Name: "singlefield"},
		},
		{
			name:  "only type byte",
			input: "0",
			want:  Item{Type: '0'},
		},
		{
			name:  "name with spaces",
			input: "0A file with spaces.txt\t/files/file with spaces.txt\thost\t70",
			want:  Item{Type: '0', Name: "A file with spaces.txt", Selector: "/files/file with spaces.txt", Host: "host", Port: 70},
		},
		{
			name:  "port zero",
			input: "0test\t/test\tgopher.example.com\t0",
			want:  Item{Type: '0', Name: "test", Selector: "/test", Host: "gopher.example.com", Port: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLine(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLine() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSplitTabs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "multiple fields",
			input: "About\t/about\tgopher.floodgap.com\t70",
			want:  []string{"About", "/about", "gopher.floodgap.com", "70"},
		},
		{
			name:  "single field",
			input: "About",
			want:  []string{"About"},
		},
		{
			name:  "empty string",
			input: "",
			want:  []string{""},
		},
		{
			name:  "leading tab",
			input: "\tfield",
			want:  []string{"", "field"},
		},
		{
			name:  "trailing tab",
			input: "field\t",
			want:  []string{"field", ""},
		},
		{
			name:  "consecutive tabs",
			input: "a\t\tb",
			want:  []string{"a", "", "b"},
		},
		{
			name:  "empty fields",
			input: "\t",
			want:  []string{"", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitTabs(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitTabs() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAtoi(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{name: "positive number", input: "70", want: 70},
		{name: "zero", input: "0", want: 0},
		{name: "empty string", input: "", want: 0},
		{name: "non-numeric", input: "abc", want: 0},
		{name: "mixed", input: "70abc", want: 70},
		{name: "leading non-numeric", input: "abc70", want: 70},
		{name: "large number", input: "65535", want: 65535},
		{name: "single digit", input: "7", want: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := atoi(tt.input)
			if got != tt.want {
				t.Errorf("atoi() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestProcessDir(t *testing.T) {
	tests := []struct {
		name string
		buf  []byte
		want []Item
	}{
		{
			name: "multiple items",
			buf:  []byte("0About\t/about\tgopher.floodgap.com\t70\r\n1Software\t/software\tgopher.floodgap.com\t70\r\n"),
			want: []Item{
				{Type: '0', Name: "About", Selector: "/about", Host: "gopher.floodgap.com", Port: 70},
				{Type: '1', Name: "Software", Selector: "/software", Host: "gopher.floodgap.com", Port: 70},
			},
		},
		{
			name: "empty buffer",
			buf:  []byte{},
			want: []Item{},
		},
		{
			name: "single item",
			buf:  []byte("0File\t/file\thost\t70\r\n"),
			want: []Item{
				{Type: '0', Name: "File", Selector: "/file", Host: "host", Port: 70},
			},
		},
		{
			name: "trailing newline",
			buf:  []byte("0One\t/one\thost\t70\r\n1Two\t/two\thost\t70\r\n"),
			want: []Item{
				{Type: '0', Name: "One", Selector: "/one", Host: "host", Port: 70},
				{Type: '1', Name: "Two", Selector: "/two", Host: "host", Port: 70},
			},
		},
		{
			name: "info lines mixed",
			buf:  []byte("iInformation line\r\n0File.txt\t/file.txt\thost\t70\r\n"),
			want: []Item{
				{Type: 'i', Name: "Information line"},
				{Type: '0', Name: "File.txt", Selector: "/file.txt", Host: "host", Port: 70},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			got := c.processDir(tt.buf)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processDir() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFindCRLF(t *testing.T) {
	tests := []struct {
		name string
		buf  []byte
		want int
	}{
		{
			name: "CRLF at start",
			buf:  []byte("\r\nhello"),
			want: 0,
		},
		{
			name: "CRLF in middle",
			buf:  []byte("hello\r\nworld"),
			want: 5,
		},
		{
			name: "no CRLF",
			buf:  []byte("hello"),
			want: -1,
		},
		{
			name: "empty buffer",
			buf:  []byte{},
			want: -1,
		},
		{
			name: "CR without LF",
			buf:  []byte("hello\rworld"),
			want: -1,
		},
		{
			name: "LF without CR",
			buf:  []byte("hello\nworld"),
			want: -1,
		},
		{
			name: "only CRLF",
			buf:  []byte("\r\n"),
			want: 0,
		},
		{
			name: "multiple CRLFs returns first",
			buf:  []byte("ab\r\ncd\r\n"),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findCRLF(tt.buf)
			if got != tt.want {
				t.Errorf("findCRLF() = %d, want %d", got, tt.want)
			}
		})
	}
}
