package main

import (
	"fmt"
	"gopher/gopher"

	"github.com/playdate-go/pdgo"
)

type ClientInterface interface {
	RequestAccess(callback pdgo.AccessCallback) pdgo.AccessReply
	GetDirectory(selector string, onItems func([]gopher.Item), onError func(error))
	GetFile(selector string, onData func([]byte), onError func(error))
}

type ViewIntarface interface {
	Data(d Data)
}

type SystemInterface interface {
	LogToConsole(msg string)
}

type Controller struct {
	client ClientInterface
	view   ViewIntarface
	system SystemInterface
}

func NewController(view ViewIntarface, client ClientInterface, system SystemInterface) *Controller {
	return &Controller{
		view:   view,
		client: client,
		system: system,
	}
}

func (c *Controller) PressB() {
	reply := c.client.RequestAccess(func(allowed bool) {
		c.system.LogToConsole(fmt.Sprintf("TCPAccessCallback: %v", allowed))
	})

	if reply != pdgo.AccessAsk {
		c.system.LogToConsole(fmt.Sprintf("TCP access reply immediate: %d", int(reply)))
	}

	if reply != pdgo.AccessAllow {
		c.system.LogToConsole("TCP access denied")

		return
	}

	c.client.GetDirectory("", func(items []gopher.Item) {
		c.view.Data(Data{
			Items: items,
		})

	}, func(err error) {
		c.system.LogToConsole(fmt.Sprintf("Error: %v", err))
	})
}

func (c *Controller) PressA(cursor gopher.Item) {

	// get item fom view
	// navigate to item
	c.system.LogToConsole("Pressed A: " + cursor.Selector)

	if cursor.Type == '1' {
		c.client.GetDirectory(cursor.Selector, func(items []gopher.Item) {
			c.view.Data(Data{
				Items: items,
			})

		}, func(err error) {
			c.system.LogToConsole(fmt.Sprintf("Error: %v", err))
		})
	}

	if cursor.Type == '0' {
		c.client.GetFile(cursor.Selector, func(file []byte) {
			c.view.Data(Data{
				File: string(file),
			})

		}, func(err error) {
			c.system.LogToConsole(fmt.Sprintf("Error: %v", err))
		})
	}
}
