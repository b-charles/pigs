package smartconfig_test

import (
	"github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	. "github.com/b-charles/pigs/smartconfig"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func complexSetup() {
	config.TestMap(map[string]string{
		"dic.co.rrup.tion":   "corruption",
		"dic.con.fe.ren.ce":  "conference",
		"dic.con.fe.ssion":   "confession",
		"dic.con.fu.sion":    "confusion",
		"dic.ga.la.xy":       "galaxy",
		"dic.ga.lle.ry":      "gallery",
		"dic.in.fec.tion":    "infection",
		"dic.in.qui.ry":      "inquiry",
		"dic.in.vest.ment":   "investment",
		"dic.pa.ra.dox":      "paradox",
		"dic.pa.ra.gra.ph":   "paragraph",
		"dic.per.for.ate":    "perforate",
		"dic.per.for.man.ce": "performance",
	})
}

var _ = Describe("IOC registration", func() {

	It("should parse complex config", func() {

		complexSetup()

		ioc.CallInjected(func(config NavConfig) {

			dic := config.Child("dic")

			Expect(dic.Keys()).To(Equal([]string{
				"co", "con", "ga", "in", "pa", "per",
			}))

			con := dic.Child("con")
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
