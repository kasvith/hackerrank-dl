/*
 * MIT License
 *
 * Copyright (c) 2020 Kasun Vithanage
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

var ext = map[string]string{
	"ada":          "ada",
	"bash":         "sh",
	"c":            "c",
	"clojure":      "clj",
	"coffeescript": "coffee",
	"cpp":          "cpp",
	"cpp14":        "cpp",
	"csharp":       "cs",
	"d":            "d",
	"elixir":       "ex",
	"erlang":       "erl",
	"fortran":      "for",
	"fsharp":       "fs",
	"go":           "go",
	"groovy":       "groovy",
	"haskell":      "hs",
	"java":         "java",
	"java8":        "java",
	"javascript":   "js",
	"julia":        "jl",
	"kotlin":       "kt",
	"lolcode":      "lol",
	"lua":          "lua",
	"pascal":       "pas",
	"perl":         "perl",
	"php":          "php",
	"r":            "r",
	"objectivec":   "m",
	"pypy":         "py",
	"pypy3":        "py",
	"sbcl":         "lsp",
	"racket":       "rkt",
	"smalltalk":    "st",
	"tcl":          "tcl",
	"visualbasic":  "vb",
	"ruby":         "ruby",
	"rust":         "rs",
	"scala":        "scala",
	"swift":        "swift",
	"ocaml":        "c",
	"python":       "py",
	"python3":      "py",
}

func GetExtension(lang string) string {
	if v, ok := ext[lang]; ok {
		return v
	}
	return "txt"
}
