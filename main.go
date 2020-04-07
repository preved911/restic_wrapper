package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// override restic executable path
	resticExecPath := os.Getenv("RESTIC_BINARY_PATH")
	if resticExecPath == "" {
		resticExecPath = "/usr/local/bin/restic"
	}

	godotenv.Load("/etc/default/restic")

	logFile, err := os.OpenFile("/var/log/restic/restic.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("log file open error: %\ns", err)
	}
	defer logFile.Close()

	if os.Args[1] == "backup" {
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	os.Setenv(
		"RESTIC_REPOSITORY",
		fmt.Sprintf("%s/%s",
			os.Getenv("RESTIC_REPOSITORY_BUCKET"),
			os.Getenv("RESTIC_REPOSITORY_PREFIX")))

	var stdout []byte
	var stderr error

	// 5 attempts before failed state return
	for i := 0; i < 5; i++ {
		cmd := exec.Command(resticExecPath, os.Args[1:]...)
		stdout, stderr = cmd.CombinedOutput()
		if stderr != nil {
			log.Printf("execution failed: %s\n", stderr)
			time.Sleep(time.Duration(5*rand.Intn(5)) * time.Second)
		} else {
			break
		}
	}

	log.Println(string(stdout))

	if stderr != nil {
		log.Fatalf("%s", stderr)
	}
}
