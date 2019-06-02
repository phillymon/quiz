package main

import (
	"fmt"
	"encoding/csv"
	"io"
	"log"
	"os"
	"time"
)

//adjustable var for how long quiz timer runs (in seconds)
var timeSet time.Duration = 30

func main () {

	//select and open file to be used as question bank
	csvfile, err := os.Open("problems.csv")
	if err != nil {
		log.Fatalln("fatal error: could not open file\n", err)
	}
	//create reader to read each line, we will store the first string in the array as the question
	//we will store the second string as the answer
	reader := csv.NewReader(csvfile)
	var questionBank []string
	var answerBank []string

	for {

		record, err := reader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}

		questionBank = append(questionBank, record[0])
		answerBank = append(answerBank, record[1])	
	}

	//prompt user for # of questions, then to begin quiz
	fmt.Printf("How many questions would you like to answer? (max: %v)\n", len(questionBank))
	var numQuestions int
    fmt.Scanln(&numQuestions)
	fmt.Printf("Type the answer to each problem and press the 'Enter' key\n\nPress 'Enter' to begin\n\n")
	var input string
	fmt.Scanln(&input)

	//channel to determine if user has completed the quiz before time up
	//value gets set to "true" after all questions completed
	finish := make(chan bool)
	//counter for right answers
	correct := 0

	//start goroutine so we can simultaneously run the quiz and the counter
	//finish states are handled in the "select" statement below
	go func() {
		for i:=0; i<numQuestions; i++ {

				//iterate through question array and compare user answers to answer array
				fmt.Printf("%v \n", questionBank[i])
				
				var ans string
		    	fmt.Scanln(&ans)
				if ans == answerBank[i] {
					fmt.Printf("Correct!\n")
					correct++
				} else {
					fmt.Printf("Wrong :(\n")
				}
			}
		//after loop through all questions is complete, send "true" to the channel
		finish <- true
	} ()

	//set timer with our global time variable
	timer1 := time.NewTimer(timeSet * time.Second)

	//handler for which finishes first, user completion of quiz or timer
	select {
        case <-finish:
	        
	    case <-timer1.C:
        	fmt.Printf("Time up!\n")  	    
    }

    //results!
	fmt.Printf("Quiz finished! %v/%v correct\n", correct, numQuestions)
	

}
