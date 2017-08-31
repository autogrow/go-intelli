# go-intelli

A NATS and REST gateway for the Intelli range of devices.  This allows hackers and tinkerers to do cool stuff
with IntelliDose or IntelliClimate devices connected via USB.  Update events can be subscribed to via NATS or
read via the HTTP API.

Current support is limited to Linux and [exprimentally] MAC.

This can be coupled with the [Jelly SDK](https://github.com/AutogrowSystems/go-jelly) to get programmatical access
to an IntelliDose.

## Installation

Install using go:

    go get github.com/AutogrowSystems/go-intelli

Build using go:

    go build github.com/AutogrowSystems/go-intelli/cmd/intellid

Optionally install and run a [NATS](https://github.com/nats-io/gnatsd/releases) server.

## Usage

To use it, simply run the binary (as sudo to get access to the USB files):

    sudo ./intellid

If you have a NATS server running you will see JSON being published to the subject `intelli.*` or `intelli.ASLID06030112` every 15 seconds.  The JSON is formatted like so: [example.json](https://github.com/AutogrowSystems/go-intelli/blob/master/example.json)

You also have two endpoints available, `/devices/count` and `/devices`.  By calling the latter you will see output like in the file
[example.json](https://github.com/AutogrowSystems/go-intelli/blob/master/example.json)

## TODO

* [ ] add ability to change settings
