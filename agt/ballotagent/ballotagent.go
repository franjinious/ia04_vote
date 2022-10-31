package ballotagent

import (
	"encoding/json"
	"fmt"
	"gitlab.utc.fr/wanhongz/ia04-vote/agt/sponsoragent"
	"gitlab.utc.fr/wanhongz/ia04-vote/agt/voteragent"
	"gitlab.utc.fr/wanhongz/ia04-vote/comsoc"
	"log"
	"net/http"
	"sync"
	"time"
)
//
// agent pour gérer les vote
//

type Ballotagent struct {
	sync.Mutex
	Sponsor sponsoragent.Sponsorinfo
	Voterinfos []voteragent.Voterinfo
	Voters map[string]bool
	p comsoc.Profile
	ID string
	Isfinish bool
	Expiration int
	seuil []int
}

const (
	VoteCreateSuccess = 201
	BadRequest = 400
	NotImplemented = 501

	VoteSuccess = 200
	UselessVote = 403
	TimeOut = 503

	OK = 200
	TooEarly = 425
	Notfind = 404
)

// tous les vote algorithme
var method_scf = map[string]interface{} {
	"condorcet"     : comsoc.CondorcetWinner,
	"majority"      : comsoc.MajoritySCF,
	"borda"         : comsoc.BordaSCF,
	"kramersimpson" : comsoc.KramerSimpsonSCF,
	"approval"      : comsoc.ApprovalSCF,
	"copeland"      : comsoc.CopelandSCF,
	"coombs"        : comsoc.CoombsSCF,
	"stv"           : comsoc.STV_SCF,
	"kemeny"        : comsoc.Kemeny_SCF,
	"singlepeak"    : comsoc.SinglePeakedSCF,
}

var method_swf = map[string]interface{} {
	"majority"      : comsoc.MajoritySWF,
	"borda"         : comsoc.BordaSWF,
	"kramersimpson" : comsoc.KramerSimpsonSWF,
	"copeland"      : comsoc.CopelandSWF,
	"coombs"        : comsoc.CoombsSWF,
	"stv"           : comsoc.STV_SWF,
	"kemeny"        : comsoc.Kemeny_SWF,
}

func (b *Ballotagent) getNewVoteRequest(v voteragent.Voterinfo,w http.ResponseWriter){
	b.Lock()
	var resp voteragent.Response

	if b.Isfinish == true {
		resp.Status = 503
		w.WriteHeader(http.StatusOK)
		log.Println(": Get a new vote request of " + v.Vote_ID + ", from " + v.Agent_ID + ", " +
			"vote failed because it has finished")
	} else {
		if b.Voters[v.Agent_ID] == false {
			resp.Status = 403
			w.WriteHeader(http.StatusOK)
			log.Println(": Get a new vote request of " + v.Vote_ID + ", from " + v.Agent_ID + ", " +
				"vote failed because the voter do not exist or he has already voted")
		} else if len(v.Prefs) != b.Sponsor.Alts {
			resp.Status = 400
			w.WriteHeader(http.StatusOK)
			log.Println(": Get a new vote request of " + v.Vote_ID + ", from " + v.Agent_ID + ", " +
				"vote failed because the prefer list is not valid")
		} else {
			flag := true
			for i:=0; i < len(v.Prefs); i++{
				pt := (*int)(&v.Prefs[i])
				if v.Prefs[i] <= 0 || *pt > b.Sponsor.Alts {
					resp.Status = 400
					w.WriteHeader(http.StatusOK)
					flag = false
					log.Println(": Get a new vote request of " + v.Vote_ID + ", from " + v.Agent_ID + ", " +
						"vote failed because the prefer list is not valid")
					break
				}
			}

			if flag == true {
				b.Voters[v.Agent_ID] = false
				b.Voterinfos = append(b.Voterinfos,v)
				resp.Status = 200
				b.p = append(b.p,v.Prefs)
				if v.Options!=nil && b.Sponsor.Rule=="approval" {
					b.seuil = append(b.seuil,v.Options[0])
				}
				w.WriteHeader(http.StatusOK)
				log.Println(": Get a new vote request of " + v.Vote_ID + ", from " + v.Agent_ID + ", " +
					"vote successfully")
			}
		}
	}

	serial, _ := json.Marshal(resp)
	w.Write(serial)
	b.Unlock()
}

func (b *Ballotagent) getNewResultRequest(ID string,w http.ResponseWriter){
	b.Lock()
	var resp voteragent.Response_Result

	if b.Isfinish == false {
		b.Isfinish = true
		for _,j := range b.Voters {
			if j != false {
				b.Isfinish = false
			}
		}
	}

	if b.Isfinish != true {
		resp.Status = TooEarly
		resp.Winner = -1
		resp.Ranking = nil
		log.Println(": Get a new result request of " + b.ID +
			" result failed because the vote has not finished")
	} else {
		fun_scf := method_scf[b.Sponsor.Rule]
		switch f := fun_scf.(type) {
		case func(comsoc.Profile)([]comsoc.Alternative,error):
			ans,e := f(b.p)
			if e != nil {
				log.Println(": Get a new result request of " + b.ID +
					 ", but it has no result bacause error: " + e.Error())
				resp.Status = 404
				resp.Winner = -1
				resp.Ranking = nil
			} else {
				resp.Status = 200
				resp.Winner = ans[0]
				resp.Ranking = nil
			}
		case func(comsoc.Profile,[]int)([]comsoc.Alternative,error):
			ans,e := f(b.p,b.seuil)
			if e != nil {
				log.Println(": Get a new result request of " + b.ID +
					", but it has no result bacause error: " + e.Error())
				resp.Status = 404
				resp.Winner = -1
				resp.Ranking = nil
			} else {
				resp.Status = 200
				resp.Winner = ans[0]
				resp.Ranking = nil
			}
		}


		if _, ok := method_swf[b.Sponsor.Rule]; ok {
			switch f := method_swf[b.Sponsor.Rule].(type) {
			case func(comsoc.Profile)(comsoc.Count,error):
				ans,e := f(b.p)
				if e != nil {
					resp.Ranking = nil
				} else {
					resp.Ranking = comsoc.SortByCount(ans)
				}
			case func(comsoc.Profile,[]int)(comsoc.Count,error):
				ans,e := f(b.p,b.seuil)
				if e != nil {
					resp.Ranking = nil
				} else {
					resp.Ranking = comsoc.SortByCount(ans)
				}
			case func(comsoc.Profile)([]comsoc.Alternative,error):
				ans,e := f(b.p)
				if e != nil {
					resp.Ranking = nil
				} else {
					resp.Ranking = ans
				}
			}
		}
	}

	serial, _ := json.Marshal(resp)
	w.Write(serial)
	b.Unlock()
}

func (b* Ballotagent) SetFinished() {
	timer := time.After(time.Duration(b.Expiration) * time.Second)
	<-timer
	fmt.Println(b.ID + " has finished")
	b.Isfinish = true
}