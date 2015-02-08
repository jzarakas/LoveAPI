package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

)

type queuedSong struct {
    Que int64
    SongId string
    StartTime time.Time 
    TimeDuration time.Duration
}

type responseNode struct {
    Que int64
    SongId string
    StartAt time.Duration
}


var queue = make([]queuedSong, 100)



func main() {

    seedData()

    http.HandleFunc("/api", welcome)
    http.ListenAndServe(":8000", nil)

}

func welcome(w http.ResponseWriter, r * http.Request){
    
    //res = responseNode(Que: queue[0].Que, SongId: queue[0].SongId, StartAt: time.Now().Sub(queue[0].StartTime))
    zero, _ := time.ParseDuration("0m")

    var res  = responseNode{Que: queue[0].Que, SongId: queue[0].SongId, StartAt: zero}

    js, err := json.Marshal(res)

    fmt.Println(js);

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
    
}


func seedData() {
    durOne, _ := time.ParseDuration("5m")
    queue[0] = queuedSong{Que: 0, SongId: "songOne", StartTime: time.Now(), TimeDuration: durOne}

    durTwo, _ := time.ParseDuration("10m")
    queue[1] = queuedSong{Que: 1, SongId: "songTwo", StartTime: queue[0].StartTime.Add(queue[0].TimeDuration), TimeDuration: durTwo}
    queue[2] = queuedSong{Que: 2, SongId: "songThree", StartTime: queue[1].StartTime.Add(queue[1].TimeDuration), TimeDuration: durOne}
}
