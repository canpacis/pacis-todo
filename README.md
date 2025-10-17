# Pacis Template

## Development

Install dependencies.

```sh
bun install
```

```sh
go mod install
```

> This project uses [air](https://github.com/air-verse/air) and [taskfile](https://taskfile.dev/) for development & build processes.

See installation guides:

- [Air](https://github.com/air-verse/air?tab=readme-ov-file#installation)
- [Taskfile](https://taskfile.dev/docs/installation)

Start the development server.

```sh
task dev
```

This will run a vite server for your frontend assets at port `:5173` and an app server for your webpages at port `:8081`.

> **Hot Reloading**: The `task dev` command will also proxy your server to the `:8080` port with air where you will get auto reload upon file changes.

Visit [http://localhost:8080](http://localhost:8080) and see your website.

While in development, the app will run on port `:8081` with a `:8080` proxy and in production it will default to `:8080`. You can change these values in `dev.go` & `prod.go` files.

## Building

### Building Manually

Build your server manually with taskfile.

```sh
task build
```

This will build your assets and your final executable. You can find it at `build/server`. It is a statically linked executable and will work on its own out of the box. The rest of the artefacts inside the build folder are irrelevant.

### Docker

The repository includes a `Dockerfile` to containerize your app. It uses bun and go images to build your server.

```sh
docker build -t pacis-template .
```