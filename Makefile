GOCMD   = go
GOBUILD = $(GOCMD) build
GORUN   = $(GOCMD) run
GOTEST  = $(GOCMD) test
GOCLEAN = $(GOCMD) clean
BINARY  = buddy
RM      = rm

all: run

build:
	CGO_ENABLED=0 $(GOBUILD) -gcflags -m -o $(BINARY) -v .

run:
	CGO_ENABLED=0 $(GORUN) -gcflags -m -v . --port 3000

test:
	CGO_ENABLED=0 $(GOTEST) -v .

clean:
	$(GOCLEAN)
	$(RM) -f $(BINARY)
