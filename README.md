# Go egor

## Introduction

- A commandline helper utility for competitive programmers to parse competitive programming tasks on online judges (codeforces, codechef ...)
and help execute tests locally via their favorite programming language.

## Installation

- There are two ways for installing egor.

### Download compiled binary
- You can download the binary that corresponds to your operating system and add it to your `PATH` variable, you can find the binaries in the [releases](https://github.com/chermehdi/go-egor/releases) page.
- For people running Mac OSX Catalina, this won't work and you will be forced to go with the Build from source solution.

### Build from source
- You can clone the repository, and have go installed in your machine
- Navigate to the directory of the cloned project and run `go build` to build the project, and `go install` to install it to your local machine, the binary can then be found in `$GOPATH/bin/egor`

## Features

- The current supported command list is outlined here, and you can find out more details in the docs page.
    - `egor parse`: Starts listening for competitive companion chrome plugin to parse the task
    - `egor test`: Runs the tests of the current task and outputs the results.
    - `egor config`: Read/Change global configuration parameters.
    - `egor testcase`: Add a custom test case to this egor task.
    - `egor showcases`: list meta data about the tests in the current task 
    - `egor printcase`: Print input and or output of a given test case.
    - `egor copy`: Copies the current task to the clipboard.
    - `egor help`: Display help for command.
    - `egor batch (in developement)`: Tests the main solution for this task against another solution (probably written in another language than the main one)
    by running both of them and feeding them the tests for this task, this is useful if you have a brute force solution and an efficient solution
    and you want to validate the optimal solution against the one you are sure that is working.

## Contribution

- Contribution to the project can be done in multiple ways, you can report bugs and issues or you can discuss new features that you like being added in future versions by creating a new issue
and tagging one of the maintainers, and if you have the time and have ideas about how you can integrate your feature, it's always a good to back your feature request with a PR (but it's fine if you can't 😅)

- PR's opened should mention (at least) one of the maintainers as a reviewer, and shouldn't break the current test suit.
