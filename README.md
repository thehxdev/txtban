# Txtban
Txtban is a server-side application for txt content sharing written in golang. It does not include a front-end.


## Build
To bulid Txtban, you need go compiler, `make` and `sqlite3` installed (sqlite is for running txtban).

```bash
make
```
This will build `tb` executable.


### Docker image
You can build txtban docker image with the following command:
```bash
make docker
```

### Use Docker to build `tb`
If you dont want to install go compiler, you can build the `tb` executable inside golang docker image:
```bash
make docker_exe
```
Then you can see `tb` executable in the project root directory.


## Configuration
Txtban uses [TOML](https://toml.io/en/) file format for it's config file. When you run `tb`,
it will search current working directory and `/etc/txtban` directory to find `config.toml` file.
If you want to specify the `config.toml` file path, use `-c` command-line flag.

See [config.toml](config.toml) file as an example.


## API Endpoints
Here is the full API documentation for txtban (They are like unix commands):

### TODO!
