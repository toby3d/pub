package common

const charsetUTF8 = "charset=UTF-8"

const (
	HeaderAcceptLanguage      string = "Accept-Language"
	HeaderContentType         string = "Content-Type"
	HeaderLocation            string = "Location"
	HeaderXContentTypeOptions string = "X-Content-Type-Options"
	HeaderLink                string = "Link"
)

const (
	MIMEApplicationForm            string = "application/x-www-form-urlencoded"
	MIMEApplicationFormCharsetUTF8 string = MIMEApplicationForm + "; " + charsetUTF8
	MIMEApplicationJSON            string = "application/json"
	MIMEApplicationJSONCharsetUTF8 string = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEMultipartForm              string = "multipart/form-data"
	MIMEMultipartFormCharsetUTF8   string = MIMEMultipartForm + "; " + charsetUTF8
	MIMETextHTML                   string = "text/html"
	MIMETextHTMLCharsetUTF8        string = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                  string = "text/plain"
	MIMETextPlainCharsetUTF8       string = MIMETextPlain + "; " + charsetUTF8
)
