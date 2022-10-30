package voteragent

import (
	"bytes"
	"encoding/json"
	"errors"
	"ia04-vote/comsoc"
	"log"
	"net/http"
	"strconv"
	"sync"
)

//
// agent pour voter
//

const (
	// for api /vote
	VoteSuccess = 200
	BadRequest = 400
	UselessVote = 403
	NotImplemented = 501
	TimeOut = 503

	// for api /result
	OK = 200
	TooEarly = 425
	Notfind = 404
)

type Voterinfo struct {
	Agent_ID string `json:"agent_id"`// e.g. "ag_id1"
	Vote_ID string `json:"vote_id"`// e.g. "vote12"
	Prefs []comsoc.Alternative `json:"prefs"` // e.g. [1, 2, 4, 3]
	Options []int `json:"options"`
}

type Voteragent struct {
	sync.Mutex
	ServerAddress string
	Voterinfo
}

func newVoteragent(mutex sync.Mutex, serverAddress string, voterinfo Voterinfo) *Voteragent {
	return &Voteragent{Mutex: mutex, ServerAddress: serverAddress, Voterinfo: voterinfo}
}

type Request struct {
	Info Voterinfo `json:"info"`
}

type Response struct {
	Status int `json:"status"`
}

func (v *Voteragent) Vote() error{
	req := Request{
		Info: Voterinfo{v.Agent_ID,v.Vote_ID,v.Prefs,v.Options},
	}

	url := "http://" + v.ServerAddress + "/vote"
	data, e := json.Marshal(req)
	if e != nil {
		return e
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	var re Response
	json.Unmarshal(buf.Bytes(), &re)
	log.SetFlags(log.Ldate | log.Ltime )
	if re.Status == VoteSuccess {
		log.Println(": " + req.Info.Agent_ID + " vote successfully for " + req.Info.Vote_ID)
	}else if re.Status == BadRequest {
		log.Println(": " +req.Info.Agent_ID + " request failed")
		return errors.New("request failed")
	}else if re.Status == UselessVote {
		log.Println(": " +req.Info.Agent_ID + " you have already voted")
		return errors.New("vote exist")
	}else if re.Status == NotImplemented {
		log.Println(": " +req.Info.Agent_ID + " function has no implemented")
		return errors.New("not implemented")
	}else {
		log.Println(": " +req.Info.Agent_ID + " vote " + req.Info.Vote_ID + " has finished")
		return errors.New("time out")
	}

	return nil
}

type Request_Result struct {
	Ballot_Id string `json:"ballot_Id"`
}

type Response_Result struct {
	Status int `json:"status"`
	Winner comsoc.Alternative `json:"winner"`
	Ranking []comsoc.Alternative `json:"ranking"`
}

func (v *Voteragent) Result() error{
	req := Request_Result{
		Ballot_Id: v.Vote_ID,
	}
	url := "http://" + v.ServerAddress + "/result"
	data, e := json.Marshal(req)
	if e != nil {
		return e
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	var re Response_Result
	json.Unmarshal(buf.Bytes(), &re)
	log.SetFlags(log.Ldate | log.Ltime )

	if re.Status == OK {
		log.Println(": " + "get vote result, " + strconv.Itoa(*(*int)(&re.Winner)) + " win")
		if re.Ranking != nil {
			log.Print(": ranking is ")
			log.Println(re.Ranking)
		}
	}else if re.Status == TooEarly {
		log.Println(": " + "vote has not finished")
		return errors.New("Too Early")
	}else if re.Status == Notfind {
		log.Println(": " + "not find this function")
		return errors.New("Not find")
	}

	return nil
}