# Changelog

## [3.0.0] UNRELEASED
- `GetMyIP` now returns `net.IP` instead of `string`.
- `CalcHoneyScore` now accepts `net.IP` instead of `string`.
- `GetDNSResolve` now returns `map[string]*net.IP` instead of `map[string]string`.
- `GetDNSReverse` now accepts `[]net.IP` instead of `[]string`.
