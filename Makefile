SRC := $(wildcard app/*.go)
BIN := site
BUILDDIR := .

.PHONY: all clean test

all: $(BUILDDIR)/$(BIN)

test:
	cd app/; go test ./...

$(BUILDDIR)/$(BIN):$(SRC)
	cd app/; go build -o $(BIN)
	mv app/$(BIN) $(BUILDDIR)/$(BIN)

clean:
	rm -f $(BUILDDIR)/$(BIN)
