package shapes

import "github.com/ByteArena/box2d"

// Handler is the shapes handler
type Handler struct {
	World      box2d.B2World
	WorldScale float64
}

// CreatePlatform creates a platform fixture for the world declared on the handler
func (h *Handler) CreatePlatform(xPx float64, yPx float64, widthPx float64, heightPx float64) *box2d.B2Body {
	bodyDef := &box2d.B2BodyDef{
		Type:     1,
		Position: box2d.B2Vec2{X: xPx / h.WorldScale, Y: yPx / h.WorldScale},
	}
	platform := h.World.CreateBody(bodyDef)
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox((widthPx/h.WorldScale)/2, (heightPx/h.WorldScale)/2)
	platform.CreateFixture(&shape, 1)
	platform.SetMassData(&box2d.B2MassData{Mass: 1, Center: box2d.B2Vec2{}, I: 1})
	return platform
}
