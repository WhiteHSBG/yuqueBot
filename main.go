package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strings"
)

func gatherResponse(response *http.Response) (string, error) {
	contentType := response.Header.Get("content-type")
	if strings.Contains(contentType, "application/json") {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return string(bodyBytes), nil
	}
	return "", nil
}

//func sendMdMsg(botKey, content string) (string, error) {
//	baseUrl := "https://qyapi.weixin.qq.com/cgi-bin/webhook/"
//	url := baseUrl + "send?key=" + botKey
//
//	msg := map[string]interface{}{
//		"msgtype": "markdown",
//		"markdown": map[string]interface{}{
//			"content": content,
//		},
//	}
//
//	body, err := json.Marshal(msg)
//	if err != nil {
//		return "", err
//	}
//
//	response, err := http.Post(url, "application/json", bytes.NewReader(body))
//	if err != nil {
//		return "", err
//	}
//	defer response.Body.Close()
//
//	return gatherResponse(response)
//}

func sendMdMsg(botKey, content string) (string, error) {
	baseUrl := "https://qyapi.weixin.qq.com/cgi-bin/webhook/"
	url1 := baseUrl + "send?key=" + botKey

	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": content,
		},
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	request, err := http.NewRequest("POST", url1, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	return gatherResponse(response)
}

type RequestBody struct {
	Data struct {
		Book struct {
			Name string `json:"name"`
		} `json:"book"`
		Actor struct {
			Name string `json:"name"`
		} `json:"actor"`
		Title      string `json:"title"`
		Path       string `json:"path"`
		ActionType string `json:"action_type"`
	} `json:"data"`
}

func handleMdMsg(botKey string, reqBody *RequestBody) (string, error) {
	//fmt.Printf("Received Request Body: %+v\n", reqBody)
	bookName := reqBody.Data.Book.Name
	userName := reqBody.Data.Actor.Name
	articleName := reqBody.Data.Title
	articleUrl := "https://www.yuque.com/" + reqBody.Data.Path

	actionWords := map[string]string{
		"update":  "更新",
		"publish": "发布",
		"delete":  "删除",
	}

	mdMsg := fmt.Sprintf("%s %s了《[%s](%s)》 来自%s", userName, actionWords[reqBody.Data.ActionType], articleName, articleUrl, bookName)
	return sendMdMsg(botKey, mdMsg)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			botKey := r.URL.Query().Get("key")
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			var reqBody RequestBody
			if err := jsoniter.Unmarshal(bodyBytes, &reqBody); err != nil {
				http.Error(w, "Failed to parse request body", http.StatusBadRequest)
				return
			}

			result, err := handleMdMsg(botKey, &reqBody)
			if err != nil {
				http.Error(w, "Failed to handle message", http.StatusInternalServerError)
				return
			}

			w.Write([]byte(result))
		} else {
			w.Write([]byte("Hello World!"))
		}
	})

	http.ListenAndServe("0.0.0.0:5990", nil)
}
