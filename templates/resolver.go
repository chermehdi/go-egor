package templates

import (
	"errors"
	"fmt"
)

const CppTemplate = `//
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
const JavaTemplate = `import java.util.*;
import java.io.*;

/**
 * Made by egor https://github.com/chermehdi/egor.
 * {{if .Author }}
 * @author {{ .Author }}
 * {{end}}
 */
public class Main {

	void solve(InputReader in, PrintWriter out) {

	}

	public static void main(String[] args) throws Exception {
		InputReader in = new InputReader(System.in);
		PrintWriter out = new PrintWriter(System.out);
		Main solver = new Main();
		solver.solve(in, out);
		out.close();
	}
	
	static class InputReader {
		BufferedReader in;
		StringTokenizer st;
	
		public InputReader(InputStream is) {
			in = new BufferedReader(new InputStreamReader(is));
		}
	
		public String next() {
			try {
				while (st == null || !st.hasMoreTokens()) {
					st = new StringTokenizer(in.readLine());
				}
				return st.nextToken();
			} catch (Exception e) {
				throw new RuntimeException(e);
			}
		}
	
		public int nextInt() {
			return Integer.parseInt(next());
		}
	
		public long nextLong() {
			return Long.parseLong(next());
		}
	}
}
`

const PythonTemplate = `#
# Created by egor http://github.com/chermehdi/egor
# {{if .Author }}
# @author {{ .Author }}
# {{end}}
`

func ResolveTemplateByLanguage(lang string) (string, error) {
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
