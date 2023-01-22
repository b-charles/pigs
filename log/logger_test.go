package log_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/json"
	. "github.com/b-charles/pigs/log"
	"github.com/benbjohnson/clock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Log Test Suite")
}

type BytesAppender struct {
	buffer bytes.Buffer
}

func (self *BytesAppender) Append(node json.JsonNode) {
	self.buffer.WriteString(node.String())
	self.buffer.WriteByte('\n')
}

func (self *BytesAppender) String() string {
	return self.buffer.String()
}

var _ = Describe("Log", func() {

	var fixedTime = "1984-07-31T13:13:00Z"

	BeforeEach(func() {

		mockClock := clock.NewMock()
		if t, err := time.Parse(time.RFC3339, fixedTime); err != nil {
			panic(err)
		} else {
			mockClock.Set(t)
		}
		ioc.TestPut(mockClock, func(clock.Clock) {})

		ioc.TestPut(&BytesAppender{}, func(Appender) {})

	})

	It("should display context.", func() {

		ioc.CallInjected(func(logger Logger, appender *BytesAppender) {

			logger.ErrorLog("what", "something")

			out := appender.String()

			Expect(out).To(ContainSubstring("\"level\":\"ERROR\""))
			Expect(out).To(ContainSubstring("\"time\":\"%s\"", fixedTime))

		})

	})

	It("should log only the correct levels.", func() {

		ioc.CallInjected(func(logger Logger, appender *BytesAppender) {

			logger.TraceLog("what", "something very small")
			logger.DebugLog("what", "something small")
			logger.InfoLog("what", "something current")
			logger.WarnLog("what", "something dangerous")
			logger.ErrorLog("what", "something problematic")
			logger.FatalLog("what", "something deadly")

			out := appender.String()

			Expect(out).ToNot(ContainSubstring("something very small"))
			Expect(out).ToNot(ContainSubstring("something small"))
			Expect(out).To(ContainSubstring("something current"))
			Expect(out).To(ContainSubstring("something dangerous"))
			Expect(out).To(ContainSubstring("something problematic"))
			Expect(out).To(ContainSubstring("something deadly"))

		})

	})

	It("should set the correct level.", func() {

		config.Set("log.level.first.second", "ERROR")

		type FirstLogger Logger
		type SecondLogger Logger
		type ThirdLogger Logger

		ioc.TestPutFactory(func(loggerFactory LoggerFactory) (FirstLogger, error) {
			return loggerFactory.NewLogger("first"), nil
		})
		ioc.TestPutFactory(func(loggerFactory LoggerFactory) (SecondLogger, error) {
			return loggerFactory.NewLogger("first.second"), nil
		})
		ioc.TestPutFactory(func(loggerFactory LoggerFactory) (ThirdLogger, error) {
			return loggerFactory.NewLogger("first.second.third"), nil
		})

		ioc.CallInjected(func(
			root Logger,
			first FirstLogger,
			second SecondLogger,
			third ThirdLogger,
			appender *BytesAppender) {

			Expect(root.GetLevel()).To(Equal(Info))
			Expect(first.GetLevel()).To(Equal(Info))
			Expect(second.GetLevel()).To(Equal(Error))
			Expect(third.GetLevel()).To(Equal(Error))

		})

	})

})
