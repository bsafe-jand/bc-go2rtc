# Fork notes

This is a fork of [AlexxIT/go2rtc](https://github.com/AlexxIT/go2rtc). It carries a
fix for upstream issue [#2443](https://github.com/AlexxIT/go2rtc/issues/2443)
(concurrent map access in the `/api` and `/api/streams` handlers) until an
equivalent fix is released upstream.

## Branches

| Branch | Purpose |
| --- | --- |
| `develop` | Integration branch. Feature and sync PRs target it. |
| `master` | Mirrors upstream `master`. |

`develop` starts from upstream `master` at `fc13c41c`.

## Upstream sync

```sh
git remote add upstream https://github.com/AlexxIT/go2rtc
git fetch upstream --tags
```

Check upstream weekly. Merge each upstream release tag into `develop` through a PR:

1. `git switch -c sync/upstream-vX.Y.Z develop`
2. `git merge --no-ff vX.Y.Z` (merge, do not rebase)
3. Resolve conflicts and run `go test ./internal/... ./pkg/...`, plus
   `go test -race -run Concurrent ./internal/api/ ./internal/app/ ./internal/streams/`
   where cgo is available.
4. Open a PR into `develop`.

If upstream releases an equivalent fix, take upstream's version in the sync PR and
drop the fork change.

`docs/fork/upstream-sync.yml.example` is an optional workflow that opens such a
sync PR. To enable it, copy it to `.github/workflows/`. It only pushes a branch and
opens a PR.
