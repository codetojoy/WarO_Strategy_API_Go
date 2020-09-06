package main

import (
    "bufio"
    "fmt"
    "log"
    "net/http"
    // "os"
    "strings"
)

func callServer(ch chan string, url string) {
    resp, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("TRACER response status:", resp.Status)
    result := strings.Builder{}

    // this is weird:
    scanner := bufio.NewScanner(resp.Body)
    for i := 0; scanner.Scan() && i < 5; i++ {
        result.WriteString(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }

    ch <- result.String()
}

func buildUrl(port int) string {
    baseUrl := fmt.Sprintf("http://localhost:%d/waro/strategy", port)

    req, err := http.NewRequest("GET", baseUrl, nil)

    if err != nil {
        log.Fatal(err)
    }

    q := req.URL.Query()
    q.Add("mode", "max")
    q.Add("prize_card", "10")
    q.Add("max_card", "20")
    q.Add("cards", "3")
    q.Add("cards", "7")
    q.Add("cards", "11")
    req.URL.RawQuery = q.Encode()

    return req.URL.String()
}

func main() {
    const port = 6160
    url := buildUrl(port)
    myChannel := make(chan string)

    go callServer(myChannel, url)
    fmt.Println("TRACER result:")
    fmt.Println(<- myChannel)
    fmt.Println("Ready.")
}
