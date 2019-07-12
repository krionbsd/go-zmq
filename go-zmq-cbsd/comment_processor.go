package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		_, err := fmt.Fprintf(b, " %s=\"%s\"", key, value)
		if err != nil {
			panic(err)
		}
	}
	return b.String()
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func DoProcess(comment *Comment) error {
	dt := time.Now()

	CreateDirIfNotExist("log")

	filePath := fmt.Sprintf("log/%s_%s_%d.txt", dt.Format(time.RFC3339), comment.Command, comment.JobID)
	commentFile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		_ = commentFile.Close()
	}()

	fmt.Printf("JobID %d\n", comment.JobID)

	cbsdArgs := createKeyValuePairs(comment.CommandArgs)

	cmdstr := fmt.Sprintf("env NOCOLOR=1 /usr/local/bin/cbsd %s inter=0 %s", comment.Command, cbsdArgs)
	_, err = fmt.Fprintf(commentFile, "%s\n", cmdstr)
	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh", filePath)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	filePath = fmt.Sprintf("log/%d.txt", comment.JobID)

	stdoutFile, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = stdoutFile.Close()
	}()

	cmd.Stdout = stdoutFile

	err = cmd.Run()

	if err != nil {
		fmt.Printf("%v\n", err)
	}

	return nil
}
