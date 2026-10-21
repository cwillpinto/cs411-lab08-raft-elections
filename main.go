// Lab 8: Raft Elections -- main.go.
//
// Run: go run -race . -nodes=5 -kill-leader=true
package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	numNodes := flag.Int("nodes", 5, "number of nodes in the simulated cluster")
	killLeader := flag.Bool("kill-leader", false, "kill the stable leader once elected and verify re-election")
	flag.Parse()

	nodes := make([]*Node, *numNodes)
	for i := 0; i < *numNodes; i++ {
		nodes[i] = NewNode(i)
	}
	for _, n := range nodes {
		n.peers = nodes
	}

	leaderStableCh := make(chan *Node, 1)
	for _, n := range nodes {
		node := n
		running := false
		node.onBecomeLeader = func() {
			// Guard against overlapping invocations: recordVote can call
			// this again if the node loses and re-wins leadership in a
			// later term before a prior invocation's ticker loop has
			// noticed (via stillLeader) that it should exit. Without this
			// guard, two goroutines race on a shared heartbeat counter.
			node.mu.Lock()
			if running {
				node.mu.Unlock()
				return
			}
			running = true
			node.mu.Unlock()
			defer func() {
				node.mu.Lock()
				running = false
				node.mu.Unlock()
			}()

			// A "stable" leader is one that has sent at least 3
			// consecutive heartbeats without another node winning an
			// election in the meantime (Part B, Requirement 3).
			// heartbeatsSent lives inside this invocation so a fresh
			// leadership term always starts its stability count at 0.
			heartbeatsSent := 0
			ticker := time.NewTicker(50 * time.Millisecond)
			defer ticker.Stop()
			for range ticker.C {
				node.mu.Lock()
				stillLeader := node.role == Leader && !node.killed
				term := node.currentTerm
				node.mu.Unlock()
				if !stillLeader {
					return
				}
				for _, peer := range node.peers {
					if peer.id != node.id {
						go peer.ReceiveHeartbeat(term)
					}
				}
				heartbeatsSent++
				if heartbeatsSent == 3 {
					select {
					case leaderStableCh <- node:
					default:
					}
				}
			}
		}
	}

	for _, n := range nodes {
		n.resetElectionTimer()
	}

	firstLeader := <-leaderStableCh
	firstLeader.mu.Lock()
	fmt.Printf("ELECTION_STABLE term=%d leader=%d\n", firstLeader.currentTerm, firstLeader.id)
	firstLeader.mu.Unlock()

	if !*killLeader {
		return
	}

	firstLeader.Kill()

	// Drain any stale signal from before the kill, then wait for the
	// remaining nodes to elect a new stable leader.
	select {
	case <-leaderStableCh:
	default:
	}

	secondLeader := <-leaderStableCh
	secondLeader.mu.Lock()
	fmt.Printf("ELECTION_STABLE term=%d leader=%d\n", secondLeader.currentTerm, secondLeader.id)
	secondLeader.mu.Unlock()

	if secondLeader.id == firstLeader.id {
		fmt.Println("ERROR: re-elected the same (killed) leader -- something is wrong with Kill() or the election logic")
		return
	}

	fmt.Println("SUCCESS TOKEN: RAFT_LAB8_OK")
}
