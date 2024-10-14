package main

import (
    "log"
    "net/http"
    // "path/filepath"
    "github.com/gorilla/websocket"
    "fmt"
)

var (
    gows *websocket.Conn
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

    gows = wsgo
    defer wsgo.Close()


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




// godotにEmoteとuser名を乗っけたjsonを送る
func sendEmoteUrlToGodot(ws *websocket.Conn, user string){
    err := ws.WriteJSON(map[string]string {"user": user, "emote_url": GetEmoteUrl()})
    if err != nil{
        log.Println(err)
    }

}


func sendEmoteTVUrlToGodot(ws *websocket.Conn, user string){
    err := ws.WriteJSON(map[string]string {"user": user, "emote_url": GetEmoteTVUrl()})
    if err != nil{
        log.Println(err)
    }

}

func sendEmote7TVUrlToGodot(ws *websocket.Conn, user string){
    err := ws.WriteJSON(map[string]string {"user": user, "emote_url": GetEmote7TVUrl()})
    if err != nil{
        log.Println(err)
    }

}



// godotにmessageを乗っけたjsonを送る
func MsgNotifyforGodot(user, msg string){
    err := gows.WriteJSON(map[string]string {"user": user, "msg": msg})
    if err != nil{
        log.Println(err)
    }
}

func EmoteNotifyforGodot(user string){
    // TODO: ここをarrayで送れるようにする
    err := gows.WriteJSON(map[string]string {"user": user, "emote_url": GetEmoteUrl()})
    if err != nil{
        log.Println(err)
    }
}

func HydrateNotifyforGodot(user string){
    err := gows.WriteJSON(map[string]string {"user": user, "hydrate": "hydrate"})
    if err != nil{
        log.Println(err)
    }
}
