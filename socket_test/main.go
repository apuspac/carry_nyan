package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

    // 変なdomainからの接続をチェック
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    // connectionをwebsocketにupgrade(内部でhttpのハンドシェイクを行なってくれる。)
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil{
        log.Println(err)
    }

    log.Println("Client Connected")
    err = ws.WriteMessage(1, []byte("Hi Client!"))
    if err != nil{
        log.Println(err)
    }

    render(ws)
}


// upgrader からのpointerを受け取る
func render(conn *websocket.Conn){
    for {
        // read in a message
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        // print out that message
        fmt.Println(string(p))


        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }
    }
}




func setupRoutes() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/ws", wsEndpoint)

}

func main() {
    fmt.Println("Go websockets")
    setupRoutes()
    log.Fatal(http.ListenAndServe(":8080", nil))

}
