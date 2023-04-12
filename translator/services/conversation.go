package services

import (
	"fmt"
	"time"
)

type Character struct{
	Id string
	Prompt string
	Gpt *ChatGPT
	Msg []Messages
}

func NewCharacter(name, prompt string, gpt *ChatGPT)Character{
	msgs := []Messages{
		Messages{
			"system",
			prompt,
		},
	}
	return Character{
		name,
		prompt,
		gpt,
		msgs,
	}
}

func (c *Character)sendRequest(msgStr string, temperature float32)(string, error){
	c.Msg = append(c.Msg, Messages{
		"user",
		msgStr,
	})
	resp, err := c.Gpt.Completions(c.Msg, temperature)
	if err != nil{
		return "", err
	}
	c.Msg = append(c.Msg, Messages{
		"assistant",
		resp.Content,
	})
	return resp.Content, nil
}

func SimulateConversation(char1, char2 Character, start string, loopCount int, temperature float32)[]Messages{
	lastResponse := start
	var err error
	history := []Messages{
		Messages{
			char2.Id, 
			start,
		},
	}
	char2.Msg = append(char2.Msg, Messages{
		"assistant",
		start,
	})
	for i:=0;i<loopCount;i++{
		time.Sleep(1*time.Second)
		//sendMsg := fmt.Sprintf("%s:%s", char2.Id, lastResponse)
		lastResponse, err = char1.sendRequest(lastResponse, temperature)
		if err != nil {
			fmt.Println("simulationConversation failed", err)
			return history
		}
		time.Sleep(1*time.Second)
		history = append(history, Messages{
			char1.Id, 
			lastResponse,
		})
		//sendMsg = fmt.Sprintf("%s:%s", char1.Id, lastResponse)
		lastResponse, err = char2.sendRequest(lastResponse, temperature)
		if err != nil{
			fmt.Println("simulation Conversation failed", err)
			return history
		}
		history = append(history, Messages{
			char2.Id, 
			lastResponse,
		})
	}
	return history
	
}