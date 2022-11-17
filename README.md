# ep-proxy-server

A forward proxy server that allows clients to use random and temporary proxies from [Ephemeral Proxies API service](https://www.ephemeral-proxies.net/)

:warning: ep-proxy-server requires a valid Rapid API key, this key can be obtained by signing up on https://rapidapi.com/.

## Quick Start

1. [Download](https://github.com/rpi4gx/ep-proxy-server/releases) and unpack the latest release for your system.

2. Run it
```
$ ./ep-proxy-server -mode=rotation -port=9090 -rapidApiKey=<USER'S RAPIDAPI KEY> -type=residential
````

4. Test it

Configure your browser to use a proxy on address http://localhost:9090 or you can test from command line:
```
$ curl -x localhost:9090 https://ifconfig.co
```

### Supported options
```
$ ./ep-proxy-server -h
Usage of ./ep-proxy-server:
  -mode="static": Mode <static|rotational>.
  -port=8090: Listening port
  -rapidApiKey="": Rapid API key.
  -size=50: Pool size for rotational mode.
```

## Compiling
### Requirements
* [Golang](https://go.dev/doc/install)

### Steps
1. Clone the repository
```
$ git clone https://github.com/rpi4gx/ep-proxy-server.git
```
2. Compile it
```
$ cd ep-proxy-server
$ go build .
```
3. Run it
```
$ ./ep-proxy-server -mode=rotation -port=9090 -rapidApiKey=<USER'S RAPIDAPI KEY> -type=datacenter
```
4. Test it
```
$ curl -x localhost:9009 https://ifconfig.co
```

## Support

Feel free to contact us on support@ephemeral-proxies.net.
