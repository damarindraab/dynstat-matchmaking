# Custom MatchMaker - Dynamic Stat Selection

## Summary

The matchmaker.go file implements a dynamic stat-based matchmaking system. Players select which stat to use for matching (e.g., character-specific MMR), and the server normalizes this into a standard attribute for AGS matching.

## Functions

### GetStatCodes()

Returns the list of stat codes configured in `statistics_config.statistics`. AGS uses this to know which player statistics to fetch.

### EnrichTicket()

For each player in the ticket:
1. Gets the selected stat code from `TicketAttributes[playerID]`
2. Extracts the stat value from `Player.Attributes[selectedStat]`
3. Sets `Player.Attributes[enrichedKey]` to the stat value
4. Removes all configured statistics from `Player.Attributes`
5. Cleans up player ID mappings from `TicketAttributes`

If a player is missing the selected stat or it has an invalid type, the enriched key is not set (validation will fail).

### ValidateTicket()

Checks that each player has the enriched attribute in their `Player.Attributes`. This is a post-enrichment validation - if any player is missing the enriched key, validation fails.

### RulesFromJSON()

Unmarshals the JSON rules string to `GameRules` struct and returns it.

### MakeMatches()

Returns `nil` to signal `UNIMPLEMENTED`. AGS uses its default matching logic based on the enriched player attributes.

### BackfillMatches()

Returns `nil` to signal `UNIMPLEMENTED`. AGS uses its default backfill logic.
