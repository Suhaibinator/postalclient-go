# Changelog

## Unreleased

This release fixes response decoding against real Postal servers and adds
context support. It contains breaking API changes; the module path stays at v1.
Most of the changed types could never decode a real Postal response, so code
that relied on them could not have worked.

### Breaking changes

- Every API method takes a `context.Context` as its first parameter:
  `SendMessage(ctx, req)`, `SendRaw(ctx, req)`, `GetMessage(ctx, id, expansions...)`
  and `GetMessageDeliveries(ctx, id)`. Pass `context.Background()` for the old behavior.
- `models.SendMessageResponse.MessageID` is now a `string` (was `int`). Postal
  returns the RFC 5322 Message-ID, such as `"<uuid>@<return-path-domain>"`, so
  every successful send failed to decode before.
- `models.SendMessageResponse.Token` is removed (Postal never returned it). The new
  `Messages` field maps each recipient to a `models.RecipientMessage` holding the
  numeric `ID` (the one `GetMessage` and `GetMessageDeliveries` take) and `Token`.
- `models.Delivery.LogID` is now a `string` (was `int`). Postal log IDs are
  strings such as `"A1B2C3D4"`.
- `models.Message.Headers` is now `map[string][]string` (was `map[string]string`).
  Postal returns every header as a list of values.
- `models.Message.Attachments` is now `[]models.MessageAttachment`, with Postal's
  `Filename`, `ContentType`, `Data`, `Size` and `Hash` fields.
  `models.Attachment` is now the send-request attachment only: its `Size` field is
  removed (Postal ignored it) and an empty `ContentType` is omitted.
- `Error.ErrorCode` and `Error.Message` now come from the response's
  `data.code` and `data.message` (they were read from top-level fields Postal
  never sends) and are no longer JSON fields of `Error`. `Error.Error()` includes
  the code when there is one.
- A non-200 response whose body is not a Postal API response now returns a
  `*postalclient.HTTPError` with the status code and the first 64 KiB of the body,
  instead of a JSON decoding error.

### Fixes

- `GetMessageDeliveries` decodes Postal's Unix-seconds `timestamp` into
  `Delivery.Timestamp` (it expected an RFC 3339 string, so any delivery failed to decode).
- `GetMessage` accepts expansions and sends them as `_expansions`; without them
  Postal returns only the ID and token. Constants `models.ExpansionStatus`,
  `ExpansionDetails`, and so on name them.
- `models.MessageStatus` and `models.MessageDetails` now have fields; they were empty.
- `models.Message` gains `ActivityEntries` for the `activity_entries` expansion.
- A trailing slash on `Client.BaseURL` no longer produces a `//` in request paths.
- A `Client` with a nil `HTTPClient` uses `http.DefaultClient` instead of panicking.
