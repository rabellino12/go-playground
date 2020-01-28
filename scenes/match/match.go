package match

import "github.com/ByteArena/box2d"

import shapes "github.com/rabellino12/go-playground/scenes"

// Scene is the match scene handler
type Scene struct {
	MatchID       string
	World         box2d.B2World
	WorldScale    float64
	Platforms     []*box2d.B2Body
	shapesHandler *shapes.Handler
}

// Initialize starts a new match scene instance
func Initialize(matchID string) Scene {
	gravity := box2d.B2Vec2{X: 0, Y: 5}
	world := box2d.MakeB2World(gravity)
	shapesHandler := &shapes.Handler{
		World:      world,
		WorldScale: 30,
	}
	matchScene := Scene{
		MatchID:       matchID,
		World:         shapesHandler.World,
		WorldScale:    shapesHandler.WorldScale,
		Platforms:     make([]*box2d.B2Body, 0),
		shapesHandler: shapesHandler,
	}
	matchScene.initializeEnvironment()
	return matchScene
}

func (s *Scene) initializeEnvironment() {
	ground := s.shapesHandler.CreatePlatform(400, 578, 800, 64)
	plat1 := s.shapesHandler.CreatePlatform(600, 400, 400, 32)
	s.Platforms = []*box2d.B2Body{ground, plat1}
}
