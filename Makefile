style = expanded

compile: src/tidgrid.scss
	npx sass src/tidgrid.scss build/tidgrid.css --style=$(style) --load-path=src

compile-minimized: src/tidgrid.scss
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
