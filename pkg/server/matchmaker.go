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

	// Validate each player in the ticket
	for _, player := range matchTicket.Players {
		playerLog := log.With("playerID", player.PlayerID)

		// Get selected stat from ticket attributes using player ID
		selectedStatRaw, exists := matchTicket.TicketAttributes[string(player.PlayerID)]
		if !exists {
			playerLog.Error("player missing selected stat mapping", "key", player.PlayerID)

			return false, status.Errorf(codes.InvalidArgument,
				"player %s missing required stat mapping", player.PlayerID)
		}

		selectedStat, ok := selectedStatRaw.(string)
		if !ok {
			playerLog.Error("selected stat is not a string", "value", selectedStatRaw)

			return false, status.Errorf(codes.InvalidArgument,
				"player %s: selected stat must be a string", player.PlayerID)
		}

		// Validate stat is in allowed list
		if !rule.Statistics.IsValidStat(selectedStat) {
			playerLog.Error("invalid stat selected",
				"stat", selectedStat,
				"allowed", rule.Statistics.Statistics)

			return false, status.Errorf(codes.InvalidArgument,
				"player %s: invalid stat '%s'", player.PlayerID, selectedStat)
		}

		// Check if player has the selected stat (unless default is configured)
		_, hasStat := player.Attributes[selectedStat]
		if !hasStat && rule.Statistics.DefaultValue == 0 {
			playerLog.Error("player missing selected stat value",
				"stat", selectedStat)

			return false, status.Errorf(codes.InvalidArgument,
				"player %s: missing stat value for '%s'", player.PlayerID, selectedStat)
		}

		playerLog.Info("player validation passed", "selectedStat", selectedStat)
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

	// Initialize ticket attributes if nil
	if matchTicket.TicketAttributes == nil {
		matchTicket.TicketAttributes = make(map[string]interface{})
	}

	// For party tickets, aggregate stat values (average)
	var totalValue float64
	var selectedStat string
	playerCount := 0

	for i, player := range matchTicket.Players {
		playerLog := log.With("playerID", player.PlayerID)

		// Get selected stat from ticket attributes using player ID
		selectedStatRaw := matchTicket.TicketAttributes[string(player.PlayerID)]
		currentStat, _ := selectedStatRaw.(string) // Already validated

		if i == 0 {
			selectedStat = currentStat
		}

		// Get value for the selected stat
		valueRaw, exists := player.Attributes[currentStat]

		var value float64
		if exists {
			switch v := valueRaw.(type) {
			case float64:
				value = v
			case int:
				value = float64(v)
			case int64:
				value = float64(v)
			default:
				playerLog.Warn("unexpected stat value type, using default", "type", fmt.Sprintf("%T", valueRaw))
				value = rule.Statistics.DefaultValue
			}
		} else {
			value = rule.Statistics.DefaultValue
		}

		totalValue += value
		playerCount++
		playerLog.Info("extracted stat value",
			"stat", currentStat,
			"value", value)
	}

	// Calculate average value for the ticket
	var averageValue float64
	if playerCount > 0 {
		averageValue = totalValue / float64(playerCount)
	}

	// Add enriched attributes to ticket
	enrichedKey := rule.Statistics.GetEnrichedKey()
	matchTicket.TicketAttributes[enrichedKey] = averageValue

	log.Info("ticket enriched",
		enrichedKey, averageValue,
		"selected_stat", selectedStat)

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
