# `memq`

[![build](https://github.com/sergeykonkin/memq/actions/workflows/build.yml/badge.svg)](https://github.com/sergeykonkin/memq/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sergeykonkin/memq.svg)](https://pkg.go.dev/github.com/sergeykonkin/memq)

A simple thread-safe in-memory message broker fo Go.

## Quick Start

1. Get it:

    ```bash
    go get github.com/sergeykonkin/memq
    ```

1. Use it:

    ```go
    // Import:
    import (
        "github.com/sergeykonkin/memq"
    )

    // Create new Broker:
    b := memq.NewBroker()

    // Subscribe to topic:
    b.Subscribe("my-topic", func(msg interface{}) {
        fmt.Printf("New message: %v", msg)
    })

    // Publish to topic:
	b.Publish("my-topic", "Hello World!")

    // You can also save a subscription to unsubscribe later:
    sub := b.Subscribe("my-topic", anotherHandler)
	defer sub.Unsubscribe()
    ```

Note that for now Broker supports only fan-out mode, so if there are 2 or more subscribers for a topic - they all will be receiving new messages.

## License

[MIT](LICENSE)
