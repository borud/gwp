
   
all: lint gen

lint:
	@buf lint

gen: clean
	@buf generate

clean:
	@rm -rf pkg

breaking:
	@buf breaking --against "https://github.com/borud/chat/archive/main.zip#strip_components=1"


publish:
	@buf push

dep:
	@go install github.com/bufbuild/buf/cmd/buf@latest