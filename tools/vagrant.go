package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type App struct {
	VagrantPath          string
	VagrantAction        string
	CurrentWorkDirectory string
	FileInfoArr          []os.FileInfo
	WaitGroup            sync.WaitGroup
}

func (a *App) Run() error {

	a.WaitGroup = sync.WaitGroup{}

	for _, fi := range a.FileInfoArr {

		if !fi.IsDir() {
			continue
		}

		err := os.Chdir(a.CurrentWorkDirectory)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fullPath := a.VagrantPath + "/" + fi.Name() + "/Vagrantfile"
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return err
		}

		if err := os.Chdir(a.VagrantPath + "/" + fi.Name()); err != nil {
			fmt.Println(err.Error())
			return err
		}

		go func() {
			switch a.VagrantAction {
			case "start":
				if err := a.start(fullPath); err != nil {
					fmt.Println(err.Error())
				}

			case "halt":
				if err := a.halt(fullPath); err != nil {
					fmt.Println(err.Error())
				}

			default:
				fmt.Println("no action :", a.VagrantAction)

			}
		}()

	}

	a.WaitGroup.Wait()

	return nil
}

func (a *App) start(fullPath string) error {

	a.WaitGroup.Add(1)
	defer a.WaitGroup.Done()

	fmt.Println("STARTING...  ", fullPath)
	cmd := exec.Command("vagrant", "up")
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (a *App) halt(fullPath string) error {

	a.WaitGroup.Add(1)
	defer a.WaitGroup.Done()

	fmt.Println("HALTING...  ", fullPath)
	cmd := exec.Command("vagrant", "halt")
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

var vagrantPath = flag.String("v", "../vagrant", "-v [vagrant path]")
var action = flag.String("a", "start", "-a [start|halt]")
var app *App

func main() {
	flag.Parse()
	path := strings.TrimSuffix(*vagrantPath, "/")

	fileInfoArr, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	app = &App{VagrantPath: path, FileInfoArr: fileInfoArr, CurrentWorkDirectory: cwd, VagrantAction: *action}
	err = app.Run()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("ALL DONE")
}
