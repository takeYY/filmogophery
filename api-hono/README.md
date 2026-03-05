## How to use it

To install dependencies:

```sh
bun install
```

To run:

```sh
bun run dev
```

open http://localhost:3000

### How to update drizzle schemas

```sh
# Run db container
make up TARGET_COMPOSE=compose.hono.yml -d

# Update drizzle schemas
npx drizzle-kit pull
```
