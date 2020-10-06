package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
	"time"
	"fmt"
)

const arrayLen = 5

type Counter struct {
	currentIndex uint8
	list [arrayLen]int64
}

func (c *Counter) Append(unixTime int64) {
	c.list[c.currentIndex] = unixTime
	c.incIndex()

}

var t = struct {
	sync.RWMutex
	bots map[string]*Counter
}{bots: make(map[string]*Counter)}


func (c *Counter) incIndex() {
	c.currentIndex++
	if c.currentIndex == arrayLen {
		c.currentIndex = 0
		c.list = [arrayLen]int64{}
	}
}

func (c *Counter) IsMore100reqInMinute() bool {
	return time.Now().Unix() - c.list[arrayLen	 - c.currentIndex] <= 60
}

func IncUserVisit(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	t.RLock()
	bot, ok := t.bots[userId]
	t.RUnlock()
	if !ok {
		t.Lock()
		t.bots[userId] =  &Counter{currentIndex: 0, list: [arrayLen]int64{time.Now().Unix()}}
		t.Unlock()
		fmt.Println(t.bots[userId])
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Println(t.bots[userId])
	bot.Append(time.Now().Unix())
	w.WriteHeader(http.StatusOK)
	return
}

func Count(w http.ResponseWriter, r *http.Request) {
	var userCnt int = 0
	for _, val := range t.bots {
		if val.IsMore100reqInMinute() {
			userCnt++
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(userCnt)))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", IncUserVisit).Methods("GET")
	router.HandleFunc("/count", Count).Methods("GET")
	http.ListenAndServe(":8080", router)
}
