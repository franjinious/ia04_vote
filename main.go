package main

import "ia04-vote/agt/ballotagent"

/**
 * main
 * @Description: Démarrer le serveur de vote
 */
func main() {
	ballotagent.StartVoteServer("127.0.0.1","8082")
}
