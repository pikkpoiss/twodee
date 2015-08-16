package twodee

type GameEventType int

type GETyper interface {
	GEType() GameEventType
}

type BasicGameEvent struct {
	geType GameEventType
}

func (e *BasicGameEvent) GEType() GameEventType {
	return e.geType
}

func NewBasicGameEvent(t GameEventType) *BasicGameEvent {
	return &BasicGameEvent{
		geType: t,
	}
}

type GameEventCallback func(GETyper)

type GameEventTypeObservers map[int]GameEventCallback

type GameEventHandler struct {
	gameEvents     chan GETyper
	eventObservers []GameEventTypeObservers
	nextObserverId int
}

func NewGameEventHandler(numGameEventTypes int) (h *GameEventHandler) {
	h = &GameEventHandler{
		gameEvents:     make(chan GETyper, 100),
		eventObservers: make([]GameEventTypeObservers, numGameEventTypes),
		nextObserverId: 0,
	}
	return
}

func (h *GameEventHandler) Poll() {
	var (
		e    GETyper
		loop = true
	)
	for loop {
		select {
		case e = <-h.gameEvents:
			for _, observer := range h.eventObservers[e.GEType()] {
				observer(e)
			}
		default:
			loop = false
		}
	}
}

func (h *GameEventHandler) Enqueue(e GETyper) {
	select {
	case h.gameEvents <- e:
		// Added to game events pool.
	default:
		// Game events pool too full; not being read quickly enough.
		// Drop the event on the floor.
		// TODO: Warn?
	}
}

func (h *GameEventHandler) AddObserver(t GameEventType, c GameEventCallback) (id int) {
	if h.eventObservers[t] == nil {
		h.eventObservers[t] = make(GameEventTypeObservers)
	}
	id = h.nextObserverId
	h.eventObservers[t][id] = c
	h.nextObserverId++
	return
}

func (h *GameEventHandler) RemoveObserver(t GameEventType, id int) {
	delete(h.eventObservers[t], id)
}
