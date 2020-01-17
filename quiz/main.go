package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
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

func StartQuiz(problems []Problem) {
	var correctAnswers int

	scanner := bufio.NewScanner(os.Stdin)
	for _, p := range problems {
		fmt.Printf("Question: %s\n", p.Question)
		fmt.Print("Your answer: ")
		scanner.Scan()
		answer := scanner.Text()
		if answer == p.Answer {
			correctAnswers++
		}
	}
	fmt.Println("Finished!")
	fmt.Printf("Correct answers: %d/%d\n", correctAnswers, len(problems))
}

func main() {
	fmt.Println("Welcome to quiz!")
	fmt.Println("..importing problems")

	problems, err := ReadProblemsFromFile("problems.csv")
	if err != nil {
		log.Fatal(err)
	}

	StartQuiz(problems)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Try again? (y/n)")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "y" {
			StartQuiz(problems)
		} else if input == "n" {
			break
		} else {
			fmt.Printf("Unrecognized input %s\n", input)
		}
		fmt.Println("Try again? (y/n)")
	}
}
