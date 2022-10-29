package voteragent

//
// agent pour voter
//

type Voteragent struct {

}

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