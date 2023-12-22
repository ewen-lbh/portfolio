dev: 
    ENV=development air

db:
    ortfodb ~/projects --scattered build to database.json --config ortfodb.yaml

clean: 
    rm -rf dist/
    rm -rf media/
    rm database.json
