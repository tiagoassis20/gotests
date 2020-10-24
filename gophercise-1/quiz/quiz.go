package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func getQuestions(r io.Reader) (total chan int, score chan int) {
	total = make(chan int)
	score = make(chan int)
	go func(r io.Reader) {
		reader := csv.NewReader(bufio.NewReader(r))
		questions, err := reader.ReadAll()
		if err != nil {
			log.Fatalln("Couldn't read the csv file", err)
		}
		total <- len(questions)
		for i, question := range questions {
			fmt.Printf("%d) %s = ", i, question[0])
			answer := ""
			fmt.Scanln(&answer)
			if answer == question[1] {
				score <- 1
			} else {
				score <- 0
			}
		}
	}(r)
	return
}

func main() {
	fmt.Println("Hello World.")
	var file = flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	var limit = flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	csvfile, err := os.Open(*file)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	totalCh, scoreCh := getQuestions(csvfile)

	total := <-totalCh

	score := 0
	t := time.NewTimer(time.Duration(*limit) * time.Second)
l:
	for i := 0; i < total; i++ {
		select {
		case s := <-scoreCh:
			score += s
		case <-t.C:
			fmt.Println("\n Time is over!")
			break l
		}
	}

	fmt.Printf("you scored %d out of %d.\n", score, total)
}
