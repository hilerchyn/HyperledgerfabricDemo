package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
			//fmt.Println(err.Error())
			continue
		}

		go func() {
			a.WaitGroup.Add(1)
			defer a.WaitGroup.Done()

			wd := a.VagrantPath + "/" + fi.Name()

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
		}()

	}

	a.WaitGroup.Wait()

	return nil
}

func (a *App) start(fullPath string) error {

	fmt.Println("STARTING...  ", fullPath)



	cmd := exec.Command("D:\\Program Files\\Vagrant\\bin\\vagrant.exe", "up")
	cmd.Dir = fullPath
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (a *App) halt(fullPath string) error {

	fmt.Println("HALTING...  ", fullPath)
	cmd := exec.Command("D:\\Program Files\\Vagrant\\bin\\vagrant.exe", "halt")
	cmd.Dir = fullPath
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

	cwd, err = filepath.Abs(cwd)

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
