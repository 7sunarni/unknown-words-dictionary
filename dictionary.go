package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RawDictionary struct {
	Word     string       `json:"word,omitempty"`
	Meanings []RawMeaning `json:"meanings,omitempty"`
	Phonetic string       `json:"phonetic,omitempty"`
}

type RawMeaning struct {
	Speech      string          `json:"partOfSpeech,omitempty"`
	Definitions []RawDefinition `json:"definitions,omitempty"`
	Synonyms    []string        `json:"synonyms,omitempty"`
	Antonyms    []string        `json:"antonyms,omitempty"`
}

type RawDefinition struct {
	Def      string   `json:"definition,omitempty"`
	Synonyms []string `json:"synonyms,omitempty"`
	Antonyms []string `json:"antonyms,omitempty"`
	Example  string   `json:"example,omitempty"`
}

var dictionaryAPI = "https://api.dictionaryapi.dev/api/v2/entries/en/%s"

func getWord(w string) (MyDictionary, error) {
	resp, err := http.Get(fmt.Sprintf(dictionaryAPI, w))
	if err != nil {
		log.Println(err)
		return MyDictionary{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return MyDictionary{}, fmt.Errorf("get api response code %d", resp.StatusCode)
	}

	dics := make([]RawDictionary, 0)
	if err := json.NewDecoder(resp.Body).Decode(&dics); err != nil {
		log.Println(err)
		return MyDictionary{}, err
	}

	return NewMyDictionary(dics), nil
}

func NewMyDictionary(raw []RawDictionary) MyDictionary {
	if len(raw) == 0 {
		return MyDictionary{}
	}
	dictionry := MyDictionary{
		Word:     raw[0].Word,
		Phonetic: raw[0].Phonetic,
	}
	defs := make([]MyDefinition, 0)
	for _, dic := range raw {
		for i := range dic.Meanings {
			rawMeaning := dic.Meanings[i]
			for j := range rawMeaning.Definitions {
				rawDefinition := rawMeaning.Definitions[j]
				defs = append(defs, MyDefinition{
					Example:    rawDefinition.Example,
					Definition: rawDefinition.Def,
				})
			}
		}
	}
	dictionry.Defs = defs
	return dictionry
}

type MyDictionary struct {
	Word     string         `json:"word,omitempty"`
	Phonetic string         `json:"phonetic,omitempty"`
	Defs     []MyDefinition `json:"defs,omitempty"`
}

type MyDefinition struct {
	Definition string `json:"definition,omitempty"`
	Example    string `json:"example,omitempty"`
}
