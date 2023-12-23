render:
    just build
    ENV=static ./tmp/main

dev: 
    ENV=development air

start:
    just build
    ENV=production ./tmp/main

build:
    templ generate
    go build -o ./tmp/main .

db:
    ortfodb ~/projects --scattered build to database.json --config ortfodb.yaml

clean: 
    rm -rf dist/
    rm -rf media/
    rm database.json
