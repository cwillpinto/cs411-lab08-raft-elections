# CS411 Lab 8: Raft Elections

Read [the assignment](Lab_08_Raft_Elections.md) before editing. The instructor-approved due time is **Wednesday, October 21, 2026, 11:59 PM America/Chicago**. Check the Blackboard Ultra assignment for its actual setting.

Instructor starter: https://github.com/cwillpinto/cs411-lab08-raft-elections

## Work in your own repository

If GitHub offers **Use this template**, create a repository under your own GitHub account, then clone that repository to your computer. You may also create an empty repository under your account and use the direct-clone workflow:

```sh
git clone https://github.com/cwillpinto/cs411-lab08-raft-elections.git
cd cs411-lab08-raft-elections
git remote rename origin upstream
git remote add origin <your-repository-url>
git remote -v
```

Verify that `origin` names **your** repository before any push. The instructor repository is an upstream source. Use your actual branch name when you push; for a `main` branch, run `git push -u origin main`.

## Run

Use Go 1.22 or newer. This starter has no third-party dependencies.

```sh
go run -race main.go node.go -nodes=5 -kill-leader=true
```

Complete the randomized election-timeout jitter and `RequestVote` in `node.go`. The starter intentionally does not pass the election check until those tasks are complete. The success token is `RAFT_LAB8_OK`, printed after two `ELECTION_STABLE` lines. The original starter's `main.go` already contains the failure-injection flag; follow the assignment's stated requirements and show its behavior in your run evidence.

## Submit in Blackboard Ultra

Submit `Lab8_Theory.pdf` or `.md`, your modified `node.go` and `main.go`, and the `RAFT_LAB8_OK` terminal evidence. Keep your work in your own GitHub repository; a repository link can supplement the original file submission, but it does not replace those graded deliverables unless the instructor updates Blackboard.
