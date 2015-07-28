package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s \n", name, elapsed)
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))

}

func PrintError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))

	}

}

func PrintOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))

	}

}

type spectroBuilder struct {
	filename string
}

func removeFile(filename string) {
	rm := exec.Command("rm", "upload/"+filename)
	output, err := rm.CombinedOutput()
	PrintError(err)
	PrintOutput(output)
}

func runSox(filename string) {
	cmd := exec.Command("sox", "upload/"+filename, "-n", "spectrogram", "-o", "processed/"+filename+".png")
	output, err := cmd.CombinedOutput()
	PrintError(err)
	PrintOutput(output)

}

// ProcessFile implements Worker interface
func (s *spectroBuilder) ProcessFile() {
	fmt.Printf("Processing %s \n", s.filename)
	runSox(s.filename)
	removeFile(s.filename)
	fmt.Printf("Completed %s \n", s.filename)
}

func BuildSpectrograms() {
	defer TimeTrack(time.Now(), "building spectrograms")

	files, _ := ioutil.ReadDir("upload/")
	p := New(len(files))
	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, f := range files {
		sB := spectroBuilder{filename: f.Name()}
		go func() {
			// Submit work to the pool
			p.Run(&sB)
			wg.Done()
		}()
	}
	wg.Wait()
	p.Shutdown()
}

type JsonResponse struct {
	Spectrograms map[string]string `json:"results"`
}
