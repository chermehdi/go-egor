# Go egor

## Introduction

- A commandline helper utility for competitive programmers to parse competitive programming tasks on online judges (codeforces, codechef ...)
and help execute tests locally via their favorite programming language.

## Installation

- The installation requires that you have `go` installed, and you just type `go get github.com/chermehdi/go-egor`
and you will have the go-egor command available to you.

## Features

- The current supported command list is outlined here, and you can find out more details in the docs page.
    - `go-egor parse`: Starts listening for competitive companion chrome plugin to parse the task
    - `go-egor test`: Runs the tests of the current task and outputs the results
    - `go-egor copy`: Copies the current task to the clipboard.
    - `go-egor batch`: Tests the main solution for this task against another solution (probably written in another language than the main one)
    by running both of them and feeding them the tests for this task, this is usefull if you have a brute force solution and an efficient solution
    and you want to validate the optimal solution against the one you are sure that is working.

## Contribution

- Contribution of any kind is welcome, Testing, Issues and Pull requests 😄

