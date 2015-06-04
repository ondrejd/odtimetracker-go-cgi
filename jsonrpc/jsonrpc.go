// Copyright 2015 Ondřej Doněk. All rights reserved.
// See LICENSE file for more details about licensing.

/*
This file contains all what we need to consume and produce JSON.
*/
package main

// Define error for our JSON-RPC response
type ResponseError struct {
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

// JSON-RPC "Parse error" error definition
var ParseError = &ResponseError{
	Code:    32700,
	Message: "Parse error"
}

// JSON-RPC "Invalid request" error definition
var InvalidRequest = &ResponseError{
	Code:    32600,
	Message: "Invalid request"
}

// JSON-RPC "Method not found" error definition
var MethodNotFound = &ResponseError{
	Code:    32601,
	Message: "Method not found"
}

// JSON-RPC "Invalid params" error definition
var InvalidParams = &ResponseError{
	Code:    32602,
	Message: "Invalid params"
}

// JSON-RPC "Internal error" error definition
var InternalError = &ResponseError{
	Code:    32603,
	Message: "Internal error"
}

// JSON-RPC "Server error" error definition
var ServerError = &ResponseError{
	Code:    32000,
	Message: "Server error",
}

// Define simple JSON-RPC response
type JsonRpcResponse struct {
	JsonRpc string // Contains version of JSON-RPC protocol. MUST be exactly "2.0."
	Result interface{} // This member is REQUIRED on success. This member MUST NOT exist if there was an error invoking the method.
	Error *ResponseError // This member is REQUIRED on error. This member MUST NOT exist if there was no error triggered during invocation.
	Id string // This member is REQUIRED. It MUST be the same as the value of the id member in the Request Object. If there was an error in detecting the id in the Request object (e.g. Parse error/Invalid Request), it MUST be Null.
}
