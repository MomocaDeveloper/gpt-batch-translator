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
	prompt := fmt.Sprintf("对以下json串中出现的中文词汇，尽量使用我给出的英文作为翻译结果，保证不影响上述规则的基础上进行翻译，不要向我解释说明，也不要给我原文。\n%s", jsonStr)
	//prompt := fmt.Sprintf("\nFor the Chinese vocabulary that appears in the following JSON string, please use the English translations that I provide.\n%s\n", jsonStr)
	return prompt, nil
}

