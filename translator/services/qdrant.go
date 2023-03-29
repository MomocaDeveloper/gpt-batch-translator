package services

import (
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
)

type Answer struct{
	Title string `json:"title"`
	Text string `json:"text"`
}
	
type SearchResult  struct{
	Code int `json:"code"`
	Data struct{
		Results []Answer `json:"result"`
	} `json:"data"`
}

func CreateAdditionPrompt1(){
	return
}

func GetSearchResult(sentences string)([]Answer, error){
	data := map[string]interface{}{
		"search": sentences,
	}

	body, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		fmt.Println("GetSearchResult on marshal json fail", marshalErr)
		return []Answer{}, marshalErr
	}

	req, searchErr := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/search", bytes.NewReader(body))
	if searchErr != nil {
		fmt.Println("GetSearchResult on search fail", searchErr)
		return []Answer{}, searchErr
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, httpErr := client.Do(req)
	if httpErr != nil {
		fmt.Println("GetSearchResult on send request fail", httpErr)
		return []Answer{}, httpErr
	}
	defer resp.Body.Close()

	var result SearchResult
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil{
		fmt.Println("GetSearchResult on decode response fail", decodeErr)
		return []Answer{}, decodeErr
	}
	return result.Data.Results, nil
}
