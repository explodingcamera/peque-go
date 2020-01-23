## PEQUE-GO
# The simple & universal message/task/job queue

**NOTE: Peque is currently still being worked on and not production ready!**

Peque (pronounced `[ˈpeke]`) is the easiest way to get a message/task/job queue running.

You don't even have to host a dedicated program for you queue, Peque can start out using your existing PostgreSQL Database for storing and distributing messages, tasks and jobs.

With a unified API message format, PEQUE can be even used across multiple different programing languages.

## Message Format

To keep everything human-readable, messages are transmitted as JSON objects by default, with support for msgpack and brotli compression in performance critical applications.

## Current Roadmap

### Work in progress
* [ ] go Client

### To be determined
* [ ] node.js Client
* [ ] clients for other languages
* [ ] messagepkg support
* [ ] brotli compression support

## Inspiration

I have been heavily inspired by [machinery](https://github.com/RichardKnop/machinery), [celery](https://github.com/celery/celery) and [eventide](https://eventide-project.org/). This projects goal is to create a unified solution that can be used across many different programing languages.
