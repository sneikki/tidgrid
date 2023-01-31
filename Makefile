compile: src/tidgrid.scss
	npx sass src/tidgrid.scss build/tidgrid.css

compile-minimized:
	make compile -- --style=compressed

.PHONY: lint-scss
lint-scss:
	npx stylelint "src/**/*.scss"

.PHONY: fix-scss
fix-scss:
	make lint-scss -- --fix

.PHONY: clean
clean:
	rm -rf build/
