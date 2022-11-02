package voteragent

import (
	"bytes"
	"encoding/json"
	"errors"
	"gitlab.utc.fr/wanhongz/ia04-vote/comsoc"
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


func (v *Voteragent) Vote() error{
	req := Voterinfo{
		 v.Agent_ID,v.Vote_ID,v.Prefs,v.Options,
	}

	// fmt.Println(req.Info)
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

	log.SetFlags(log.Ldate | log.Ltime )
	if resp.StatusCode == VoteSuccess {
		log.Println(": " + req.Agent_ID + " vote successfully for " + req.Vote_ID)
	}else if resp.StatusCode == BadRequest {
		log.Println(": " +req.Agent_ID + " request failed")
		return errors.New("request failed")
	}else if resp.StatusCode == UselessVote {
		log.Println(": " +req.Agent_ID + " you have already voted")
		return errors.New("vote exist")
	}else if resp.StatusCode == NotImplemented {
		log.Println(": " +req.Agent_ID + " function has no implemented")
		return errors.New("not implemented")
	}else {
		log.Println(": " +req.Agent_ID + " vote " + req.Vote_ID + " has finished")
		return errors.New("time out")
	}

	return nil
}

type Request_Result struct {
	Ballot_Id string `json:"ballot_Id"`
}

type Response_Result struct {
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

	if resp.StatusCode == OK {
		ou := (": " + "get vote result, " + strconv.Itoa(*(*int)(&re.Winner)) + " win")
		if re.Ranking != nil {
			out := ", ranking is [ "
			for _,j := range re.Ranking {
				temp := *(*int)(&j)
				out += strconv.Itoa(temp)
				out += " "
			}
			out += "]"
			log.Print(ou + out)
		}
	}else if resp.StatusCode == TooEarly {
		log.Println(": " + "vote has not finished")
		return errors.New("Too Early")
	}else if resp.StatusCode == Notfind {
		log.Println(": " + "not find this function")
		return errors.New("Not find")
	}

	return nil
}