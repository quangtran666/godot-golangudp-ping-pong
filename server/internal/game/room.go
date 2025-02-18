package game

import "sync"

type Room struct {
	ID      string
	Player1 *Player
	Player2 *Player
	mutex   sync.RWMutex // Mutex để thread-safe
}

func NewRoom(id string) *Room {
	return &Room{
		ID: id,
	}
}

func (room *Room) AddPlayer(player *Player) bool {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	if room.Player1 == nil {
		room.Player1 = player
		player.Room = room
		return true
	}

	if room.Player2 == nil {
		room.Player2 = player
		player.Room = room
		return true
	}

	return false
}

func (room *Room) RemovePlayer(player *Player) {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	if room.Player1 == player {
		room.Player1 = nil
	}

	if room.Player2 == player {
		room.Player2 = nil
	}

	player.Room = nil
}

func (room *Room) IsFull() bool {
	room.mutex.RLock()
	defer room.mutex.RUnlock()

	return room.Player1 != nil && room.Player2 != nil
}
