package test

import (
	"fmt"
	"ia04-vote/agt/sponsoragent"
	"time"
)

func Test_newballot(){
	time.Sleep(2*time.Second)
	s1 := []string{"ag_id1","ag_id2","ag_id3"}
	g := sponsoragent.Sponsorinfo{"majority","Mon Jan 15:04:05 UTC 2006",s1,3}
	p := sponsoragent.Sponsoragent{g,"127.0.0.1:8082","none"}
	p.New_ballot()
	fmt.Println(p.ID)

	time.Sleep(2*time.Second)
	s2 := []string{"ag_id1","ag_id2","ag_id3"}
	g1 := sponsoragent.Sponsorinfo{"majority","Mon Jan 15:04:05 UTC 2006",s2,3}
	p1 := sponsoragent.Sponsoragent{g1,"127.0.0.1:8082","none"}
	p1.New_ballot()
	fmt.Println(p1.ID)
}