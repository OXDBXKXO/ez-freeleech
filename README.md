# ez-freeleech
A simple private tracker ratio cheater written in Golang.

`ez-freeleech` acts as a HTTP proxy between your BitTorrent client and private trackers to tamper download reports. Every announce requests made to trackers will report 0 downloaded bytes. Based on [libtorrent tracker protocol documentation](https://libtorrent.org/udp_tracker_protocol.html#announcing).

## Installation

```bash
git clone https://github.com/OXDBXKXO/ez-freeleech
cd ez-freeleech
go get
```

## Usage

```bash
go run ez-freeleech
```

After starting `ez-freeleech`, set up your BitTorrent client to use 127.0.0.1:8888 (or any other port you choose) as HTTP proxy.

## Options

`-port`: By default, ez-freeleech will run on port 8888. Alternatively, you can specify a port using option `port`.

`-locale`: As a proud member of the octet gang, I could not bring myself to provide only byte representation. You can use `-locale fr` to use octet representation.
