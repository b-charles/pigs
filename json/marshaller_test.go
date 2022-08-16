package json_test

import (
	"github.com/b-charles/pigs/ioc"
	. "github.com/b-charles/pigs/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testMarshall(json Json, value any, expected string) {
	node, err := json.Marshal(value)
	Expect(err).WithOffset(1).To(Succeed())
	Expect(node.String()).WithOffset(1).To(Equal(expected))
}

var _ = Describe("Json marshallers", func() {

	BeforeEach(func() {
		ioc.TestPut("ioc test flag")
	})

	It("should marshall simple values", func() {

		ioc.CallInjected(func(json Json) {

			testMarshall(json, "Praise You", `"Praise You"`)

			testMarshall(json, 414.15, `414.15`)
			testMarshall(json, 42, `42`)

			testMarshall(json, true, `true`)
			testMarshall(json, false, `false`)

		})

	})

	It("should marshall simple struct", func() {

		type mySub struct {
			MyValue int `json:"value"`
		}

		type myStruct struct {
			MyString string
			MySub1   mySub  `json:"sub1"`
			MySub2   *mySub `json:"sub2"`
		}

		ioc.CallInjected(func(json Json) {
			testMarshall(json, myStruct{"Road Trippin'", mySub{42}, &mySub{21}}, `{"MyString":"Road Trippin'","sub1":{"value":42},"sub2":{"value":21}}`)
		})

	})

	It("should marshall recursive struct", func() {

		type recStruct struct {
			Value string     `json:"v"`
			Sub   *recStruct `json:"s"`
		}

		ioc.CallInjected(func(json Json) {
			testMarshall(json,
				recStruct{"Fatboy Slim", &recStruct{"You've Come a Long Way Baby", &recStruct{"Praise You", nil}}},
				`{"v":"Fatboy Slim","s":{"v":"You've Come a Long Way Baby","s":{"v":"Praise You","s":null}}}`)
		})

	})

	It("should marshall slices and maps", func() {

		type completeStruct struct {
			Slice []string
			Map   map[string]int
		}

		ioc.CallInjected(func(json Json) {
			testMarshall(json,
				completeStruct{[]string{"Wild Cherry", "Play That Funky Music"}, map[string]int{"Daft Punk": 2}},
				`{"Slice":["Wild Cherry","Play That Funky Music"],"Map":{"Daft Punk":2}}`)
		})

	})

})
