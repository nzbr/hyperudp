# <div align=center>HyperUDP</div>

<p align=center>
	<img alt="GitHub" src="https://img.shields.io/github/license/nzbr/hyperudp?label=License">
	<a href="https://goreportcard.com/report/github.com/nzbr/hyperudp"><img src="https://goreportcard.com/badge/github.com/nzbr/hyperudp" /></a>
	<a href="https://actions-badge.atrox.dev/nzbr/hyperudp/goto"><img alt="Build Status" src="https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fnzbr%2Fhyperudp%2Fbadge&style=flat" /></a>
	<a href="https://libraries.io/github/nzbr/hyperudp"><img alt="Libraries.io dependency status for GitHub repo" src="https://img.shields.io/librariesio/github/nzbr/hyperudp?label=Dependencies"></a>
</p>
A program that converts raw UDP streams to <a href="https://github.com/hyperion-project/hyperion.ng">Hyperion's</a> Protocol Buffer format

Hyperion.ng currently has no way to receive a raw UDP stream itself.  
For example this can be used to use [ColorChord](https://github.com/cnlohr/colorchord) with Hyperion.ng

# Building

Just run `go build .`  

# Usage

Basically it's just `./hyperudp <Hyperion IP>`  
This, by default, opens a UDP socket on port `1337`. Any color data sent there should appear as a video feed in Hyperion and thus on your connected RGB devices

For more advanced settings please see `./hyperudp --help`
