package main

import (
	"fmt"
	"encoding/csv"
	"io"
	"log"
	"os"
	"time"
	"math/rand"
	"regexp"
	"strings"
)

//adjustable var for how long quiz timer runs (in seconds)
var timeSet time.Duration = 30
//flag for randomizing questions from .csv
var random = true

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

	//A little regex magic to "clean up" the answer in the .csv
	//we want to only accept letters and numbers
	//making the regex statement here
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
    if err != nil {
        log.Fatal(err)
    }

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

    	//apply the regex statement to the answer
    	//we only do this to the answer string 
    	//because it's what we compare to the user input later
    	record[1] = reg.ReplaceAllString(record[1], "")
    	//store answer as lowercase, we will do the same to the user input
    	strings.ToLower(record[1])

    	answerBank = append(answerBank, record[1])
    }	
		


	//this is our randomizer
	//if the flag is set to true, we iterate backward through the array
	//pick a random int from the range, then swap the values
	if random {
		for i := len(questionBank) - 1; i > 0; i-- {
	        j := rand.Intn(i + 1)
	        questionBank[i], questionBank[j] = questionBank[j], questionBank[i]
	        answerBank[i], answerBank[j] = answerBank[j], answerBank[i]
	    }
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
	    	//calling back to the regex to clean up user errors in input
	    	//also force string to be lowercase in the event of capitalization errors
	    	ans = reg.ReplaceAllString(ans, "")
	    	strings.ToLower(ans)
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
