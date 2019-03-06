default: build

build:
	hugo

install:
	cd themes/ferret && npm i && cd ../..

deploy:
	cd public && git add --all && git commit -m "Publishing to gh-pages" && cd .. && git push origin gh-pages