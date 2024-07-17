set dotenv-load := true

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
    ortfodb --scattered build database-regular.json
    sed -i 's/scattered mode folder: .ortfo$/scattered mode folder: .ortfo.personal/' ortfodb.yaml
    ortfodb --scattered build database-personal-overrides.json
    sed -i 's/scattered mode folder: .ortfo.personal$/scattered mode folder: .ortfo/' ortfodb.yaml
    jq -s '.[0] * .[1]' database-regular.json database-personal-overrides.json > database.json
    just diff-with-online-db

diff-with-online-db:
    scp ewen@ewen.works:~/portfolio/database.json database.online.json
    jq . --sort-keys database.json > database.keys.json
    jq . --sort-keys database.online.json > database.online.keys.json
    difft database.online.keys.json database.keys.json

clean: 
    rm -f */*_templ.go
    rm -rf dist/
    rm -f database.json

clean-media:
    rm -rf media/

deploy ssh='$SSH_SERVER':
    git push
    just upload-media {{ ssh }}
    just upload-assets {{ ssh }}
    rsync -av database.json {{ ssh }}:~/portfolio/
    just git-pull {{ ssh }}
    ssh {{ ssh }} "tmux send-keys -t 0:0.0 C-c 'git pull --autostash --rebase' Enter 'just start' Enter"

upload-assets ssh:
    rsync -avz public/* {{ ssh }}:~/www/assets.ewen.works/

upload-media ssh:
    rsync -av media/* {{ ssh }}:~/www/media.ewen.works/

git-pull ssh:
    ssh {{ ssh }} "bash -c 'cd ~/portfolio && git stash && git pull && git stash apply'"
