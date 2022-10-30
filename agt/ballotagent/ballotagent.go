package ballotagent

import (
	"encoding/json"
	"ia04-vote/agt/sponsoragent"
	"ia04-vote/agt/voteragent"
	"ia04-vote/comsoc"
	"net/http"
	"sync"
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
}

// tous les vote algorithme
var method_scf = map[string]interface{} {
	"condorcet"     : comsoc.CondorcetWinner,
	"majority"      : comsoc.MajoritySCF,
	"borda"         : comsoc.BordaSCF,
	"kramersimpson" : comsoc.KramerSimpsonSCF,
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

func (b *Ballotagent) getNewVoteRequest(v voteragent.Voterinfo,w http.ResponseWriter){
	b.Lock()
	var resp voteragent.Response

	if b.Isfinish == true {
		resp.Status = 503
		w.WriteHeader(http.StatusOK)
	} else {
		if b.Voters[v.Agent_ID] == false {
			resp.Status = 403
			w.WriteHeader(http.StatusOK)
		} else if len(v.Prefs) != b.Sponsor.Alts {
			resp.Status = 400
			w.WriteHeader(http.StatusOK)
		} else {
			flag := true
			for i:=0; i < len(v.Prefs); i++{
				pt := (*int)(&v.Prefs[i])
				if v.Prefs[i] <= 0 || *pt > b.Sponsor.Alts {
					resp.Status = 400
					w.WriteHeader(http.StatusOK)
					flag = false
					break
				}
			}

			if flag == true {
				b.Voters[v.Agent_ID] = false
				b.Voterinfos = append(b.Voterinfos,v)
				resp.Status = 200
				b.p = append(b.p,v.Prefs)
				w.WriteHeader(http.StatusOK)
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

	b.Isfinish = true
	for _,j := range b.Voters {
		if j != false {
			b.Isfinish = false
		}
	}

	if b.Isfinish != true {
		resp.Status = TooEarly
		resp.Winner = -1
		resp.Ranking = nil
	} else {
		fun_scf := method_scf[b.Sponsor.Rule]
		switch f := fun_scf.(type) {
		case func(comsoc.Profile)([]comsoc.Alternative,error):
			ans,e := f(b.p)
			if e != nil {
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
			}
		}
	}

	serial, _ := json.Marshal(resp)
	w.Write(serial)
	b.Unlock()
}