package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var vagrantPath = flag.String("v", ".."+string(os.PathSeparator)+"vagrant", "-v [vagrant path]")
var action = flag.String("a", "start", "-a [start|halt]")
var exe = flag.String("e", "vagrant", "-e \\path\\of\\vagrant.exe")
var app *App
var globalMap = sync.Map{}

type App struct {
	VagrantPath          string
	VagrantAction        string
	VagrantExe           string
	CurrentWorkDirectory string
	FileInfoArr          []os.FileInfo
	WaitGroup            sync.WaitGroup
}

func (a *App) Run() error {

	// get all boxes
	cmd := exec.Command(a.VagrantExe, "global-status")
	out := &bytes.Buffer{}
	cmd.Stdout = out
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	// split list
	all := strings.Split(out.String(), "\r\n")
	if len(all) < 5 {
		return nil
	}

	reg, err := regexp.Compile("^([0-9]|[a-z]){7,7} ")
	if err != nil {
		return err
	}

	// get box ID and execute COMMAND
	count := 0
	for i:= 0; i < len(all); i++ {
		line := all[i]

		if !reg.MatchString(line) {
			continue
		}

		id := reg.FindString(line)
		fmt.Println("get ID with regex: ",id)
		count++

		l := strings.Split(line, " ")


		if len(l[0]) == 7  {
			if (l[0][6] >=48 && l[0][6] <=57) || (l[0][6] >=97 && l[0][6] <=122) {
				fmt.Println(l[0])
				fmt.Println(l[4])


				cmd := exec.Command(a.VagrantExe, a.VagrantAction, l[0])
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Start(); err != nil {
					return err
				}
				if err := cmd.Wait(); err != nil {
					return err
				}
			}

		}

	}

	return nil


}


func main() {
	flag.Parse()

	// vagrant path
	path := strings.TrimSuffix(*vagrantPath, "/")
	fileInfoArr, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// current work directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cwd, err = filepath.Abs(cwd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	app = &App{VagrantPath: path, FileInfoArr: fileInfoArr, CurrentWorkDirectory: cwd, VagrantAction: *action}

	app.VagrantExe = *exe

	err = app.Run()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("ALL DONE")
}
