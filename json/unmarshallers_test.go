package json_test

import (
	"github.com/b-charles/pigs/ioc"
	. "github.com/b-charles/pigs/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testUnmarshall[T any](jsons Jsons, source string, expected T) {
	err := jsons.UnmarshalFromString(source, func(v T) {
		Expect(v).WithOffset(1).To(Equal(expected))
	})
	Expect(err).WithOffset(1).To(Succeed())
}

var _ = Describe("Json unmarshallers", func() {

	BeforeEach(func() {
		ioc.TestPut("ioc test flag")
	})

	It("should unmarshall simple values", func() {

		ioc.CallInjected(func(jsons Jsons) {

			testUnmarshall(jsons, `"Praise You"`, "Praise You")

			testUnmarshall(jsons, `414.15`, 414.15)
			testUnmarshall(jsons, `42`, 42)

			testUnmarshall(jsons, `true`, true)
			testUnmarshall(jsons, `false`, false)

		})

	})

	It("should unmarshall simple struct", func() {

		type mySub struct {
			MyValue int `json:"value"`
		}

		type myStruct struct {
			MyString string
			MySub1   mySub  `json:"sub1"`
			MySub2   *mySub `json:"sub2"`
		}

		ioc.CallInjected(func(jsons Jsons) {
			testUnmarshall(jsons,
				`{"MyString":"Road Trippin'","sub1":{"value":42},"sub2":{"value":21}}`,
				myStruct{"Road Trippin'", mySub{42}, &mySub{21}})
		})

	})

	It("should unmarshall recursive struct", func() {

		type recStruct struct {
			Value string     `json:"v"`
			Sub   *recStruct `json:"s"`
		}

		ioc.CallInjected(func(jsons Jsons) {
			testUnmarshall(jsons,
				`{"v":"Fatboy Slim","s":{"v":"You've Come a Long Way Baby","s":{"v":"Praise You","s":null}}}`,
				recStruct{"Fatboy Slim", &recStruct{"You've Come a Long Way Baby", &recStruct{"Praise You", nil}}})
		})

	})

	It("should marshall slices and maps", func() {

		type completeStruct struct {
			Slice []string
			Map   map[string]int
		}

		ioc.CallInjected(func(jsons Jsons) {
			testUnmarshall(jsons,
				`{"Slice":["Wild Cherry","Play That Funky Music"],"Map":{"Daft Punk":2}}`,
				completeStruct{[]string{"Wild Cherry", "Play That Funky Music"}, map[string]int{"Daft Punk": 2}})
		})

	})

})
