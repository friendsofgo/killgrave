# Acceptance tests

Run `make acceptance` from the repository root. It builds Killgrave and runs each
`tests/<scenario>.txtar` against a separate process with an isolated working
directory and port. No external services are needed.

Keep each scenario in one archive containing its setup and named HTTP steps:

```text
-- imposters/routes.imp.json --
[{"request":{"method":"GET","endpoint":"/example"},"response":{"status":200,"headers":{"Content-Type":"text/plain"},"body":"example\n"}}]
-- http/matched/req.http --
GET /example HTTP/1.1

-- http/matched/res.http --
HTTP/1.1 200 OK
Content-Length: 8
Content-Type: text/plain


example
```

Each `http/<step>/` group is a named subtest containing `req.http` and `res.http`.
Steps run in archive order, so response sequences are visible together without
numbering separate files. Scenarios run concurrently. Expected responses include
all headers except `Date`; redirects are checked directly.

All entries outside `http/`, except `scenario.json`, are extracted as setup files.
By default Killgrave uses the extracted `imposters` directory and an allocated
port. An optional `scenario.json` entry can replace startup arguments:

```text
-- scenario.json --
{"args":["--config","settings/config.yml"]}
-- settings/config.yml --
host: localhost
port: ${PORT}
imposters_path: imposters
```

Arguments and setup files support `${PORT}` and `${UPSTREAM}` substitution.
Arguments run directly without a shell. Config-file imposter paths resolve
relative to the config file; `config_precedence.txtar` records precedence over
host, port, and imposter flags.

Optional `scenario.json` settings:

| Setting | Meaning |
| --- | --- |
| `args` | Replace the default startup arguments. |
| `error` | Require an unsuccessful exit and this diagnostic substring. No HTTP steps are needed. |
| `occupyPort` | Hold the allocated port open to exercise bind failure. |
| `upstreamCalls` | Start a local upstream and assert this exact request count, including zero. |
| `upstreamDown` | Close that upstream before startup to exercise a connection failure. |
| `https` | Send fixture requests over HTTPS. Startup arguments or configuration must also enable secure mode. |

The upstream returns status 202, `X-Upstream: yes`, and a body echoing the method,
request URI, `X-Test` header, and body. TCP readiness does not call the upstream
or advance response sequences.

Watcher steps can include `http/<step>/write/<relative-path>` entries to atomically
create or replace files before the request. An `http/<step>/eventually` entry
containing `true` retries the response check for up to five seconds. Writes happen
once. Use retries only for stateless GET requests without bodies; sequence tests
must check every response without retries. See `watcher_flag.txtar` for an example.

Delay steps may include `http/<step>/timing.json`, for example
`{"min":"180ms","max":"2s"}`. The harness measures the complete response, with a
lower bound to detect missing delays and a generous upper allowance for CI load.
Fixed and ranged delays are covered in both JSON and YAML. Timing checks cannot
use retries.

HTTPS clients accept Killgrave's bundled dummy certificate only within the test
transport. HTTPS proxy fixtures still connect to Killgrave before reaching the
local upstream; the unit suite separately checks a trusted local TLS upstream.

`TestGracefulShutdown` is a Go table test for SIGINT and SIGTERM. It uses channels
to hold an upstream request open, checks that new connections are refused during
shutdown, then releases the response and verifies clean exit. This coordination
stays in Go rather than adding concurrency commands to the archive format.

HTTP requests time out after three seconds; readiness and startup-error commands
have five-second limits. Cleanup allows three seconds for an interrupt before
killing and reaping the process, and verifies that the listening port can be bound
again. Failed scenarios include server logs.
