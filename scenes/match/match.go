package scenes

import (
	"errors"

	"github.com/ByteArena/box2d"
	game "github.com/rabellino12/go-playground/db/collections"
	ioMatch "github.com/rabellino12/go-playground/ioclient/match"
	"github.com/rabellino12/go-playground/scenes"
)

// WorldScene is the match scene handler
type WorldScene struct {
	MatchID       string
	World         *box2d.B2World
	WorldScale    float64
	Platforms     []*box2d.B2Body
	Edges         []*box2d.B2Body
	Players       map[string]*box2d.B2Body
	shapesHandler *scenes.Handler
}

// MakeMatch starts a new match scene instance
func MakeMatch(gameObj game.Game) *WorldScene {
	gravity := box2d.B2Vec2{X: 0, Y: 5}
	world := box2d.MakeB2World(gravity)
	shapesHandler := &scenes.Handler{
		World:      &world,
		WorldScale: 30,
	}
	matchScene := WorldScene{
		MatchID:       gameObj.ID.Hex(),
		World:         shapesHandler.World,
		WorldScale:    shapesHandler.WorldScale,
		Platforms:     make([]*box2d.B2Body, 0),
		Edges:         make([]*box2d.B2Body, 0),
		Players:       make(map[string]*box2d.B2Body, 0),
		shapesHandler: shapesHandler,
	}
	matchScene.initializeEnvironment()
	for _, player := range gameObj.Players {
		userData := ioMatch.PlayerUserData{
			Player: player,
			Action: "stop",
		}
		matchScene.AddPlayer(player.Position.X, player.Position.Y, userData)
	}
	return &matchScene
}

func (s *WorldScene) initializeEnvironment() {
	ground := s.shapesHandler.CreatePlatform(400, 578, 800, 64)
	plat1 := s.shapesHandler.CreatePlatform(600, 400, 400, 32)
	s.Platforms = []*box2d.B2Body{ground, plat1}
	edgeLeft := s.shapesHandler.CreateEdge(0, 0, 0, 0, 0, 568)
	edgeRight := s.shapesHandler.CreateEdge(0, 0, 800, 0, 800, 568)
	s.Edges = []*box2d.B2Body{edgeLeft, edgeRight}
}

// AddPlayer creates a new player at the given
func (s *WorldScene) AddPlayer(x float64, y float64, userData ioMatch.PlayerUserData) *box2d.B2Body {
	player := s.shapesHandler.CreatePlayer(x, y)
	player.SetUserData(userData)
	s.Players[userData.ID] = player
	return player
}

// AddMove moves a player in the world and returns the new player state
func (s *WorldScene) AddMove(move ioMatch.Move) (*box2d.B2Body, error) {
	player := s.Players[move.UserID]
	if player == nil {
		return &box2d.B2Body{}, errors.New("Player not found")
	}
	userData := player.GetUserData().(ioMatch.PlayerUserData)
	userData.Jumping = move.Jumping
	userData.Action = move.Action
	player.SetUserData(userData)
	s.shapesHandler.MovePlayer(move, player)
	return player, nil
}

// GetSnapshot returns a world snapshot
func (s *WorldScene) GetSnapshot() map[string][]ioMatch.Move {
	moves := make(map[string][]ioMatch.Move, 0)
	for userID, player := range s.Players {
		userData := player.GetUserData().(ioMatch.PlayerUserData)
		position := player.GetPosition()
		x := position.X
		y := position.Y
		moves[userID] = append(moves[userID], ioMatch.Move{
			UserID: userID,
			Action: userData.Action,
			Position: game.Position{
				X: x,
				Y: y,
			},
			Jumping: userData.Jumping,
		})
	}
	return moves
}
