default: serve

serve:
	rm -rf public/* && hugo server

build:
	rm -rf public/* && hugo

install:
	cd themes/ferret && npm i && cd ../..

publish:
	sh publish.sh