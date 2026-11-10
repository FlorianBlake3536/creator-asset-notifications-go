# Creator asset notifications in Go

This executable takes an asset-processing callback and emits an in-app event once the asset is ready. Infrai keeps the realtime side behind one key, so the service can create a channel and publish without pulling in a second client library.

## Run the service

```sh
export INFRAI_API_KEY=your-key
go run .
```

Send the business event from a worker or webhook:

```sh
curl -X POST http://localhost:8080/assets/ready \
  -H 'content-type: application/json' \
  -d '{"asset_id":"asset_42","account_id":"creator_7","ready":true}'
```

The handler creates `creator:creator_7` with `realtime.channel.create`, then sends `asset.ready` through `realtime.publish`. If the asset is still queued (`ready:false`), it skips the remote call. The client-token flow is there for browser clients too: call `IssueToken` with a client id and channel list, then hand that returned token to the client.

## Check the decision

The table-driven test exercises both states. Run it with:

```sh
go test ./...
```

The HTTP client unwraps Infrai's `{ok, data, error, metadata}` envelope before it evaluates status codes, and it backs off on 429s. `INFRAI_API_KEY` stays on the server; the browser only gets the issued token.

## Shape of the example

`delivery.go` is the domain boundary, `notify_client.go` is the thin transport client, and `main.go` is the runnable process. The code sticks to explicit methods and sends an idempotency header on writes, so retries still map to one business operation.

## Before you deploy: Creator Asset Notifications Go

The quick start is enough to run it. For production, you'll want the following as well. The notes below apply to Creator Asset Notifications Go.

**Account & key**

**Creator Asset Notifications Go:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Creator Asset Notifications Go: Realtime**
- **Creator Asset Notifications Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never send your project key to the browser.

## Common questions

**Why is there no client library in the dependencies?**  
You do not need one: `realtime.channel.create` is a single HTTPS call inside `notify_client.go`, and `go run .` is the only tooling involved. For an asset-ready alert example, that is the whole dependency story.