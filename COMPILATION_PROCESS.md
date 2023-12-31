# Compilation process

## Compile

```
   +------------------------------+
   | ~/projects/*/.portfoliodb/** |
   +---------------+--------------+
                   |
                   |   pnpm database:crawl
                   v
+------------------+----------------+
|  ~/projects/portfolio/database/** |
+------------------+----------------+
                   |
                   |   pnpm database:build
                   v
     +-------------+----------+     +-----------+                             +----------+          +-----------+         +-----------+
     | database/database.json |     | src/*.pug |                             | src/*.ls |          | src/*.css |         | assets/** |
     +-------------+----------+     +-----+-----+                             +----+-----+          +-----+-----+         +-----+-----+
                   |                      |                                        |                      |                     |
                   +--------------------->+   pnpm pug:hydrate                     |                      |                     |
                                          v                                        |                      |                     |
                         +----------------+---------------+                        |                      |                     |
                         | artifacts/phase_1/{lang}/*.pug |                        |                      |                     |
                         +----------------+---------------+                        |                      |                     |
                                          |                                        |                      |                     |
                                          |                                        |                      |                     |
                                          v                                        |                      |                     |
                              +-----------+-----------------+                      |                      |                     |
      pnpm messages:extract   |                             |                      |                      |                     |
      pnpm messages:combine   |                             |   pnpm pug:build     |                      |                     |
                              v                             v                      |                      |                     |
                  +-----------+--------+  +-----------------+---------------+      |                      |                     |
                  | messages/{lang}.po |  | artifacts/phase_2/{lang}/*.html |      |                      |                     |
                  +-----------+--------+  +-----------------+---------------+      |                      |                     |
                              |                             |                      |                      |                     |
        pnpm messages:build   |                             |                      |                      |                     |
                              v                             |                      |                      |                     |
                  +-----------+--------+                    |                      |                      |                     |
                  | messages/{lang}.mo |                    |                      |                      |                     |
                  +-----------+--------+                    |                      |                      |                     |
                              |                             |                      |                      |                     |
                              |                             |                      |                      |                     |
                              +-----------+-----------------+                      |                      |                     |
                                          |                                        |                      |                     |
                                          |   pnpm html:translate                  |                      |                     |
                                          v                                        |   pnpm ls:build      |   pnpm stylus:build |   pnpm assets:build
                               +----------+---------+                              |                      |                     |
                               | dist/{lang}/*.html +<-----------------------------+----------------------+---------------------+
                               +--------------------+
```

## Macro-commands

<dl>
  <dt><code>database:update</code></dt>
  <dd><code>database:crawl</code> then <code>database:build</code></dd>
  <dt><code>prepare:i18n</code></dt>
  <dd>from <code>pug:hydrate</code> to <code>messages:combine</code></dd>
  <dt><code>make</code></dt>
  <dd><code>*:build</code> and <code>html:translate</code></dd>
</dl>

## Cleanup


## Workflow

1. Edit `{src,}/**` files
2. Update database if needed: `pnpm database:update` (i.e. if you changed some description.md files and/or added a new project)
3. Run `pnpm prepare:i18n`
4. Add missing translations to `messages/*.po`
5. Run `pnpm make`
6. (Optionally) Run `pnpm clean`

Or:

1. Edit source files
2. Run `./imake`
