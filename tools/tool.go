package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var action = flag.String("a", "up", "-a [up|halt]")
var exe = flag.String("e", "vagrant", "-e \\path\\of\\vagrant.exe")
var app *App

type App struct {
	VagrantPath          string
	VagrantAction        string
	VagrantExe           string
	CurrentWorkDirectory string
	FileInfoArr          []os.FileInfo
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
	start := time.Now()
	for i := 0; i < len(all); i++ {
		line := all[i]

		if !reg.MatchString(line) {
			continue
		}

		vmId := strings.TrimRight(reg.FindString(line), " ")
		count++

		fmt.Println(a.VagrantAction, "... ", fmt.Sprintf("%s [%d]", vmId, count))
		cmd := exec.Command(a.VagrantExe, a.VagrantAction, vmId)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}

	}

	fmt.Println("TAKE TIME (minutes): ", time.Now().Sub(start).Minutes())

	return nil

}

func main() {
	flag.Parse()

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

	app = &App{CurrentWorkDirectory: cwd, VagrantAction: *action, VagrantExe: *exe}

	err = app.Run()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("ALL DONE")
}
