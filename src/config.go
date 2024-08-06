package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)


func Load_file(filePath string) map[string]string{
    // 読み込むファイルのパス

    var config map[string]string = make(map[string]string)
    var tmp [4]string

    // ファイルを開く
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Println("エラー:", err)
    }
    defer file.Close()

    // スキャナを使ってファイルを読み込む
    scanner := bufio.NewScanner(file)
    i := 0
    for scanner.Scan() {
        // 1行ずつ読み込む
        line := scanner.Text()

        fields := strings.Fields(line)

        if len(fields) >= 2 {
            tmp[i] = fields[1]
        } else {
            fmt.Println("2つ目のフィールドが存在しません")
        }

        i++

    }

    // エラーがあった場合
    if err := scanner.Err(); err != nil {
        fmt.Println("エラー:", err)
    }

    config["client_id"] = tmp[0]
    config["access_token"] = tmp[1]
    config["broadcaster_id"] = tmp[2]
    config["sender_id"] = tmp[3]

    return config
}
