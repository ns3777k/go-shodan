# Changelog

## [3.1.0] Unreleased
- Add ssl-related data to `HostData`.

## [3.0.2]
- Implement data bulk api
- Make use of dep instead of glide
- `SetDebug` is goroutine-safe now
- Drop support for go 1.6
- Lots of refactoring to support `context.Context` in all methods.
- Rewritten streaming methods and implemented some more.
- `GetMyIP` now returns `net.IP` instead of `string`.
- `CalcHoneyScore` now accepts `net.IP` instead of `string`.
- `GetDNSResolve` now returns `map[string]*net.IP` instead of `map[string]string`.
- `GetDNSReverse` now accepts `[]net.IP` instead of `[]string`.

## [2.0.4]
- Dropped support for old golang versions (`1.1` - `1.5`).
- `GetHttpHeaders` is renamed to `GetHTTPHeaders`.
- Invalid url will now `panic`.
- Stream methods no longer return error.
