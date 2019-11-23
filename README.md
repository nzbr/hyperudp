# <div align=center>HyperUDP</div>

<p align=center>
	<img alt="GitHub" src="https://img.shields.io/github/license/nzbr/hyperudp">
	<a href="https://actions-badge.atrox.dev/nzbr/hyperudp/goto"><img alt="Build Status" src="https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fnzbr%2Fhyperudp%2Fbadge&style=flat" /></a>
    <img alt="Libraries.io dependency status for GitHub repo" src="https://img.shields.io/librariesio/github/nzbr/hyperudp">
</p>
A program that converts raw UDP streams to [Hyperion](https://github.com/hyperion-project/hyperion.ng)'s Protocol Buffer format

Hyperion.ng currently has no way to receive a raw UDP stream itself.  
For example this can be used to use [ColorChord](https://github.com/cnlohr/colorchord) with Hyperion.ng

# Building

Just run `go build .`  

# Usage

Basically it's just `./hyperudp <Hyperion IP>`  
This, by default, opens a UDP socket on port `1337`. Any color data sent there should appear as a vide feed in Hyperion and thus on your connected RGB devices

For more advanced settings please see `./hyperudp --help`