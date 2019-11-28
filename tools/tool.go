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
	"time"
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

	// regexp for ID
	reg, err := regexp.Compile("^([0-9]|[a-z]){7,7} ")
	if err != nil {
		return err
	}

	// get box ID and execute COMMAND
	count := 0
	start := time.Now()
	for i := 0; i < len(all); i++ {
		line := all[i]

		// match ID
		if !reg.MatchString(line) {
			continue
		}

		// get ID
		id := reg.FindString(line)
		fmt.Println("get ID with regex: ", id)
		count++

		// display info
		fmt.Println(a.VagrantAction, "... ", id, "(", count, ")")

		// execute command
		cmd := exec.Command(a.VagrantExe, a.VagrantAction, id)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}

	}

	// print statics
	fmt.Println("TAKE TIME (minute): ", time.Now().Sub(start).Minutes())

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
