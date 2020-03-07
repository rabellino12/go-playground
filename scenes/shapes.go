package scenes

import (
	"github.com/ByteArena/box2d"
	ioclient "github.com/rabellino12/go-playground/ioclient/match"
)

// Handler is the shapes handler
type Handler struct {
	World      *box2d.B2World
	WorldScale float64
}

// CreatePlatform creates a platform fixture for the world declared on the handler
func (h *Handler) CreatePlatform(xPx float64, yPx float64, widthPx float64, heightPx float64) *box2d.B2Body {
	bodyDef := &box2d.B2BodyDef{
		Type:     0,
		Position: box2d.B2Vec2{X: xPx / h.WorldScale, Y: yPx / h.WorldScale},
	}
	platform := h.World.CreateBody(bodyDef)
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox((widthPx/h.WorldScale)/2, (heightPx/h.WorldScale)/2)
	platform.CreateFixture(&shape, 1)
	platform.GetFixtureList().SetFriction(0)
	platform.SetMassData(&box2d.B2MassData{Mass: 1, Center: box2d.B2Vec2{}, I: 1})
	return platform
}

// CreateEdge returns a box2d edge shaped body
func (h *Handler) CreateEdge(positionX float64, positionY float64, x1Px float64, y1Px float64, x2Px float64, y2Px float64) *box2d.B2Body {
	bodyDef := &box2d.B2BodyDef{
		Type:     0,
		Position: box2d.B2Vec2{X: positionX / h.WorldScale, Y: positionY / h.WorldScale},
	}
	edge := h.World.CreateBody(bodyDef)
	shape := box2d.MakeB2EdgeShape()
	shape.Set(box2d.B2Vec2{X: x1Px / h.WorldScale, Y: y1Px / h.WorldScale}, box2d.B2Vec2{X: x2Px / h.WorldScale, Y: y2Px / h.WorldScale})
	edge.CreateFixture(&shape, 1)
	edge.SetMassData(&box2d.B2MassData{Mass: 1, Center: box2d.B2Vec2{}, I: 1})
	return edge
}

// CreatePlayer returns a box fixture with player attributes
func (h *Handler) CreatePlayer(x float64, y float64) *box2d.B2Body {
	bodyDef := &box2d.B2BodyDef{
		Type:          2,
		FixedRotation: true,
		Position:      box2d.B2Vec2{X: x / h.WorldScale, Y: y / h.WorldScale},
	}
	player := h.World.CreateBody(bodyDef)
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox((28/h.WorldScale)/2, (48/h.WorldScale)/2)
	player.CreateFixture(&shape, 1)
	player.SetMassData(&box2d.B2MassData{Mass: 1, Center: box2d.B2Vec2{}, I: 1})
	return player
}

// MovePlayer handles player movement in some direction
func (h *Handler) MovePlayer(move ioclient.Move, player *box2d.B2Body, elapsed int64) *box2d.B2Body {
	vel := player.GetLinearVelocity()
	force := 0.0
	if move.Action == "up" {
		h.jump(player)
	}
	switch move.Action {
	case "left":
		if (vel.X * float64(elapsed)) > -5 {
			force = -50
		}
	case "right":
		if (vel.X * float64(elapsed)) < 5 {
			force = 50
		}
	default:
		if (vel.X * float64(elapsed)) != 0 {
			force = vel.X * -10
		}
	}
	player.ApplyForce(box2d.B2Vec2{X: force, Y: 0}, player.GetWorldCenter(), true)
	return player
}

func (h *Handler) jump(player *box2d.B2Body) {
	f := player.GetWorldVector(box2d.B2Vec2{X: 0, Y: -1})
	p := player.GetWorldPoint(box2d.B2Vec2{X: 0, Y: 0.1})
	player.ApplyLinearImpulse(f, p, true)
}
