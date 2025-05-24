# ctx Package

This package provides utilities for managing operation IDs within Go's `context.Context`, enabling traceability and correlation across service boundaries.

## Features (Implemented)

- **Operation ID Management**: Easily set and retrieve an `operationID` in a context for tracking operations across service calls.
- **Type-Safe Context Keys**: Uses custom types for context keys to avoid collisions and ensure type safety.
- **Helper Functions**: Provides `WithOperationID` and `GetOperationID` for convenient operation ID handling.

## Possible Future Enhancements

- **Request ID Support**: Add functions to set and get a request ID in context for tracking individual requests.
- **User Context**: Support storing user information (e.g., user ID, roles) in context.
- **Correlation ID**: Add correlation ID utilities for distributed tracing across services.
- **Timeout Management**: Provide helpers for creating contexts with common timeout presets.
- **Context Cancellation**: Add utilities for creating and propagating cancellable contexts.
- **Context Value Validation**: Type-safe validation and presence checks for required context values.
- **Context Metadata**: Support for storing and managing arbitrary metadata in context.
- **Context Chain Utilities**: Functions to merge or inherit values between contexts.

---

Feel free to contribute or suggest additional features!
