package main

import (
    "log"
    "net/http"
    // "path/filepath"
    "github.com/gorilla/websocket"
    "fmt"
)

var (
    godotclients = make(map[*websocket.Conn]bool)
    godot_websockets_upgrader = websocket.Upgrader{}
)


func LaunchServerForGodot(){
    http.HandleFunc("/godot", godot_wsEndpoint)
    http.ListenAndServe(":3001", nil)
}



// godotとwebsocketの通信。
func godot_wsEndpoint(w http.ResponseWriter, r *http.Request) {
    fmt.Println("viewEmoteHandler")
    var upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }

    // 変なdomainからの接続をチェック
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    wsgo, err := upgrader.Upgrade(w, r, nil)
    if err != nil{
        log.Println(err)
    }

    defer wsgo.Close()


    clients[wsgo] = true
    defer delete(clients, wsgo)


    log.Println("Client Connected")
    err = wsgo.WriteMessage(1, []byte("Hi Client!"))
    if err != nil{
        log.Println(err)
    }


    godot_communication(wsgo)
}



func godot_communication(conn *websocket.Conn){
    for {
        // read in a message
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        // print out that message
        fmt.Println(string(p))


        emoteUrl := GetEmoteUrl()
        p = []byte(emoteUrl)


        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }

    }

}



// godotにmessageを乗っけたjsonを送る
func sendMessageToGodot(ws *websocket.Conn, user, msg string){
    err := ws.WriteJSON(map[string]string {"user": user, "msg": msg})
    if err != nil{
        log.Println(err)
    }

}

// godotにEmoteとuser名を乗っけたjsonを送る
func sendEmoteUrlToGodot(ws *websocket.Conn, user string){
    err := ws.WriteJSON(map[string]string {"user": user, "emote_url": GetEmoteUrlStatic()})
    if err != nil{
        log.Println(err)
    }

}


// すべてのclientに通知(clientは一つだからこれ一つでよくないか？
func MsgNotifyforGodot(user, msg string){
    for client := range clients {
        sendMessageToGodot(client, user, msg)
    }

}

func EmoteNotifyforGodot(user, msg string){
    for client := range clients {
        sendEmoteUrlToGodot(client, user)
    }

}


