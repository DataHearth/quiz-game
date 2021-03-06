package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type game struct {
	Question string `json:"question" xml:"question"`
	Result   string `json:"result" xml:"result"`
}

// Debug - activate development logs
var Debug bool

func main() {
	var data []game

	// Define flags
	filePath := flag.String("path", "problems.csv", "File path with file ending with [json, csv, xml]")
	questionSeparator := flag.String("separator", ",", "Question separator")
	debug := flag.Bool("debug", false, "Need output ?")
	timerDuration := flag.Int("duration", 1000, "Time duration in seconds")

	// Setup flags for output
	flag.Parse()

	// Set debug variable
	Debug = *debug

	// Debug values
	if Debug {
		fmt.Printf("file path: %v\n", *filePath)
		fmt.Printf("question separator: %v\n", *questionSeparator)
	}

	fileTypeEnding := strings.Split(*filePath, ".")
	fileType := fileTypeEnding[len(fileTypeEnding)-1]

	if Debug {
		log.Printf("File type array: %v\n", fileTypeEnding)
		log.Printf("File type: %v\n", fileType)
	}

	// Check type for different format
	switch fileType {
	case "csv":
		_, ioReader := readData(*filePath, true)
		records, err := csv.NewReader(ioReader).ReadAll()
		if err != nil {
			log.Fatalln(err.Error())
		}

		for _, value := range records {
			data = append(data, game{Question: value[0], Result: value[1]})
		}
		break
	case "json":
		jsonData, _ := readData(*filePath, false)
		json.Unmarshal(jsonData, &data)
		break
	case "xml":
		// Doesn't work
		xmlData, _ := readData(*filePath, false)
		xml.Unmarshal(xmlData, &data)
		break
	default:
		log.Fatalln(errors.New("Your file needs to end with .[csv, json, xml]"))
		break
	}

	if Debug {
		log.Println(data)
	}

	fmt.Println("Game starting!")
	score := askQuestions(data, *timerDuration)
	fmt.Printf("Your score is %d out of %d\n", score, len(data))
}

func readData(filePath string, ioReader bool) ([]byte, io.Reader) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// If io.Reader is needed, return file
	if ioReader {
		return nil, file
	}

	// Read data from file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Couldn't read data\nError: %v", err)
	}

	return data, nil
}

func askQuestions(questions []game, duration int) int {
	var score int
	answerChan := make(chan string)

	timer := time.NewTimer(time.Duration(duration) * time.Second)

	for i := 0; i < len(questions); i++ {
		fmt.Printf("Question %d: %v\n", i+1, questions[i].Question)
		if Debug {
			log.Printf("question: %v, answer: %v", questions[i].Question, questions[i].Result)
		}

		// Make the user input timer depends (if the timer ends, the problem is not stuck in scanf)
		go func() {
			var userInput string

			fmt.Scanf("%v\n", &userInput)
			answerChan <- userInput
		}()

		// Check for user response OR timer's end (no other option allowed)
		select {
		case <-timer.C:
			return score
		case answer := <-answerChan:
			if answer == questions[i].Result {
				score++
			}
		}
	}

	return score
}
