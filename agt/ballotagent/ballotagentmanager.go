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
	sync.Mutex
	IP string
	Port string
	Ballotagents map[string]*Ballotagent
	NowID int
}

func (bs *Ballotagentmanager)handlerNewBallot(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime )
	log.Println(": get a new ballot request")

	var resp sponsoragent.Response
	var re sponsoragent.Request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &re)
	if err != nil {
		resp.ID = "none"
		resp.Status = 400
		w.WriteHeader(http.StatusBadRequest)
	}else {
		var a sync.Mutex
		id := "vote"

		bs.Lock()
		id += strconv.Itoa(bs.NowID)
		b := Ballotagent{a,re.Info,make([]voteragent.Voterinfo,0),
			make(map[string]bool),make(comsoc.Profile,0),id, false}
		for i := 1; i <= re.Info.Alts; i++ {
			id := "ag_id"
			id += strconv.Itoa(i)
				b.Voters[id] = true
		}
		bs.Ballotagents[id] = &b
		resp.ID = id
		bs.NowID++
		bs.Unlock()

		resp.Status = 201
		w.WriteHeader(http.StatusOK)
		log.Println(": Create a new ballot " + resp.ID)
	}

	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (bs *Ballotagentmanager)handlerVoteRequest(w http.ResponseWriter, r *http.Request){
	bs.Lock()
	log.SetFlags(log.Ldate | log.Ltime )
	var resp voteragent.Response
	var re voteragent.Request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &re)

	if err != nil {
		resp.Status = 400
		w.WriteHeader(http.StatusBadRequest)
		serial, _ := json.Marshal(resp)
		w.Write(serial)
		bs.Unlock()
		return
	}
	log.Println(": Get a new vote request of " + re.Info.Vote_ID + ", from " + re.Info.Agent_ID)

	agent := bs.Ballotagents[re.Info.Vote_ID]
	agent.getNewVoteRequest(re.Info,w)
	bs.Unlock()
}

func (bs *Ballotagentmanager)handlerResultRequest(w http.ResponseWriter, r *http.Request){
	bs.Lock()
	log.SetFlags(log.Ldate | log.Ltime )
	var resp voteragent.Response_Result
	var re voteragent.Request_Result
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &re)
	if err != nil {
		resp.Status = 400
		w.WriteHeader(http.StatusBadRequest)
		serial, _ := json.Marshal(resp)
		w.Write(serial)
		bs.Unlock()
		return
	}
	log.Println(": Get a new result request of " + re.Ballot_Id)
	var agent *Ballotagent
	if _, ok := bs.Ballotagents[re.Ballot_Id]; ok {
		agent = bs.Ballotagents[re.Ballot_Id]
		agent.getNewResultRequest(re.Ballot_Id,w)
		bs.Unlock()
	} else {
		resp.Status = 404
		w.WriteHeader(http.StatusOK)
		serial, _ := json.Marshal(resp)
		w.Write(serial)
		bs.Unlock()
		return
	}
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
	mux.HandleFunc("/vote", bs.handlerVoteRequest)
	mux.HandleFunc("/result", bs.handlerResultRequest)

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

func StartVoteServer(IP string,Port string){
	var mutex sync.Mutex
	bs := Ballotagentmanager{mutex,IP,Port,make(map[string]*Ballotagent),0}
	go bs.Start()
}
