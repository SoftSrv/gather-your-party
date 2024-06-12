run:
	templ generate
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css
	go run main.go

templ:
	templ generate

tw:
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css
