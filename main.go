package main

import "ia04-vote/agt/ballotagent"
import "ia04-vote/test"

/**
 * main
 * @Description: Démarrer le serveur de vote
 */
func main() {
	// fmt.Println(time.Now().Format("Mon Jan 15:04:05 UTC 2006"))\
	bs := ballotagent.Ballotagentmanager{"127.0.0.1","8082",make(map[string]ballotagent.Ballotagent),0}
	go bs.Start()
	test.Test_newballot()
}
