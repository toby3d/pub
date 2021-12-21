// Package provides a media HTTP endpoints.
//
// To upload a file to the Media Endpoint, the client sends a
// `multipart/form-data` request with one part named file. The Media Endpoint
// MAY ignore the suggested filename that the client sends.
//
// The Media Endpoint processes the file upload, storing it in whatever backend
// it wishes, and generates a URL to the file. The URL SHOULD be unguessable,
// such as using a UUID in the path. If the request is successful, the endpoint
// MUST return the URL to the file that was created in the HTTP Location header,
// and respond with HTTP 201 Created. The response body is left undefined.
//
// The Micropub client can then use this URL as the value of e.g. the "photo"
// property of a Micropub request.
//
// The Media Endpoint MAY periodically delete files uploaded if they are not
// used in a Micropub request within a specific amount of time.
package http
