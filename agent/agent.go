package agent

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sacOO7/gowebsocket"
	"github.com/splable/agent/v1/conf"
	"github.com/splable/agent/v1/logger"
)

// https://leopard.in.ua/2012/07/08/using-cors-with-rails#.X2DJnJMzbUI
// https://github.com/NullVoxPopuli/action_cable_client#the-action-cable-protocol

const (
	// UserAgent for Rails CORS whitelist.
	UserAgent = "agent"
	// ActionName must match the ActionCable method name.
	ActionName = "broadcast"
	// CommandTypeSubscribe used to subscribe to a channel.
	CommandTypeSubscribe = "subscribe"
	// CommandTypeMessage used after subscribe.
	CommandTypeMessage = "message"
)

// Client holds our agent connections.
type Client struct {
	common service

	socket  *gowebsocket.Socket
	Channel *ChannelService
}

type service struct {
	client *Client
}

// ChannelService provides a client for performing web socket operations.
type ChannelService service

// ChannelSubscribe returns an ActionCable compliant subscription.
type ChannelSubscribe struct {
	Command    string `json:"command"`
	Identifier string `json:"identifier"`
}

// ChannelMessage returns an ActionCable compliant message.
type ChannelMessage struct {
	Command    string `json:"command"`
	Identifier string `json:"identifier"`
	Data       string `json:"data"`
}

type channelIdentifier struct {
	Channel string `json:"channel"`
}

type channelData struct {
	Action  string         `json:"action"`
	Content channelContent `json:"content"`
}

type channelContent struct {
	Datetime string `json:"datetime"`
	Message  string `json:"message"`
}

// NewSocket returns a websocket.
func NewSocket(l logger.Logger, conf conf.File) *Client {
	socket := gowebsocket.New(conf.Hostname)
	socket.RequestHeader.Set("Origin", UserAgent)
	socket.RequestHeader.Set("Authorization", fmt.Sprintf("Token: %s", conf.Token))
	socket.Connect()

	socket.OnConnected = func(socket gowebsocket.Socket) {
		l.Info("Connected to server  %s", conf.Hostname)
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		l.Fatal("Recieved connect error  %s", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		l.Debug("Recieved message %s", message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		l.Debug("Recieved binary data %s", data)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		l.Debug("Recieved ping %s", data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		l.Debug("Recieved pong %s", data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		l.Info("Disconnected from %s", conf.Hostname)
		return
	}

	c := &Client{
		socket: &socket,
	}
	c.common.client = c
	c.Channel = (*ChannelService)(&c.common)

	return c
}

// CloseSocket closes a client socket.
func (c *Client) CloseSocket() {
	c.socket.Close()
}

// Subscribe joins a particular ActionCable channel.
func (c *ChannelService) Subscribe(l logger.Logger, channelName string) {
	identifier := channelIdentifier{
		Channel: channelName,
	}

	encodedIdentifier, err := json.Marshal(identifier)
	if err != nil {
		l.Error("Error subscribing to %s channel: %s", channelName, err)
	}

	subscribe := ChannelSubscribe{
		Command:    CommandTypeSubscribe,
		Identifier: string(encodedIdentifier),
	}

	encodedSubscribe, err := json.Marshal(subscribe)
	if err != nil {
		l.Error("Error subscribing to %s channel: %s", channelName, err)
	}

	// TODO: Need to handle subscription failures.
	c.client.socket.SendText(string(encodedSubscribe))

	// TODO: Related to the above. Waiting isn't nessesary if the subscription successful response is handled.
	time.Sleep(1 * time.Second)
}
