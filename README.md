# Interactions

Fetch a rocket for a given channel using:
````sh
make rocket id=<channel_id>
````

Fetch all rockets:
````sh
make rockets
````

You can apply optional query parameters to filter the results:
- `filterExploded` (`bool`): If true, only non-exploded rockets are returned.
- `limit` (`int`): Maximum number of events per channel. Use `-1` to return all.

Examples:
````sh
make rockets filterExploded=true limit=2
make rockets limit=2

````

These values are appended as query parameters to the API call.

# Design

The project follows the following design principles:
- Responsibilities are separated using controllers, models, and repository.
- Domain-Driven Design (DDD) patterns are applied:
    - `IRepository` abstracts the storage logic.
    - Rehydration logic lives close to the domain entity (`RocketState`).
- Dependency injection is used to provide `IRepository` to the `Controller`, adhering to the Dependency Inversion Principle.
- A simple GitHub Actions workflow runs unit tests on every `push` and `merge` to the `main` branch, ensuring code is validated automatically.

# Scalability
To improve scalability:
- Add observability: logs, metrics, and tracing.
- Introduce read models or snapshots to avoid reconstructing `RocketState` from events on every query.

This optimization should be done only after observability is in place and metrics confirm a bottleneck. Premature optimization should be avoided.

Read models or snapshots can be created:
- Inline (in the same transaction), or
- Using an outbox pattern (eventual consistency).

# Short-Time Improvements
- Add integration tests for `Controller` endpoints.
- Add documentation comments to all public functions and interfaces.
- `MemoryRepository` uses a Go map, which does not guarantee deterministic order when listing rockets.
- Use persistent storage to avoid data loss on restart.
