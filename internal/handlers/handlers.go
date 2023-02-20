package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

type WebSocketConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action string `json:"action"`
	Message string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsPayload struct {
	Action string `json:"action"`
	Username string `json:"username"`
	Message string `json:"message"`
	Conn WebSocketConnection `json:"-"`
}

func WsEndpoint(w http.ResponseWriter, r *http.Request)  {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap)  error {
	view, err := views.GetTemplate(tmpl)

	if err != nil {
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		return err
	}

	return nil
}
