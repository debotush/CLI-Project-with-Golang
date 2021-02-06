package pkg

import (
	"bufio"
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"strings"
	"sync"
	"time"
)

// "Title,Message 1,Message 2,Stream Delay,Run Times\nCLI Invoker Name,First Message,Second Msg,2,10"
type CliStreamerRecord struct {
	Title       string `csv:"Title"`
	Message1    string `csv:"Message 1"`
	Message2    string `csv:"Message 2"`
	StreamDelay int    `csv:"Stream Delay"`
	RunTimes    int    `csv:"Run Times"`
}

type CliRunnerRecord struct {
	// How many streamer will run.
	Run         string `csv:"Run"`
	Title       string `csv:"Title"`
	Message1    string `csv:"Message 1"`
	Message2    string `csv:"Message 2"`
	StreamDelay int    `csv:"Stream Delay"`
	RunTimes    int    `csv:"Run Times"`
}

func (cliRunnerRecord CliRunnerRecord) CliStreamerRecord() CliStreamerRecord {
	return CliStreamerRecord{
		Title:       cliRunnerRecord.Title,
		Message1:    cliRunnerRecord.Message1,
		Message2:    cliRunnerRecord.Message2,
		StreamDelay: cliRunnerRecord.StreamDelay,
		RunTimes:    cliRunnerRecord.RunTimes,
	}
}

func Csv(cliRunners *[]CliRunnerRecord) string {
	out, err := gocsv.MarshalString(cliRunners)
	if err != nil {
		panic(err)
	}
	return out
}

func PrintMassage(cliStreamerRecord CliStreamerRecord, wg  *sync.WaitGroup, m *sync.Mutex, f *os.File) {
	defer wg.Done()
	for idx := 1 ;idx <= cliStreamerRecord.RunTimes; idx++ {
		massageOne := cliStreamerRecord.Title + "->" + cliStreamerRecord.Message1
		massageTwo := cliStreamerRecord.Title + "->" + cliStreamerRecord.Message2
		fmt.Println(massageOne)

		m.Lock()
		w := bufio.NewWriter(f)
		_, err := fmt.Fprintf(w, "%v\n", massageOne)
		if err != nil {
			panic(err)
		}
		w.Flush()
		m.Unlock()

		time.Sleep(time.Duration(cliStreamerRecord.StreamDelay)*time.Millisecond)

		fmt.Println(massageTwo)
		m.Lock()
		w2 := bufio.NewWriter(f)
		_, err = fmt.Fprintf(w2, "%v\n", massageTwo)
		if err != nil {
			panic(err)
		}
		w2.Flush()
		m.Unlock()

	}
}



func CommandLineInterface() {
	args := strings.Join(os.Args[1:], " ")
	multiLineArgs := strings.Replace(args, `\n`, "\n", len(args))
	var mutex sync.Mutex
	f, err := os.Create("CLIStreamAndCLIRunner/output/log.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var w sync.WaitGroup

	var cliRunners []CliRunnerRecord
	gocsv.UnmarshalString(multiLineArgs, &cliRunners)

	for _, runner := range cliRunners {
		w.Add(1)
		go PrintMassage(runner.CliStreamerRecord(),&w ,&mutex, f)

	}
	w.Wait()
}