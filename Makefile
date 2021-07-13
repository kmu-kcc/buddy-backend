GOCMD   = go
GOBUILD = $(GOCMD) build
GORUN   = $(GOCMD) run
GOTEST  = $(GOCMD) test
GOCLEAN = $(GOCMD) clean
BINARY  = buddy
RM      = rm

all: run

build:
	$(GOBUILD) -gcflags -m -o $(BINARY) -v .

run:
	$(GORUN) -gcflags -m -v . --port 3000

test:
	$(GOTEST) -v .

clean:
	$(GOCLEAN)
	$(RM) -f $(BINARY)
