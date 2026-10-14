package server

import (
	"context"
	"fmt"

	"log/slog"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
)

func NewPlayerMailServiceHandler() craig_starsv1connect.PlayerMailServiceHandler {
	return &playerMailService{}
}

type playerMailService struct{}

// GetPlayerMail returns the mail the current player has written and the mail delivered to them
func (s *playerMailService) GetPlayerMail(ctx context.Context, req *connect.Request[craig_starsv1.GetPlayerMailRequest]) (*connect.Response[craig_starsv1.GetPlayerMailResponse], error) {
	dbClient := contextDb(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	mail, err := dbClient.GetPlayerMailForPlayer(ctx, game.ID, gamePlayer.Num, game.Year)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get mail from db"))
	}

	return connect.NewResponse(&craig_starsv1.GetPlayerMailResponse{
		Mail: converter.C.ConvertCSPlayerMails(mail),
	}), nil
}

// SavePlayerMail writes mail to be delivered with the next turn, or changes mail that hasn't been delivered
func (s *playerMailService) SavePlayerMail(ctx context.Context, req *connect.Request[craig_starsv1.SavePlayerMailRequest]) (*connect.Response[craig_starsv1.SavePlayerMailResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	if req.Msg.Mail == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing mail"))
	}
	// mail written while the turn is generating would miss its delivery
	if game.State != cs.GameStateWaitingForPlayers {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("game is not waiting for players"))
	}

	mail := converter.C.ConvertPlayerMail(req.Msg.Mail)
	mail.GameID = game.ID
	mail.Year = game.Year
	mail.FromPlayerNum = gamePlayer.Num

	if err := mail.Validate(len(game.Players)); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if mail.ID != 0 {
		existing, err := s.getPendingMail(ctx, mail.ID)
		if err != nil {
			return nil, err
		}
		mail.CreatedAt = existing.CreatedAt
	} else {
		pending, err := dbClient.GetPendingPlayerMailForGame(ctx, game.ID, game.Year)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get mail from db"))
		}
		count := 0
		for _, m := range pending {
			if m.FromPlayerNum == gamePlayer.Num {
				count++
			}
		}
		if count >= cs.PlayerMailMaxPerYear {
			return nil, connect.NewError(connect.CodeResourceExhausted, fmt.Errorf("a player can send at most %d messages a year", cs.PlayerMailMaxPerYear))
		}
	}

	if err := dbWriteClient.SavePlayerMail(ctx, mail); err != nil {
		slog.Error("save player mail in database", slog.Any("error", err), slog.Int64("ID", mail.ID))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.SavePlayerMailResponse{
		Mail: converter.C.ConvertCSPlayerMail(mail),
	}), nil
}

// DeletePlayerMail deletes mail the current player wrote that hasn't been delivered
func (s *playerMailService) DeletePlayerMail(ctx context.Context, req *connect.Request[craig_starsv1.DeletePlayerMailRequest]) (*connect.Response[craig_starsv1.DeletePlayerMailResponse], error) {
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)

	if game.State != cs.GameStateWaitingForPlayers {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("game is not waiting for players"))
	}

	mail, err := s.getPendingMail(ctx, req.Msg.MailId)
	if err != nil {
		return nil, err
	}

	if err := dbWriteClient.DeletePlayerMail(ctx, mail.ID); err != nil {
		slog.Error("delete player mail from database", slog.Any("error", err), slog.Int64("ID", mail.ID))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.DeletePlayerMailResponse{}), nil
}

// getPendingMail returns mail the current player wrote that they can still change
func (s *playerMailService) getPendingMail(ctx context.Context, id int64) (*cs.PlayerMail, error) {
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	mail, err := contextDb(ctx).GetPlayerMail(ctx, id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get mail from db"))
	}
	if mail == nil || mail.GameID != game.ID || mail.FromPlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("mail not found"))
	}
	if mail.Delivered(game.Year) {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("mail was already delivered"))
	}
	return mail, nil
}
