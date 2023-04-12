package services

import (
		"fmt"
		"strings"
		"strconv"
		)

type Translator = Content

var BASE_REQUIRE = []string{
	"请尽可能使用古英语进行翻译，在其中适量加入拉丁文，同时符合角色性格。可以对内容和句子结构进行符合英文本土化方面的优化调整，译文台词的自然性、趣味性、生动性的优先级是最高的。",
	"翻译内容的每一行都必须保留阿拉伯数字序列，并可接受非连续数字。如果翻译的段落中有空行，你可以跳过它们。如果要翻译的内容已经是英文或超出了翻译能力，请保留原文。",
}

func GetBaseRequire(roleRequire string)[]string{
	return []string{
		// "请尽可能使用古英语进行翻译，在其中适量加入拉丁文，同时符合角色性格。可以对内容和句子结构进行符合英文本土化方面的优化调整，译文台词的自然性、趣味性、生动性的优先级是最高的。",
		roleRequire,
		"原文每一行的阿拉伯数字复制粘贴到翻译结果的最前面",
	}
}

func GetPromptHead()string{
	return "请将一款异世界动漫题材游戏中的角色台词翻译成英文，同时在翻译过程中需要你遵循以下要求：\n"
}

func GetPromptTail()string{
	return fmt.Sprintf("需要翻译的内容如下：")
}

func GetFullPrompt(roleRequire string, requires ...string)string{
	baseRequire := GetBaseRequire(roleRequire)
	full_require := append(baseRequire, requires...)
	prompt := GetPromptHead()
	ind := 'a'
	for _, require := range full_require{
		if len(require) == 0{
			continue
		}
		prompt += fmt.Sprintf("%c. %s\n", ind, require)
		ind += 1
	}
	prompt_tail := GetPromptTail()
	prompt += prompt_tail
	return prompt
}
func GetTempPrompt()string{
	//return "You are now a remarkable translator, responsible for translating the script of a game with a theme set in a different world. I will provide you with some Chinese vocabulary and sentences, with each line marked with a number in Arabic numerals. Your task is to translate the text into English while preserving the original numbering. You are not allowed to add any additional content or modify the original text in any way. Please skip any empty lines I may have included. If the text I provide is already in English or contains content that you cannot translate, please leave it as is. Do not provide any explanations for your translations."
	return ""
}

func (t *Translator) BuildMessageText()(int, string) {
	text := ""
	length := len(t.Cells)
	if length == 0{
		return 0, text
	}

	cells := t.Cells
	for _, cell := range cells{
		i := cell.Id
		source := cell.Source
		line := fmt.Sprintf("%d. %s\n", i, source)
		text = text + line
	}
	return len(cells), text
}

func (t *Translator) ParseTranslateResponse(respStr string, isSec bool){
	curIndex := 0
	//fmt.Println("response string:", respStr)
	resultSlice := strings.Split(respStr, "\n")
	defer func(){
		if r:=recover();r!=nil{
			fmt.Println("Recovered:", r)
		}
	}()
	for _, line := range(resultSlice){
		if len(line)==0{
			continue
		}
		if curIndex > len(t.Cells) {
			panic("something was wrong, index out of range")
			break
		}
		dotIndex := strings.Index(line, ".")
		if dotIndex == -1 {
			continue
		}
		id, err := strconv.Atoi(line[:dotIndex])
		if err != nil{
			fmt.Println("parse translate response failed", err)
			continue
		}
		target := line[dotIndex+1:]
		if id!= t.Cells[curIndex].Id{
			fmt.Println("parse translate result id not equal", id, t.Cells[curIndex].Id)	
		}
		if isSec == true{
			t.Cells[curIndex].SecTarget = strings.TrimSpace(target)
		}else{
			t.Cells[curIndex].Target = strings.TrimSpace(target)
		}
		curIndex += 1
	}	
}

func (t *Translator) TranslateMultiLines(gpt *ChatGPT, baseNoun *map[string]string, character *map[string]characterDes, temperatureStr string)(int,error){
	temperature, _ := strconv.ParseFloat(temperatureStr, 32)
	if temperature > 1{
		temperature = 1
	}
	if temperature < 0{
		temperature = 0
	}
	line, command := t.BuildMessageText()
	if len(command)==0{
		fmt.Println("translate finish~")
		return line, nil
	}

	mentionedNoun := GetMentionedNoun(command, baseNoun)
	require := t.Require
	if len(require)>0{
		require = fmt.Sprintf("额外要求:%s.", require)
	}
	keyword, err := CreateNounPrompt(&mentionedNoun)
	if err != nil{
		return line, err
	}

	charInfo := (*character)[t.Role].roleDes

	prompt := GetFullPrompt(charInfo, keyword, require)

	msgs := []Messages{
		{Role: "system", Content:prompt},
		{Role: "user", Content:command},
	}
	resp, err := gpt.Completions(msgs, float32(temperature))
	if err != nil {
		fmt.Println("try translate multi lines failed", err)
		return line, err
	}

	fmt.Println("resp",resp.Content)
	t.ParseTranslateResponse(resp.Content, false)
	roleTranslateDes := (*character)[t.Role].roleTranslateDes
	if len(roleTranslateDes) == 0 {
		return line, nil
	}

	msgs = []Messages{
		{Role: "system", Content:prompt},
		{Role: "user", Content:command},
		{Role: "assistant", Content:resp.Content},
		{Role: "user", Content:roleTranslateDes},
	}

	secResp, secErr := gpt.Completions(msgs, float32(temperature))
	if secErr != nil {
		fmt.Println("try translate sec multi lines failed", secErr)
		return line, secErr
	}

	fmt.Println("resp sec:",secResp.Content)
	t.ParseTranslateResponse(secResp.Content, true)
	return line, nil
}
