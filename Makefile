OUT = fawrwebservice

$(OUT):
	CGO_ENABLED=1 \
	go build -v -o $(OUT)

.PHONY: clean
clean:
	rm -f $(OUT)