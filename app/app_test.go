package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var appRouterObj *appRouter

func Test_checkIfPalindrome(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
	}{
		{"palindrome string",
			args{"malayalam"},
		},
		{"non-palindrome string",
			args{"sample"},
		},
	}

	for id, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkIfPalindrome(tt.args.str)
			if id == 0 {
				if !reflect.DeepEqual(*got, true) {
					t.Errorf("checkIfPalindrome() = %v, want %v", *got, true)
				}
			} else {
				if reflect.DeepEqual(got, true) {
					t.Errorf("checkIfPalindrome() = %v, want %v", got, true)
				}
			}
		})
	}
}
func Test_appRouter_getMessage_Neg(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/messages/{id}", nil)
	if err != nil {
		t.Error("test failed with error: ", err)
		return
	}
	w := httptest.NewRecorder()

	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"id": "1",
	}

	req = mux.SetURLVars(req, vars)

	appRouterObj.getMessage(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, `{"error":"input text ID 1 not found "}`, strings.TrimSuffix(string(w.Body.Bytes()), "\n"))
}

func Test_appRouter_createMessage_getMessage(t *testing.T) {
	var b bytes.Buffer
	b.WriteString(`{"text":"sample"}`)
	req, err := http.NewRequest("POST", "/v1/messages", &b)
	if err != nil {
		t.Error("test failed with error: ", err)
		return
	}
	w := httptest.NewRecorder()

	appRouterObj.createMessage(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"id":1,"text":"sample"}`, strings.TrimSuffix(string(w.Body.Bytes()), "\n"))

	// Test Get as well
	getReq, err := http.NewRequest("GET", "/v1/messages/{id}", nil)
	if err != nil {
		t.Error("test failed with error: ", err)
		return
	}
	w = httptest.NewRecorder()

	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"id": "1",
	}

	getReq = mux.SetURLVars(getReq, vars)

	appRouterObj.getMessage(w, getReq)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"id":1,"text":"sample"}`, strings.TrimSuffix(string(w.Body.Bytes()), "\n"))
}

func Test_appRouter_palindrome_createMessage_getMessage(t *testing.T) {
	var b bytes.Buffer
	b.WriteString(`{"text":"malayalam"}`)
	req, err := http.NewRequest("POST", "/v1/messages", &b)
	if err != nil {
		t.Error("test failed with error: ", err)
		return
	}
	w := httptest.NewRecorder()

	appRouterObj.createMessage(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"id":2,"text":"malayalam"}`, strings.TrimSuffix(string(w.Body.Bytes()), "\n"))

	// Test Get as well
	getReq, err := http.NewRequest("GET", "/v1/messages/{id}", nil)
	if err != nil {
		t.Error("test failed with error: ", err)
		return
	}
	w = httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"id": "2",
	}
	getReq = mux.SetURLVars(getReq, vars)

	param := getReq.URL.Query()
	param.Add("is-palindrome", "")
	getReq.URL.RawQuery = param.Encode()

	appRouterObj.getMessage(w, getReq)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"id":2,"text":"malayalam","is-palindrome":true}`, strings.TrimSuffix(string(w.Body.Bytes()), "\n"))
}

func Test_appRouter_getAllMessages(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/messages", nil)
	if err != nil {
		t.Error("test failed with error: ", err)
		return
	}
	w := httptest.NewRecorder()

	appRouterObj.getAllMessages(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `[{"id":1,"text":"sample"},{"id":2,"text":"malayalam"}]`, strings.TrimSuffix(string(w.Body.Bytes()), "\n"))
}

func Test_appRouter_deleteMessage(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/v1/messages/{id}", nil)
	if err != nil {
		t.Error("test failed with error: ", err)
		return
	}
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"id": "2",
	}
	req = mux.SetURLVars(req, vars)

	appRouterObj.deleteMessage(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestMain(m *testing.M) {
	appRouterObj = NewAppRouter(200)
	err := appRouterObj.InitTracing("jaeger", ":8080")
	if err != nil {
		print("Initialization of tests failed with error: ", err)
		os.Exit(-1)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}
