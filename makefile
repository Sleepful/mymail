# makefile per https://templ.guide/developer-tools/live-reload-with-other-tools

# basic way of generating templ pages and running the dev server:
# NOTE:
#		For a more complete command, look at `live`
gen:
	templ generate --watch --proxy="http://localhost:8000" --open-browser=false --cmd="go run ."

# Default url: http://localhost:7331
live/templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o tmp/bin/main && sleep 4 && tput bel" --build.bin "tmp/bin/main serve" --build.delay "100" \
	--build.exclude_dir "pulumi,node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# run tailwindcss to generate the styles.css bundle in watch mode.
live/tailwind:
	npx --yes tailwindcss -i ./app/input.css -o ./assets/styles.css --minify --watch

# run esbuild to generate the index.js bundle in watch mode.
live/esbuild:
	npx --yes esbuild js/index.ts --bundle --outdir=assets/ --watch

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "pulumi,node_modules" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

# start all 5 watch processes in parallel.
live: 
	make -j5 live/templ live/server live/tailwind live/esbuild live/sync_assets

# Docker
dockerrun:
	docker run -p 8090:8090 "imap"

dockerbuild:
	docker build  . --tag "imap"  --progress=plain --no-cache

# env
envsample:
	sed 's/=.*/=/' .env > env.example
