VERSION ?= v1.21

all: kwebsite

clean:
	rm -rf kwebsite/content/en/docs/* kwebsite/public

kwebsite: clean
	mkdir -p kwebsite/content/en/docs
	go run cmd/main.go kwebsite --config-dir config/$(VERSION)/ --file api/$(VERSION)/swagger.json --output-dir kwebsite/content/en/docs --templates ./templates

copy_files: kwebsite
	cp -r kwebsite/content/en/docs/* ../docsy-example/content/en/docs/Reference/

# Create english PO files from swagger and sources
potfiles:
	mkdir -p po
	go run cmd/main.go po -f api/$(VERSION)/swagger.json -p po/
	# go install github.com/leonelquinteros/gotext/cli/xgotext
	xgotext -default gen-resourcesdocs -in pkg/outputs/kwebsite -out po

# Initialize translation files for a new language, using LANG env variable
initlang:
	scripts/initlang.sh
