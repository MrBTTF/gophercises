package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
)

type Problem struct {
	Question string
	Answer   string
}

func ReadProblemsFromFile(filename string) ([]Problem, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	var problems []Problem

	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		problem := Problem{
			Question: record[0],
			Answer:   record[1],
		}
		problems = append(problems, problem)
	}

	return problems, nil
}

func StartQuiz(problems []Problem, timeoutChan chan bool, quizFinished chan int) {
	var correctAnswers int

	for _, p := range problems {
		fmt.Printf("Question: %s\n", p.Question)
		fmt.Println("Your answer: ")

		answerChan := make(chan string)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			var answer string
			for len(answer) == 0 {
				select {
				case <-timeoutChan:
					fmt.Println("\nOops! Time is up")
					quizFinished <- correctAnswers
					return
				default:
				}
				answer, _ = reader.ReadString('\n')
			}
			answerChan <- answer
		}()
		answer := <-answerChan
		answer = strings.TrimSpace(answer)
		if answer == p.Answer {
			correctAnswers++
		}
	}
	quizFinished <- correctAnswers
}

func main() {
	fd := int(os.Stdin.Fd())

	var problemsFile = flag.String("p", "problems.csv", "path to CSV file containing problems")
	var timeout = flag.Int("t", 30, "total time to finish a quiz")
	flag.Parse()

	fmt.Println("Welcome to quiz!")
	fmt.Println("..importing problems")

	problems, err := ReadProblemsFromFile(*problemsFile)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Press enter to start")
	scanner.Scan()

	start := func() {
		syscall.SetNonblock(fd, true)
		defer syscall.SetNonblock(fd, false)
		timeoutChan := make(chan bool)
		quizFinished := make(chan int)
		var correctAnswers int
		go StartQuiz(problems, timeoutChan, quizFinished)
		select {
		case <-time.After(time.Duration(*timeout) * time.Second):
			timeoutChan <- true
			correctAnswers = <-quizFinished
		case correctAnswers = <-quizFinished:
		}

		fmt.Println("Finished!")
		fmt.Printf("Correct answers: %d/%d\n", correctAnswers, len(problems))
	}
	start()
	fmt.Println("Try again? (y/n)")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "y" {
			start()
		} else if input == "n" {
			break
		} else {
			fmt.Printf("Unrecognized input %s\n", input)
		}
		fmt.Println("Try again? (y/n)")
	}
}
