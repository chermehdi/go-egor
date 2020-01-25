package config

type TestCase struct {
	Input  string
	Output string
}

type IOType struct {
	Type string
}

type LanguageDescription struct {
	MainClass string
	TaskClass string
}

// Task representation
// Struct mapping to the json object received by competitive companion.
type Task struct {
	Name        string
	Group       string
	Url         string
	Interactive bool
	MemoryLimit float64
	TimeLimit   float64
	Tests       []TestCase
	TestType    string
	Input       IOType
	Output      IOType
	Languages   map[string]LanguageDescription
}
