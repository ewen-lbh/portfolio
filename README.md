
# ewen.works
My personal website, built with [ortfo/db](https://ortfo.org/db) and [go-templ](https://templ.guide)
## Setup
Run this before any command
```sh
cp .env.example .env
# fill values in .env
source .env
```

Below are how to do various things like start/build/etc
([generated](./readme_from_justfile.rb) from [the Justfile](./Justfile))
## Dev

```sh
ENV=development air
```

## Start

```sh
# Run Build's commands
ENV=production ./tmp/main
```

## Build

```sh
templ generate
go build -o ./tmp/main .
```

## Db

```sh
ortfodb --scattered build database.json
```

## Clean

```sh
rm -f */*_templ.go
rm -rf dist/
rm -f database.json
```

## Deploy

```sh
rsync -av media/* YOUR_SSH:~/www/media.ewen.works/
rsync -avz public/* YOUR_SSH:~/www/assets.ewen.works/
rsync -av database.json YOUR_SSH:~/portfolio/
ssh YOUR_SSH "tmux send-keys -t 0:0.0 C-c 'git pull --autostash --rebase' Enter 'just start' Enter"
```
