# myhttp

`myhttp` is a CLI tool that performs HTTP GET requests for urls provided and prints the MD5 Hash for their response body.

## Prerequisites
### Install Go

[Here](https://go.dev/doc/install) you will find instructions on how to setup up the `Go Compiler` locally

## Building the Application

To build the CLI app, you'll run:
```bash
go build -o myhttp main.go
```
Or if you have GNU Make locally:
```bash
make build
```

## Testing

There are unit tests written for the application. To run the tests, run:
```bash
go test -v
```
Or
```bash
make test
```

## Using The Application

Below are example snippets on how to use the binary and their sample output. First make sure you [build](#building-the-application) the binary.

The `-parallel` flag is used to specify the limit of the number of parallel requests that will be scheduled.

```bash
./myhttp http://www.adjust.com http://google.com
http://www.adjust.com d1b40e2a2ba488a054186e4ed0733f9752f66949
http://google.com 9d8ec921bdd275fb2a605176582e08758eb60641
```

```bash
./myhttp adjust.com
http://adjust.com d1b40e2a2ba488a054186e4ed0733f9752f66949
```

```bash
./myhttp -parallel 3 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com
reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com
http://google.com 8ff1c478ccca08cca025b028f68b352f
http://adjust.com 6b2560b9a5262571258cc173248b7492
http://yandex.com 4baab01ff9ff0f793bf423aeef539c9d
http://facebook.com ccae5ffa91c4936aef3efd5091a43f3e
http://twitter.com 857efe81a54c8b5c2241846ac4f08e66
http://reddit.com/r/funny ff3b2b7dcd9e716ca0adcbd208061c9a
http://reddit.com/r/notfunny ff3b2b7dcd9e716ca0adcbd208061c9a
http://yahoo.com e2d50a30b7bfbfda097d72e32578c6a6
http://baroquemusiclibrary.com 8e5138a0111364f08b10d37ed3371b11
```