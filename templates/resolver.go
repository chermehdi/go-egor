package templates

import (
	"errors"
	"fmt"
	"github.com/chermehdi/egor/config"
	"io/ioutil"
)

// This is embeded as a literal string for easy shipping with the binary.
// We could consider using some new Go feature to embed it as a static resource.
// At the time of creation of this, this is not a priority.
const RandH = `
class Rand {
public:
	Rand(int argc, char** argv) {
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
