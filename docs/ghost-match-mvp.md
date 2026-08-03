# Ghost Match MVP

## Product hypothesis

For a selected Deadlock match, find a successful match played on the same hero with a similar duration, economy, KDA and rank. Show measurable differences and generate exactly one objective for the next match.

The core promise is deliberately narrow:

> We do not claim to know why the player made a decision. We show where measurable outcomes differ from a comparable successful game.

## What is feasible with the repository today

The current Deadlock API client already loads player match history with:

- match ID;
- hero;
- kills, deaths and assists;
- net worth;
- duration;
- result/team;
- start time.

That is sufficient for the first backend slice:

1. Select a target match.
2. Filter successful matches on the same hero.
3. Rank candidates by duration, net worth per minute, KDA and rank proximity.
4. Compare economy, KDA and death rate.
5. Generate one measurable training mission.
6. Return confidence and explicit limitations.

The first implementation uses the player's own successful history. This avoids pretending that we already have a complete global cohort dataset.

## What is not yet supported honestly

The following claims require richer match metadata or demo parsing:

- exact moment where two matches diverged;
- movement and lane route comparison;
- purchase timings;
- teamfight participation by timestamp;
- objective decisions;
- positioning mistakes;
- causal statements such as “this action lost the match.”

They must not be generated from aggregate match history alone.

## Next integration steps

### Stage 1 — expose the service

Add an endpoint such as:

```text
GET /api/v1/players/{steamId}/matches/{matchId}/ghost
```

The handler should:

1. load recent matches through `PlayerProfileService`;
2. call `GhostMatchService.BuildReport`;
3. return `404` when the target match is missing;
4. return `422` when no comparable victory exists.

### Stage 2 — global candidate pool

Create a stored candidate table populated by background jobs. Partition candidates by:

- patch/build version;
- hero;
- rank bucket;
- duration bucket.

Only then should the service compare a player with broader cohorts.

### Stage 3 — metadata/demo enrichment

Add item timings, ability order and timeline events from verified endpoints. Keep every recommendation linked to source metrics and include a confidence score.

## Acceptance criteria for the first public experiment

- A report loads for at least 70% of users with ten or more matches on one hero.
- The selected reference match is understandable to the user.
- The report contains no unsupported causal language.
- The mission can be checked automatically after the next match.
- At least 20 test users complete two report → mission → next-match cycles.

## Monetization gate

Do not add payment before measuring repeat usage. A reasonable initial signal is that at least 20–30% of test users return after another match to check whether they completed the mission.
