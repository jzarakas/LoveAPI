package main

import (
	"encoding/json"
	_ "fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type queuedSong struct {
	Que          int64
	SongId       string
	StartTime    time.Time
	TimeDuration time.Duration
}

type responseNode struct {
	Que     int64
	SongId  string
	StartAt float64
}

var queue = make([]queuedSong, 100)
var queLen int64 = -1

var sendError = responseNode{Que: -1, SongId: "error", StartAt: 0}

func main() {

	http.HandleFunc("/api/", welcome)
	http.HandleFunc("/getNext/", getNext)
	http.HandleFunc("/addSong/", addSong)
	http.HandleFunc("/getCurr/", getCurr)
	http.ListenAndServe(":8000", nil)

}

func getCurr(w http.ResponseWriter, r *http.Request) {
	if queLen < 0 || queue[queLen].StartTime.Add(queue[queLen].TimeDuration).Before(time.Now()) {
		js, _ := json.Marshal(sendError)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		for i := queLen; i >= 0; i-- {
			if queue[i].StartTime.Add(queue[i].TimeDuration).After(time.Now()) {

				var res = responseNode{Que: i, SongId: queue[i].SongId, StartAt: time.Now().Sub(queue[i].StartTime).Seconds()}
				js, _ := json.Marshal(res)

				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
				return

			}

		}
		js, _ := json.Marshal(sendError)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func addSong(w http.ResponseWriter, r *http.Request) {
	// /addSong/songId/
	if len(r.URL.Path) < 4 {
		return
	}
	song := strings.Split(r.URL.Path, "/")[2]
	dur := strings.Split(r.URL.Path, "/")[3]

	if queLen < 0 || queue[queLen].StartTime.Add(queue[queLen].TimeDuration).Before(time.Now()) {
		queLen = 0
		timedur, _ := time.ParseDuration(dur)
		queue[queLen] = queuedSong{Que: queLen, SongId: song, StartTime: time.Now(), TimeDuration: timedur}
	} else {
		queLen = queLen + 1
		timedur, _ := time.ParseDuration(dur)
		queue[queLen] = queuedSong{Que: queLen, SongId: song, StartTime: queue[queLen-1].StartTime.Add(queue[queLen-1].TimeDuration), TimeDuration: timedur}
	}

	w.Write([]byte("hello!"))

}

func welcome(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(sendError)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func getNext(w http.ResponseWriter, r *http.Request) {

	key := strings.Split(r.URL.Path, "/")[2]

	num, err := strconv.ParseInt(key, 10, 0)

	if err != nil || num+1 > queLen {
		js, _ := json.Marshal(sendError)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}

	zero, _ := time.ParseDuration("0m")
	var res = responseNode{Que: queue[num+1].Que, SongId: queue[num].SongId, StartAt: zero.Seconds()}
	js, err := json.Marshal(res)
	if err != nil {
		js, _ := json.Marshal(sendError)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func getCurrent(w http.ResponseWriter, r *http.Request) {

}

//func seedData() {
//	durOne, _ := time.ParseDuration("5m")
//	queue[0] = queuedSong{Que: 0, SongId: "songOne", StartTime: time.Now(), TimeDuration: durOne}
//
//	durTwo, _ := time.ParseDuration("10m")
//	queue[1] = queuedSong{Que: 1, SongId: "songTwo", StartTime: queue[0].StartTime.Add(queue[0].TimeDuration), TimeDuration: durTwo}
//	queue[2] = queuedSong{Que: 2, SongId: "songThree", StartTime: queue[1].StartTime.Add(queue[1].TimeDuration), TimeDuration: durOne}
//
//	queLen = 3
//}
