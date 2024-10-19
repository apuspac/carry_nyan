package main

import (
    // "encoding/json"
    "strings"
    "log"
    "net/http"

    // "path/filepath"
    "fmt"

    "github.com/gorilla/websocket"
)

var (
    gows *websocket.Conn
)


func LaunchServerForGodot(){
    http.HandleFunc("/godot", godot_wsEndpoint)
    http.ListenAndServe(":3001", nil)
}


type EmoteJSON struct {
        User string
        // 上のurlとidは同じindexで対応
        EmoteUrl []string
        EmoteId []string
}



// godotとwebsocketの通信。
func godot_wsEndpoint(w http.ResponseWriter, r *http.Request) {
    log.Println("godot_wsendpoint create")
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


    log.Println("godot Client Connected")
    err = wsgo.WriteMessage(1, []byte("Hi Client!"))
    if err != nil{
        log.Println(err)
    }


    godot_communication(wsgo)
}

func testEmote() {
        testEmoteArray := []string{
            "https://static-cdn.jtvnw.net/emoticons/v2/emotesv2_dad5dd0bf0484a9185f3019fd332baa3/animated/light/4.0", 
            "https://static-cdn.jtvnw.net/emoticons/v2/emotesv2_35dbfad152f8403dbfa7197d095cb960/static/light/4.0",
        }
        testEmoteIdArray := []string{
            "yunivoMunimuniL",
            "x37nuiStar",
        }

        testEmote := EmoteJSON{
            User: "test",
            EmoteUrl: testEmoteArray,
            EmoteId: testEmoteIdArray,
        }

        data := map[string]string {"user": testEmote.User, "emote_url": strings.Join(testEmote.EmoteUrl, ","), "emote_id": strings.Join(testEmote.EmoteId, ",")}


        if err := gows.WriteJSON(data); err != nil {
            log.Println(err)
            return
        }
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
        log.Println("from godot " + string(p))

        // testEmote()


        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }
    }
}

// godotにmessageを乗っけたjsonを送る
func MsgNotifyforGodot(user, msg string){
    if gows != nil{
        err := gows.WriteJSON(map[string]string {"user": user, "msg": msg})
        if err != nil{
            log.Println(err)
        }
    }
}

func EmoteNotifyforGodot(user string){
    emoteJSON := EmoteJSON{
        User: user,
        EmoteUrl: GetEmoteUrl(),
    }

    data := map[string]string {"user": emoteJSON.user , "emote_url": strings.Join(emoteJSON.emote_url, ",")}
    // data := map[string]string {"user": testEmote.User, "emote_url": strings.Join(testEmote.EmoteUrl, ","), "emote_id": strings.Join(testEmote.EmoteId, ",")}
    
    if gows != nil{
        err := gows.WriteJSON(data)
        if err != nil{
            log.Println(err)
        }
    }

    fmt.Println("emoteJSON:", emoteJSON)
}

func HydrateNotifyforGodot(user string){
    if gows != nil{
        err := gows.WriteJSON(map[string]string {"user": user, "hydrate": "hydrate"})
        if err != nil{
            log.Println(err)
        }
    }
}
