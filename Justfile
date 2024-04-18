render:
    just build
    REMOVE_UNUSED_MESSAGES=1 ENV=static ./tmp/main

dev: 
    ENV=development air

start:
    just build
    ENV=production ./tmp/main

build:
    templ generate
    go build -o ./tmp/main .

db:
    ortfodb --scattered build database.json

clean: 
    rm -f */*_templ.go
    rm -rf dist/
    rm -f database.json

clean-media:
    rm -rf media/

deploy ssh='$SSH_SERVER':
    rsync -av media/* {{ ssh }}:~/www/media.ewen.works/
    rsync -avz public/* {{ ssh }}:~/www/assets.ewen.works/
    rsync -av database.json {{ ssh }}:~/portfolio/
    ssh {{ ssh }} "tmux send-keys -t 0:0.0 C-c 'git pull --autostash --rebase' Enter 'just start' Enter"
