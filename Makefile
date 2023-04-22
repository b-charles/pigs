
.PHONY=test watch tidy

test:
	ginkgo -r

watch:
	ginkgo watch -r -depth 25

tidy:
	go mod tidy

