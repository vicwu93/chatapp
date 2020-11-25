package main

import (
	"example.com/constants"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strings"
)

type wsConn struct {
	*websocket.Conn        // websocket connection
	User            string // username
}

// interface struct that holds a response
type resp struct {
	User    string // name of user sending the message back through a response
	Message string // msg user is sending back
	Option  string // receiving end option to match (case) to see what to do
}

// interface struct that holds a message to send
type msg struct {
	Msg string
}

// slice of connections
var listOfConns = make([]*wsConn, 0)

// websocket handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// server will call this upgrader from the http request
	upgrader := websocket.Upgrader{
		CheckOrigin:     func(*http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// upgrades to a websocket connection
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// can not open a websocket connection
		log.Println(err)
	}

	fmt.Println(r.URL.Query(), "reached query")

	// get username from URL
	user := r.URL.Query()["user"][0]

	// create a struct of websocket connection
	// to add to a list of connections
	conn := wsConn{Conn: c, User: user}
	listOfConns = append(listOfConns, &conn)

	// start a new go routine to handle this new connection
	// and to see what we want to do with this connection
	go echo(listOfConns, &conn)
}

// helper func to send a message to every websocket connection
func sendToAll(conn *wsConn, option string, msg string) {
	fmt.Printf("connection with ptr: %p with message: %s\n", &conn, msg)
	for _, c := range listOfConns {
		if c == conn {
			continue
		} else {
			// WriteJSON marshals and sends json over io.writer
			// our io.writer is the websocket (obvious)
			c.WriteJSON(resp{
				User:    conn.User,
				Message: msg,
				Option:  option,
			})
		}
	}
}

func connect(conn *wsConn) {
	fmt.Printf("connection made: %p\n", &conn)
	for _, c := range listOfConns {
		c.WriteJSON(resp{
			User:    conn.User,
			Message: constants.EMPTY_STRING,
			Option:  constants.CONNECT,
		})
	}
	return
}

/**
 * Go has no filter func, built my own to
 * filter out the connection that gets called to close
 */
func filter(arr []*wsConn, cond func(*wsConn) bool) []*wsConn {
	result := []*wsConn{}
	for i := range arr {
		if cond(arr[i]) {
			result = append(result, arr[i])
		}
	}
	return result
}

/**
 * Remove connection by filtering it out with helper func
 * An anonymouse func to check if tmp == is equal
 * to the iterated list, if so skip and do not append
 * filter out the current connection from the list
 */
func disconnect(conn *wsConn) {
	fmt.Printf("connection removed: %p\n", &conn)
	for _, c := range listOfConns {
		c.WriteJSON(resp{
			User:    conn.User,
			Message: constants.EMPTY_STRING,
			Option:  constants.DISCONNECT,
		})
	}
	return
}

/**
 * A while loop used to call this func to read through received connection
 * and appending it to Msg{}
 * An if check is used for err when user navigates away from page or server goes down
 */
func storeMsg(conn *wsConn) {
	for {
		// storing msg here
		m := msg{}
		// readJSON and sets msg struct to user msg being sent
		err := conn.ReadJSON(&m)
		if err != nil && strings.Contains(err.Error(), constants.GOING_AWAY) {
			return
		} else {
			sendToAll(conn, "msg", m.Msg)
		}
	}
}

// func which handles all i/o connections
func echo(listOfConns []*wsConn, conn *wsConn) {
	connect(conn)
	storeMsg(conn)
	disconnect(conn)
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/ws", wsHandler)
	fmt.Printf("Listening on port : " + port)
	// http.ListenAndServe(":8888", nil)
	http.ListenAndServe(":"+port, nil)
}
