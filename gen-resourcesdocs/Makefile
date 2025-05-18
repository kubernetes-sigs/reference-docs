VERSION ?= v1.33

all: kwebsite

clean:
	rm -rf kwebsite/content/en/docs/* kwebsite/public

kwebsite: clean
	mkdir -p kwebsite/content/en/docs
	go run cmd/main.go kwebsite --config-dir config/$(VERSION)/ --file api/$(VERSION)/swagger.json --output-dir kwebsite/content/en/docs --templates ./templates


copy_files: kwebsite
	cp -r kwebsite/content/en/docs/* ../docsy-example/content/en/docs/Reference/
