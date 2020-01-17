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

	scanner := bufio.NewScanner(os.Stdin)
	for _, p := range problems {
		fmt.Printf("Question: %s\n", p.Question)
		fmt.Println("Your answer: ")
		scanner.Scan()
		answer := scanner.Text()
		answer = strings.TrimSpace(answer)
		if answer == p.Answer {
			correctAnswers++
		}
		select {
		case <-timeoutChan:
			fmt.Println("Oops! Time is up")
			quizFinished <- correctAnswers
			return
		default:
		}
	}
	quizFinished <- correctAnswers
}

func main() {
	var problemsFile = flag.String("p", "problems.csv", "path to CSV file containing problems")
	var timeout = flag.Int("t", 5, "total time to finish a quiz")
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
		timeoutChan := make(chan bool, 1)
		quizFinished := make(chan int, 1)
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
