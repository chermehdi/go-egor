package config

import (
	"bytes"
	json2 "encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

type IoFile struct {
	Name   string
	Path   string
	Custom bool
}

func NewIoFile(fileName, filePath string, customCase bool) IoFile {
	return IoFile{
		Name:   fileName,
		Path:   filePath,
		Custom: customCase,
	}
}

func (ioFile *IoFile) GetId() int {
	tokens := strings.Split(ioFile.Name, "-")
	id, err := strconv.Atoi(tokens[1])
	if err != nil {
		return 0
	}
	return id
}

// Type mapping to the `egor-meta.json` file.
// The egor meta configuration is the source of truth for the task runner
// so an update to it (either from the outside, or by invoking egor commands) can change
// the behavior of execution of the egor cli.
type EgorMeta struct {
	TaskName  string
	TaskLang  string
	Inputs    []IoFile
	Outputs   []IoFile
	TaskFile  string
	TimeLimit float64
}

// Resolves the task file given the default language.
func GetTaskName(config Config) (string, error) {
	if config.Lang.Default == "cpp" {
		return "main.cpp", nil
	} else if config.Lang.Default == "java" {
		return "Main.java", nil
	} else if config.Lang.Default == "python" {
		return "main.py", nil
	} else {
		return "", errors.New(fmt.Sprintf("Unknown default language %s, please edit your settings", config.Lang.Default))
	}
}

// Creates a new Egor meta object from the parsed task, and the configuration values.
func NewEgorMeta(task Task, config Config) EgorMeta {
	testCount := len(task.Tests)
	inputs := make([]IoFile, testCount)
	outputs := make([]IoFile, testCount)
	for i := 0; i < testCount; i++ {
		fileName := fmt.Sprintf("test-%d", i)
		inputs[i] = NewIoFile(fileName, path.Join("inputs", fileName+".in"), false)
		outputs[i] = NewIoFile(fileName, path.Join("outputs", fileName+".ans"), false)
	}
	taskFile, err := GetTaskName(config)
	if err != nil {
		panic(err)
	}
	if task.TimeLimit == 0 {
		task.TimeLimit = 10000
	}
	return EgorMeta{
		TaskName:  task.Name,
		TaskLang:  config.Lang.Default,
		Inputs:    inputs,
		Outputs:   outputs,
		TaskFile:  taskFile,
		TimeLimit: task.TimeLimit,
	}
}

// count the number of tests cases in the metadata. The number of tests is
// the number of inputs
func (egor *EgorMeta) CountTestCases() int {
	return len(egor.Inputs)
}

func (egor *EgorMeta) toJson() (string, error) {
	var buffer bytes.Buffer
	encoder := json2.NewEncoder(&buffer)
	if err := encoder.Encode(egor); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (egor *EgorMeta) Save(w io.Writer) error {
	jsonContent, err := egor.toJson()
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, jsonContent)
	return err
}

func (egor *EgorMeta) SaveToFile(filePath string) error {
	file, _ := CreateFile(filePath)
	return egor.Save(file)
}

// TODO(chermehdi): probably this should't be a member function
func (egor *EgorMeta) Load(r io.Reader) error {
	decoder := json2.NewDecoder(r)
	err := decoder.Decode(egor)
	return err
}

// Load egor meta data from a given reader
func LoadMeta(r io.Reader) (EgorMeta, error) {
	var egor_meta EgorMeta
	decoder := json2.NewDecoder(r)
	err := decoder.Decode(&egor_meta)
	return egor_meta, err
}

// Load egor meta data form a filepath
func LoadMetaFromPath(filePath string) (EgorMeta, error) {
	file, _ := OpenFileFromPath(filePath)
	return LoadMeta(file)
}

// TODO(Eroui): this is a duplicate function from parse.go
// consider moving this somewhere common or use the other one
func CreateFile(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0777)
}

// Open file with a given file path
func OpenFileFromPath(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	return file, err
}
