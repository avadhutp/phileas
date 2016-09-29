deps:
	@echo Getting dependencies for Phileas
	@go get github.com/mattn/gom
	@gom install

test:
	@echo Testing Phileas
	@(go list ./... | grep -v -e /vendor/ | xargs -L1 gom test -cover || exit;)

lint:
	@echo Linting Phileas sources
	@(go list ./... | grep -v -e /vendor/ | xargs -L1 golint || exit;)