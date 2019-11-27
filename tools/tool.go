package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var vagrantPath = flag.String("v", ".."+string(os.PathSeparator)+"vagrant", "-v [vagrant path]")
var action = flag.String("a", "start", "-a [start|halt]")
var exe = flag.String("e", "D:\\Program Files\\Vagrant\\bin\\vagrant.exe", "-e \\path\\of\\vagrant.exe")
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


	cmd := exec.Command("vagrant", "global-status")

	out := &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Start()
	cmd.Wait()

	all := strings.Split(out.String(), "\r\n")
	fmt.Println(len(all))
	if len(all) < 5 {
		return nil
	}


	for i:= 0; i < len(all); i++ {
		line := all[i]
		l := strings.Split(line, " ")

		if len(l[0]) == 7 {
			if (l[0][6] >=48 && l[0][6] <=57) || (l[0][6] >=97 && l[0][6] <=122) {
				fmt.Println(l[0])

				cmd := exec.Command("vagrant", a.VagrantAction, l[0])


				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Start()
				cmd.Wait()

			}

		}

	}




	return nil

	a.WaitGroup = sync.WaitGroup{}

	for _, fi := range a.FileInfoArr {

		if !fi.IsDir() {
			continue
		}

		fullPath := a.CurrentWorkDirectory + string(os.PathSeparator) + a.VagrantPath + string(os.PathSeparator) + fi.Name() + string(os.PathSeparator) + "Vagrantfile"
		fullPath, err := filepath.Abs(fullPath)
		if err != nil {
			continue
		}

		if _, err := os.Stat("file:\\" + fullPath); os.IsNotExist(err) {
			fmt.Println(err.Error())
			continue
		}

		wd := a.VagrantPath + string(os.PathSeparator) + fi.Name()

		switch a.VagrantAction {
		case "start":
			if err := a.start(wd); err != nil {
				fmt.Println(err.Error(), fullPath)
			}

		case "halt":
			if err := a.halt(wd); err != nil {
				fmt.Println(err.Error(), fullPath)
			}

		default:
			fmt.Println("no action :", a.VagrantAction)

		}

	}

	return nil
}

func (a *App) start(fullPath string) error {
	return a.exe(fullPath, "up", "STARTING...")
}

func (a *App) halt(fullPath string) error {
	return a.exe(fullPath, "up", "STOPPING...")
}

func (a *App) exe(fullPath, command, info string) error {

	// check task queue
	if _, ok := globalMap.Load(fullPath); ok {
		return nil
	}

	// show info
	fmt.Println(info, fullPath)

	// add to task queue
	globalMap.Store(fullPath, true)

	// create command
	cmd := exec.Command("vagrant", command)

	// set command
	dir, err := filepath.Abs(fullPath)
	if err != nil {
		return err
	}

	fmt.Println(dir)
	os.Chdir(dir)

	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdout = os.Stdin

	// start
	if err := cmd.Start(); err != nil {
		return err
	}

	// wait
	if err := cmd.Wait(); err != nil {
		return err
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
