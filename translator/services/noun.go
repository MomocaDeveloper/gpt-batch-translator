package services

import (
	"fmt"
	"strings"
	"encoding/json"
)

func GetMentionedNoun(content string, nounMap *map[string]string)map[string]string{
	result := make(map[string]string)	
	for key, val := range *nounMap{
		if strings.Contains(content, key){
			result[key] = val
		}
	}
	return result
}

func CreateNounPrompt(nounMap *map[string]string)(string, error){
	if len(*nounMap) == 0{
		return "", nil
	}	
	jsonStr, marshalErr := json.Marshal(*nounMap)
	if marshalErr != nil {
		fmt.Println("CreateAdditionPrompt try marshal nounMap fail", marshalErr)
		return "", marshalErr
	}
	prompt := fmt.Sprintf("\n对以下json串中出现的中文词汇，使用我给出的英文作为翻译结果\n%s\n", jsonStr)
	//prompt := fmt.Sprintf("\nFor the Chinese vocabulary that appears in the following JSON string, please use the English translations that I provide.\n%s\n", jsonStr)
	return prompt, nil
}

