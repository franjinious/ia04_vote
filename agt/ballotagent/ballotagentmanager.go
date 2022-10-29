package ballotagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ia04-vote/agt/sponsoragent"
	"ia04-vote/agt/voteragent"
	"ia04-vote/comsoc"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Ballotagentmanager struct {
	IP string
	Port string
	Ballotagents map[string]Ballotagent
	NowID int
}

func (bs *Ballotagentmanager)handlerNewBallot(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime )
	log.Println(": Get a new ballot request")

	var resp sponsoragent.Response
	var re sponsoragent.Request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &re)
	if err != nil {
		resp.ID = "none"
		resp.Status = 400
		w.WriteHeader(http.StatusOK)
	}else {
		var a sync.Mutex
		var c comsoc.Profile
		id := "vote"
		id += strconv.Itoa(bs.NowID)

		b := Ballotagent{a,re.Info,make([]voteragent.Voteragent,0),
			c,id, false}
		bs.Ballotagents[id] = b
		resp.ID = id
		bs.NowID++
		resp.Status = 201
		w.WriteHeader(http.StatusOK)
		log.Println(": Create a new ballot " + resp.ID)
	}

	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (bs *Ballotagentmanager) Start(){
	banner := "  ___    _    ___  _  _      __     __    _       \n " +
		   "|_ _|  / \\  / _ \\| || |     \\ \\   / /__ | |_ ___ \n  " +
		   "| |  / _ \\| | | | || |_ ____\\ \\ / / _ \\| __/ _ \\\n  " +
		  "| | / ___ \\ |_| |__   _|_____\\ V / (_) | ||  __/\n " +
		"|___/_/   \\_\\___/   |_|        \\_/ \\___/ \\__\\___|"

	fmt.Println(banner)
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", bs.handlerNewBallot)

	s := &http.Server{
		Addr:           "127.0.0.1:8082",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.SetFlags(log.Ldate | log.Ltime )
	log.Println(": start listen on \""+ bs.IP + ":" + bs.Port +"\"")
	go log.Fatal(s.ListenAndServe())
}