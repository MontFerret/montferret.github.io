---
title: "Wire security"
weight: 31
draft: false
description: "Configure verified TLS, mutual TLS, and host-supplied token authentication for a Wire endpoint."
---

# Secure a Wire endpoint

Wire exposes the host's configured runtime. Protect that endpoint with verified transport identities and host-defined access policy. The examples below are compiled and tested in the [Wire demo module](https://github.com/MontFerret/montferret.github.io/tree/main/examples/wire); follow the [setup instructions]({{< ref "docs/embedding/wire" >}}#run-a-host-and-client) first. They use the same published Wire `v{{< data "versions.wire.v1" >}}` and API `v{{< data "versions.wire.api" >}}` consumable-output contract as the greeting demo.

Each host function accepts an existing `api.Runtime` and returns a configured `*server.Server`. Start it with blocking `srv.Run(ctx, address)`, or use `srv.Serve(ctx, listener)` for a caller-created listener. Both paths apply the same security policy. The tests use ephemeral loopback listeners; production addresses, certificate loading and renewal, trust stores, and token validation are host inputs. Wire provides no certificate-management or built-in authentication service.

## Verify the server with TLS

Supply the host's certificate and the client's trusted CA roots. The client verifies both trust and the expected server identity: `serverName` must match a DNS subject alternative name in the certificate, even when the TCP target is a loopback address. For an IP identity, use a certificate with the corresponding IP subject alternative name. Never bypass verification.

{{< tabs >}}
{{< tab title="Host" >}}
{{< code lang="go" name="security/tls_host" file="examples/wire/security/tls_host.go" >}}{{< /code >}}
{{< /tab >}}
{{< tab title="Client" >}}
{{< code lang="go" name="security/tls_client" file="examples/wire/security/tls_client.go" >}}{{< /code >}}
{{< /tab >}}
{{< /tabs >}}

Without a custom transport option, `client.New` uses TLS with system trust and normal identity verification. A failed TLS handshake never falls back to plaintext. `client.WithInsecure()` is an explicit plaintext choice for a local development endpoint; it does not mean TLS with verification disabled. A TLS host rejects plaintext clients.

TLS here authenticates the server. It does not identify or authorize the client.

## Require client identities with mTLS

Mutual TLS also verifies the client's certificate. The host requires a client certificate chaining to `clientCAs`; missing or untrusted certificates fail the handshake. The client still verifies the server using its own trust roots and expected identity.

{{< tabs >}}
{{< tab title="Host" >}}
{{< code lang="go" name="security/mtls_host" file="examples/wire/security/mtls_host.go" >}}{{< /code >}}
{{< /tab >}}
{{< tab title="Client" >}}
{{< code lang="go" name="security/mtls_client" file="examples/wire/security/mtls_client.go" >}}{{< /code >}}
{{< /tab >}}
{{< /tabs >}}

Host middleware and runtime calls can inspect `peer.FromContext(ctx)` and `credentials.TLSInfo`, including verified certificate chains. A trusted certificate establishes an identity; the host still decides what that identity may do.

## Add token authentication over TLS

The host supplies unary and stream middleware, inside Wire's sanitized panic-recovery boundary, for every registered RPC. This example rejects absent, malformed, and invalid bearer metadata with a generic `Unauthenticated` status. Its validation callback is the host's policy, not a built-in token implementation.

{{< tabs >}}
{{< tab title="Host" >}}
{{< code lang="go" name="security/token_host" file="examples/wire/security/token_host.go" >}}{{< /code >}}
{{< /tab >}}
{{< tab title="Client" >}}
{{< code lang="go" name="security/token_client" file="examples/wire/security/token_client.go" >}}{{< /code >}}
{{< /tab >}}
{{< /tabs >}}

`client.WithPerRPCCredentials` installs credentials at the **connection level**, covering the streaming startup handshake, subsequent unary and streaming operations (including `ReadExecutionOutput`), and resource cleanup. `RequireTransportSecurity` returns `true`; keep transport credentials and call credentials separate, as in [gRPC's authentication model](https://grpc.io/docs/guides/auth/). Never put bearer tokens in targets, option diagnostics, or logs. A credential source may refresh tokens and must support concurrent calls, including cleanup after an operation's context is canceled.

Stream middleware authenticates at establishment. A stream already open when a token expires or is revoked does not automatically run that middleware again. The host must define and enforce expiry and revocation for established streams, including Wire's connection-lifetime stream.

Authentication identifies a caller. Authorization decides which compilation, execution, debugging, and host capabilities that caller can use. Neither creates tenant isolation. Wire's logical connection and resource IDs define lifetimes; they do not bind resources to a principal. Keep principal checks and any required separation of runtimes, host services, data, and tenant capabilities in the host's policy.

## Execute and close resources

Each client function returns a runtime that owns its channel. After connecting, pass it to this query function, then close the remote runtime and report both operation and cleanup errors. The query collects and closes its output before closing the session and plan, then returns detached `*api.Content`; the caller closes the runtime. Preparation errors come from `Run`; terminal errors and available content come from `Collect`. Preserve both the returned content and error.

{{< code lang="go" name="security/query" file="examples/wire/security/query.go" >}}{{< /code >}}

Keep the query's invocation context alive through collection. `Consume` is the alternative for forwarding borrowed chunks without collecting the entire result; see [output consumption]({{< ref "docs/embedding/wire" >}}#consume-output-and-observe-completion). Message limits apply per gRPC message, not to the aggregate result. Budget memory for `Collect` and materialized debugger output.

`client.New`'s context controls construction only. Once startup succeeds, cancel individual operations through their own contexts. Runtime closure rejects new root operations, while existing plans and admitted work may keep the owned channel alive until their final release. `client.From(ctx, conn)` always borrows a caller-owned channel; the caller must close it after descendants settle. On the host, a managed shutdown timeout can leave retained cleanup in progress: await actual `Shutdown` settlement before closing the borrowed runtime. See [startup and shutdown ownership]({{< ref "docs/embedding/wire" >}}#choose-startup-and-shutdown-ownership).

## Inspect option validation

`server.New` and `Server.Run` invoke every non-nil option once in registration order, collect nil-option and validation failures with `errors.Join`, and abort before server construction or startup listening. A rejected option never applies its setter. Nil credentials, nil interceptor entries, non-positive limits, and non-positive shutdown timeouts are invalid. Empty interceptor lists are valid; repeated lists preserve registration order, while valid credential and timeout overrides use the last value. `WithLimits` replaces the complete limit set, so start from `server.DefaultLimits()` when changing one field.

Factory errors expose `go-options` `ValidationError` through `errors.As`:

{{< code lang="go" name="security/validation" file="examples/wire/security/validation.go" >}}{{< /code >}}

`errors.As` finds one matching error in a wrapped or joined tree; preserve the full error to inspect all failures and their underlying causes. Collection errors have named outer fields such as `unary interceptors` or `limits`, then relative labels such as `[0]` or `["max connections"]`. Interceptor failures appear in ascending index order. Limit failures come from a private map with trusted field names; their order is unspecified.

Collection aggregate and location wrappers set `OmitValue`, omitting whole slices and maps from diagnostics. Child errors retain rejected values and causes, such as a nil entry or a non-positive integer. This is not general redaction: pass only non-secret configuration values to validators and keep authentication secrets out of validation errors.

{{< docs-related tiles="embedding-wire,embedding-go-custom-functions,embedding-go" >}}
