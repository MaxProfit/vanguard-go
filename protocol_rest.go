// Copyright 2023 Buf Technologies, Inc.
//
// All rights reserved.

//nolint:forbidigo,revive,gocritic // this is temporary, will be removed when implementation is complete
package vanguard

import (
	"io"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type restClientProtocol struct{}

var _ clientProtocolHandler = restClientProtocol{}
var _ clientBodyPreparer = restClientProtocol{}
var _ clientProtocolEndMustBeInHeaders = restClientProtocol{}

// restClientProtocol implements the REST protocol for
// processing RPCs received from the client.
func (r restClientProtocol) protocol() Protocol {
	return ProtocolREST
}

func (r restClientProtocol) acceptsStreamType(op *operation, streamType connect.StreamType) bool {
	switch streamType {
	case connect.StreamTypeUnary:
		return true
	case connect.StreamTypeClient:
		return requestIsSpecialHTTPBody(op)
	case connect.StreamTypeServer:
		// TODO: support server streams even when body is not google.api.HttpBody
		return responseIsSpecialHTTPBody(op)
	default:
		return false
	}
}

func (r restClientProtocol) endMustBeInHeaders() bool {
	// TODO: when we support server streams over REST, this should return false when streaming
	return true
}

func (r restClientProtocol) extractProtocolRequestHeaders(op *operation, headers http.Header) (requestMeta, error) {
	var reqMeta requestMeta
	reqMeta.compression = headers.Get("Content-Encoding")
	headers.Del("Content-Encoding")
	// TODO: A REST client could use "q" weights in the `Accept-Encoding` header, which
	//       would currently cause the middleware to not recognize the compression.
	//       We may want to address this. We'd need to sort the values by their weight
	//       since other protocols don't allow weights with acceptable encodings.
	reqMeta.acceptCompression = parseMultiHeader(headers.Values("Accept-Encoding"))
	headers.Del("Accept-Encoding")

	reqMeta.codec = CodecJSON // if actually a custom content-type, handled by body preparer methods
	contentType := headers.Get("Content-Type")
	if contentType != "" &&
		contentType != "application/json" &&
		contentType != "application/json; charset=utf-8" &&
		!requestIsSpecialHTTPBody(op) {
		// invalid content-type
		reqMeta.codec = contentType + "?"
	}
	headers.Del("Content-Type")

	return reqMeta, nil
}

func (r restClientProtocol) addProtocolResponseHeaders(meta responseMeta, headers http.Header, allowedCompression []string) int {
	//TODO implement me
	panic("implement me")
}

func (r restClientProtocol) encodeEnd(codec Codec, end *responseEnd, writer io.Writer, wasInHeaders bool) http.Header {
	//TODO implement me
	panic("implement me")
}

func (r restClientProtocol) requestNeedsPrep(o *operation) bool {
	//TODO implement me
	panic("implement me")
}

func (r restClientProtocol) prepareUnmarshalledRequest(op *operation, src []byte, target proto.Message) error {
	//TODO implement me
	panic("implement me")
}

func (r restClientProtocol) responseNeedsPrep(o *operation) bool {
	//TODO implement me
	panic("implement me")
}

func (r restClientProtocol) prepareMarshalledResponse(op *operation, base []byte, src proto.Message, headers http.Header) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (r restClientProtocol) String() string {
	return protocolNameREST
}

// restServerProtocol implements the REST protocol for
// sending RPCs to the server handler.
type restServerProtocol struct{}

var _ serverProtocolHandler = restServerProtocol{}
var _ requestLineBuilder = restServerProtocol{}
var _ serverBodyPreparer = restServerProtocol{}

func (r restServerProtocol) protocol() Protocol {
	return ProtocolREST
}

func (r restServerProtocol) addProtocolRequestHeaders(meta requestMeta, headers http.Header, allowedCompression []string) {
	//TODO implement me
	panic("implement me")
}

func (r restServerProtocol) extractProtocolResponseHeaders(statusCode int, headers http.Header) (responseMeta, responseEndUnmarshaler, error) {
	//TODO implement me
	panic("implement me")
}

func (r restServerProtocol) extractEndFromTrailers(o *operation, headers http.Header) (responseEnd, error) {
	//TODO implement me
	panic("implement me")
}

func (r restServerProtocol) requestNeedsPrep(o *operation) bool {
	//TODO implement me
	panic("implement me")
}

func (r restServerProtocol) prepareMarshalledRequest(op *operation, base []byte, src proto.Message, headers http.Header) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (r restServerProtocol) responseNeedsPrep(o *operation) bool {
	//TODO implement me
	panic("implement me")
}

func (r restServerProtocol) prepareUnmarshalledResponse(op *operation, src []byte, target proto.Message) error {
	//TODO implement me
	panic("implement me")
}

func (r restServerProtocol) requiresMessageToProvideRequestLine(o *operation) bool {
	//TODO implement me
	panic("implement me")
}

func (r restServerProtocol) requestLine(op *operation, req proto.Message) (urlPath, queryParams, method string, includeBody bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r restServerProtocol) String() string {
	return protocolNameREST
}

func requestIsSpecialHTTPBody(op *operation) bool {
	return isSpecialHTTPBody(op.method.Input(), op.restTarget.requestBodyFields)
}

func responseIsSpecialHTTPBody(op *operation) bool {
	return isSpecialHTTPBody(op.method.Output(), op.restTarget.responseBodyFields)
}

func isSpecialHTTPBody(msg protoreflect.MessageDescriptor, bodyPath []protoreflect.FieldDescriptor) bool {
	if len(bodyPath) > 0 {
		field := bodyPath[len(bodyPath)-1]
		msg = field.Message()
	}
	return msg != nil && msg.FullName() == "google.api.HttpBody"
}
