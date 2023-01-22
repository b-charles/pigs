package log

import (
	"bufio"
	"io"
	"os"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
)

type Appender interface {
	Append(json.JsonNode)
}

// Default: Buffered stdout

type BufferedWriterAppender struct {
	writer *bufio.Writer
}

func NewBufferedWriterAppender(w io.Writer) *BufferedWriterAppender {
	return &BufferedWriterAppender{bufio.NewWriter(w)}
}

func (self *BufferedWriterAppender) Append(node json.JsonNode) {
	if _, err := self.writer.WriteString(node.String()); err != nil {
		panic(err)
	}
}

func (self *BufferedWriterAppender) Close() error {
	return self.writer.Flush()
}

func init() {

	ioc.DefaultPutFactory(func() (*BufferedWriterAppender, error) {
		return NewBufferedWriterAppender(os.Stdout), nil
	}, func(Appender) {})

}
