#!/bin/bash
ragel -Z -T0 -o stomp.go stomp.rl && go build -o stomp stomp.go main.go command.go && ./stomp
