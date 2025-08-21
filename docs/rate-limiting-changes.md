# Rate Limiting Handling Changes

## Overview

This document describes changes made to handle rate limiting (HTTP 429 - Too Many Requests) errors from the Oracle Cloud Infrastructure (OCI) API.

## Issue

The application was encountering "Too Many Requests" errors (HTTP 429) when making multiple API calls to the OCI Compute Service, specifically when trying to enrich instance data with network information by calling the ListVnicAttachments API endpoint.

Error message example:
```
Error: failed to execute root command: list instances: listing instances from repository: enriching instance ocid1.instance.oc2.us-luke-1.anwhkljrlkxuqyqcb7ywrsz3qk57xsvf4kw6dgsxyozxcqhyc5r3vynb25fa with network: Error returned by Compute Service. Http Status Code: 429. Error Code: TooManyRequests.
```

## Solution

Implemented a retry mechanism with exponential backoff for all API calls in the instance adapter that might encounter rate limiting. This includes:

1. `getPrimaryVnic` - Calls `ListVnicAttachments` and `GetVnic`
2. `getSubnet` - Calls `GetSubnet`
3. `getVcn` - Calls `GetVcn`
4. `getImage` - Calls `GetImage`
5. `getRouteTable` - Calls `GetRouteTable`

### Retry Strategy

For each API call:
- Maximum of 5 retry attempts
- Initial backoff of 1 second
- Exponential backoff with a maximum of 32 seconds
- Only retry on HTTP 429 (Too Many Requests) errors
- For other errors, fail immediately

## Implementation Details

The retry mechanism follows this pattern:
1. Make the API call
2. If successful, proceed normally
3. If error is HTTP 429 (Too Many Requests):
   - If maximum retries reached, return error
   - Otherwise, wait with exponential backoff and retry
4. If any other error, return immediately

## Future Considerations

Consider implementing a more general retry mechanism that can be applied to all API calls in the application, possibly as a middleware or wrapper around the OCI client.