// Copyright (c) 2022 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"matchmaking-function-grpc-plugin-server-go/pkg/common"
	"matchmaking-function-grpc-plugin-server-go/pkg/matchmaker"
	matchfunctiongrpc "matchmaking-function-grpc-plugin-server-go/pkg/pb"
)

// MatchFunctionServer is for the handler (upper level of match logic)
type MatchFunctionServer struct {
	matchfunctiongrpc.UnimplementedMatchFunctionServer
	MM MatchLogic
}

// matchTicketProvider contains the go channel of matchmaker tickets needed for making matches
type matchTicketProvider struct {
	channelTickets         chan matchmaker.Ticket
	channelBackfillTickets chan matchmaker.BackfillTicket
}

// GetTickets will return the go channel of tickets from the matchTicketProvider
func (m matchTicketProvider) GetTickets() chan matchmaker.Ticket {
	return m.channelTickets
}

// GetBackfillTickets returns the go channel of backfill tickets
func (m matchTicketProvider) GetBackfillTickets() chan matchmaker.BackfillTicket {
	return m.channelBackfillTickets
}

// GetStatCodes uses the assigned MatchMaker to get the stat codes of the ruleset
func (m *MatchFunctionServer) GetStatCodes(ctx context.Context, req *matchfunctiongrpc.GetStatCodesRequest) (*matchfunctiongrpc.StatCodesResponse, error) {
	scope := common.ChildScopeFromRemoteScope(ctx, "MatchFunctionServer.GetStatCodes")
	defer scope.Finish()

	rules, err := m.MM.RulesFromJSON(scope, req.Rules.Json)
	if err != nil {
		scope.Log.Error("could not get rules from json", "error", err)

		return nil, err
	}

	statCodes := m.MM.GetStatCodes(scope, rules)

	return &matchfunctiongrpc.StatCodesResponse{Codes: statCodes}, nil
}

// ValidateTicket uses the assigned MatchMaker to validate the ticket
func (m *MatchFunctionServer) ValidateTicket(ctx context.Context, req *matchfunctiongrpc.ValidateTicketRequest) (*matchfunctiongrpc.ValidateTicketResponse, error) {
	scope := common.ChildScopeFromRemoteScope(ctx, "MatchFunctionServer.ValidateTicket")
	defer scope.Finish()

	scope.Log.Info("validating ticket")

	rules, err := m.MM.RulesFromJSON(scope, req.Rules.Json)
	if err != nil {
		scope.Log.Error("could not get rules from json", "error", err)

		return nil, err
	}

	matchTicket := matchfunctiongrpc.ProtoTicketToMatchfunctionTicket(req.Ticket)

	validTicket, err := m.MM.ValidateTicket(scope, matchTicket, rules)

	return &matchfunctiongrpc.ValidateTicketResponse{ValidTicket: validTicket}, err
}

// EnrichTicket uses the assigned MatchMaker to enrich the ticket
func (m *MatchFunctionServer) EnrichTicket(ctx context.Context, req *matchfunctiongrpc.EnrichTicketRequest) (*matchfunctiongrpc.EnrichTicketResponse, error) {
	scope := common.ChildScopeFromRemoteScope(ctx, "MatchFunctionServer.EnrichTicket")
	defer scope.Finish()

	scope.Log.Info("enriching ticket")

	rules, err := m.MM.RulesFromJSON(scope, req.Rules.Json)
	if err != nil {
		scope.Log.Error("could not get rules from json", "error", err)

		return nil, err
	}

	matchTicket := matchfunctiongrpc.ProtoTicketToMatchfunctionTicket(req.Ticket)
	enrichedTicket, err := m.MM.EnrichTicket(scope, matchTicket, rules)
	if err != nil {
		return nil, err
	}
	newTicket := matchfunctiongrpc.MatchfunctionTicketToProtoTicket(enrichedTicket)

	response := &matchfunctiongrpc.EnrichTicketResponse{Ticket: newTicket}
	scope.Log.Info("ticket enriched successfully")

	return response, nil
}

// MakeMatches returns UNIMPLEMENTED to let AGS use default matching algorithm
func (m *MatchFunctionServer) MakeMatches(server matchfunctiongrpc.MatchFunction_MakeMatchesServer) error {
	scope := common.ChildScopeFromRemoteScope(context.Background(), "MatchFunctionServer.MakeMatches")
	defer scope.Finish()

	scope.Log.Info("MakeMatches returning UNIMPLEMENTED - using AGS default matching")

	return status.Error(codes.Unimplemented, "MakeMatches not implemented - using AGS default matching")
}

// BackfillMatches returns UNIMPLEMENTED to let AGS use default backfill logic
func (m *MatchFunctionServer) BackfillMatches(server matchfunctiongrpc.MatchFunction_BackfillMatchesServer) error {
	scope := common.ChildScopeFromRemoteScope(context.Background(), "MatchFunctionServer.BackfillMatches")
	defer scope.Finish()

	scope.Log.Info("BackfillMatches returning UNIMPLEMENTED - using AGS default backfill")

	return status.Error(codes.Unimplemented, "BackfillMatches not implemented - using AGS default backfill")
}
