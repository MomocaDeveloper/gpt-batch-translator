package services

import (
		"fmt"
		"strings"
		"strconv"
		)

type Translator = Content

func GetTempPrompt()string{
	return "You are now a remarkable translator, responsible for translating the script of a game with a theme set in a different world. I will provide you with some Chinese vocabulary and sentences, with each line marked with a number in Arabic numerals. Your task is to translate the text into English while preserving the original numbering. You are not allowed to add any additional content or modify the original text in any way. Please skip any empty lines I may have included. If the text I provide is already in English or contains content that you cannot translate, please leave it as is. Do not provide any explanations for your translations."
	//return "你现在是一名出色的翻译，为一款异世界题材的游戏翻译文案部分。接下来我会提供一些中文词汇和句子，我会用阿拉伯数字为每行要翻译的内容标上序号，你需要将我给出的文字翻译成英文。注意：序号需要原封不动地在结果中保留，你只能翻译我给出的内容，不能自己创作额外的东西。我的序号可能不连续。如果我给出的翻译段落有空行，你只需要跳过这些空行。如果我给出的内容本身就是英文或者是其他你无法翻译的内容，那么保留原文。无论你做了任何操作，都不要向我解释。"
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

func (t *Translator) ParseTranslateResponse(respStr string){
	curIndex := 0
	//fmt.Println("response string:", respStr)
	resultSlice := strings.Split(respStr, "\n")
	for _, line := range(resultSlice){
		if len(line)==0{
			continue
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
		t.Cells[curIndex].Target = target
		curIndex += 1
	}	
}

func (t *Translator) TranslateMultiLines(gpt *ChatGPT, baseNoun *map[string]string, character *map[string]string)(int,error){
	prompt := GetTempPrompt()
	line, command := t.BuildMessageText()
	if len(command)==0{
		fmt.Println("translate finish~")
		return line, nil
	}

	mentionedNoun := GetMentionedNoun(command, baseNoun)
	require := t.Require
	keyword, err := CreateNounPrompt(&mentionedNoun)
	if err != nil{
		return line, err
	}
	
	charInfo, _ := (*character)[t.Role]

	additionPrompt := charInfo + keyword + require
	//fmt.Println("additionPrompt", additionPrompt)
	//fmt.Println("command",command)
	msgs := []Messages{
		{Role: "system", Content:prompt + additionPrompt},
		{Role: "user", Content:command},
	}
	resp, err := gpt.Completions(msgs)
	if err != nil {
		fmt.Println("try translate multi lines failed", err)
		return line, err
	}

	fmt.Println("resp",resp.Content)
	t.ParseTranslateResponse(resp.Content)
	return line, nil
}
