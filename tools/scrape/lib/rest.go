package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Restpost(url string, payload []byte, result any) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("get http status %s", resp.Status)
	}
	return nil
}
