# Miniredis

**miniredis** is a minimal Redis implementation written in Go with ~500 lines of code only.

## Overview

This was written as a learning project for evening to understand better how Redis works.

It supports following commands:
- PING
- SET
- GET
- HSET
- HGET
- HGETALL

It uses AOF (append only file) data persistence method to store values

## Running

- Clone this repository
- Run ```make run``` (runs on port 6379)
- Connect using ```redis-cli```
