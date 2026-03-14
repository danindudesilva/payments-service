# ADR 0001: Use Stripe first, but hide it behind an interface

## Status
Accepted

## Decision
Implement Stripe first for testability and learning speed, but isolate it behind our own gateway abstraction.

## Why
- Free test mode
- Realistic 3DS testing support
- Strong documentation
- Easier future provider swap
