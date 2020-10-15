package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/sacOO7/gowebsocket"
	"github.com/splable/agent/v1/conf"
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
func NewSocket(conf conf.File) *Client {
	socket := gowebsocket.New(conf.Hostname)
	socket.RequestHeader.Set("Origin", UserAgent)
	socket.RequestHeader.Set("Authorization", fmt.Sprintf("Token: %s", conf.Token))
	socket.Connect()

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		log.Println("Recieved binary data ", data)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
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
func (c *ChannelService) Subscribe(channelName string) {
	identifier := channelIdentifier{
		Channel: channelName,
	}

	encodedIdentifier, err := json.Marshal(identifier)
	if err != nil {
		log.Panic(err)
	}

	subscribe := ChannelSubscribe{
		Command:    CommandTypeSubscribe,
		Identifier: string(encodedIdentifier),
	}

	encodedSubscribe, err := json.Marshal(subscribe)
	if err != nil {
		log.Panic(err)
	}

	// TODO: Need to handle subscription failures.
	c.client.socket.SendText(string(encodedSubscribe))

	// TODO: Related to the above. Waiting isn't nessesary if the subscription successful responce is handled.
	time.Sleep(1 * time.Second)
}
