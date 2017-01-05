package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type data struct {
	url      string
	revision string
	author   string
	status   string
}

func main() {

	argsWithoutProg := os.Args[1:]

	towrite := data{}

	for i, v := range argsWithoutProg {
		switch i {
		case 0:
			if v == "--help" {
				writehelp()
				return
			}
			towrite.author = `"` + v + `"`
		case 1:
			towrite.url = `"` + v + `"`
		case 2:
			towrite.revision = `"` + v + `"`
		case 3:
			towrite.status = `"` + v + `"`
		}
	}
	if len(towrite.author) == 0 {
		towrite.author = author()
	}
	if len(towrite.url) == 0 {
		towrite.url = url()
	}
	if len(towrite.revision) == 0 {
		towrite.revision = revision()
	}

	writefile(towrite)

}

func url() string {

	url, err := exec.Command(`git`, `remote`, `-v`).Output() //`| sed -e '1!d;s/origin\t//g;s/ (fetch)//g'`
	if err != nil {
		fmt.Print(err)
		panic("Could not get scm url. Execute this only inside git repos.")
	}
	res := strings.Split(string(url), "\n")[0]
	res = strings.Replace(res, "origin\t", "", -1)
	res = strings.Replace(res, " (fetch)", "", -1)
	res = strings.TrimSpace(res)
	return `"` + res + `"`

}
func author() string {

	for _, v := range os.Environ() {
		res := strings.Split(v, "=")
		if res[0] == "USER" {
			return `"` + res[1] + `"`
		}
	}
	return ""
}

func revision() string {
	revision, err := exec.Command(`git`, `log`).Output() // | sed '1!d; s/commit //g'`
	res := strings.Split(string(revision), "\n")[0]
	res = strings.Replace(res, "commit", "", -1)
	res = strings.TrimSpace(res)
	if err != nil {
		fmt.Print(err)
		panic("Could not get revision. Execute this only inside git repos.")
	}

	return `"` + res + `"`
}

func writefile(d data) {

	rawdata := []byte{}
	rawdata = append(rawdata, []byte("{\n\t")...)
	rawdata = append(rawdata, []byte(`"url": `)...)
	rawdata = append(rawdata, d.url...)
	rawdata = append(rawdata, []byte("\n\t")...)
	rawdata = append(rawdata, []byte(`"revision": `)...)
	rawdata = append(rawdata, d.revision...)
	rawdata = append(rawdata, []byte("\n\t")...)
	rawdata = append(rawdata, []byte(`"author": `)...)
	rawdata = append(rawdata, d.author...)
	rawdata = append(rawdata, []byte("\n\t")...)
	rawdata = append(rawdata, []byte(`"status": `)...)
	rawdata = append(rawdata, d.status...)
	rawdata = append(rawdata, []byte("\n}")...)

	fmt.Println(string(rawdata))

	err := ioutil.WriteFile("scm-source.json", rawdata, 0644)
	if err != nil {
		fmt.Print(err)
		panic("Could not write file.")
	}
}

func writehelp() {
	fmt.Println("usage: gscm [author] [url] [revision] [status]")
	fmt.Println()
	fmt.Println("all args are optional if not provided will use the following:")
	fmt.Println()
	fmt.Println("\tauthor = $USER  (", author(), ")")
	fmt.Println("\turl = $(git remote -v | sed -e '1!d;s/origin\t//g;s/ (fetch)//g')  (", url(), ")")
	fmt.Println("\trevision = $(git log | sed '1!d; s/commit //g')  (", revision(), ")")
	fmt.Println("\tstatus wont be filled")
}
