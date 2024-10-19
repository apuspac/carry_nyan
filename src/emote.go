package main

import (
    // "log"
    "net/http"
    "path/filepath"
    "strings"
    "fmt"
    // "image/gif"
    // "bytes"
    // "io/ioutil"
)


type Emote struct {
    Platform string
    Url string
    Id string
}

var (
    // EmoteWebUrlStatic string = "https://static-cdn.jtvnw.net/emoticons/v2/emotesv2_d878fe2c4fe4463c8a6cdd5257d6a0ed/static/light/4.0"
    // EmoteWebUrlAnimated string = "https://static-cdn.jtvnw.net/emoticons/v2/emotesv2_d878fe2c4fe4463c8a6cdd5257d6a0ed/animated/light/4.0"
    // EmoteTVUrl string = "https://cdn.betterttv.net/emote/5ba6d5ba6ee0c23989d52b10/3x"
    // Emote7TVUrl string = "https://cdn.7tv.app/emote/63a2d49d7ace4044b8d02681/4x.gif"

    EmoteArray []Emote
)



func GetEmoteUrl() []string {
    var emoteUrls []string


    for _, url := range EmoteArray {
        emoteUrls = append(emoteUrls, url.Url)
    }

    return emoteUrls
}

func GetEmoteId() []string {
    var emoteUrls []string
    for _, emote_id := range EmoteArray {
        emoteUrls = append(emoteUrls, emote_id.Id)
    }

    return emoteUrls
}



func SetEmoteUrl(url_id, format, platform, emote_id string) {
    var emoteWebUrl string

    switch platform {
    case "twitch":
        emoteWebUrl = "https://static-cdn.jtvnw.net/emoticons/v2/" + url_id + "/" + "static" + "/light/4.0"

        if format == "animated" {
            emoteWebUrl = "https://static-cdn.jtvnw.net/emoticons/v2/" + url_id + "/" + "animated" + "/light/4.0"
        }
    case "betterttv":
        emoteWebUrl = "https://cdn.betterttv.net/emote/" + url_id + "/3x"
        if format == "animated" {
            emoteWebUrl = "https://cdn.betterttv.net/emote/" + url_id + "/3x"
        }

    case "7tv":
        emoteWebUrl = "https://cdn.7tv.app/emote/" + url_id + "/3x.png"

        if format == "animated" {
            emoteWebUrl = "https://cdn.7tv.app/emote/" + url_id + "/3x.gif"
        }
    }



    EmoteArray = append(EmoteArray, Emote{platform,emoteWebUrl, emote_id})
}


func ClearEmoteArray(){
    EmoteArray = nil
}


func emoteServeHandler(w http.ResponseWriter, r *http.Request){
    // urlからemote名をtrimして、emoteDIRからemoteを取得
    emoteName := strings.TrimPrefix(r.URL.Path, "emote/")
    emotePath := filepath.Join("emote", emoteName + ".gif")

    http.ServeFile(w, r, emotePath)


}

func ServerforEmote(){
    http.HandleFunc("/emote/", emoteServeHandler)

    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server", err)
    }
}
