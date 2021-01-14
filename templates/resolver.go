package templates

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/chermehdi/egor/config"
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
const RandH = `#pragma once

#include <algorithm>
#include <cassert>
#include <climits>
#include <ctime>
#include <random>
#include <string>

const std::string LOWERCASE_ALPHABET = "abcdefghijklmnopqrstuvwyxz";
const std::string UPPERCASE_ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWYXZ";
const std::string DIGITS_SET = "0123456789";

// Used to determine the set of allowed characters
// for the String method of the Rand class
const int LOWER = 1;
const int UPPER = 1 << 1;
const int DIGITS = 1 << 2;

// Returns a random value of the following types
// int, long long, double, char, string
// given main arguments
// Example:
// Rand rand(argc, argv);
// int randomInt = rand.Int();
// string randomString = rand.String(10, DIGITS);
//
// @author NouemanKHAL
//
class Rand {
  unsigned int seed;

  unsigned int ReadSeedFromArgs(int argc, char* argv[]) {
    if (argc > 0) {
      return std::stoi(argv[1]);
    }
    return 0;
  }

 public:
  Rand(unsigned int _seed = 0) : seed(_seed) { srand(seed); }

  Rand(int argc, char* argv[]) {
    unsigned int seed = ReadSeedFromArgs(argc, argv);
    Rand{seed};
  }

  // Returns an random int value in the range [from, to] inclusive
  int Int(int from, int to) {
    if (from == to) return from;
    assert(from < to);
    return rand() % (to - from) + from;
  }

  // Returns an random long long value in the range [from, to] inclusive
  long long Long(long long from, long long to) {
    if (from == to) return from;
    assert(from < to);
    return rand() % (to - from) + from;
  }

  // Returns an random double value in the range [from, to] inclusive
  double Double(double from, double to) {
    assert(from <= to);
    double tmp = (double)rand() / RAND_MAX;
    return from + tmp * (to - from);
  }

  // Returns an random char value in the range [from, to] inclusive
  // Parameters are optional
  char Char(char from = CHAR_MIN, char to = CHAR_MAX) {
    assert(from <= to);
    return static_cast<char>(Int(from, to));
  }

  // Returns an random char value in the range [from, to] inclusive
  // Parameters are optional, by default returns a random lowercase letter
  char Lower(char from = 'a', char to = 'z') {
    assert('a' <= from && from <= to && to <= 'z');
    return Char(from, to);
  }

  // Returns an random char value in the range [from, to] inclusive
  // Parameters are optional, by default returns a random uppercase letter
  char Upper(char from = 'A', char to = 'Z') {
    assert('A' <= from && from <= to && to <= 'Z');
    return Char(from, to);
  }

  // Returns an random letter either in lowercase or uppercase
  char Alpha() {
    int k = rand() % 52;
    if (k < 26) return Lower('a', 'a' + k);
    return Upper('a', 'a' + k % 26);
  }

  // Returns an random digit character
  // Parameters are optional, by default returns a random digit character in the
  // range ['0', '9'] inclusive
  char Digit(char from = '0', char to = '9') {
    assert(from <= to);
    return Char(from, to);
  }

  // Returns a random alphanumerical character
  char AlphaNum() {
    if (rand() & 1) return Alpha();
    return Digit();
  }

  // Returns a random boolean value.
  bool Bool() { return bool(rand() & 1); }

  // Returns an std::string of length size consisting only of characters allowed
  // in the given mask using the constants LOWER, UPPER, DIGITS Example: Rand
  // rand(argc, argv); std::string str = rand.String(10, LOWER | DIGITS); str is
  // an std::string of size 10 consisting only of lowercase letters and digits
  std::string String(size_t size, const int mask = LOWER | UPPER) {
    std::string charset;

    // Building the set of allowed charset from the mask
    if (mask & LOWER) {
      charset += LOWERCASE_ALPHABET;
    }
    if (mask & UPPER) {
      charset += UPPERCASE_ALPHABET;
    }
    if (mask & DIGITS) {
      charset += DIGITS_SET;
    }

    std::random_shuffle(charset.begin(), charset.end());

    std::string res(size, ' ');

    int len = charset.size();

    for (char& c : res) {
      size_t randomIndex = Int(0, len - 1);
      c = charset[randomIndex];
    }

    return res;
  }
}
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

import java.io.*;
import java.util.*;

/**
 * Made by egor https://github.com/chermehdi/egor.
 * {{if .Author }}
 * @author {{ .Author }}
 * {{end}}
 */
public class Main {
  void solve(Scanner in, PrintWriter out) {}
  public static void main(String[] args) {
    try (Scanner in = new Scanner(System.in); PrintWriter out = new PrintWriter(System.out)) {
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
