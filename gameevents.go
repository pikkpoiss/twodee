package twodee

type GameEventType int

type GameEvent interface {
	geType() GameEventType
}

type GameEventCallback func(GameEvent)

type GameEventTypeObservers map[int]GameEventCallback

type GameEventHandler struct {
	gameEvents        chan GameEvent
	eventObservers    []GameEventTypeObservers
	nextObserverId    int
	nextGameEventType GameEventType
}

func NewGameEventHandler() (h *GameEventHandler) {
	h = &GameEventHandler{
		gameEvents:        make(chan GameEvent, 100),
		eventObservers:    make([]GameEventTypeObservers),
		nextObserverId:    0,
		nextGameEventType: 0,
	}
	return
}

func (h *GameEventHandler) Poll() {
	for e := range h.gameEvents {
		for _, observer := range h.eventObservers[e.geType()] {
			observer(e)
		}
	}
}

func (h *GameEventHandler) Enqueue(e GameEvent) {
	select {
	case h.gameEvents <- e:
		// Added to game events pool.
	default:
		// Game events pool too full; not being read quickly enough.
		// Drop the event on the floor.
		// TODO: Warn?
	}
}

func (h *GameEventHandler) RegisterNewEventType() (t GameEventType) {
	t = nextGameEventType
	h.eventObservers[t] = make(GameEventTypeObservers)
	h.nextObserverId++
	return
}

func (h *GameEventHandler) AddObserver(t GameEventType, c GameEventCallback) (id int) {
	id = h.nextObserverId
	h.eventObservers[t][id] = c
	h.nextObserverId++
	return
}

func (h *GameEventHandler) RemoveObserver(t GameEventType, id int) {
	delete(h.eventObservers[t], id)
}
