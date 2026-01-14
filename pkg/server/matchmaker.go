// Copyright (c) 2022 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"matchmaking-function-grpc-plugin-server-go/pkg/common"
	"matchmaking-function-grpc-plugin-server-go/pkg/matchmaker"
)

// New returns a MatchMaker of the MatchLogic interface
func New() MatchLogic {
	return MatchMaker{}
}

// GetStatCodes returns the stat codes configured in the rules
func (b MatchMaker) GetStatCodes(scope *common.Scope, matchRules interface{}) []string {
	log := scope.Log.With("method", "MatchMaker.GetStatCodes")

	rule, ok := matchRules.(GameRules)
	if !ok {
		log.Error("unexpected game rule type", "type", fmt.Sprintf("%T", matchRules))

		return []string{}
	}

	// If no statistics configured, return empty
	if len(rule.Statistics.Statistics) == 0 {
		log.Info("no statistics configured, returning empty stat codes")

		return []string{}
	}

	log.Info("returning stat codes", "codes", rule.Statistics.Statistics)

	return rule.Statistics.Statistics
}

// ValidateTicket validates that the ticket has a valid selected stat
func (b MatchMaker) ValidateTicket(scope *common.Scope, matchTicket matchmaker.Ticket, matchRules interface{}) (bool, error) {
	log := scope.Log.With("method", "MatchMaker.ValidateTicket", "ticketID", matchTicket.TicketID)
	log.Info("validating ticket")
	log.Info("ticket/rules snapshot", "ticket", matchTicket, "rules", matchRules)

	rule, ok := matchRules.(GameRules)
	if !ok {
		log.Error("unexpected game rule type", "type", fmt.Sprintf("%T", matchRules))

		return false, status.Error(codes.Internal, "invalid game rules type")
	}

	// If no statistics configured, skip validation
	if len(rule.Statistics.Statistics) == 0 {
		log.Info("no statistics config, skipping validation")

		return true, nil
	}

	enrichedKey := rule.Statistics.GetEnrichedKey()

	// Validate each player has the enriched attribute
	for _, player := range matchTicket.Players {
		playerLog := log.With("playerID", player.PlayerID)

		if _, exists := player.Attributes[enrichedKey]; !exists {
			playerLog.Error("player missing enriched attribute", "key", enrichedKey)

			return false, status.Errorf(codes.InvalidArgument,
				"player %s: missing enriched attribute '%s'", player.PlayerID, enrichedKey)
		}

		playerLog.Info("player validation passed")
	}

	log.Info("ticket validation successful")

	return true, nil
}

// EnrichTicket extracts the selected stat value and adds it to ticket attributes
func (b MatchMaker) EnrichTicket(scope *common.Scope, matchTicket matchmaker.Ticket, ruleSet interface{}) (matchmaker.Ticket, error) {
	log := scope.Log.With("method", "MatchMaker.EnrichTicket", "ticketID", matchTicket.TicketID)
	log.Info("enriching ticket")
	log.Info("ticket/rules snapshot", "ticket", matchTicket, "rules", ruleSet)

	rule, ok := ruleSet.(GameRules)
	if !ok {
		log.Error("unexpected game rule type", "type", fmt.Sprintf("%T", ruleSet))

		return matchTicket, status.Error(codes.Internal, "invalid game rules type")
	}

	// If no statistics configured, skip enrichment
	if len(rule.Statistics.Statistics) == 0 {
		log.Info("no statistics config, skipping enrichment")

		return matchTicket, nil
	}

	enrichedKey := rule.Statistics.GetEnrichedKey()

	// For each player, set enriched attribute and remove configured stats
	for i, player := range matchTicket.Players {
		playerLog := log.With("playerID", player.PlayerID)

		// Get selected stat from ticket attributes using player ID
		selectedStatRaw := matchTicket.TicketAttributes[string(player.PlayerID)]
		selectedStat, _ := selectedStatRaw.(string)

		// Try to get and set the enriched value
		enriched := false
		valueRaw, exists := player.Attributes[selectedStat]

		if exists {
			var value float64

			switch v := valueRaw.(type) {
			case float64:
				value = v
				enriched = true
			case int:
				value = float64(v)
				enriched = true
			case int64:
				value = float64(v)
				enriched = true
			default:
				playerLog.Warn("unexpected stat value type", "type", fmt.Sprintf("%T", valueRaw))
			}

			if enriched {
				// Initialize player attributes if nil
				if matchTicket.Players[i].Attributes == nil {
					matchTicket.Players[i].Attributes = make(map[string]interface{})
				}

				matchTicket.Players[i].Attributes[enrichedKey] = value
				playerLog.Info("player enriched", "selectedStat", selectedStat, "value", value)
			}
		} else {
			playerLog.Warn("player missing selected stat", "stat", selectedStat)
		}

		// Always remove configured statistics from player attributes
		for _, stat := range rule.Statistics.Statistics {
			delete(matchTicket.Players[i].Attributes, stat)
		}
	}

	// Clean up player ID mappings - no longer needed after enrichment
	for _, player := range matchTicket.Players {
		delete(matchTicket.TicketAttributes, string(player.PlayerID))
	}

	log.Info("ticket enriched", "enrichedKey", enrichedKey)

	return matchTicket, nil
}

// RulesFromJSON returns the ruleset from the Game rules JSON
func (b MatchMaker) RulesFromJSON(scope *common.Scope, jsonRules string) (interface{}, error) {
	var ruleSet GameRules
	err := json.Unmarshal([]byte(jsonRules), &ruleSet)
	if err != nil {
		return nil, err
	}

	return ruleSet, nil
}

// MakeMatches returns nil to signal UNIMPLEMENTED - AGS will use default matching
func (b MatchMaker) MakeMatches(scope *common.Scope, ticketProvider TicketProvider, matchRules interface{}) <-chan matchmaker.Match {
	scope.Log.Info("MakeMatches not implemented, delegating to AGS default matching")

	return nil
}

// BackfillMatches returns nil to signal UNIMPLEMENTED - AGS will use default backfill
func (b MatchMaker) BackfillMatches(scope *common.Scope, ticketProvider TicketProvider, matchRules interface{}) <-chan matchmaker.BackfillProposal {
	scope.Log.Info("BackfillMatches not implemented, delegating to AGS default backfill")

	return nil
}
