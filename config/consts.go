package config

// Status corresponding to various verdicts for running a solution
var (
	AC int8 = 0
	SK int8 = 1
	RE int8 = 2
	WA int8 = 3
	TL int8 = 4
)

var (
	OK int8 = 0 // OK execution status
	TO int8 = 1 // Timed out status
)

// WorkDir is the default work directory name
const WorkDir = "work"

// Time out delta in miliseconds
const TimeOutDelta float64 = 25.0
