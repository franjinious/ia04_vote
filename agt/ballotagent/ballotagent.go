package ballotagent

import (
	"ia04-vote/agt/sponsoragent"
	"sync"
)
import "ia04-vote/agt/voteragent"
import "ia04-vote/comsoc"
//
// agent pour gérer les vote
//

type Ballotagent struct {
	sync.Mutex
	Sponsor sponsoragent.Sponsorinfo
	Voters []voteragent.Voteragent
	p comsoc.Profile
	ID string
	Isfinish bool
}

// tous les vote algorithme
var method = map[string]int {
	"condorcet"     :0,
	"majority"      :1,
	"borda"         :2,
	"kramersimpson" :3,
	"copeland"      :4,
	"coombs"        :5,
	"stv"           :6,
	"kemeny"        :7,
	"singlepeak"    :8,
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







