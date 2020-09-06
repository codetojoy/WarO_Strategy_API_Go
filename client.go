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
    baseUrl := fmt.Sprintf("http://localhost:%d%v", port, BASE_URL)

    req, err := http.NewRequest("GET", baseUrl, nil)

    if err != nil {
        log.Fatal(err)
    }

    q := req.URL.Query()
    q.Add(MODE_PARAM, "min")
    q.Add(PRIZE_CARD_PARAM, "10")
    q.Add(MAX_CARD_PARAM, "20")
    q.Add(CARDS_PARAM, "3")
    q.Add(CARDS_PARAM, "7")
    q.Add(CARDS_PARAM, "11")
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
