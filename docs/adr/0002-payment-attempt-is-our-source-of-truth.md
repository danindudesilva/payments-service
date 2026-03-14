# ADR 0002: Payment attempt is our source of truth

## Status
Accepted

## Decision
Persist and manage a local payment attempt lifecycle instead of trusting the last client response.

## Why
This makes the flow recoverable after browser refreshes, kiosk disconnects, webhook retries, or provider-side asynchronous completion.
