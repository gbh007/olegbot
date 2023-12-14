package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type ExportData struct {
	Messages []struct {
		TextEntities []struct {
			Text string `json:"text"`
		} `json:"text_entities"`
	} `json:"messages"`
}

var rgx = regexp.MustCompile(`.+\(\s*[cCсС]\s*\).*`)

func main() {
	fileIn, err := os.Open("in.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer fileIn.Close()

	out := make([]string, 0)

	exp := new(ExportData)
	err = json.NewDecoder(fileIn).Decode(exp)
	if err != nil {
		log.Fatalln(err)
	}

	for _, msg := range exp.Messages {
		text := ""
		for _, entr := range msg.TextEntities {
			text += entr.Text
		}

		if !strings.Contains(text, "\n") && rgx.MatchString(text) {
			out = append(out, text)
			fmt.Println(text)
		}
	}

	fileOut, err := os.Create("out.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer fileOut.Close()

	err = json.NewEncoder(fileOut).Encode(out)
	if err != nil {
		log.Fatalln(err)
	}
}
