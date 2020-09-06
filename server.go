package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "strconv"
    "time"
)

// cribbed from https://dev.to/moficodes/build-your-first-rest-api-with-go-2gcj

func getTime() string {
    t := time.Now()
    result := fmt.Sprintf("%d-%v-%d %v:%v.%v", t.Year(), t.Month(), t.Day(),
                                                t.Hour(), t.Minute(), t.Second())
    return result
}

func getIntegersParam(param string, req *http.Request) ([]int, error) {
    result := []int{}
    var err error
    values, ok := req.URL.Query()[param]

    if ok {
        for _, value := range values {
            if err == nil {
                card, e1 := strconv.Atoi(value)

                if e1 == nil {
                    result = append(result, card)
                } else {
                    err = e1
                }
            }
        }
    } else {
        err = errors.New("string param error for: " + param)
    }

    return result, err
}

func getStringParam(param string, req *http.Request) (string, error) {
    var result string
    var err error
    values, ok := req.URL.Query()[param]

    if ok && len(values) == 1 {
        result = values[0]
    } else {
        err = errors.New("string param error for: " + param)
    }

    return result, err
}

func getIntParam(param string, req *http.Request) (int, error) {
    var result int
    var err error
    resultStr, e1 := getStringParam(param, req)

    if e1 == nil {
        result, err = strconv.Atoi(resultStr)
    } else {
        err = e1
    }

    return result, err
}

func handleError(writer http.ResponseWriter) {
    writer.Header().Set("Content-Type", "application/json")
    writer.WriteHeader(http.StatusOK)

    const errorValue = -99
    now := getTime()

    // TODO: encode JSON
    body := strings.Builder{}
    body.WriteString("{");
    body.WriteString(fmt.Sprintf(`"message": "TRACER %v error",`, now))
    body.WriteString(fmt.Sprintf(`"card": %d`, errorValue))
    body.WriteString("},");

    bodyStr := body.String()
    fmt.Println(bodyStr)
    writer.Write([]byte(bodyStr))
}

const CARDS_PARAM = "cards"
const MAX_CARD_PARAM = "max_card"
const MODE_PARAM = "mode"
const PRIZE_CARD_PARAM = "prize_card"

type MyParams struct {
    cards []int
    maxCard int
    mode string
    prizeCard int
}

func (params *MyParams) String() string {
    return fmt.Sprintf("mode: %v prizeCard: %d maxCard: %d cards: %v",
            params.mode, params.prizeCard, params.maxCard, params.cards)
}

func getParams(req *http.Request) (MyParams, error) {
    result := MyParams{}
    var err error

    maxCard, e1 := getIntParam(MAX_CARD_PARAM, req)
    mode, e2 := getStringParam(MODE_PARAM, req)
    prizeCard, e3 := getIntParam(PRIZE_CARD_PARAM, req)
    cards, e4 := getIntegersParam(CARDS_PARAM, req)

    if e1 == nil && e2 == nil && e3 == nil && e4 == nil {
        result.maxCard = maxCard
        result.mode = mode
        result.prizeCard = prizeCard
        result.cards = cards
    } else {
        if e1 != nil { err = e1 }
        if e2 != nil { err = e2 }
        if e3 != nil { err = e3 }
        if e4 != nil { err = e4 }
    }

    return result, err
}

func handleGet(writer http.ResponseWriter, req *http.Request) {
    writer.Header().Set("Content-Type", "application/json")

    params, err := getParams(req)

    if err == nil {
        writer.WriteHeader(http.StatusOK)
        now := getTime()
        fmt.Printf("TRACER %v %v\n", now, params.String())
        body := strings.Builder{}

        // TODO: encode JSON
        body.WriteString("{");
        body.WriteString(fmt.Sprintf(`"card": 10`))
        body.WriteString("}");

        bodyStr := body.String()
        fmt.Println(bodyStr)
        writer.Write([]byte(bodyStr))
    } else {
        handleError(writer)
    }
}

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "GET":
        handleGet(w, r)
    case "POST":
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"message": "TRACER POST"}`))
    case "PUT":
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte(`{"message": "TRACER PUT"}`))
    case "DELETE":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "TRACER DELETE"}`))
    default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "TRACER not found"}`))
    }
}

func main() {
    numArgs := len(os.Args)

    if numArgs > 1 {
        port, err := strconv.Atoi(os.Args[1])

        if err != nil {
            log.Fatal(err)
        }

        portStr := fmt.Sprintf(":%d", port)
        fmt.Printf("TRACER running on port %d\n", port)

        http.HandleFunc("/waro/strategy", home)
        log.Fatal(http.ListenAndServe(portStr, nil))
    } else {
        fmt.Println("Usage.")
    }
}
