package main


import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "bytes"
    // "time"
    // "io/ioutil"

    "github.com/gorilla/websocket"
)

const (
    wsUrl = "wss://eventsub.wss.twitch.tv/ws"
)

type StreamToken struct {
    AccsessToken string
    CliendId string
    BroadcasterId string
    SenderId string
    SessionId string
}

type Received struct {
    Metadata Metadata `json:"metadata"`
    Payload Payload `json:"payload"`
}

type Metadata struct {
    MessageID         string `json:"message_id"`
    MessageType       string `json:"message_type"`
    MessageTimestamp  string `json:"message_timestamp"`
    Subscriptiontype  interface{} `json:"subscription_type"` // subscription後
    Subscriptionversion  interface{} `json:"subscription_version"` // subscription後
}

type Payload struct {
    Session Session `json:"session"`
    Subscription Subscription `json:"subscription"`
    Event ChatMessageEvent `json:"event"`
}

type Session struct {
    ID                     string      `json:"id"`
    Status                 string      `json:"status"`
    ConnectedAt            string      `json:"connected_at"`
    KeepaliveTimeoutSeconds int       `json:"keepalive_timeout_seconds"`
    ReconnectURL           interface{} `json:"reconnect_url"` // これはnull値が来る可能性があるためinterface{}にします
    RecoveryURL            interface{} `json:"recovery_url"`  // これもnull値が来る可能性があるためinterface{}にします
}

type Subscription struct {
    Id  string `json:"id"`
    Status string `json:"status"`
    Type string `json:"type"`
    Version string `json:"version"`
}

type ChatMessageEvent struct {
    BroadcasterUserID    string `json:"broadcaster_user_id"`
    BroadcasterUserLogin string `json:"broadcaster_user_login"`
    BroadcasterUserName  string `json:"broadcaster_user_name"`
    ChatterUserID        string `json:"chatter_user_id"`
    ChatterUserLogin     string `json:"chatter_user_login"`
    ChatterUserName      string `json:"chatter_user_name"`
    MessageID            string `json:"message_id"`
    Message              struct {
        Text string `json:"text"`
        Fragments []MessageFragment `json:"fragments"`
    } `json:"message"`
    Color string `json:"color"`
}


type MessageFragment struct {
    Type string `json:"type"`
    Text string `json:"text"`
    Cheermote bool `json:"cheermote"`
    Emote struct{
        Id string `json:"id"`
        EmoteSetId string `json:"emote_set_id"`
    }`json:"emote"`
    Mention struct{
        UserId string `json:"user_id"`
        UserName string `json:"user_name"`
        UserLogin string `json:"user_login"`

    }`json:"mention"`
}





type EventSubRequest struct {
    Type string `json:"type"`
    Version string `json:"version"`
    Condition interface{} `json:"condition"`
    Transport struct {
        Method string `json:"method"`
        SessionID string `json:"session_id"`
    } `json:"transport"`
}


func setStreamToken(filePath string) StreamToken{
    var config map[string]string = make(map[string]string)
    config = Load_file(filePath)
    var streamToken StreamToken

    streamToken.AccsessToken = config["access_token"]
    streamToken.CliendId = config["client_id"]
    streamToken.BroadcasterId = config["broadcaster_id"]
    streamToken.SenderId = config["sender_id"]

    return streamToken
}


func createSubscription(ws *websocket.Conn, streamToken *StreamToken)error {
    eventSubURL := "https://api.twitch.tv/helix/eventsub/subscriptions"
    reqBody := EventSubRequest{
        Type:    "channel.chat.message",
        Version: "1",
        Condition: map[string]string{
            "broadcaster_user_id": streamToken.BroadcasterId,
            "user_id": streamToken.BroadcasterId,
        },
    }
    reqBody.Transport.Method = "websocket"
    reqBody.Transport.SessionID = streamToken.SessionId


    body, err_r := json.Marshal(reqBody)
    if err_r != nil {
        return err_r
    }


    req, err := http.NewRequest("POST", eventSubURL, bytes.NewBuffer(body))
    if err != nil {
        return err
    }

    req.Header.Set("Client-Id", streamToken.CliendId)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", streamToken.AccsessToken))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusAccepted {
        return fmt.Errorf("failed to create subscription, status code: %d", resp.StatusCode)
    }

    log.Println("Subscription created successfully")

    return nil
}

func chatCommand(ws *websocket.Conn, streamToken *StreamToken, msg string) {
    if msg == "!nya" {
        SendMessage(streamToken, "にゃーん")
    }

}



func handleSessionWelcome(ws *websocket.Conn, streamToken *StreamToken) {
    fmt.Println("Received Session Welcome")
    createSubscription(ws, streamToken)
}


func handleNotification(ws *websocket.Conn, received Received, streamToken *StreamToken) {
    switch received.Metadata.Subscriptiontype {
    case "channel.chat.message":
        fmt.Println(received.Payload.Event.ChatterUserName, ": ", received.Payload.Event.Message.Text)
        if received.Payload.Event.Message.Text[0:1] == "!"{
            chatCommand(ws, streamToken, received.Payload.Event.Message.Text)
        }
    default:
        fmt.Println("Received Other Notification")
    }
}


func handleSessionKeepalive(ws *websocket.Conn) {
    // fmt.Println("Received Session Keepalive")
}

func handleOtherMessage(ws *websocket.Conn, received Received) error{
    fmt.Println("Received Other Message")
    // fmt.Println(received.Metadata.MessageType)
    body, err := json.MarshalIndent(received, "", "    ")
    if err != nil {
        return err
    }
    fmt.Println(body)

    return nil
}


// Listen for the messages
func listenForMessages(ws *websocket.Conn, streamToken *StreamToken) {
    for {
        _, msg, err := ws.ReadMessage()
        if err != nil {
            log.Printf("Read Error: %+v\n", err)
            continue
        }

        var received Received

        err_json := json.Unmarshal([]byte(msg), &received)
        if err_json != nil {
            fmt.Println("Error decoding JSON:", err)
            return
        }



        switch received.Metadata.MessageType {
        case "session_welcome":
            streamToken.SessionId = received.Payload.Session.ID
            handleSessionWelcome(ws,streamToken)
        case "notification":
            handleNotification(ws, received, streamToken)
        case "session_keepalive":
            handleSessionKeepalive(ws)
        default:
            handleOtherMessage(ws, received)
        }


    }
}


func main() {
    filePath := "../config/config.txt"
    streamToken := setStreamToken(filePath)
    
    ws, _, err := websocket.DefaultDialer.Dial(wsUrl, http.Header{})
    if err != nil {
        log.Fatal(err)
    }

    defer ws.Close()
    

    listenForMessages(ws, &streamToken)
}

