package main

import (
	"encoding/json"
	"log"

	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

type MessageObject struct {
	Data string `json:"data"`
	From string `json:"from"`
	To   string `json:"to"`
}

func main() {

	clients := make(map[string]string)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	ikisocket.On(ikisocket.EventConnect, func(ep *ikisocket.EventPayload) {
		log.Printf("Connected: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	ikisocket.On(ikisocket.EventDisconnect, func(ep *ikisocket.EventPayload) {
		delete(clients, ep.Kws.GetStringAttribute("user_id"))
		log.Printf("Disconnected: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	ikisocket.On(ikisocket.EventMessage, func(ep *ikisocket.EventPayload) {
		log.Printf("Message: %s", ep.Data)
		message := MessageObject{}
		err := json.Unmarshal(ep.Data, &message)
		if err != nil {
			log.Println(err)
			return
		}

		if message.To == "all" {
			ikisocket.Broadcast(ep.Data)
			return
		}

		err = ep.Kws.EmitTo(clients[message.To], ep.Data)
		if err != nil {
			log.Println(err)
		}
	})

	ikisocket.On(ikisocket.EventClose, func(ep *ikisocket.EventPayload) {
		delete(clients, ep.Kws.GetStringAttribute("user_id"))
		log.Printf("Closed: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	ikisocket.On(ikisocket.EventError, func(ep *ikisocket.EventPayload) {
		log.Printf("Error: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	app.Get("/ws/:id", ikisocket.New(func(kws *ikisocket.Websocket) {
		userId := kws.Params("id")
		clients[userId] = kws.UUID
		kws.SetAttribute("user_id", userId)
		kws.Broadcast([]byte("New user connected: "+userId), true)
		kws.Emit([]byte("Welcome " + userId))
	}))

	app.Listen(":3000")
}
