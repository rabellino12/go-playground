package scenes

import (
	"github.com/ByteArena/box2d"
	game "github.com/rabellino12/go-playground/db/collections"
	"github.com/rabellino12/go-playground/scenes"
)

// Scene is the match scene handler
type Scene struct {
	MatchID       string
	World         *box2d.B2World
	WorldScale    float64
	Platforms     []*box2d.B2Body
	Edges         []*box2d.B2Body
	Players       []*box2d.B2Body
	shapesHandler *scenes.Handler
}

// MakeMatch starts a new match scene instance
func MakeMatch(gameObj game.Game) Scene {
	gravity := box2d.B2Vec2{X: 0, Y: 5}
	world := box2d.MakeB2World(gravity)
	shapesHandler := &scenes.Handler{
		World:      &world,
		WorldScale: 30,
	}
	matchScene := Scene{
		MatchID:       gameObj.ID.Hex(),
		World:         shapesHandler.World,
		WorldScale:    shapesHandler.WorldScale,
		Platforms:     make([]*box2d.B2Body, 0),
		Edges:         make([]*box2d.B2Body, 0),
		Players:       make([]*box2d.B2Body, 0),
		shapesHandler: shapesHandler,
	}
	matchScene.initializeEnvironment()
	return matchScene
}

func (s *Scene) initializeEnvironment() {
	ground := s.shapesHandler.CreatePlatform(400, 578, 800, 64)
	plat1 := s.shapesHandler.CreatePlatform(600, 400, 400, 32)
	s.Platforms = []*box2d.B2Body{ground, plat1}
	edgeLeft := s.shapesHandler.CreateEdge(0, 0, 0, 0, 0, 568)
	edgeRight := s.shapesHandler.CreateEdge(0, 0, 800, 0, 800, 568)
	s.Edges = []*box2d.B2Body{edgeLeft, edgeRight}
}

// AddPlayer creates a new player at the given
func (s *Scene) AddPlayer(x float64, y float64, userData game.Player) *box2d.B2Body {
	player := s.shapesHandler.CreatePlayer(x, y)
	s.Players = append(s.Players, player)
	return player
}
