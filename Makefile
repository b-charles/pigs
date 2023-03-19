
.PHONY=test watch

test:
	ginkgo -r

watch:
	ginkgo watch -r -p

