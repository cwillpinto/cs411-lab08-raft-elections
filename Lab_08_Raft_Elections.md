# Lab 8: Raft Elections
## CS 411: Computing Paradigms (Fall 2026)

**Instructor starter:** https://github.com/cwillpinto/cs411-lab08-raft-elections

**Due:** Wednesday, October 21, 2026, 11:59 PM America/Chicago. The live Blackboard Ultra setting must be checked there.

### ⚡ The Mission
It's 11:47 PM during quarter close. Your five-node payments ledger cluster just lost its leader node to a rack power failure, and two other nodes are dropping packets on a flaky switch. The audit log **cannot** fork into two competing "truths" — a single bad election could double-commit a transaction. You've been handed a bare-bones Raft simulation framework where every node's election timer fires at the exact same instant, which means the cluster is currently deadlocked in a permanent split-vote livelock. Fix the timing model, verify the cluster still elects exactly one leader per term even when you kill the sitting leader mid-run, and prove it holds together before the next close cycle.

---

### Part A: Brief Systems Checks (50 Points)
**Instructions:** Submit concise answers in `Lab8_Theory.pdf` or `Lab8_Theory.md`: one short vote table, one calculation with a brief interpretation, and a 2–3 sentence safety explanation.

1. **Vote trace.** N4 crashes in term 6. N2 times out at 180ms, N5 at 185ms, N1 at 310ms, and N3 at 340ms. N5's log is older than N2's; N1 and N3 have logs at least as up to date as N2's. In a short table, show possible votes and terms under the `RequestVote` term and log checks. State which RPC timing or log facts are still needed before naming a leader.
2. **Split vote and retry.** A five-node election splits 2–2–1 in term 7. State what candidates do next. For independent timeouts redrawn from [150ms, 300ms), estimate the chance that two draws fall within 5ms and explain in one sentence why another split is possible but less likely.
3. **Committed-entry safety.** Entry 7 is committed on 3 of 5 nodes before the leader crashes. In 2–3 sentences, explain how the up-to-date-log voting rule keeps a newly elected leader from missing that committed entry, including why an up-to-date voter rejects a stale candidate even if its term is higher.

---

### Part B: The Implementation (50 Points)
**Instructions:** Clone the provided `raft-sim` starter (a runnable multi-goroutine Raft simulation, *not* a production consensus library — that's Project 2). As shipped, every node's `resetElectionTimer()` uses a hardcoded 150ms timeout, so all five nodes time out in lockstep and the cluster livelocks in permanent split votes. Your job is to fix the timing model and add controlled failure injection so the cluster provably self-heals.

**Requirements:**
1. In `node.go`, replace the fixed 150ms election timeout with a randomized jitter drawn uniformly from `[150ms, 300ms)` per node, re-rolled every time the timer resets (on election timeout, on granting a vote, and on receiving a valid heartbeat).
2. Implement the `RequestVote` RPC handler so each node casts at most one vote per term, rejects requests from stale terms, and enforces the log up-to-date check from Part A, Question 3.
3. Implement the `-kill-leader` CLI flag: once a stable leader has emitted at least 3 consecutive heartbeats, the simulation must terminate that leader's goroutine (simulating a crash) and let the remaining 4 nodes re-elect.
4. **Verification:** Run `go run -race . -nodes=5 -kill-leader=true`. The harness prints `ELECTION_STABLE term=<N> leader=<id>` twice — once for the initial election, once after the induced crash — followed by the **Success Token** `RAFT_LAB8_OK`. If you see it, your implementation is correct.

---

### 🚀 Submission
1. `Lab8_Theory.pdf` (or `.md`)
2. `node.go` and `main.go` (your modified source)
3. The **Success Token** (`RAFT_LAB8_OK`) captured as text or a terminal screenshot
