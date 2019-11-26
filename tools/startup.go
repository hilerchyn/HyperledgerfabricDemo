package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var vagrantPath = flag.String("v", "../vagrant", "-v [vagrant path]")

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

	for _, fi := range fileInfoArr {

		if !fi.IsDir() {
			continue
		}

		os.Chdir(cwd)

		fullPath := path + "/" + fi.Name() + "/Vagrantfile"
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			continue
		}

		fmt.Println("STARTING********")
		fmt.Println(fullPath)
		os.Chdir(path + "/" + fi.Name())
		cmd := exec.Command("vagrant", "up")
		cmd.Start()
		cmd.Wait()

	}

	fmt.Println("ALL DONE")
}
