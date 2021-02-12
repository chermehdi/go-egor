# Go egor

## Introduction

- A commandline helper utility for competitive programmers to parse competitive programming tasks on online judges (codeforces, codechef ...)
and help execute tests locally via their favorite programming language.

## Installation

- There are two ways for installing egor.

### Download compiled binary
- You can download the binary corresponding to your operating system and add it to your `PATH` variable, you can find the binaries in the [releases](https://github.com/chermehdi/go-egor/releases) page.
- For people running Mac OSX Catalina, this won't work and you will be forced to go with the Build from source solution.

### Build from source
- You can clone the repository, and have go installed in your machine
- Navigate to the directory of the cloned project and run `go build` to build the project, and `go install` to install it to your local machine, the binary can then be found in `$GOPATH/bin/egor`

## Features

- The current supported command list is outlined here, and you can find out more details in the docs page.
    - `egor parse`: Starts listening for competitive companion chrome plugin to parse a task or a contest.
    - `egor test`: Runs the tests of the current task and outputs the results.
    - `egor config`: Read/Change global configuration parameters.
    - `egor testcase`: Add a custom test case to this egor task.
    - `egor showcases`: list meta data about the tests in the current task 
    - `egor printcase`: Print input and or output of a given test case.
    - `egor copy`: Copies the current task to the clipboard.
    - `egor help`: Display help for command.
    - `egor batch`: Tests the main solution for this task against another solution (probably written in another language than the main one)
    by running both of them and feeding them the tests for this task, this is useful if you have a brute force solution and an efficient solution
    and you want to validate the optimal solution against the one you are sure that is working.
    
## C++ Library code
- If you are using C++ to do competitive programming, and you have your own library that you use,
you can just tell egor the location of your library, and it will use it when finding tasks that rely on it.

### Steps to make it work
1- run `egor config set cpp.lib.location /path/to/root/of/library` to tell egor the path to use to resolve your includes.

2- Create your tasks normally and use your library:

```cpp
#include <iostream>
#include <vector>
// notice the include is using " and not <
#include "kratos/graphs/tree.h"

using namespace std;

int main() {
    int n; cin >> n;
    kratos::Tree tree;
    for(int i = 0; i < n - 1; ++i) {
        int u, v; cin >> u >> v;
        tree.add_edge(u, v):
    } 
    cerr << tree.debug() << endl;
}
```

3- In this example, the configuration i have is:
```
$ egor config set cpp.lib.location
/home/directory/include

$ ls /home/directory/include
kratos
``` 

## Custom templates

- Egor also support the use of custom templates per each language, powered by the golang template engine. to take full advantage of it you can
take a look at all the details in the official [documentation](https://golang.org/pkg/text/template/).

- A typical configuration for a template, is the following: 
```
//
// {{ if .Problem}} {{ .Problem }} {{ end }} {{ if .Url }} {{ .Url }} {{ end }}
{{- if .Author }}
// @author {{ .Author }}
// created {{ .Time }}
{{- end }}
// 
#include <iostream>
#include <vector>
#include <set>
#include <algorithm>
#include <map>

using namespace std;
{{ if .MultipleTestCases }}
void solve() {
}
{{ end }}

int main() {
  {{- if .FastIO }}
  ios_base::sync_with_stdio(false);
  cin.tie(0);
  {{- end}}

  {{- if .MultipleTestCases }}
  int t; cin >> t;
  while(t--) {
    solve();
  }
  {{- end}}
}
```
Each template is provided by a model containing some basic information about the task that is going to be generated
to help create dynamic templates, the model reference is
```go
type TemplateContext struct {
	Author            string
	Time              string
	MultipleTestCases bool
	Interactive       bool
	FastIO            bool
	Problem           string
	Url               string
}
```
You can access any of the model fields in the template, and make changes accordingly dependending on your preferences.
- Running the command `egor config set custom.templates.{lang} /path/to/template/file` will register the given template and will use it for the given language for future tasks.
You can also see what is the current custom template for your language by running `egor config get config.templates`

## Batch testing

- This one is usefull if you are stuck on a problem, your solution gets `WA` and
you can't figure out a testcase that will make it break.
- If you know how to write a brute force solution to the problem, egor can help
  with running the solution that you know gives the correct answer against
  your (optimal) solution, and figure out when will they differ.

### Steps to make a batch test
- Create a batch directory structure: 

```
egor batch create 
```

Running this will create 3 additional files: 

- `gen.cpp`: which is the random input generator, this include a small random
  generation library `"rand.h"`, feel free to use it to generate testcases for
  the given problem.
- `main_brute.cpp`: which will act as the brute force solution.
- `rand.h`: the random generation library, contains a bunch of helper methods
  to generate random input, read the source and docs on public methods for
  additional info.

Once all the files are filled according to their purpose, you can start running
your batch test using

- Run the batch test:

```
egor batch run --tests=100
```

the parameter `--tests` tells egor how many times it should run your code,
a `100` is a good start (and it's the default).

If you can't find a failing testcase after this, you can increase the
`--tests` parameter, or change the generation strategy, (ex: try bigger
/ smaller numbers), egor will give your program a random seed every time it runs
the generator, so that it generates a new testcase with every run.

## Contribution

- Contribution to the project can be done in multiple ways, you can report bugs and issues or you can discuss new features that you like being added in future versions by creating a new issue
and tagging one of the maintainers, and if you have the time and have ideas about how you can integrate your feature, it's always a good to back your feature request with a PR (but it's fine if you can't 😅)

- PR's opened should mention (at least) one of the maintainers as a reviewer, and shouldn't break the current test suit.
