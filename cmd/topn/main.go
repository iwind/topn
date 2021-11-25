// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	ProcDir = "/proc"
)

func main() {
	topPath, err := exec.LookPath("top")
	if err != nil {
		fmt.Println("'top' command not found")
		return
	}

	var args = os.Args[1:]
	var cmdArgs = []string{}
	if len(args) > 0 {
		var keyword = args[0]
		if !strings.HasPrefix(keyword, "-") { // not option
			if regexp.MustCompile(`^\d+$`).MatchString(keyword) { // pid
				cmdArgs = append(cmdArgs, "-p", keyword)
			} else {
				// process name
				commFiles, err := filepath.Glob(ProcDir + "/*/comm")
				if err != nil {
					fmt.Println("[ERROR]travel " + ProcDir + " files failed: " + err.Error())
					return
				}

				var found = false
				for _, commFile := range commFiles {
					data, err := ioutil.ReadFile(commFile)
					if err != nil {
						continue
					}
					if strings.Contains(string(data), keyword) {
						var pieces = strings.Split(commFile, "/")
						var pid = pieces[len(pieces)-2]
						cmdArgs = append(cmdArgs, "-p", pid)

						found = true
					}
				}

				if !found {
					fmt.Println("can not find process with keyword '" + keyword + "'")
					return
				}

				cmdArgs = append(cmdArgs, args[1:]...)
			}
		} else {
			cmdArgs = append(cmdArgs, args...)
		}
	}

	var cmd = exec.Command(topPath, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
