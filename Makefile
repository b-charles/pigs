
.PHONY=test watch tidy

test:
	ginkgo -r

watch:
	ginkgo watch -r -p

tidy:
	go mod tidy

