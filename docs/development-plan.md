# Development-plan

## Slice 1
- Go service skeleton
- health endpoint
- logging
- Dockerfile
- CI
- ADRs

## Slice 2
- payment attempt entity
- status model
- repository and gateway interfaces
- state transition tests

## Slice 3
- create payment attempt endpoint
- in-memory repository
- fake provider

## Slice 4
- Stripe test mode integration
- PaymentIntent creation
- requires_action mapping

## Slice 5
- webhook endpoint
- idempotent event processing
- duplicate event tests

## Slice 6
- browser client
- reconnect / resume flow
- Cloud Run deployment
