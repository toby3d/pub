// Code generated by qtc from "baseof.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line web/template/baseof.qtpl:1
package template

//line web/template/baseof.qtpl:1
import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//line web/template/baseof.qtpl:6
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line web/template/baseof.qtpl:6
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line web/template/baseof.qtpl:6
type Page interface {
//line web/template/baseof.qtpl:6
	body() string
//line web/template/baseof.qtpl:6
	streambody(qw422016 *qt422016.Writer)
//line web/template/baseof.qtpl:6
	writebody(qq422016 qtio422016.Writer)
//line web/template/baseof.qtpl:6
	dir() string
//line web/template/baseof.qtpl:6
	streamdir(qw422016 *qt422016.Writer)
//line web/template/baseof.qtpl:6
	writedir(qq422016 qtio422016.Writer)
//line web/template/baseof.qtpl:6
	head() string
//line web/template/baseof.qtpl:6
	streamhead(qw422016 *qt422016.Writer)
//line web/template/baseof.qtpl:6
	writehead(qq422016 qtio422016.Writer)
//line web/template/baseof.qtpl:6
	lang() string
//line web/template/baseof.qtpl:6
	streamlang(qw422016 *qt422016.Writer)
//line web/template/baseof.qtpl:6
	writelang(qq422016 qtio422016.Writer)
//line web/template/baseof.qtpl:6
	title() string
//line web/template/baseof.qtpl:6
	streamtitle(qw422016 *qt422016.Writer)
//line web/template/baseof.qtpl:6
	writetitle(qq422016 qtio422016.Writer)
//line web/template/baseof.qtpl:6
	t(format message.Reference, a ...any) string
//line web/template/baseof.qtpl:6
	streamt(qw422016 *qt422016.Writer, format message.Reference, a ...any)
//line web/template/baseof.qtpl:6
	writet(qq422016 qtio422016.Writer, format message.Reference, a ...any)
//line web/template/baseof.qtpl:6
}

//line web/template/baseof.qtpl:16
type BaseOf struct {
	language language.Tag
	printer  *message.Printer
}

func NewBaseOf(lang language.Tag) *BaseOf {
	return &BaseOf{
		language: lang,
		printer:  message.NewPrinter(lang),
	}
}

//line web/template/baseof.qtpl:29
func (b *BaseOf) streamlang(qw422016 *qt422016.Writer) {
//line web/template/baseof.qtpl:29
	qw422016.N().S(`
`)
//line web/template/baseof.qtpl:30
	qw422016.E().S(b.language.String())
//line web/template/baseof.qtpl:30
	qw422016.N().S(`
`)
//line web/template/baseof.qtpl:31
}

//line web/template/baseof.qtpl:31
func (b *BaseOf) writelang(qq422016 qtio422016.Writer) {
//line web/template/baseof.qtpl:31
	qw422016 := qt422016.AcquireWriter(qq422016)
//line web/template/baseof.qtpl:31
	b.streamlang(qw422016)
//line web/template/baseof.qtpl:31
	qt422016.ReleaseWriter(qw422016)
//line web/template/baseof.qtpl:31
}

//line web/template/baseof.qtpl:31
func (b *BaseOf) lang() string {
//line web/template/baseof.qtpl:31
	qb422016 := qt422016.AcquireByteBuffer()
//line web/template/baseof.qtpl:31
	b.writelang(qb422016)
//line web/template/baseof.qtpl:31
	qs422016 := string(qb422016.B)
//line web/template/baseof.qtpl:31
	qt422016.ReleaseByteBuffer(qb422016)
//line web/template/baseof.qtpl:31
	return qs422016
//line web/template/baseof.qtpl:31
}

//line web/template/baseof.qtpl:33
func (b *BaseOf) streamdir(qw422016 *qt422016.Writer) {
//line web/template/baseof.qtpl:33
	qw422016.N().S(`
`)
//line web/template/baseof.qtpl:34
	for _, tag := range []language.Tag{
		language.Arabic,
		language.Hebrew,
		language.Persian,
		language.Urdu,
	} {
//line web/template/baseof.qtpl:39
		qw422016.N().S(`
`)
//line web/template/baseof.qtpl:40
		if b.language != tag {
//line web/template/baseof.qtpl:40
			qw422016.N().S(`
`)
//line web/template/baseof.qtpl:41
			continue
//line web/template/baseof.qtpl:42
		}
//line web/template/baseof.qtpl:42
		qw422016.N().S(`
rtl
`)
//line web/template/baseof.qtpl:44
		return
//line web/template/baseof.qtpl:45
	}
//line web/template/baseof.qtpl:45
	qw422016.N().S(`
ltr
`)
//line web/template/baseof.qtpl:47
}

//line web/template/baseof.qtpl:47
func (b *BaseOf) writedir(qq422016 qtio422016.Writer) {
//line web/template/baseof.qtpl:47
	qw422016 := qt422016.AcquireWriter(qq422016)
//line web/template/baseof.qtpl:47
	b.streamdir(qw422016)
//line web/template/baseof.qtpl:47
	qt422016.ReleaseWriter(qw422016)
//line web/template/baseof.qtpl:47
}

//line web/template/baseof.qtpl:47
func (b *BaseOf) dir() string {
//line web/template/baseof.qtpl:47
	qb422016 := qt422016.AcquireByteBuffer()
//line web/template/baseof.qtpl:47
	b.writedir(qb422016)
//line web/template/baseof.qtpl:47
	qs422016 := string(qb422016.B)
//line web/template/baseof.qtpl:47
	qt422016.ReleaseByteBuffer(qb422016)
//line web/template/baseof.qtpl:47
	return qs422016
//line web/template/baseof.qtpl:47
}

//line web/template/baseof.qtpl:49
func (b *BaseOf) streamtitle(qw422016 *qt422016.Writer) {
//line web/template/baseof.qtpl:49
	qw422016.N().S(`
Micropub
`)
//line web/template/baseof.qtpl:51
}

//line web/template/baseof.qtpl:51
func (b *BaseOf) writetitle(qq422016 qtio422016.Writer) {
//line web/template/baseof.qtpl:51
	qw422016 := qt422016.AcquireWriter(qq422016)
//line web/template/baseof.qtpl:51
	b.streamtitle(qw422016)
//line web/template/baseof.qtpl:51
	qt422016.ReleaseWriter(qw422016)
//line web/template/baseof.qtpl:51
}

//line web/template/baseof.qtpl:51
func (b *BaseOf) title() string {
//line web/template/baseof.qtpl:51
	qb422016 := qt422016.AcquireByteBuffer()
//line web/template/baseof.qtpl:51
	b.writetitle(qb422016)
//line web/template/baseof.qtpl:51
	qs422016 := string(qb422016.B)
//line web/template/baseof.qtpl:51
	qt422016.ReleaseByteBuffer(qb422016)
//line web/template/baseof.qtpl:51
	return qs422016
//line web/template/baseof.qtpl:51
}

//line web/template/baseof.qtpl:53
func (b *BaseOf) streamhead(qw422016 *qt422016.Writer) {
//line web/template/baseof.qtpl:53
}

//line web/template/baseof.qtpl:53
func (b *BaseOf) writehead(qq422016 qtio422016.Writer) {
//line web/template/baseof.qtpl:53
	qw422016 := qt422016.AcquireWriter(qq422016)
//line web/template/baseof.qtpl:53
	b.streamhead(qw422016)
//line web/template/baseof.qtpl:53
	qt422016.ReleaseWriter(qw422016)
//line web/template/baseof.qtpl:53
}

//line web/template/baseof.qtpl:53
func (b *BaseOf) head() string {
//line web/template/baseof.qtpl:53
	qb422016 := qt422016.AcquireByteBuffer()
//line web/template/baseof.qtpl:53
	b.writehead(qb422016)
//line web/template/baseof.qtpl:53
	qs422016 := string(qb422016.B)
//line web/template/baseof.qtpl:53
	qt422016.ReleaseByteBuffer(qb422016)
//line web/template/baseof.qtpl:53
	return qs422016
//line web/template/baseof.qtpl:53
}

//line web/template/baseof.qtpl:54
func (b *BaseOf) streambody(qw422016 *qt422016.Writer) {
//line web/template/baseof.qtpl:54
}

//line web/template/baseof.qtpl:54
func (b *BaseOf) writebody(qq422016 qtio422016.Writer) {
//line web/template/baseof.qtpl:54
	qw422016 := qt422016.AcquireWriter(qq422016)
//line web/template/baseof.qtpl:54
	b.streambody(qw422016)
//line web/template/baseof.qtpl:54
	qt422016.ReleaseWriter(qw422016)
//line web/template/baseof.qtpl:54
}

//line web/template/baseof.qtpl:54
func (b *BaseOf) body() string {
//line web/template/baseof.qtpl:54
	qb422016 := qt422016.AcquireByteBuffer()
//line web/template/baseof.qtpl:54
	b.writebody(qb422016)
//line web/template/baseof.qtpl:54
	qs422016 := string(qb422016.B)
//line web/template/baseof.qtpl:54
	qt422016.ReleaseByteBuffer(qb422016)
//line web/template/baseof.qtpl:54
	return qs422016
//line web/template/baseof.qtpl:54
}

//line web/template/baseof.qtpl:56
func (b *BaseOf) streamt(qw422016 *qt422016.Writer, format message.Reference, a ...any) {
//line web/template/baseof.qtpl:56
	qw422016.N().S(`
`)
//line web/template/baseof.qtpl:57
	qw422016.E().S(b.printer.Sprintf(format, a...))
//line web/template/baseof.qtpl:57
	qw422016.N().S(`
`)
//line web/template/baseof.qtpl:58
}

//line web/template/baseof.qtpl:58
func (b *BaseOf) writet(qq422016 qtio422016.Writer, format message.Reference, a ...any) {
//line web/template/baseof.qtpl:58
	qw422016 := qt422016.AcquireWriter(qq422016)
//line web/template/baseof.qtpl:58
	b.streamt(qw422016, format, a...)
//line web/template/baseof.qtpl:58
	qt422016.ReleaseWriter(qw422016)
//line web/template/baseof.qtpl:58
}

//line web/template/baseof.qtpl:58
func (b *BaseOf) t(format message.Reference, a ...any) string {
//line web/template/baseof.qtpl:58
	qb422016 := qt422016.AcquireByteBuffer()
//line web/template/baseof.qtpl:58
	b.writet(qb422016, format, a...)
//line web/template/baseof.qtpl:58
	qs422016 := string(qb422016.B)
//line web/template/baseof.qtpl:58
	qt422016.ReleaseByteBuffer(qb422016)
//line web/template/baseof.qtpl:58
	return qs422016
//line web/template/baseof.qtpl:58
}

//line web/template/baseof.qtpl:60
func StreamTemplate(qw422016 *qt422016.Writer, p Page) {
//line web/template/baseof.qtpl:60
	qw422016.N().S(`
<!DOCTYPE html>
<html lang="`)
//line web/template/baseof.qtpl:62
	p.streamlang(qw422016)
//line web/template/baseof.qtpl:62
	qw422016.N().S(`"
      dir="`)
//line web/template/baseof.qtpl:63
	p.streamdir(qw422016)
//line web/template/baseof.qtpl:63
	qw422016.N().S(`">

  <head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, initial-scale=1.0">
    `)
//line web/template/baseof.qtpl:69
	p.streamhead(qw422016)
//line web/template/baseof.qtpl:69
	qw422016.N().S(`
    <title>`)
//line web/template/baseof.qtpl:70
	p.streamtitle(qw422016)
//line web/template/baseof.qtpl:70
	qw422016.N().S(`</title>
  </head>

  <body>
    `)
//line web/template/baseof.qtpl:74
	p.streambody(qw422016)
//line web/template/baseof.qtpl:74
	qw422016.N().S(`
  </body>
</html>
`)
//line web/template/baseof.qtpl:77
}

//line web/template/baseof.qtpl:77
func WriteTemplate(qq422016 qtio422016.Writer, p Page) {
//line web/template/baseof.qtpl:77
	qw422016 := qt422016.AcquireWriter(qq422016)
//line web/template/baseof.qtpl:77
	StreamTemplate(qw422016, p)
//line web/template/baseof.qtpl:77
	qt422016.ReleaseWriter(qw422016)
//line web/template/baseof.qtpl:77
}

//line web/template/baseof.qtpl:77
func Template(p Page) string {
//line web/template/baseof.qtpl:77
	qb422016 := qt422016.AcquireByteBuffer()
//line web/template/baseof.qtpl:77
	WriteTemplate(qb422016, p)
//line web/template/baseof.qtpl:77
	qs422016 := string(qb422016.B)
//line web/template/baseof.qtpl:77
	qt422016.ReleaseByteBuffer(qb422016)
//line web/template/baseof.qtpl:77
	return qs422016
//line web/template/baseof.qtpl:77
}
