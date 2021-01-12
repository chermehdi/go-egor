package templates

import (
	"errors"
	"fmt"
	"github.com/chermehdi/egor/config"
	"io/ioutil"
)

const BruteH = `
#include <iostream>
#include <vector>
#include <set>

using namespace std;

int main() {
}
`

// This is embeded as a literal string for easy shipping with the binary.
// We could consider using some new Go feature to embed it as a static resource.
// At the time of creation of this, this is not a priority.
const RandH = `
#pragma once

#include <ctime>
#include <string>

enum CASE {
	LOWERCASE,
	UPPERCASE,
	MIXEDCASE
};

class Rand {
private:
	int seed;

	void init() {
		srand(seed);
	}

	void readArgs(int argc, char * argv[]) {
		if (argc > 0) {
			seed = std::stoi(argv[1]);
		} else {
			seed = 0;
		}
	}

public:
	Rand(int _seed) : seed(_seed) {
		init();
	}

	Rand(int argc, char* argv[]) {
		readArgs(argc, argv);
		init();
	}

	// Numbers

	int Int(int from, int to) {
		return rand() % (to - from) + from;
	}

	long Long(long from, long to) {
		return rand() % (to - from) + from;
	}

	long long LongLong(long long from, long long to) {
		return rand() % (to - from) + from;
	}

	double Double(double from, double to) {
		double tmp = (double) rand() / RAND_MAX;
		return from + tmp * (to - from);
	}

	// Chars

	char Lower(char from = 'a', char to = 'z') {
		return Char(from, to);
	}

	char Upper(char from = 'A', char to = 'Z') {
		return Char(from, to);
	}

	char Alpha(char from = 'a', char to = 'z') {
		int k = rand() % 52;
		if (k < 26) return Lower('a', 'a' + k);
		return Upper('a', 'a' + k % 26);
	}

	char Digit(char from = '0', char to = '9') {
		return Char(from, to);
	}

	char Char(char from, char to) {
		return static_cast <char> (Int(from, to));
	}

	// Strings

	std::string String(size_t size, CASE mode = MIXEDCASE) {
		std::string res(size, ' ');

		switch (mode) {
			case LOWERCASE:
				for (char& c : res) {
					c = Lower();
				}
				break;
			case UPPERCASE:
				for (char& c : res) {
					c = Upper();
				}
				break;
			case MIXEDCASE:
				for (char &c : res) {					
					c = Alpha();
				}
				break;
			default:
				break;
		}

		return res;
	}

};
`

const GeneratorTemplate = `
//
// Created by egor http://github.com/chermehdi/egor
// {{if .Author }}
// @author {{ .Author }}
{{end}}
#include <iostream>
#include <vector>
#include "rand.h"

using namespace std;

int main(int argc, char** argv) {
	// Do not remove this line
	Rand rand(argc, argv);	
}
`
const CppTemplate = `
//
// Created by egor http://github.com/chermehdi/egor
// {{if .Author }}
// @author {{ .Author }}
{{end}}
#include <iostream>
#include <vector>
#include <set>
#include <map>
#include <algorithm>
#include <cmath>

using namespace std;

int main() {

}
`
const JavaTemplate = `
import java.util.*;
import java.io.*;

/**
 * Made by egor https://github.com/chermehdi/egor.
 * {{if .Author }}
 * @author {{ .Author }}
 * {{end}}
 */
public class Main {

    void solve(Scanner in, PrintWriter out) {

    }

    public static void main(String[] args) {
        try(Scanner in = new Scanner(System.in);
            PrintWriter out = new PrintWriter(System.out)) {
            new Main().solve(in, out);
        }
    }
}
`

const PythonTemplate = `
#
# Created by egor http://github.com/chermehdi/egor
# {{if .Author }}
# @author {{ .Author }}
# {{end}}
`

func ResolveTemplateByLanguage(config config.Config) (string, error) {
	templates := config.CustomTemplate
	path, has := templates[config.Lang.Default]
	if !has {
		return resolveWithDefaultTemplate(config.Lang.Default)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func resolveWithDefaultTemplate(lang string) (string, error) {
	switch lang {
	case "cpp":
		return CppTemplate, nil
	case "java":
		return JavaTemplate, nil
	case "python":
		return PythonTemplate, nil
	default:
		return "", errors.New(fmt.Sprintf("Unknown language %s provided", lang))
	}
}
