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


func SetEmoteUrl(id, format string) {
    EmoteWebUrlStatic = "https://static-cdn.jtvnw.net/emoticons/v2/" + id + "/" + "static" + "/light/4.0"

    if format == "animated" {
        EmoteWebUrlAnimated = "https://static-cdn.jtvnw.net/emoticons/v2/" + id + "/" + "animated" + "/light/4.0"
    }else {
        EmoteWebUrlAnimated = ""

    }
    // notifyClients(EmoteWebUrl)
}




