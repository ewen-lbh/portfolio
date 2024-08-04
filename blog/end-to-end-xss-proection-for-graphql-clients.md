---
date: 2024-08-04
title: End-to-end XSS protection for GraphQL clients
works: [churros]
---

# End-to-end XSS protection for GraphQL clients

There you are, accepting Markdown content from your users, and exposing on your API a field that resolves to a XSS-safe HTML rendering of that user-generated content.

You're feeling good about yourself, and you should. You're doing the right thing. But, within the same breath, you turn around to your client code, get that generated HTML _string_, and include it somewhere on your page, telling your templating engin / js framework “it's fine, this is safe for HTML inclusion”.

A few months pass, you add more and more features, and suddenly your API has _a lot_of fields, some that output HTML and some that output plain text.

As a documentation feature, creating a separate `HTML` scalar type to signal fields that have HTML content (and thus cannot be inserted into the page as simple text).

But we can get way better than that: _end-to-end guarantees_ of XSS safety, by type systems.

## The recipe

1. Define a separate scalar type for HTML content on your GraphQL API: we'll call it `HTML`.
2. On your client, create an [Opaque type](https://en.wikipedia.org/wiki/Opaque_data_type) that wraps a string, such that no string value can be assigned to that opaque type, but data of that type can be assigned to a string.
3. Bind the `HTML` scalar type to the opaque type on your client.
4. Create a component that accepts a value of the opaque type, and includes the HTML value without escaping
5. (bonus) Enforce that no other place on your frontend code is allowed to include HTML without escaping.

## An example tech stack: Pothos -> Houdini -> SvelteKit:


