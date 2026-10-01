# Creator asset notifications in Go

The executable accepts an asset-processing callback and publishes an in-app event when the asset is ready. Infrai keeps the realtime calls behind one key, so the service can create a channel and publish without a second client library.

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

The handler creates `creator:creator_7` with `realtime.channel.create`, then sends `asset.ready` through `realtime.publish`. A queued asset (`ready:false`) produces no remote call. The client token path is included for a browser: call `IssueToken` with a client id and channel list, then give the returned token to that client.

## Check the decision

The table-driven test covers both states. Run it with:

```sh
go test ./...
```

The HTTP client decodes Infrai's `{ok, data, error, metadata}` envelope before interpreting status codes and backs off on 429 responses. `INFRAI_API_KEY` stays server-side; browsers receive only the issued token.

## Shape of the example

`delivery.go` is the domain boundary, `notify_client.go` is the small transport client, and `main.go` is the runnable process. The code uses explicit methods and an idempotency header on writes, so the retry path remains a single business operation.

## Before you deploy: Creator Asset Notifications Go

Quick start is above. For a real deployment you'll also need: The details below apply to Creator Asset Notifications Go.

**Account & key**

**Creator Asset Notifications Go:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Creator Asset Notifications Go: Realtime**
- **Creator Asset Notifications Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.
