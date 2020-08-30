// Package classification Messaging Service.
//
// Documentation of Messaging Service API.
//
//     Schemes: http
//     BasePath: /
//     Version: v1
//     Host: localhost
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package docs

import "github.com/shailendra-k-singh/example.messaging.service/message"

// swagger:route POST /v1/messages create-messages createMessageRequest
// Creates a message record based on input text and returns the same.
// responses:
//   200: createMessageResponse

// Returns the created Message record containing system created ID and input text.
// swagger:response createMessageResponse
type createMessageResponseWrapper struct {
	// in:body
	Body struct {
		Id           int64  `json:"id"`
		Text         string `json:"text"`
	}
}

// swagger:parameters createMessageRequest
type createMessageRequestWrapper struct {
	// Accepts a string text as input
	// in:body
	Body struct {
		Text string `json:"text"`
	}
}

// swagger:route GET /v1/messages/{id}?is-palindrome get-message-pcheck checkIfPalindrome
// Retrieves a message with input id as path param, performs a palindrome check on the text
// and returns the result in field "is-palindrome" in the response, else returns error if message not found.
// responses:
//   200: getPMessageSuccResponse
//	 404: getPMessageFailResponse

// Returns the specified message record.
// swagger:response getPMessageSuccResponse
type getPMessageSuccResponseWrapper struct {
	// in:body
	Body message.MessageObj
}

// Returns error response in case of any failure.
// swagger:response getPMessageFailResponse
type getPMessageFailResponseWrapper struct {
	// in:body
	Body struct {
		Error string `json:"error"`
	}
}

// Message id and query parameter.
// swagger:parameters checkIfPalindrome
type isPalindromeCheckWrapper struct {
	// in:path
	// required:true
	Id int64	`json:"id"`
	// in:query
	// required:true
	IsPalindrome int64	`json:"is-palindrome"`
}

// swagger:route GET /v1/messages/{id} get-message getMessageID
// Retrieves a message with input id as path param. Returns error if message not found.
// responses:
//   200: getMessageSuccResponse
//	 404: getMessageFailResponse

// Returns the specified message record.
// swagger:response getMessageSuccResponse
type getMessageSuccResponseWrapper struct {
	// in:body
	Body struct {
		Id           int64  `json:"id"`
		Text         string `json:"text"`
	}
}

// Returns error response in case of any failure.
// swagger:response getMessageFailResponse
type getMessageFailResponseWrapper struct {
	// in:body
	Body struct {
		Error string `json:"error"`
	}
}

// The message id.
// swagger:parameters getMessageID
type getMessageIDWrapper struct {
	// in:path
	// required:true
	Id int64	`json:"id"`
}

// swagger:route GET /v1/messages get-all-messages getAllMessageID
// Retrieves all created messages.
// responses:
//   200: getAllMessagesSuccResponse
//	 404: getAllMessagesFailResponse

// Returns all the stored message records.
// swagger:response getAllMessagesSuccResponse
type getAllMessageSuccResponseWrapper struct {
	// in:body
	Body []struct{
		Id           int64  `json:"id"`
		Text         string `json:"text"`
	}
}

// Returns error response in case of any failure.
// swagger:response getAllMessagesFailResponse
type getAllMessageFailResponseWrapper struct {
	// in:body
	Body struct{
		Error string `json:"error"`
	}
}

// swagger:parameters getAllMessageID
type getAllMessageIDWrapper struct {
}

// swagger:route DELETE /v1/messages/{id} del-message getDelMessageID
// Deletes a message with input id as path param. Returns error if message not found.
// responses:
//   204: delMessageSuccResponse
//	 404: delMessageFailResponse


// Returns error response in case of any failure.
// swagger:response delMessageFailResponse
type getDelMessageFailResponseWrapper struct {
	// in:body
	Body struct {
		Error string `json:"error"`
	}
}

// The message id.
// swagger:parameters getDelMessageID
type getDMessageIDWrapper struct {
	// in:path
	// required:true
	Id int64	`json:"id"`
}

// Returns no content.
// swagger:response delMessageSuccResponse
type delMessageIDWrapper struct {

}