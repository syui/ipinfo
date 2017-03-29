package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	path, query, contenttype, body string
}

func TestApp(t *testing.T) {

	response := &Response{
		path:        "/v1/media/popular",
		contenttype: "application/json",
		body: `{
			"meta": {
				"code": 200
			},
			"data": [{
				"attribution": null,
				"comments": {
					"count": 0,
					"data": []
				},
				"filter": "Normal",
				"created_time": "1407830189",
				"id": "785258998994306042_21773839"
			}]
		}`,
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check request.
		if g, w := r.URL.Path, response.path; g != w {
			t.Errorf("request got path %s, want %s", g, w)
		}

		// Send response.
		w.Header().Set("Content-Type", response.contenttype)
		io.WriteString(w, response.body)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	res := fetch(server.URL + "/v1/media/popular")

	var media map[string]interface{}
	err := json.Unmarshal(res, &media)
	if err != nil {
		t.Fatal(err)
	}

	dataArr := media["data"].([]interface{})
	data := dataArr[0].(map[string]interface{})

	if data["id"] != "785258998994306042_21773839" {
		t.Fatal("Error")
	}
}
