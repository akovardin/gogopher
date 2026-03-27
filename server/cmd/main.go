package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	port    = "7070"
	host    = "localhost"
	dataDir = "data"
)

type menuEntry struct {
	Type     byte
	Display  string
	Selector string
}

var rootMenu = []menuEntry{
	{'i', "Привет всем любителям старины!", ""},
	{'i', "", ""},
	{'i', "", ""},
	{'0', "О закерах (фантастический рассказ)", "/about_zakers.txt"},
	{'1', "Ещё рассказы", "/stories"},
}

var storiesMenu = []menuEntry{
	{'i', "Фантастические рассказы", ""},
	{'i', "", ""},
	{'i', "", ""},
	{'0', "Рассказ первый", "/stories/story1.txt"},
	{'0', "Рассказ второй", "/stories/story2.txt"},
}

func main() {
	addr := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Gopher server started on %s:%s", host, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()

	sel, err := selector(conn)
	if err != nil {
		log.Printf("Read error: %v", err)
		return
	}

	log.Printf("Request: %q", sel)

	switch {
	case sel == "" || sel == "/":
		menu(conn, rootMenu)
	case sel == "/stories" || sel == "/stories/":
		menu(conn, storiesMenu)
	case strings.HasPrefix(sel, "/"):
		file(conn, sel)
	default:
		fail(conn, "Unknown selector")
	}
}

func selector(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func menu(conn net.Conn, entries []menuEntry) {
	writer := bufio.NewWriter(conn)
	for _, entry := range entries {
		switch entry.Type {
		case 'i': // info
			fmt.Fprintf(writer, "i%s\tfake\tfake\t0\r\n", entry.Display)
		default: // file
			fmt.Fprintf(writer, "%c%s\t%s\t%s\t%s\r\n",
				entry.Type, entry.Display, entry.Selector, host, port)
		}
	}
	writer.WriteString(".\r\n")
	writer.Flush()
}

func file(conn net.Conn, selector string) {
	cleanPath := filepath.Clean(strings.TrimPrefix(selector, "/"))
	fullPath := filepath.Join(dataDir, cleanPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		fail(conn, "File not found")
		return
	}

	writer := bufio.NewWriter(conn)
	writer.Write(data)
	if !strings.HasSuffix(string(data), "\n") {
		writer.WriteString("\r\n")
	}
	writer.WriteString(".\r\n")
	writer.Flush()
}

func fail(conn net.Conn, msg string) {
	writer := bufio.NewWriter(conn)
	fmt.Fprintf(writer, "3%s\tfake\tfake\t0\r\n", msg)
	writer.WriteString(".\r\n")
	writer.Flush()
}
