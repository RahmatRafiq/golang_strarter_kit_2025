build:
	go build -o main .

test:
	ginkgo -r --race