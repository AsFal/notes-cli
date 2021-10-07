package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

func notes(notesName string) {
	cmd := exec.Command("vim", getNotesPath(notesName))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func notesList() {
	getNotesList()
}

func notesClear(notesName string) {
	cmd := exec.Command("rm", getNotesSwapPath(notesName))
	cmd.Run()
}

func notesCopy(notesName string) {
	//create command
	catCmd := exec.Command("cat", getNotesPath(notesName))
	pbcopyCmd := exec.Command("pbcopy")

	//make a pipe
	reader, writer := io.Pipe()

	//set the output of "cat" command to pipe writer
	catCmd.Stdout = writer
	//set the input of the "pbcopy" command pipe reader

	pbcopyCmd.Stdin = reader

	//start to execute "cat" command
	catCmd.Start()

	//start to execute "pbcopy" command
	pbcopyCmd.Start()

	//waiting for "cat" command complete and close the writer
	catCmd.Wait()
	writer.Close()

	//waiting for the "pbcopy" command complete and close the reader
	pbcopyCmd.Wait()
	reader.Close()
}

func mustGetNotesDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	notesDir := path.Join(dir, ".notes")
	os.Mkdir(notesDir, os.ModePerm)
	return notesDir
}

func getNotesList() {
	root := mustGetNotesDir()
	files, _ := ioutil.ReadDir(root)
	notesNamesDir := make(map[string]bool, 0)
	for _, file := range files {
		fileName := file.Name()
		notesName := strings.Split(fileName, ".")[0]
		notesNamesDir[notesName] = true
	}
	for notesName := range notesNamesDir {
		fmt.Println(notesName)
	}
}

func getNotesPath(notesName string) string {
	notesFileName := "default"
	if len(notesName) != 0 {
		notesFileName = notesName
	}
	return path.Join(mustGetNotesDir(), fmt.Sprintf("%s.md", notesFileName))
}

func getNotesSwapPath(notesName string) string {
	notesFileName := "root"
	if len(notesName) != 0 {
		notesFileName = notesName
	}
	return path.Join(mustGetNotesDir(), fmt.Sprintf("%s.md.swp", notesFileName))
}

func getNotesName(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func main() {
	commandLineArgs := os.Args[1:]

	if len(commandLineArgs) == 0 {
		notes("")
	} else if commandLineArgs[0] == "list" {
		notesList()
	} else if commandLineArgs[0] == "copy" {
		notesCopy(getNotesName(commandLineArgs[1:]))
	} else if commandLineArgs[0] == "clear" {
		notesClear(getNotesName(commandLineArgs[1:]))
	} else {
		notes(getNotesName(commandLineArgs))
	}
}
