package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/shailendra-k-singh/example.messaging.service/message"
	log "github.com/sirupsen/logrus"
)

const (
	defaultBitsize = 64
)

type msgRequestBody struct {
	Text string `json:"text"`
}

type errMsg struct {
	Error string `json:"error"`
}

type appRouter struct {
	router *mux.Router
	m      *message.MessageServer
	limit  int
	t      *tracerObj
}

func NewAppRouter(limit int) *appRouter {
	return &appRouter{router: mux.NewRouter(), m: message.NewMessageServer(), limit: limit, t: &tracerObj{}}
}

func (r *appRouter) GetRouter() *mux.Router {
	return r.router
}

func (r *appRouter) SetRoutes() {
	r.router.Use(r.t.startTracing)

	r.router.Methods("GET").Path("/v1/messages").HandlerFunc(r.getAllMessages)
	r.router.Methods("POST").Path("/v1/messages").HandlerFunc(r.createMessage)
	r.router.Methods("GET").Path("/v1/messages/{id}").HandlerFunc(r.getMessage)
	r.router.Methods("DELETE").Path("/v1/messages/{id}").HandlerFunc(r.deleteMessage)

	// A default root handler
	r.router.Methods("GET").Path("/").HandlerFunc(r.root)
}

func jsonResponse(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("Error while writing JSON response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		i, err := w.Write([]byte(`{"error": "Error while writing response object"}`))
		if err != nil || i < 1 {
			log.Error("Error while writing default error to http.ResponseWriter: ", err)
		}
	}
}

func respondWithError(w http.ResponseWriter, message errMsg, code int) {
	jsonResponse(w, message, code)
}

func (r *appRouter) createMessage(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Error while reading request body: ", err)
		r.addSpan(req.Context(), http.StatusBadRequest, req)
		respondWithError(w, errMsg{"Invalid request body"}, http.StatusBadRequest)
		return
	}
	log.Debug("POST: Received request body: ", string(body))

	msg := msgRequestBody{}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		log.Error("error while unmarshalling request body: ", err)
		r.addSpan(req.Context(), http.StatusBadRequest, req)
		respondWithError(w, errMsg{"Invalid request body"}, http.StatusBadRequest)
		return
	}

	// Basic length sanity check
	l := len(msg.Text)
	if l < 1 {
		log.Error("incorrect input format or message length zero")
		r.addSpan(req.Context(), http.StatusBadRequest, req)
		respondWithError(w, errMsg{"Invalid input body, must be a non-zero length string in specified format"}, http.StatusBadRequest)
		return
	}
	if l > r.limit {
		log.Errorf("message length %d greater than limit %d", l, r.limit)
		msg := fmt.Sprintf("Input text length must be in range 1-%d", r.limit)
		r.addSpan(req.Context(), http.StatusBadRequest, req)
		respondWithError(w, errMsg{msg}, http.StatusBadRequest)
		return
	}

	resp := r.m.Add(msg.Text)
	r.addSpan(req.Context(), http.StatusOK, req)
	jsonResponse(w, resp, http.StatusOK)
	log.Infof("Added message with id %v successfully", resp.Id)
}

func (r *appRouter) validateMsgID(w http.ResponseWriter, req *http.Request) (int64, error) {
	var err error
	vars := mux.Vars(req)
	log.Infof("Received route params as: %#v", vars)
	id, ok := vars["id"]
	if !ok {
		err = fmt.Errorf("id param not present in request path")
		r.addSpan(req.Context(), http.StatusBadRequest, req)
		respondWithError(w, errMsg{"Message id not passed in the request. Retry in the format: /v1/messages/{id} "}, http.StatusBadRequest)
		return 0, err
	}

	val, err := strconv.ParseInt(id, 10, defaultBitsize)
	if err != nil || val < 1 {
		err = fmt.Errorf("invalid message id value: %s,err:%s ", id, err)
		r.addSpan(req.Context(), http.StatusBadRequest, req)
		respondWithError(w, errMsg{"Invalid message id value, should be a valid positive integer"}, http.StatusBadRequest)
		return 0, err
	}
	log.Debug("Validated message ID successfully")
	return val, err
}

func (r *appRouter) getMessage(w http.ResponseWriter, req *http.Request) {
	id, err := r.validateMsgID(w, req)
	if err != nil {
		log.Error("error validating request: ", err)
		return
	}
	resp, err := r.m.Get(id)
	if err != nil {
		log.Error("error while retrieving message: ", err)
		r.addSpan(req.Context(), http.StatusNotFound, req)
		respondWithError(w, errMsg{err.Error()}, http.StatusNotFound)
		return
	}
	// check for optional query param "is-palindrome"
	param := req.URL.Query()
	if val, ok := param["is-palindrome"]; ok {
		log.Debugf("value of query param is %#v:", val)
		if len(val) > 0 && val[0] != "" {
			log.Error("query param passed incorrectly: ", param)
			r.addSpan(req.Context(), http.StatusBadRequest, req)
			respondWithError(w, errMsg{"Incorrect URL structure, should be passed as: /v1/messages/{id}?is-palindrome "}, http.StatusBadRequest)
			return
		}
		resp.IsPalindrome = checkIfPalindrome(resp.Text)
		log.Info("Result of Palindrome check: ", *resp.IsPalindrome)
	}
	r.addSpan(req.Context(), http.StatusOK, req)
	jsonResponse(w, resp, http.StatusOK)
	log.Infof("Retrieved message %d successfully", id)
}

func (r *appRouter) getAllMessages(w http.ResponseWriter, req *http.Request) {
	resp, err := r.m.GetAll()
	if err != nil {
		log.Error("error while retrieving messages: ", err)
		r.addSpan(req.Context(), http.StatusNotFound, req)
		respondWithError(w, errMsg{err.Error()}, http.StatusNotFound)
		return
	}
	r.addSpan(req.Context(), http.StatusOK, req)
	jsonResponse(w, resp, http.StatusOK)
	log.Info("Retrieved all messages successfully")
}

func (r *appRouter) deleteMessage(w http.ResponseWriter, req *http.Request) {

	id, err := r.validateMsgID(w, req)
	if err != nil {
		log.Error("error validating request: ", err)
		r.addSpan(req.Context(), http.StatusBadRequest, req)
		respondWithError(w, errMsg{err.Error()}, http.StatusBadRequest)
		return
	}
	err = r.m.Delete(id)
	if err != nil {
		log.Error("error while deleting message: ", err)
		r.addSpan(req.Context(), http.StatusNotFound, req)
		respondWithError(w, errMsg{err.Error()}, http.StatusNotFound)
		return
	}
	r.addSpan(req.Context(), http.StatusNoContent, req)
	jsonResponse(w, struct{}{}, http.StatusNoContent)
	log.Infof("Deleted message %d successfully", id)
}

func checkIfPalindrome(str string) *bool {
	res := new(bool)
	rev := ""

	for _, val := range str {
		rev = string(val) + rev
	}
	*res = str == rev
	return res
}

// Dummy root handler
func (r *appRouter) root(w http.ResponseWriter, req *http.Request) {
	r.addSpan(req.Context(), http.StatusOK, req)
	jsonResponse(w, fmt.Sprintf("Welcome, it's %s now", time.Now()), http.StatusOK)
}

func (r *appRouter) InitTracing(lib string, host string, addr string) error {
	return r.t.initTracing(lib, host, addr)
}

func (r *appRouter) Close() {
	r.t.Close()
}

func (r *appRouter) addSpan(ctx context.Context, code int, req *http.Request) {
	// retrieve current Span from Context
	var parentCtx opentracing.SpanContext
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		parentCtx = parentSpan.Context()
	}

	// start a new Span to wrap HTTP request
	span := r.t.tracer.StartSpan(
		req.URL.Path,
		opentracing.ChildOf(parentCtx),
	)
	span.SetTag(string(ext.Component), "HTTP server")
	span.SetTag(string(ext.HTTPStatusCode), code)
	span.SetTag(string(ext.HTTPUrl), req.URL.Path)
	span.SetTag(string(ext.HTTPMethod), req.Method)
	span.Finish()
}
