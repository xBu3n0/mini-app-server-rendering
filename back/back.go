package back

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	fmt.Println("Connected")

	s.conns[ws] = true

	s.readLoop(ws)
}

type SCounter struct {
	Counter int64  `json:"counter"`
	Name    string `json:"name"`
}

func (c *SCounter) Init() {
	c.Counter = 20
	c.Name = "Nome teste"
}

func (c *SCounter) Increment() {
	c.Counter += 1
}

func (c *SCounter) Append(FormData map[string]string) {
	c.Name += FormData["name"]
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024)

	for {
		// err := websocket.JSON.Receive(ws, &content)
		n, err := ws.Read(buff)

		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println(err)
			continue
		}

		contents := strings.Split(string(buff[:n]), "\r\n")

		call := contents[0]
		body_request := contents[1]
		body_variables := contents[2]

		var FormData map[string]string
		counterS := SCounter{}

		err = json.Unmarshal([]byte(body_request), &FormData)
		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal([]byte(body_variables), &counterS)
		if err != nil {
			fmt.Println(err)
		}

		switch call {
		case "init":
			counterS.Init()
		case "increment":
			counterS.Increment()
		case "append":
			counterS.Append(FormData)
		}

		fmt.Println(counterS, call)
		response, err := json.Marshal(counterS)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(counterS, call)

		ws.Write(response)
	}
}

func main() {
	s := NewServer()
	http.Handle("/ws", websocket.Handler(s.HandleWS))

	http.ListenAndServe(":3000", nil)
}
