// send.go
package main


import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
)



func GetEmotes(streamToken *StreamToken)error {
    url := "https://api.twitch.tv/helix/chat/emotes/user?user_id=" + streamToken.SenderId

    // jsonにしたpayloadをPOSTリクエストで送信
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }


    // ヘッダーの定義
    req.Header.Set("Authorization", "Bearer "+streamToken.AccsessToken)
    req.Header.Set("Client-Id", streamToken.CliendId)

    // ここ送信してるのかな？
    // client := &http.Client{}
    // resp, err := client.Do(req)
    client := new(http.Client)
    resp, _ := client.Do(req)

    // check responce
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error:", err)
    }

    // 3. JSONデータをデコードする
    var jsonData map[string]interface{}
    if err := json.Unmarshal(body, &jsonData); err != nil {
        fmt.Println("Error:", err)
    }

    // 4. JSONデータを表示する
    fmt.Println(jsonData)

    fmt.Println("response Status:", resp.Status)


    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send message: %s", resp.Status)
    }

    return nil
}




func getUser(broadcasterID, senderID, message, token, clientID, get_user_id string) error {
    url := "https://api.twitch.tv/helix/users?login=" + get_user_id

    // jsonにしたpayloadをPOSTリクエストで送信
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }


    // ヘッダーの定義
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Client-Id", clientID)

    // ここ送信してるのかな？
    // client := &http.Client{}
    // resp, err := client.Do(req)
    client := new(http.Client)
    resp, _ := client.Do(req)

    // check responce
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error:", err)
    }

    // 3. JSONデータをデコードする
    var jsonData map[string]interface{}
    if err := json.Unmarshal(body, &jsonData); err != nil {
        fmt.Println("Error:", err)
    }

    // 4. JSONデータを表示する
    fmt.Println(jsonData)

    fmt.Println("response Status:", resp.Status)


    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send message: %s", resp.Status)
    }

    return nil
}

func getModeration(broadcasterID, senderID, message, token, clientID, get_user_id string) error {
    url := "https://api.twitch.tv/helix/moderation/mopderators?broadcaster_id=" + broadcasterID

    // jsonにしたpayloadをPOSTリクエストで送信
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }


    // ヘッダーの定義
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Client-Id", clientID)

    // ここ送信してるのかな？
    // client := &http.Client{}
    // resp, err := client.Do(req)
    client := new(http.Client)
    resp, _ := client.Do(req)

    // check responce
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error:", err)
    }

    // 3. JSONデータをデコードする
    var jsonData map[string]interface{}
    if err := json.Unmarshal(body, &jsonData); err != nil {
        fmt.Println("Error:", err)
    }

    // 4. JSONデータを表示する
    fmt.Println(jsonData)

    fmt.Println("response Status:", resp.Status)


    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send message: %s", resp.Status)
    }

    return nil
}


func SendMessage(streamToken *StreamToken, msg string) error {
    url := "https://api.twitch.tv/helix/chat/messages"

    payload := map[string]string{
        "broadcaster_id": streamToken.BroadcasterId,
        "sender_id":      streamToken.SenderId,
        "message":        msg,
    }
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    // jsonにしたpayloadをPOSTリクエストで送信
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    // ヘッダーの定義
    req.Header.Set("Authorization", "Bearer "+ streamToken.AccsessToken)
    req.Header.Set("Client-Id", streamToken.CliendId)
    req.Header.Set("Content-Type", "application/json")

    // ここ送信してるのかな？
    client := &http.Client{}
    resp, err := client.Do(req)

    // check responce
    if err != nil {
        return err
    }
    defer resp.Body.Close() // deferって遅延実行らしくよく分からない。

    fmt.Println("response Status:", resp)


    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send message: %s", resp.Status)
    }

    return nil
}




// Annoucementes形式で、POSTrequestを送信します。
// colorは、blue, green, orange, purple primaryのいずれかを指定してください。
func SendAnnouncementes(streamToken *StreamToken, msg, color string) error {
    url := "https://api.twitch.tv/helix/chat/announcements?broadcaster_id=" + streamToken.BroadcasterId+ "&moderator_id=" + streamToken.SenderId

    payload := map[string]string{
        "message": msg,
        "color": color,
    }
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    // jsonにしたpayloadをPOSTリクエストで送信
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    // ヘッダーの定義
    req.Header.Set("Authorization", "Bearer "+ streamToken.AccsessToken)
    req.Header.Set("Client-Id", streamToken.CliendId)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)

    // check responce
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp)


    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send message: %s", resp.Status)
    }

    return nil
}


//
// func main() {
//     filePath := "../config/config.txt"
//     var config map[string]string = make(map[string]string)
//     config = Load_file(filePath)
//
//
//     broadcasterID := config["broadcaster_id"]
//     senderID := config["sender_id"]
//     clientID := config["client_id"]
//     token := config["access_token"]
//     message := "aaahhh"
//
//     // err := sendMessage(broadcasterID, senderID, message, token, clientID)
//     // err := getUser(broadcasterID, senderID, message, token, clientID, "carry_nyan")
//     // err := getModeration(broadcasterID, senderID, message, token, clientID, "carry_nyan")
//     err := sendAnnouncementes(broadcasterID, senderID, message, token, clientID )
//     if err != nil {
//         fmt.Println("Error sending message:", err)
//     } else {
//         fmt.Println("Message sent successfully")
//     }
//
//     
//
//
// }
//
