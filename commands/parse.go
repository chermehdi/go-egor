package commands

import (
	"bytes"
	"context"
	json2 "encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	template2 "text/template"
	"time"

	. "github.com/chermehdi/egor/config"
	"github.com/chermehdi/egor/templates"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

const listenAddr = ":4243"

var TaskExistsError = errors.New("Task already exists, creation skipped!")

// Serialize task into a JSON string.
func SerializeTask(meta EgorMeta) (string, error) {
	var buffer bytes.Buffer
	encoder := json2.NewEncoder(&buffer)
	if err := encoder.Encode(meta); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func CreateDirectoryStructure(task Task, config Config, rootDir string, context *cli.Context) (string, error) {
	taskDir := path.Join(rootDir, task.Name)
	if err := os.Mkdir(taskDir, 0777); err != nil {
		if os.IsExist(err) {
			return "", TaskExistsError
		}
		return "", err
	}
	if err := os.Chdir(taskDir); err != nil {
		return "", err
	}
	egorMeta := NewEgorMeta(task, config)

	file, err := CreateFile(config.ConfigFileName)
	if err != nil {
		return "", err
	}
	if err = egorMeta.Save(file); err != nil {
		return "", err
	}
	if err = os.Mkdir("inputs", 0777); err != nil {
		return "", err
	}

	if err = os.Mkdir("outputs", 0777); err != nil {
		return "", err
	}
	inputs := egorMeta.Inputs
	for i := 0; i < len(inputs); i++ {
		file, err := CreateFile(inputs[i].Path)
		if err != nil {
			return "", err
		}
		_, err = file.WriteString(task.Tests[i].Input)
		if err != nil {
			return "", err
		}
	}

	outputs := egorMeta.Outputs
	for i := 0; i < len(outputs); i++ {
		file, err := CreateFile(outputs[i].Path)
		if err != nil {
			return "", err
		}
		_, err = file.WriteString(task.Tests[i].Output)
		if err != nil {
			return "", err
		}
	}
	templateContent, err := templates.ResolveTemplateByLanguage(config)
	if err != nil {
		return "", err
	}
	template := template2.New("Solution template")
	compiledTemplate, err := template.Parse(templateContent)
	if err != nil {
		return "", err
	}
	taskFile, err := CreateFile(egorMeta.TaskFile)
	if err != nil {
		return "", err
	}
	templateContext := CreateTemplateContext(config, task)
	templateContext.MultipleTestCases = context.Bool("multiple")
	templateContext.FastIO = context.Bool("fast-io")
	if err = compiledTemplate.Execute(taskFile, templateContext); err != nil {
		return "", err
	}
	return taskDir, nil
}

// Given the json string returned by competitive companion
// it will parse it as a json Task and return it.
func extractTaskFromJson(json string) (*Task, error) {
	var task Task
	task.TimeLimit = 10000 // default timelimit to 10 seconds
	err := json2.Unmarshal([]byte(json), &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func createWebServer(quit chan<- string) *http.Server {
	router := http.NewServeMux() // here you could also go with third party packages to create a router
	// Register your routes
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		quit <- string(content)
	})
	return &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
}

func waitForShutDown(server *http.Server, done chan<- string, quit <-chan string, problemsCount int) {

	readProblems := func() <-chan string {
		results := make(chan string, problemsCount)
		defer close(results)
		has := false
		for i := 0; i < problemsCount; i++ {
			if has {
				select {
				case content := <-quit:
					results <- content
				case <-time.After(5 * time.Second):
					fmt.Println("Timed Out!")
					return results
				}
			} else {
				content := <-quit
				results <- content
				has = true
			}

		}
		return results
	}

	consumeProblems := func(results <-chan string) {
		defer close(done)
		for result := range results {
			done <- result
		}
	}

	problems := readProblems()
	consumeProblems(problems)

	color.Green("Received data from CHelper companion")
	color.Magenta("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)

	if err := server.Shutdown(ctx); err != nil {
		color.Red("Could not shutdown the server")
	}

	color.Green("Server stopped")
}

func ParseAction(context *cli.Context) error {
	problemsCount := 1

	if context.Bool("contest") && context.NArg() > 0 {
		count, err := strconv.Atoi(context.Args().Get(0))
		if err != nil {
			color.Red(fmt.Sprintf("Cannot parse problems count = '%s', a number required!", context.Args().Get(0)))
			return fmt.Errorf("Failed to parse test id = %s", context.Args().Get(0))
		}
		problemsCount = count
	}

	msgReceiveChannel := make(chan string, problemsCount)
	msgReadChannel := make(chan string, problemsCount)

	server := createWebServer(msgReadChannel)

	go waitForShutDown(server, msgReceiveChannel, msgReadChannel, problemsCount)

	color.Green("Starting the server on %s\n", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		color.Red("Could not listen on %s, %v\n", listenAddr, err)
	}

	config, err := LoadDefaultConfiguration()
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for json := range msgReceiveChannel {
		// first line contains a json string
		lines := strings.Split(json, "\n")
		task, err := extractTaskFromJson(lines[1])
		if err != nil {
			return err
		}

		taskDir, err := CreateDirectoryStructure(*task, *config, cwd, context)
		if err != nil {
			if err == TaskExistsError {
				color.Magenta("Skipping creating task %s as it already exists", task.Name)
				continue
			} else {
				color.Red("Unexpected error happened %s", err.Error())
				return err
			}
		}

		color.Green("Created task directory in : %s\n", taskDir)
	}

	return nil
}

// Command to parse tasks from the `Competitive Companion` Chrome extension.
// Running this command, will start a server on the default port that the extension
// uses, and will create a directory structure containing input files, expected output files,
// and an additional `egor-meta.json` file, and finally your task file, which is usually a `main.cpp` or `Main.java`
// file depending on the default configured language.
var ParseCommand = cli.Command{
	Name:    "parse",
	Aliases: []string{"p"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "fast-io",
			Usage:   "Indicates that this task will require the use of Fast IO",
			Aliases: []string{"io", "fio"},
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "multiple",
			Usage:   "Indicates if this task has multiple test cases",
			Aliases: []string{"m", "mul"},
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "contest",
			Usage:   "Indicates if this is a contest to parse",
			Aliases: []string{"c"},
			Value:   false,
		},
	},
	Usage:     "Parse task from navigator",
	UsageText: "egor parse",
	Action:    ParseAction,
}
