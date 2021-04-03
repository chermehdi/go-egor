package config

// Status corresponding to various verdicts for running a solution
const (
	AC int8 = iota
	SK
	RE
	WA
	TL
)

const (
	OK int8 = iota // OK execution status
	TO             // Timed out status
)

// WorkDir is the default work directory name
const WorkDir = "work"

// Time out delta in miliseconds
const TimeOutDelta int64 = 25
