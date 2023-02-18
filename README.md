# FastCGI-Serve - A webserver for FastCGI
cloning from https://github.com/beberlei/fastcgi-serve

## Usage

You can keep your Apache/Nginx setup with this proxy and start a webserver (defaults to `localhost:8080`)
for a given document root (defaults to current working directory).

    $ fastcgi-serve --document-root=/var/www --listen=127.0.0.1:8080
    Listening on http://localhost:8080
    Document root is /var/www
    Press Ctrl-C to quit.

## Configuration

The following settings are available:

- `--document-root` - The document root to serve files from (default: current working directory)
- `--listen` - The webserver bind address to listen to (default:127.0.0.1)
- `--server` - The FastCGI server to listen to
- `--server-port` The FastCGI server port to listen to
- `--index` The default script to call when request path cannot be served with an existing file

