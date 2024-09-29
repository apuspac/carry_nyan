package main

import (
    // "html/template"
    // "log"
    // "net/http"
    // "path/filepath"
    // "github.com/gorilla/websocket"
    // "fmt"
)

var (
    EmoteWebUrl string = "https://static-cdn.jtvnw.net/emoticons/v2/emotesv2_d878fe2c4fe4463c8a6cdd5257d6a0ed/animated/light/4.0"
)

func GetEmoteUrl() string {
    return EmoteWebUrl
}

func SetEmoteUrl(id, format string) {
    EmoteWebUrl = "https://static-cdn.jtvnw.net/emoticons/v2/" + id + "/" + format + "/light/4.0"
    notifyClients(EmoteWebUrl)
}
