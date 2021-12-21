// Package media provide a better user experience for Micropub applications, as
// well as to overcome the limitation of being unable to upload a file with the
// JSON syntax, a Micropub server MAY support a "Media Endpoint". The role of the
// Media Endpoint is exclusively to handle file uploads and return a URL that
// can be used in a subsequent Micropub request.
//
// When a Micropub server supports a Media Endpoint, clients can start uploading
// a photo or other media right away while the user is still creating other
// parts of the post.
package media
