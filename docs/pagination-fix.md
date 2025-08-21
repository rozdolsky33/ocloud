# Pagination Fix

## Issue

The application was experiencing a panic with the following error:

```
panic: runtime error: slice bounds out of range [-300:]
```

This occurred in the `List` method of the instance service when called with `pageNum=0`. The error happened because the pagination logic was calculating a negative start index when `pageNum` was 0:

```go
start := (pageNum - 1) * limit
```

With `pageNum=0` and `limit=300`, this resulted in `start = -300`, which is an invalid slice index in Go.

## Solution

Modified the `List` method in `internal/services/compute/instance/service.go` to handle the case where `pageNum` is 0 or negative by treating it as page 1:

```go
// Handle pageNum=0 as the first page
if pageNum <= 0 {
    pageNum = 1
}

start := (pageNum - 1) * limit
```

This ensures that the start index is never negative, preventing the panic.

## Root Cause

The issue occurred because the bastion flow in `cmd/identity/bastion/flows_instance.go` was calling the instance service with `pageNum=0`:

```go
instances, _, _, err := instService.List(ctx, 300, 0, true)
```

While the service implementation expected `pageNum` to start at 1.

## Future Considerations

Consider standardizing pagination across all services to either:
1. Use 0-based indexing consistently (page 0 is the first page)
2. Use 1-based indexing consistently (page 1 is the first page)

Document this convention clearly to prevent similar issues in the future.