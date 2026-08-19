package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"fmt"
	"os"
	"io"
)


type StackJson struct {
	Items						[]StackData			`json:"items"`
}


type StackData struct {
	IsAnswered			bool						`json:"is_answered"`
	QuestionID			int							`json:"question_id"`
	Score						int							`json:"score"`
	Title						string					`json:"title"`
	Link						string					`json:"link"`
	Tags						[]string				`json:"tags"`
}


type AnswersJson struct {
	Items						[]StackAnswers	`json:"items"`
}

type StackAnswers struct {
	Body						string					`json:"body"`
}

func main() {
	query := strings.Join(os.Args[1:], " ")

	rawData, err := getData(query)
	if err != nil {
		fmt.Println("Error: Failed to get data:", err)
		return
	}

	parsedData, err := unmarshal(rawData)
	if err != nil {
		fmt.Println("Error: Failed to Parse Data:", err)
		return
	}

	for i, question := range parsedData.Items {
		if i >= 3 {
			break
		}
		rawAnswers, err := getAnswers(question.QuestionID)
		if err != nil {
			fmt.Println("Error: Failed to get answers:", err)
			continue
		}
		parsedAnswers, err := unmarshalAnswers(rawAnswers)
		if err != nil {
			fmt.Println("Error: Failed to parse answers:", err)
			continue
		}
		fmt.Printf("Title: %s\nScore: %d\nLink: %s\n", question.Title, question.Score, question.Link)
		if len(parsedAnswers.Items) > 0 {
			fmt.Printf("Answer: %s\n", parsedAnswers.Items[0].Body)
		}
		fmt.Println("=====================================================================================================================================\n")
	}
}

func getData(query string) ([]byte, error) {
	apiURL := "https://api.stackexchange.com/2.3/search/advanced?q=" + url.QueryEscape(query) + "&sort=votes&order=desc&site=stackoverflow"
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return rawData, nil
}

func unmarshal(rawData []byte) (StackJson, error) {
	var parsedData StackJson
	err := json.Unmarshal(rawData, &parsedData)
	if err != nil {
		fmt.Println("Error Unmarshaling data", err)
		return parsedData, err
	}
	return parsedData, nil
}

func getAnswers(questionID int) ([]byte, error) {
	apiURL := "https://api.stackexchange.com/2.3/questions/" + strconv.Itoa(questionID) + "/answers?sort=votes&order=desc&site=stackoverflow&filter=withbody"
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rawAnswers, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return rawAnswers, nil
}

func unmarshalAnswers(rawAnswers []byte) (AnswersJson, error) {
	var parsedAnswers AnswersJson
	err := json.Unmarshal(rawAnswers, &parsedAnswers)
	if err != nil {
		return parsedAnswers, err
	}
	return parsedAnswers, nil
}
