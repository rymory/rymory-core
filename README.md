# Rymory — Open Identity Infrastructure

> **Author & Creator:** Onur Yaşar ([@onurid](https://github.com/onurid))
> Built since 2017. All rights reserved.

---

## What is Rymory?

Rymory is an open identity infrastructure project — a federation-ready identity layer for distributed applications, multi-tenant ecosystems and modern digital trust.

It is **not a concept**. The core identity engine has been in active development since 2017, evolving through multiple production iterations. What you see here is that infrastructure, opened.

### Core capabilities

- **Federated SSO** — one identity, many applications, cross-domain session propagation
- **JWT HS512 sessions** — short-lived tokens, refresh chains, passkey/security-key support
- **Role-based access control** — six-tier hierarchy embedded directly in JWT claims
- **Multi-tenant isolation** — appId + merchantId + projectId scoping at every layer
- **Brute-force protection** — progressive lockout (3 → 6 → 9 attempts), account lock propagation
- **Multi-language** — EN, TR, RU, ES, FR, ET, IT, PL built in
- **Go backend** — modular package architecture (authenticate, account, role, project, member, validation)

---

## Architecture

```
id.rymory.org          ← single sign-on entry point
      │
      ├── notes.lemoras.com
      ├── drive.lemoras.com
      ├── [any third-party app]
      │
account.rymory.org     ← account management
```

Identity lives in Rymory. Services live anywhere. The protocol is the product.

---

## Repository Structure

```
rymory-core/           ← Go identity gateway (authenticate, account, role, project)
rymory-spec/           ← Protocol specification and RFC documents
rymory-ui/             ← Reference frontend implementation
```

---

## License

This project is licensed under **GNU AGPL v3** with a **Commercial License Exception**.

| Use case | License |
|---|---|
| Personal / academic / open source (AGPL-compliant) | Free — AGPL v3 |
| Commercial product or SaaS without source disclosure | Paid — contact author |

See [LICENSE.txt](./LICENSE.txt) for full terms.

For commercial licensing: **onxorg@proton.me**

---

## Status

Rymory is currently a project under active development by its sole author.
Long-term governance goal: independent foundation model.

> *"Infrastructure this critical should be owned by the community, not a single company."*

---

## Author

**Onur Yaşar**
Sole creator and author since 2017.

- GitHub: [@onurid](https://github.com/onurid)
- Email: onxorg@proton.me
- Web: [rymory.org](https://rymory.org)
- Trademarks: "Rymory" and "Lemoras" are registered with TÜRKPATENT

---

*© 2017–2026 Onur Yaşar. All rights reserved.*
*"Rymory" is a trademark of Onur Yaşar.*
