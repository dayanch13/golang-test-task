package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CheckSpelling(text string) error {
	resp, err := http.Get("https://speller.yandex.net/services/spellservice.json/checkText?text=" + text)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result []struct {
		Word string `json:"word"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	if len(result) > 0 {
		return fmt.Errorf("spelling errors found: %v", result)
	}
	return nil
}
