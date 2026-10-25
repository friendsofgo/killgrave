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

The upstream returns status 202, `X-Upstream: yes`, and a body echoing the method,
request URI, `X-Test` header, and body. TCP readiness does not call the upstream
or advance response sequences.

Watcher steps can include `http/<step>/write/<relative-path>` entries to atomically
create or replace files before the request. An `http/<step>/eventually` entry
containing `true` retries the response check for up to five seconds. Writes happen
once. Use retries only for stateless GET requests without bodies; sequence tests
must check every response without retries. See `watcher_flag.txtar` for an example.

HTTP requests time out after three seconds; readiness and startup-error commands
have five-second limits. Cleanup allows three seconds for an interrupt before
killing and reaping the process. Failed scenarios include server logs.
