BINARY := s3undelete

install: build test
	cp $(BINARY) $(GOPATH)/bin

.PHONY: lint
lint:
	golint cmd/$(BINARY)/*.go
	golint pkg/$(BINARY)/*.go

build: lint
	go build -o $(BINARY) cmd/$(BINARY)/*.go

.PHONY: test
.ONESHELL:
test:
	cd test
	terraform init -input=false
	terraform apply -input=false -auto-approve
	./runtest.sh
	terraform destroy -force

.ONESHELL:
clean:
	rm -f $(BINARY)
	cd test
	terraform destroy -force
