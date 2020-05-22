default: serve

serve:
	rm -rf public/* && hugo server

generate-stdlib-docs:
	sh generate-docs.sh stdlib-docs-rep.yaml

build:
	rm -rf public/* && hugo

install:
	cd themes/ferret && npm i && cd ../..

publish:
	sh publish.sh