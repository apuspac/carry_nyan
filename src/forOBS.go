package main

import (
    "html/template"
    "log"
    "net/http"
    // "path/filepath"
    "github.com/gorilla/websocket"

    // test
    "math/rand"
    "time"
    "fmt"
)

var (
    clients = make(map[*websocket.Conn]bool)
    EmoteWebUrl string = "https://static-cdn.jtvnw.net/emoticons/v2/emotesv2_d878fe2c4fe4463c8a6cdd5257d6a0ed/animated/light/4.0"
    upgrader = websocket.Upgrader{}
)


func getEmoteUrl() string {
    return EmoteWebUrl
}

func SetEmoteUrl(id, format string) {
    EmoteWebUrl = "https://static-cdn.jtvnw.net/emoticons/v2/" + id + "/" + format + "/light/2.0"
    notifyClients(EmoteWebUrl)
}



func LaunchServerForOBS(){
    http.HandleFunc("/ws", wsEndpoint)
    http.ListenAndServe(":3000", nil)

    // wsUrl := "ws://localhost:3030/ws"
    //
    // conn, _, err := websocket.DefaultDialer.Dial(wsUrl, http.Header{})
    //
    // if err != nil {
    //     log.Println("Error while upgrading connection:", err)
    //     return
    // }
    // defer conn.Close()



}




func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    fmt.Println("viewEmoteHandler")
    var upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }

    // 変なdomainからの接続をチェック
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    // connectionをwebsocketにupgrade(内部でhttpのハンドシェイクを行なってくれる。)
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil{
        log.Println(err)
    }

    defer ws.Close()


    clients[ws] = true
    defer delete(clients, ws)


    log.Println("Client Connected")
    err = ws.WriteMessage(1, []byte("Hi Client!"))
    if err != nil{
        log.Println(err)
    }


    render(ws)
    // ramdom_generate(ws)
}

func sendEmoteMessage(ws *websocket.Conn){
    err := ws.WriteMessage(1, []byte(getEmoteUrl()))
    if err != nil{
        log.Println(err)
    }
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


        // ここでオウム返し

        emoteUrl := getEmoteUrl()
        // p = []byte(emoteUrl + " " + string(p))
        p = []byte(emoteUrl)


        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }

    }

}


func ramdom_generate(conn *websocket.Conn){
    for {
        fmt.Println("ramdom_generate")
        // Generate a random number
        rand.Seed(time.Now().UnixNano())
        number := rand.Intn(100)

        // Send the number to the React client
        if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", number))); err != nil {
            log.Println("Error sending message:", err)
            return
        }

        // Wait for a while before generating the next number
        time.Sleep(5 * time.Second)
    }

}






func notifyClients(meg string){
    for client := range clients {
        sendEmoteMessage(client)
    }

}


// defaultのurl画像が表示されるはず...
func viewEmoteHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the image filename from the URL query

    // ここで、?url=を取得してて、
    // emoteUrl := r.URL.Query().Get("url")
    emoteUrl := getEmoteUrl()
    if emoteUrl == "" {
        http.Error(w, "Missing image parameter", http.StatusBadRequest)
        return
    }

    // Parse and execute the HTML template
    // htmlをparseして、imageを組み込んで表示
    tmpl, err := template.ParseFiles("templates/emote.html")
    if err != nil {
        http.Error(w, "Error parsing template", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, emoteUrl)
}


func _LaunchServerForOBS(){
    fmt.Println("LaunchServerForOBS")

    // URLが/viewの場合はこっちが判定される
    http.HandleFunc("/emote", viewEmoteHandler)
    // fs := http.FileServer(http.Dir("templates"))
    // http.Handle("/", fs)
    // http.HandleFunc("/ws", handleEmoteWebSockets)

    log.Println("Listening... on :3000...")

    // localhost:3000でserver が立ち上がる。
    err := http.ListenAndServe(":3000", nil)
    if err != nil {
        log.Fatal(err)
    }

}
