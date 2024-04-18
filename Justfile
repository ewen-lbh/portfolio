render:
    just build
    REMOVE_UNUSED_MESSAGES=1 ENV=static ./tmp/main

readme:
    echo > README.md
    echo "# ewen.works" >> README.md
    echo "My personal website, built with [ortfo/db](https://ortfo.org/db) and [go-templ](https://templ.guide)" >> README.md
    
    echo "## Setup" >> README.md
    echo "Run this before any command" >> README.md
    
    echo "\`\`\`sh" >> README.md
    echo "cp .env.example .env" >> README.md
    echo "# fill values in .env" >> README.md
    echo "source .env" >> README.md
    echo "\`\`\`" >> README.md
    echo  >> README.md
    echo "Below are how to do various things like start/build/etc" >> README.md
    echo "([generated](./readme_from_justfile.rb) from [the Justfile](./Justfile))" >> README.md
    
    ./readme_from_justfile.rb including dev start db build clean deploy >> README.md

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
