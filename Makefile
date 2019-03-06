default: build

build:
	rm -rf public && hugo

install:
	cd themes/ferret && npm i && cd ../..

publish:
	sh publish.sh