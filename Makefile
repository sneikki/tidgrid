style = expanded

build: src/index.scss
	make compile

compile: src/index.scss
	npx sass src/index.scss build/tidgrid.css --style=$(style) --load-path=src

compile-minimized: src/index.scss
	make -s compile style=compressed

.PHONY: test
test:
	npx jest

.PHONY: lint-scss
lint-scss:
	npx stylelint "{src,test}/**/*.scss"

.PHONY: fix-scss
fix-scss:
	npx stylelint "{src,test}/**/*.scss" --fix

.PHONY: clean
clean:
	rm -rf build/
