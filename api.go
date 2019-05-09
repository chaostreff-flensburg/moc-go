package moc

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"

	"github.com/chaostreff-flensburg/moc-go/models"
)

type Client struct {
	Endpoint    string
	LastMessage *models.Message
	NewMessages chan *models.Message
	Log         *logrus.Entry
}

// NewClient create a new MOC API Client
func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint:    endpoint,
		NewMessages: make(chan *models.Message, 5),
		Log:         log.WithFields(log.Fields{"component": "moc-api"}),
	}
}

// Request queries the last messages. To do this, attach to the set endpoint `/messages`.
// If a new message is found, it is additionally packed into the channel client.NewMessage.
func (c *Client) Request() (error, []*models.Message) {
	c.Log.Info("Crawl new Messages...")

	url := fmt.Sprintf("%s/messages", c.Endpoint)

	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	netClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	resp, err := netClient.Get(url)
	if err != nil {
		c.Log.Fatalln(err)
		return err, []*models.Message{}
	}

	var result []*models.Message
	json.NewDecoder(resp.Body).Decode(&result)

	if len(result) == 0 {
		return nil, result
	}

	// push new message to queue
	latestMessage := (result)[len(result)-1]
	if c.LastMessage == nil || c.LastMessage.ID != latestMessage.ID {
		c.Log.Info("New Message detected")
		c.LastMessage = latestMessage

		// if channel full, discarding first value
		select {
		case c.NewMessages <- latestMessage:
		default:
			<-c.NewMessages
		}

	}

	return nil, result
}

// Loop calls the API in the defined amount of time.
func (c *Client) Loop(tick time.Duration) {
	go func() {
		for range time.Tick(tick) {
			c.Request()
		}
	}()
}
