package gopher

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/playdate-go/pdgo"
)

type Connection interface {
	SetConnectionClosedCallback(callback func(err pdgo.PDNetErr))
	Open(callback func(err pdgo.PDNetErr)) pdgo.PDNetErr
	Read(buf []byte) int
	Write(data []byte) int
	GetBytesAvailable() int
	GetError() pdgo.PDNetErr
}

const (
	TypeFile      byte = '0'
	TypeDirectory byte = '1'
	TypeSearch    byte = '7'
	TypeError     byte = '3'
	TypeBinHex    byte = '4'
	TypeBinary    byte = '9'
	TypeInfo      byte = 'i'
	TypeTelnet    byte = 'T'
	TypeSound     byte = 's'
	TypeImage     byte = 'I'
)

type Item struct {
	Type     byte
	Name     string
	Selector string
	Host     string
	Port     int
}

type state int

const (
	stateIdle state = iota
	stateDir
	stateFile
)

type Client struct {
	pd *pdgo.PlaydateAPI

	server string
	port   int
}

func New(pd *pdgo.PlaydateAPI, server string, port int) *Client {
	return &Client{
		pd:     pd,
		server: server,
		port:   port,
	}
}

func (c *Client) RequestAccess(callback pdgo.AccessCallback) pdgo.AccessReply {
	tcp := c.pd.Network.TCP()
	return tcp.RequestAccess(c.server, c.port, true, "for testing", callback)
}

func (c *Client) GetDirectory(selector string, onItems func([]Item), onError func(error)) {
	conn, err := c.connect()

	if err != nil {
		onError(err)

		return
	}

	c.begin(conn, selector, stateDir, onItems, nil, onError)
}

func (c *Client) GetFile(selector string, onData func([]byte), onError func(error)) {
	conn, err := c.connect()

	if err != nil {
		onError(err)

		return
	}

	c.begin(conn, selector, stateFile, nil, onData, onError)
}

func (c *Client) Search(selector, query string, onItems func([]Item), onError func(error)) {
	conn, err := c.connect()

	if err != nil {
		onError(err)

		return
	}

	c.begin(conn, selector+"\t"+query, stateDir, onItems, nil, onError)
}

func (c *Client) connect() (Connection, error) {
	tcp := c.pd.Network.TCP()
	conn := tcp.NewConnection(c.server, c.port, false)

	if conn == nil {
		return nil, fmt.Errorf("failed to create connection")
	}

	return conn, nil
}

func (c *Client) begin(
	conn Connection,
	request string,
	st state,
	onItems func([]Item),
	onData func([]byte),
	onError func(error),
) {
	c.pd.System.LogToConsole("client begin")

	if conn == nil {
		return
	}

	c.pd.System.LogToConsole("set callback")

	conn.SetConnectionClosedCallback(func(err pdgo.PDNetErr) {
		c.pd.System.LogToConsole("connection close callback")

		avail := conn.GetBytesAvailable()
		if err != pdgo.NetOK {
			c.pd.System.LogToConsole(fmt.Sprintf("TCP request complete, err=%d", int(err)))

			onError(fmt.Errorf("TCP request failed, err=%d", int(err)))

			return
		} else {
			c.pd.System.LogToConsole(fmt.Sprintf("TCP request complete, %d bytes available", avail))
		}

		data := []byte{}

		// Read all available data
		for avail > 0 {
			buf := make([]byte, 256)

			n := conn.Read(buf)
			if n > 0 {
				data = append(data, buf...)
			}
			avail = conn.GetBytesAvailable()
		}

		switch st {
		case stateDir:
			items := c.processDir(data)

			if onItems != nil {
				onItems(items)
			}
		case stateFile:
			if onData != nil {
				onData(data)
			}
		}
	})

	conn.Open(func(err pdgo.PDNetErr) {
		c.pd.System.LogToConsole("connection open callback")
		if err != pdgo.NetOK {
			onError(fmt.Errorf("TCP opem failed, err=%d", int(err)))

			return
		}
		conn.Write([]byte(request + "\r\n"))
	})
}

func (c *Client) processDir(buf []byte) []Item {
	items := make([]Item, 0)

	strs := strings.Split(string(buf), "\r\n")

	for _, str := range strs {
		if len(str) > 0 {
			items = append(items, parseLine(str))
		}
	}

	return items
}

func findCRLF(buf []byte) int {
	for i := 0; i < len(buf)-1; i++ {
		if buf[i] == '\r' && buf[i+1] == '\n' {
			return i
		}
	}
	return -1
}

func parseLine(line string) Item {
	var item Item
	if len(line) == 0 {
		return item
	}
	item.Type = line[0]
	rest := line[1:]

	fields := splitTabs(rest)

	if len(fields) > 0 {
		item.Name = fields[0]
	}
	if len(fields) > 1 {
		item.Selector = fields[1]
	}
	if len(fields) > 2 {
		item.Host = fields[2]
	}
	if len(fields) > 3 {
		// TODO: remove atoi
		item.Port, _ = strconv.Atoi(fields[3])
	}
	return item
}

func splitTabs(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\t' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func atoi(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			n = n*10 + int(s[i]-'0')
		}
	}
	return n
}

type gopherError struct {
	code pdgo.PDNetErr
}

func (e *gopherError) Error() string {
	return "gopher error"
}
