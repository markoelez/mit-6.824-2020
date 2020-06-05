package raft

type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

type ServerState int

const (
	Leader ServerState = iota
	Follower
	Candidate
)

type Log struct {
	Term    int
	Command interface{}
}

type RequestVoteArgs struct {
}

type RequestVoteReply struct {
}
