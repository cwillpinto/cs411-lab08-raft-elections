// Lab 8: Raft Elections -- starter.
//
// As shipped, resetElectionTimer() uses a hardcoded 150ms timeout for
// every node, so all nodes time out in lockstep and the cluster
// livelocks in permanent split votes. Fix the timing model (Part B,
// Requirement 1), implement RequestVote correctly (Requirement 2),
// and add -kill-leader support (Requirement 3, in main.go).
package main

import (
	"sync"
	"time"
)

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

// LogEntry is a minimal stand-in log entry (term is all that matters
// for the up-to-date check this lab exercises).
type LogEntry struct {
	Term int
}

type Node struct {
	mu sync.Mutex

	id          int
	peers       []*Node
	role        Role
	currentTerm int
	votedFor    int // -1 if no vote cast this term
	log         []LogEntry

	electionTimer *time.Timer
	killed        bool
	votesReceived int
	onBecomeLeader func()
}

func NewNode(id int) *Node {
	return &Node{
		id:       id,
		role:     Follower,
		votedFor: -1,
	}
}

// resetElectionTimer (re)starts this node's election timer.
func (n *Node) resetElectionTimer() {
	// TODO (Part B, Requirement 1): replace the hardcoded 150ms below
	// with a randomized duration drawn uniformly from [150ms, 300ms).
	// Call this every time the timer needs to restart: on election
	// timeout, on granting a vote, and on receiving a valid heartbeat.
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.electionTimer != nil {
		n.electionTimer.Stop()
	}
	timeout := 150 * time.Millisecond
	n.electionTimer = time.AfterFunc(timeout, n.startElection)
}

func (n *Node) lastLogTerm() int {
	if len(n.log) == 0 {
		return 0
	}
	return n.log[len(n.log)-1].Term
}

// startElection transitions this node to Candidate and requests votes
// from every peer.
func (n *Node) startElection() {
	n.mu.Lock()
	if n.killed {
		n.mu.Unlock()
		return
	}
	n.role = Candidate
	n.currentTerm++
	n.votedFor = n.id
	n.votesReceived = 1 // votes for itself
	term := n.currentTerm
	lastLogTerm := n.lastLogTerm()
	lastLogIndex := len(n.log)
	n.mu.Unlock()

	n.resetElectionTimer()

	for _, peer := range n.peers {
		if peer.id == n.id {
			continue
		}
		go func(p *Node) {
			granted := p.RequestVote(term, n.id, lastLogIndex, lastLogTerm)
			if granted {
				n.recordVote(term)
			}
		}(peer)
	}
}

func (n *Node) recordVote(term int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.role != Candidate || n.currentTerm != term {
		return // stale vote for a term/role we've already moved past
	}
	n.votesReceived++
	if n.votesReceived > (len(n.peers)+1)/2 && n.role != Leader {
		n.role = Leader
		if n.onBecomeLeader != nil {
			go n.onBecomeLeader()
		}
	}
}

// RequestVote implements the RequestVote RPC handler. See Part A,
// Question 4 for the log up-to-date check this must enforce.
func (n *Node) RequestVote(candidateTerm, candidateID, candidateLastLogIndex, candidateLastLogTerm int) bool {
	// TODO (Part B, Requirement 2):
	// 1. If candidateTerm < n.currentTerm, reject (stale term).
	// 2. If candidateTerm > n.currentTerm, update n.currentTerm to
	//    candidateTerm, reset n.votedFor to -1, and step down to
	//    Follower if not already.
	// 3. Grant the vote (set n.votedFor = candidateID, call
	//    n.resetElectionTimer(), return true) ONLY if n.votedFor is
	//    -1 or already candidateID, AND the candidate's log is at
	//    least as up-to-date as this node's own log (compare
	//    candidateLastLogTerm/candidateLastLogIndex against
	//    n.lastLogTerm()/len(n.log) using the standard Raft rule:
	//    higher last-log-term wins; if equal, longer log wins).
	// 4. Otherwise, reject (return false) without granting.
	//
	// Locking note: resetElectionTimer() takes n.mu itself, so you
	// must n.mu.Unlock() before calling it -- calling it while you're
	// still holding the lock will deadlock. Compute your grant/deny
	// decision under the lock, unlock, then call resetElectionTimer()
	// only if you're granting, then return.
	panic("RequestVote not implemented")
}

// ReceiveHeartbeat is called when a node hears from the current
// leader -- it must suppress this node's own election timer.
func (n *Node) ReceiveHeartbeat(leaderTerm int) {
	n.mu.Lock()
	if leaderTerm >= n.currentTerm {
		n.currentTerm = leaderTerm
		n.role = Follower
	}
	n.mu.Unlock()
	n.resetElectionTimer()
}

// Kill simulates a node crash: it stops responding to everything.
func (n *Node) Kill() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.killed = true
	if n.electionTimer != nil {
		n.electionTimer.Stop()
	}
}
