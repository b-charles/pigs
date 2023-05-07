package json_test

import (
	"github.com/b-charles/pigs/ioc"
	. "github.com/b-charles/pigs/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func testMarshall(jsons Jsons, value any, expected string) {
	node, err := jsons.Marshal(value)
	Expect(err).WithOffset(1).To(Succeed())
	Expect(node.String()).WithOffset(1).To(Equal(expected))
}

type doer interface {
	Do()
}

type doer_impl struct{}

func (self *doer_impl) Do() {}

type dodoer interface {
	Do()
	DoDo()
}

type dodoInh interface {
	doer
	DoDo()
}

type dodoer_impl struct{}

func (self *dodoer_impl) Do()   {}
func (self *dodoer_impl) DoDo() {}

type dododoer interface {
	Do()
	DoDo()
	DoDoDo()
}

type dododoer_impl struct{}

func (self *dododoer_impl) Do()     {}
func (self *dododoer_impl) DoDo()   {}
func (self *dododoer_impl) DoDoDo() {}

type myJsoner struct{}

func (self *myJsoner) Json() JsonNode {
	return JsonString("ANGRY JSONER")
}

var _ = Describe("Json marshallers", func() {

	BeforeEach(func() {
		ioc.TestPut("ioc test flag")
	})

	It("should marshall simple values", func() {

		ioc.CallInjected(func(jsons Jsons) {

			testMarshall(jsons, "Praise You", `"Praise You"`)

			testMarshall(jsons, 414.15, `414.15`)
			testMarshall(jsons, 42, `42`)

			testMarshall(jsons, true, `true`)
			testMarshall(jsons, false, `false`)

		})

	})

	It("should marshall Jsoner interfaces", func() {

		ioc.CallInjected(func(json Jsons) {

			testMarshall(json, &myJsoner{}, `"ANGRY JSONER"`)

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

		ioc.CallInjected(func(jsons Jsons) {
			testMarshall(jsons, myStruct{"Road Trippin'", mySub{42}, &mySub{21}}, `{"MyString":"Road Trippin'","sub1":{"value":42},"sub2":{"value":21}}`)
		})

	})

	It("should marshall recursive struct", func() {

		type recStruct struct {
			Value string     `json:"v"`
			Sub   *recStruct `json:"s"`
		}

		ioc.CallInjected(func(jsons Jsons) {
			testMarshall(jsons,
				recStruct{"Fatboy Slim", &recStruct{"You've Come a Long Way Baby", &recStruct{"Praise You", nil}}},
				`{"v":"Fatboy Slim","s":{"v":"You've Come a Long Way Baby","s":{"v":"Praise You","s":null}}}`)
		})

	})

	It("should marshall slices and maps", func() {

		type completeStruct struct {
			Slice []string
			Map   map[string]int
		}

		ioc.CallInjected(func(jsons Jsons) {
			testMarshall(jsons,
				completeStruct{[]string{"Wild Cherry", "Play That Funky Music"}, map[string]int{"Daft Punk": 2}},
				`{"Slice":["Wild Cherry","Play That Funky Music"],"Map":{"Daft Punk":2}}`)
		})

	})

	It("should marshall interfaces", func() {

		ioc.TestPut(func(v doer) (JsonNode, error) {
			return JsonString("Do"), nil
		}, func(JsonMarshaller) {})

		ioc.CallInjected(func(jsons Jsons) {
			testMarshall(jsons, &doer_impl{}, `"Do"`)
		})

	})

	It("should choose wisely the correct marshaller in case of inheritance", func() {

		ioc.TestPut(func(v doer) (JsonNode, error) {
			return JsonString("Do"), nil
		}, func(JsonMarshaller) {})

		ioc.TestPut(func(v dodoInh) (JsonNode, error) {
			return JsonString("DoDo"), nil
		}, func(JsonMarshaller) {})

		ioc.CallInjected(func(jsons Jsons) {
			testMarshall(jsons, &doer_impl{}, `"Do"`)
			testMarshall(jsons, &dodoer_impl{}, `"DoDo"`)
		})

	})

	It("should choose wisely the correct marshaller in case of interface overlap", func() {

		ioc.TestPut(func(v dodoer) (JsonNode, error) {
			return JsonString("DoDo"), nil
		}, func(JsonMarshaller) {})

		ioc.TestPut(func(v dododoer) (JsonNode, error) {
			return JsonString("DoDoDo"), nil
		}, func(JsonMarshaller) {})

		ioc.TestPut(func(v doer) (JsonNode, error) {
			return JsonString("Do"), nil
		}, func(JsonMarshaller) {})

		ioc.CallInjected(func(jsons Jsons) {
			testMarshall(jsons, &doer_impl{}, `"Do"`)
			testMarshall(jsons, &dodoer_impl{}, `"DoDo"`)
			testMarshall(jsons, &dododoer_impl{}, `"DoDoDo"`)
		})

	})

})
