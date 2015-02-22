# mtorrent-go

Uber-hipster version of [mtorrent-node](https://github.com/martintrojer/mtorrent-node), node.js and ClojureScript is soooo 2014.

## Building

Make sure you have Docker (>=1.5.0) installed (and perhaps [boot2docker](http://boot2docker.io)).

### Building the binary

First build the development environment;

`$ ./build-dev.sh`

Then build the mtorrent binary;

```
$ ./connect-dev.sh
$ go build
```

### Build the distribution container

`$ ./build-dist.sh`

## Configuration

You might want to edit the `mtorrent.config` file before building the distribution container. That file contains settings and is hopefully self explanatory.

## Running

`$ docker run -v <DATA_FOLDER>:/data -p 1337:1337 --rm -t mtorrent-go`

## License

Copyright © 2015 Martin Trojer

Distributed under the Eclipse Public License.
