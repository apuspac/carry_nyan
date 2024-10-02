package main

import (
    // "html/template"
    // "log"
    // "net/http"
    // "path/filepath"
    // "github.com/gorilla/websocket"
    // "strings"
    // "fmt"
    // "io/ioutil"
)

var (
    EmoteWebUrlStatic string = "https://static-cdn.jtvnw.net/emoticons/v2/emotesv2_d878fe2c4fe4463c8a6cdd5257d6a0ed/static/light/4.0"
    EmoteWebUrlAnimated string = "https://static-cdn.jtvnw.net/emoticons/v2/emotesv2_d878fe2c4fe4463c8a6cdd5257d6a0ed/animated/light/4.0"

    EmoteTVUrl string = "https://cdn.betterttv.net/emote/5ba6d5ba6ee0c23989d52b10/3x"
    Emote7TVUrl string = "https://cdn.7tv.app/emote/60ae7316f7c927fad14e6ca2/3x.webp"
)

func GetEmoteUrl() string {
    if(EmoteWebUrlAnimated == ""){
        return EmoteWebUrlStatic
    }
    return EmoteWebUrlAnimated
}

func GetEmoteUrlStatic() string {
    return EmoteWebUrlStatic
}


func GetEmoteTVUrl() string {
    return EmoteTVUrl
}

func GetEmote7TVUrl() string {
    return Emote7TVUrl
}

func SetEmoteUrl(id, format string) {
    EmoteWebUrlStatic = "https://static-cdn.jtvnw.net/emoticons/v2/" + id + "/" + "static" + "/light/4.0"

    if format == "animated" {
        EmoteWebUrlAnimated = "https://static-cdn.jtvnw.net/emoticons/v2/" + id + "/" + "animated" + "/light/4.0"
    }else {
        EmoteWebUrlAnimated = ""

    }
    // notifyClients(EmoteWebUrl)
}

func SetTVEmoteUrl(id string){
    EmoteTVUrl = "https://cdn.betterttv.net/emote/" + id + "/3x"
}

func Set7TVEmoteUrl(id string){
    Emote7TVUrl = "https://cdn.7tv.app/emote/" + id + "/3x.webp"
}

