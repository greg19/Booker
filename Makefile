README.pdf: README.md
	pandoc -V geometry:margin=1in -o README.pdf README.md

