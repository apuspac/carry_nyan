package main

import (
    "net/http"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "log"
)

type betterttv_global struct {
    Id string `json:"id"`
    Code string `json:"code"`
    ImageType string `json:"imageType"`
    Animated bool `json:"animated"`
    Userid string `json:"userId"`
    Modifier bool `json:"modifier"` 
}

type betterttv_users struct {
    Id string `json:"id"`
    Bots []string `json:"bots"`
    AvatarUrl string `json:"avatar"`
    ChannelEmotes []struct{
        Id        string `json:"id"`
        Code      string `json:"code"`
        ImageType string `json:"imageType"`
        Animated  bool `json:"animated"`
        Userid    string `json:"userId"`
    }`json:"channelEmotes"`
    SharedEmotes []struct{
        Id        string `json:"id"`
        Code      string `json:"code"`
        ImageType string `json:"imageType"`
        Animated  bool `json:"animated"`
        User struct{
            Id          string `json:"id"`
            Name        string `json:"name"`
            DisplayName string `json:"displayName"`
            ProviderId  string `json:"providerId"`
        } `json:"user"`
    }`json:"sharedEmotes"`
}

type tv7_emote_sets struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Emotes []tv7_emote `json:"emotes"`
}

type tv7_emote struct {
    Id string `json:"id"`
    Code string `json:"name"`
}



var (
    // おい ここに書くんじゃありません。
    betterttv_user_id string = "139776244"
    tv7_user_id string = "6676d128ca58f5151ecd3b7e"
    tv7_emote_set_id string = "6676d198180087fd48c2a128"
    Replace_emote_list_ttv map[string]string = make(map[string]string)
    Replace_emote_list_7tv map[string]string = make(map[string]string)
)


func http_request_GET(url string, return_body *[]byte) error{

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

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
    
    fmt.Println(url)
    log.Println("response Status:", resp.Status)

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send message: %s", resp.Status)
    }

    *return_body = body

    return nil

}


func Get7tvEmoteWebp (emote_url string, webpbyte *[]byte) {
    url := emote_url

    http_request_GET(url, webpbyte)
}



func GetList7tvEmoteSets() {
    url := "https://7tv.io/v3/emote-sets/" + tv7_emote_set_id

    var body []byte
    http_request_GET(url, &body)

    var emoteData tv7_emote_sets
    if err := json.Unmarshal(body, &emoteData); err != nil {
        fmt.Println("Error:", err)
    }

    for _, emote := range emoteData.Emotes{
        Replace_emote_list_7tv[emote.Code] = emote.Id
    }
}

func GetListBetterttvUser() {
    url := "https://api.betterttv.net/3/cached/users/twitch/" + betterttv_user_id

    var body []byte
    http_request_GET(url, &body)

    var emoteData betterttv_users
    if err := json.Unmarshal(body, &emoteData); err != nil {
        fmt.Println("Error:", err)
    }

    for _, emote := range emoteData.ChannelEmotes {
        Replace_emote_list_ttv[emote.Code] = emote.Id
    }

    for _, emote := range emoteData.SharedEmotes {
        Replace_emote_list_ttv[emote.Code] = emote.Id
    }
}


func GetListBetterttvGlobal() {
    url := "https://api.betterttv.net/3/cached/emotes/global"

    var body []byte
    http_request_GET(url, &body)

    var emoteData []betterttv_global
    if err := json.Unmarshal(body, &emoteData); err != nil {
        fmt.Println("Error:", err)
    }

    for _, emote := range emoteData {
        Replace_emote_list_ttv[emote.Code] = emote.Id
    }
}
