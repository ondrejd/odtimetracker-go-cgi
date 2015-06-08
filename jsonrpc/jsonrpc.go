// Copyright 2015 Ondřej Doněk. All rights reserved.
// See LICENSE file for more details about licensing.

/*
This file contains all what we need to consume and produce JSON.
*/
package jsonrpc

// Define error for our JSON-RPC response
type Error struct {
	Code    int
	Message string
	Data    interface{}
}

// Here goes error defined according to JSON-RPC specification
// See: http://www.jsonrpc.org/specification#error_object
/*
-32602 	Invalid params 	Invalid method parameter(s).
-32603 	Internal error 	Internal JSON-RPC error.
-32000 to -32099 	Server error 	Reserved for implementation-defined server-errors.
*/

// JSON-RPC "Parse error" error definition.
var ParseError = &Error{
	Code:    32700,
	Message: "Parse error",
}

// JSON-RPC "Invalid request" error definition.
var InvalidRequest = &Error{
	Code:    32600,
	Message: "Invalid request",
}

// JSON-RPC "Method not found" error definition.
var MethodNotFound = &Error{
	Code:    32601,
	Message: "Method not found",
}

// JSON-RPC "Invalid params" error definition.
var InvalidParams = &Error{
	Code:    32602,
	Message: "Invalid params",
}

// JSON-RPC "Internal error" error definition.
var InternalError = &Error{
	Code:    32603,
	Message: "Internal error",
}

// JSON-RPC "Server error" error definition.
var ServerError = &Error{
	Code:    32000,
	Message: "Server error",
}

// Our custom error: "Storage initialization failed".
var InitStorageError = &Error{Code: 32001, Message: "Storage initialization failed"}

// Our custom error: "No running activity error".
var NoRunningActivityError = &Error{Code: 32002, Message: "There is no running activity"}

// Our custom error: "Updating activity failed".
var UpdateActivityError = &Error{Code: 32003, Message: "Updating activity failed"}

// Our custom error: "Another activity is running".
var AnotherRunningActivityError = &Error{Code: 32004, Message: "Another activity is running"}

// Our custom error: "Inserting new project failed".
var NewProjectError = &Error{Code: 32005, Message: "Inserting new project failed"}

// Our custom error: "Inserting new activity failed".
var NewActivityError = &Error{Code: 32006, Message: "Inserting new activity failed"}

// Define simple JSON-RPC response.
type Response struct {
	JsonRpc string      // Contains version of JSON-RPC protocol. MUST be exactly "2.0.".
	Result  interface{} // This member is REQUIRED on success. This member MUST NOT exist if there was an error invoking the method.
	Id      string      // This member is REQUIRED. It MUST be the same as the value of the id member in the Request Object. If there was an error in detecting the id in the Request object (e.g. Parse error/Invalid Request), it MUST be Null.
}

// Define simple JSON-RPC response.
type ErrorResponse struct {
	JsonRpc string // Contains version of JSON-RPC protocol. MUST be exactly "2.0.".
	Error   interface{}  // This member is REQUIRED on error. This member MUST NOT exist if there was no error triggered during invocation.
	Id      string // This member is REQUIRED. It MUST be the same as the value of the id member in the Request Object. If there was an error in detecting the id in the Request object (e.g. Parse error/Invalid Request), it MUST be Null.
}

// Create new JsonRpcResponse.
func NewResponse(result interface{}, id string) *Response {
	return &Response{JsonRpc: "2.0", Result: result, Id: id}
}

// Create new JsonRpcResponse.
func NewErrorResponse(err *Error, id string) *ErrorResponse {
	return &ErrorResponse{JsonRpc: "2.0", Error: err, Id: id}
}
