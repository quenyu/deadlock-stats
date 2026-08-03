# Architecture decision: clean-slate v2

## Product boundary

The first release answers one question only:

> Which victory on the same hero is most comparable to this match, what summary metrics differ, and what one measurable target should the player try next?

## Runtime path

```text
HTTP request
  -> validate account and match IDs
  -> Deadlock API match-history request
  -> pure Ghost Match selection and metric calculation
  -> JSON response with confidence and limitations
```

## Deliberately absent

- PostgreSQL and migrations;
- Redis and multi-layer caches;
- authentication and user accounts;
- old rank conversion tables;
- static hero asset loading;
- search and fuzzy search;
- crosshair/build/social features;
- metrics infrastructure beyond structured logs;
- LLM-generated coaching.

These are not forbidden forever. They require an observed product need first.

## Data rules

1. Upstream DTOs live only in `internal/deadlock`.
2. `match_result` is the winning team; player victory is `player_team == match_result`.
3. The domain report contains only values used by Ghost Match.
4. No rank is inferred from an obsolete local conversion table.
5. No causal claim is made from end-of-match aggregates.
6. A missing comparable victory is a valid product result, not an internal server error.

## When storage becomes justified

Add PostgreSQL only when at least one of these becomes necessary:

- global candidate pools across players;
- patch-aware historical cohorts;
- mission completion tracking;
- product analytics requiring durable events.

Add Redis only after measuring an actual upstream-rate or latency problem that cannot be handled by process-local caching and request coalescing.

## Next technical milestone

1. Validate the endpoint against real account histories.
2. Add an adapter for match metadata and item timing behind a separate interface.
3. Build a single report page.
4. Recruit initial testers before adding accounts, payments or AI text generation.
