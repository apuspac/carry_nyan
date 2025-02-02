package main


import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "bytes"
    "strings"
    // "time"
    // "io/ioutil"

    "regexp"
    "unicode"

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
    Event Event`json:"event"`
}

type Session struct {
    ID                     string      `json:"id"`
    Status                 string      `json:"status"`
    ConnectedAt            string      `json:"connected_at"`
    KeepaliveTimeoutSeconds int       `json:"keepalive_timeout_seconds"`
    //これが あまり良くないらしい。
    ReconnectURL           interface{} `json:"reconnect_url"` // これはnull値が来る可能性があるためinterface{}にします
    RecoveryURL            interface{} `json:"recovery_url"`  // これもnull値が来る可能性があるためinterface{}にします
}

type Subscription struct {
    Id  string `json:"id"`
    Status string `json:"status"`
    Type string `json:"type"`
    Version string `json:"version"`
}

type Event interface{}

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
        Owner_Id string `json:"owner_id"`
        Format []string `json:"format"`
    }`json:"emote"`
    Mention struct{
        UserId string `json:"user_id"`
        UserName string `json:"user_name"`
        UserLogin string `json:"user_login"`

    }`json:"mention"`
}

type CustomRewardRedemptionAddEvent struct {
    Id string `json:"id"`
    BroadcasterUserID    string `json:"broadcaster_user_id"`
    BroadcasterUserLogin string `json:"broadcaster_user_login"`
    BroadcasterUserName  string `json:"broadcaster_user_name"`
    UserID        string `json:"user_id"`
    UserLogin     string `json:"user_login"`
    UserName      string `json:"user_name"`
    UserInput     string `json:"user_input"`
    Status        string `json:"status"`
    Reward       struct{
        Id      string `json:"id"`
        Title   string `json:"title"`
        Cost    int `json:"cost"`
        Prompt  string `json:"prompt"`

    }`json:"reward"`
    RedeemedAt    string `json:"redeemed_at"`
}

func IsAllSpace(s string) bool {
    for _, r := range s {
        if !unicode.IsSpace(r) {
            return false
        }
    }
    return true
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


// for create subscription
type EventSubRequest struct {
    Type string `json:"type"`
    Version string `json:"version"`
    Condition interface{} `json:"condition"`
    Transport struct {
        Method string `json:"method"`
        SessionID string `json:"session_id"`
    } `json:"transport"`
}



// twitchAPIとのsubscriptionを作成
// event_sub_typeはここ。: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/
func createSubscription(ws *websocket.Conn, streamToken *StreamToken, event_sub_type string)error {
    eventSubURL := "https://api.twitch.tv/helix/eventsub/subscriptions"
    reqBody := EventSubRequest{
        Type:    event_sub_type,
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

    log.Println(event_sub_type + ": Subscription created successfully")

    return nil
}

// コマンド処理の実行。
// 最初のメッセージが!はチャットコマンドとして扱う。
func chatCommand(ws *websocket.Conn, streamToken *StreamToken, msg string) {
    if msg == "!nya" {
        SendMessage(streamToken, "にゃーん")
        // GetChatters(streamToken)
    }
    if msg == "!dis" {
        SendAnnouncementes(streamToken, "https://discord.gg/HZwVQXPPwM", "blue")
    }
    if msg == "!ocha" {
        SendMessage(streamToken, "~(=^･ω･^)_旦")
    }

    if msg == "!nyan" {
        SendMessage(streamToken, "▒▒▒▒▒█▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀█ ▒▒▒▒▒█░▒▒▒▒▒▒▓▒▓▒▒▒▒▒▒▒░█ ▒▒▒▒▒█░▒▒▓▒▒▒▒▒▒▒▄▄▒▓▒▒░█░▄▄ ▄▀▀▄▄█░▒▒▒▒▒▒▓▒▒█░░▀▄▄▄▄▄▀░░█ █░░░░█░▒▒▒▒▒▒▒▒▒█░░░░░░░░░░░█ ▒▀▀▄▄█░▒▒▒▓▒▒▓▒█░░░█░░░░░█░░░█ ▒▒▒▒▒█░▒▓▒▒▒▓▒▒█░░░░░░░▀░░░░░█ ▒▒▒▄▄█░▒▒▒▓▒▒▒▒▒█░░█▄▄█▄▄█░░█ ▒▒█░░░█▄▄▄▄▄▄▄▄█░█▄▄▄▄▄▄▄▄▄█ ▒▒▒█▄▄█░░█▄▄█░░░░█▄▄█░░█▄▄")
    }


}

// jsonファイルから読みだして、emoteのURLをsetする
func loadEmoteData(received *ChatMessageEvent, frg_index int) {
    MessageFragment := received.Message.Fragments[frg_index]
    fmt.Println(MessageFragment.Emote.Id)
    fmt.Println(MessageFragment.Emote.Format[0])

    id := MessageFragment.Emote.Id
    var format string

    if len(MessageFragment.Emote.Format) == 2 {
        format = "animated"
    }else {
        format = "static"
    }
    SetEmoteUrl(id, format, "twitch", MessageFragment.Text)

    // emoteの文字列削除
}

// subscriptionを確立させる
func handleSessionWelcome(ws *websocket.Conn, streamToken *StreamToken) {
    log.Println("Received Session Welcome from Twitch")
    createSubscription(ws, streamToken, "channel.chat.message")
    createSubscription(ws, streamToken, "channel.channel_points_custom_reward_redemption.add")
}


// messageなどを受け取ったときの処理
func handleNotification(ws *websocket.Conn, received Received, streamToken *StreamToken) {
    switch received.Metadata.Subscriptiontype {
    case "channel.chat.message":
        // EvnetのinterfaceでChatMessageEventの型を指定
        var rcv_event ChatMessageEvent
        eventData, _ := json.Marshal(received.Payload.Event)
        
        // rcv_eventに格納
        json.Unmarshal(eventData, &rcv_event)
        fmt.Println(rcv_event.ChatterUserName, ": ", rcv_event.Message.Text)

        if rcv_event.Message.Text[0:1] == "!"{
            chatCommand(ws, streamToken, rcv_event.Message.Text)
        }

        // messageの中のemoteの分だけ、loadEmoteDataを呼び出す。
        // emoteがないときは、すっとばされる。
        for index, msg_frag := range rcv_event.Message.Fragments {
            if msg_frag.Type == "emote" {
                loadEmoteData(&rcv_event, index)
                rcv_event.Message.Text = strings.Replace(rcv_event.Message.Text, msg_frag.Text, "", 1)
            }
        }

        // check betterttv emote
        for _, ttv_emote:= range Replace_emote_list_ttv {
            re := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(ttv_emote.Code)))
            if re.MatchString(rcv_event.Message.Text){
            // if strings.Contains(rcv_event.Message.Text, ttv_emote.Code) {
            fmt.Println("Replace Emote: " + ttv_emote.Id, "Code: " + ttv_emote.Code)


                for range(strings.Count(rcv_event.Message.Text, ttv_emote.Code)) {
                    var format = "static"
                    if ttv_emote.Animated { format = "animated" }
                    SetEmoteUrl(ttv_emote.Id, format, "betterttv", ttv_emote.Code)
                }

                // rcv_event.Message.Text = strings.ReplaceAll(rcv_event.Message.Text, ttv_emote.Code, "")
                rcv_event.Message.Text = re.ReplaceAllString(rcv_event.Message.Text, "")
            }
        }

        for _, tv7_emote := range Replace_emote_list_7tv {
            // emote文字列をreで検索するためにあらかじめcompile
            re := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(tv7_emote.Code)))
            if re.MatchString(rcv_event.Message.Text){
            // if strings.Contains(rcv_event.Message.Text, tv7_emote.Code) {
                // fmt.Println("Replace Emote: " + tv7_emote.Id)
                fmt.Println("Replace Emote: " + tv7_emote.Id, "Code: " + tv7_emote.Code)

                for range(strings.Count(rcv_event.Message.Text, tv7_emote.Code)) {
                    var format = "static"
                    if tv7_emote.Animated { format = "animated" }
                    SetEmoteUrl(tv7_emote.Id, format, "7tv", tv7_emote.Code)
                }

                // rcv_event.Message.Text = strings.ReplaceAll(rcv_event.Message.Text, tv7_emote.Code, "")
                rcv_event.Message.Text = re.ReplaceAllString(rcv_event.Message.Text, "")
            }
        }


        if len(rcv_event.Message.Text) != 0 && IsAllSpace(rcv_event.Message.Text) == false{
            MsgNotifyforGodot(rcv_event.ChatterUserName, rcv_event.Message.Text)
        }


        if len(EmoteArray) != 0 {
            EmoteNotifyforGodot(rcv_event.ChatterUserName)
        }

        ClearEmoteArray()


    // channel points引き換え
    case "channel.channel_points_custom_reward_redemption.add":
        // Evnet interfaceでCustomRewardRedemptionAddEvnetの型を指定
        var rcv_event CustomRewardRedemptionAddEvent
        eventData, _ := json.Marshal(received.Payload.Event)
        json.Unmarshal(eventData, &rcv_event)

        fmt.Println(rcv_event)

        switch rcv_event.Reward.Title {
        case "Hydrate!":
            HydrateNotifyforGodot(rcv_event.UserName)
        default:
            fmt.Println("Received Other Reward")
        }





    default:
        fmt.Println("Received Other Notification")
    }
}

// 接続が維持されているかが返ってきたときの処理
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


// websocketからメッセージを受信
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


        // MetaDataに乗っかってくるmessagetypeによって処理分け。
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

    filePath := "config/config.txt"
    streamToken := setStreamToken(filePath)

    // GetChatters(&streamToken)

    // getemote list 
    emoteFilePath := "config/config_emote.txt"
    LoadExtentionEmotesList(emoteFilePath)



    ws, _, err := websocket.DefaultDialer.Dial(wsUrl, http.Header{})
    if err != nil {
        log.Fatal(err)
    }


    // mainが終了されたら、実行される。
    defer ws.Close()

    go listenForMessages(ws, &streamToken)
    go LaunchServerForOBS()
    go LaunchServerForGodot()
    go ServerforEmote()

    select {}
}

