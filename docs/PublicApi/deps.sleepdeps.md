# `deps.SleepDeps`

`sandbox/deps/sleepdeps`

## `Contract`

Contract is the pause library injected whole as the Deps.SleepDeps field.

| Field | Type | Description |
| --- | --- | --- |
| `Sleep` | `func(nanoseconds int64)` | Sleep pauses the calling goroutine for at least nanoseconds nanoseconds. A zero or negative duration returns at once. The sandbox may not name a `time.Duration`, so a duration crosses this boundary as a plain integer. |

[every contract](doc.md)
