# Acceptance tests

Run `make acceptance` from the repository root. This builds the binary before
testing it as a separate process. Scenarios run concurrently without external
network services.

Each directory in `tests/` contains:

- `config.txtar`: files to extract into an isolated temporary working directory.
  Paths are relative to that directory. The archive can contain only a comment
  when a startup error needs no files.
- `http/*.txtar`: a `req.http` and `res.http` pair per exchange. These run in
  filename order within a scenario. Expected responses include all headers
  except `Date`, which the harness removes. Redirects are asserted directly.
- Optional `scenario.json`: startup arguments and additional expectations.

Without `scenario.json`, Killgrave starts with `--imposters <temp>/imposters`
and an available port. If `args` is present, it replaces those default arguments.
Arguments and extracted file contents support `${PORT}` and `${UPSTREAM}`
substitution. Arguments run directly without a shell.

For example, to exercise configuration from a different directory:

```json
{
  "args": ["--config", "settings/config.yml"]
}
```

The extracted config must include `port: ${PORT}`. Its imposter path should be
relative to the config's directory, rather than the process working directory.
The `config_precedence` scenario records the current behavior: config-file host,
port, and imposter settings take precedence over the corresponding flags.

Other settings:

| Setting | Meaning |
| --- | --- |
| `error` | Expected diagnostic substring; require an unsuccessful exit within five seconds. No HTTP fixtures are used. |
| `occupyPort` | Hold the allocated port open to exercise bind failure. |
| `upstreamCalls` | Start a local upstream and assert this exact request count. Zero also verifies that traffic stays local. |
| `upstreamDown` | Close that upstream before launching Killgrave to exercise a connection failure. |

The upstream returns status 202, `X-Upstream: yes`, and a text body echoing the
method, request URI, `X-Test` header, and request body. Proxy fixtures use these
to check forwarding as well as upstream response propagation. Readiness uses a
TCP connection, so it neither calls the upstream nor advances a response sequence.

HTTP calls time out after three seconds; readiness and startup-error commands
have five-second limits. Cleanup allows three seconds for an interrupt before
killing and reaping the process. Server logs are attached to failed scenarios.

Current groups cover config-file paths and precedence, matching/schema fallback,
proxy modes (including schema-rejected bodies and unavailable upstreams), startup
errors, and the original simple resource example. Malformed imposter handling is
not yet specified by this suite; malformed startup configuration is covered.
