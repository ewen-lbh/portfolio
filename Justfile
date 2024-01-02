render:
    just build
    WAKATIME_API_KEY=$(rbw get 'wakatime api key') REMOVE_UNUSED_MESSAGES=1 ENV=static ./tmp/main

dev: 
    WAKATIME_API_KEY=$(rbw get 'wakatime api key') MAIL_PASSWORD=$(rbw get mail.ewen.works) ENV=development air

start:
    just build
    ENV=production ./tmp/main

build:
    templ generate
    go build -o ./tmp/main .

db:
    ortfodb ~/projects --scattered build to database.json --config ortfodb.yaml

clean: 
    rm -f */*_templ.go
    rm -rf dist/
    rm -f database.json

clean-media:
    rm -rf media/

deploy:
    rsync -av media/* ewen@ewen.works:~/www/media.ewen.works/
    rsync -avz public/* ewen@ewen.works:~/www/assets.ewen.works/
    ssh ewen@ewen.works "tmux send-keys -t 0:0.0 C-c 'git pull --autostash --rebase' Enter 'just start' Enter"
