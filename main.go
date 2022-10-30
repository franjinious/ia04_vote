package main

import (
	"fmt"
	"ia04-vote/agt/ballotagent"
	"ia04-vote/test"
)

/**
 * main
 * @Description: Démarrer le serveur de vote
 */
func main() {
	// fmt.Println(time.Now().Format("Mon Jan 15:04:05 UTC 2006"))\
	ballotagent.StartVoteServer("127.0.0.1","8082")
	test.Test_newballot()
	test.Test_vote()
	test.Test_vote2()
	test.Test_vote3()
	test.Test_vote2()
	test.Test_vote4()

	fmt.Scan()
}
