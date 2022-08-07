package smartconf_test

import (
	"github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	. "github.com/b-charles/pigs/smartconf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func complexSetup() {
	config.SetTest(map[string]string{
		"co.rrup.tion":   "corruption",
		"con.fe.ren.ce":  "conference",
		"con.fe.ssion":   "confession",
		"con.fu.sion":    "confusion",
		"ga.la.xy":       "galaxy",
		"ga.lle.ry":      "gallery",
		"in.fec.tion":    "infection",
		"in.qui.ry":      "inquiry",
		"in.vest.ment":   "investment",
		"pa.ra.dox":      "paradox",
		"pa.ra.gra.ph":   "paragraph",
		"per.for.ate":    "perforate",
		"per.for.man.ce": "performance",
	})
}

var _ = Describe("IOC registration", func() {

	It("should parse empty config", func() {

		config.SetTest(map[string]string{})

		ioc.CallInjected(func(config NavConfig) {
			Expect(config.Keys()).To(BeEmpty())
		})

	})

	It("should parse complex config", func() {

		complexSetup()

		ioc.CallInjected(func(config NavConfig) {

			Expect(config.Keys()).To(Equal([]string{
				"co", "con", "ga", "in", "pa", "per",
			}))

			con := config.Child("con")
			Expect(con.Keys()).To(Equal([]string{"fe", "fu"}))

			confu := con.Child("fu")
			Expect(confu.Keys()).To(Equal([]string{"sion"}))

			confusion := confu.Child("sion")
			Expect(confusion.Keys()).To(BeEmpty())

			Expect(confusion.Value()).To(Equal("confusion"))

		})

	})

	It("should return empty value for undefined key", func() {

		complexSetup()

		ioc.CallInjected(func(config NavConfig) {

			unk := config.Child("unk")
			Expect(unk.Keys()).To(Equal([]string{}))

			nown := unk.Child("nown")
			Expect(nown.Keys()).To(Equal([]string{}))

		})

	})

})
