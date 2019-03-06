default: build

build:
	hugo

install:
	cd themes/ferret && npm i && cd ../..

deploy:
	git subtree push --prefix public origin gh-pages