package memfun

import (
	"sync"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mem Fun Test Suite")
}

var _ = Describe("Mem fun", func() {

	It("should work in single goroutine", func() {

		call := make(map[int]int)

		fibo := NewMemFun(func(n int, rec func(n int) (int, error)) (int, error) {

			call[n] += 1

			if n == 0 || n == 1 {
				return n, nil
			}
			if n1, err1 := rec(n - 1); err1 != nil {
				return 0, err1
			} else if n2, err2 := rec(n - 2); err2 != nil {
				return 0, err2
			} else {
				return n1 + n2, nil
			}

		})

		Expect(fibo.Get(15)).To(Equal(610))

		Expect(len(call)).To(Equal(16))
		for k, v := range call {
			Expect(v).To(Equal(1), "Calls for %d", k)
		}

	})

	It("should not block on different key", func() {

		nb := 3

		wg := sync.WaitGroup{}
		wg.Add(nb)

		fun := NewMemFun(func(n int, rec func(n int) (int, error)) (int, error) {

			wg.Done()
			wg.Wait()

			return 4, nil

		})

		c := make(chan int)

		for i := 0; i < nb; i++ {
			go func(i int) {
				if v, err := fun.Get(i); err != nil {
					panic(err)
				} else {
					c <- v
				}
			}(i)
		}

		for i := 0; i < nb; i++ {
			Expect(<-c).To(Equal(4))
		}

		close(c)

	})

	It("should handle panics", func() {

		call := make(map[int]int)

		fun := NewMemFun(func(n int, rec func(n int) (int, error)) (int, error) {

			call[n] += 1

			if n == 4 {
				panic("I'm panicking.")
			} else {
				return n, nil
			}
		})

		Expect(func() { fun.Get(4) }).Should(Panic())
		Expect(func() { fun.Get(5) }).ShouldNot(Panic())
		Expect(func() { fun.Get(4) }).Should(Panic())

		Expect(len(call)).To(Equal(2))
		for k, v := range call {
			Expect(v).To(Equal(1), "Calls for %d", k)
		}

	})

	It("should throw cycle errors", func() {

		fun := NewMemFun(func(v int, rec func(v int) (int, error)) (int, error) {
			return rec((v + 1) % 10)
		})

		_, err := fun.Get(0)

		exp := CyclicLoopError[int]{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}}
		Expect(err).Should(MatchError(exp))

	})

	It("should allow cyclic lazy loading", func() {

		call := make(map[int]int)

		fun := NewMemFun(func(v int, rec func(v int) (func(int) int, error)) (func(int) int, error) {

			call[v] += 1

			return func(s int) int {
				if s == 0 {
					return v
				} else {
					if n, e := rec((v + 1) % 10); e != nil {
						panic(e)
					} else {
						return v + n(s-1)
					}
				}
			}, nil

		})

		sumfun, err := fun.Get(7)
		Expect(err).Should(Succeed())

		// sumfun(s) = sum_{v=0}^{s}{v+7 (mod 10)}
		Expect(sumfun(254)).Should(Equal(1150))

		Expect(len(call)).To(Equal(10))
		for k, v := range call {
			Expect(v).To(Equal(1), "Calls for %d", k)
		}

	})

})
