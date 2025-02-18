package game

import "net"

type Player struct {
	ID       string
	Addr     *net.UDPAddr // Địa chỉ UDP của player để gửi message
	Room     *Room
	Position float64
}

func NewPlayer(id string, addr *net.UDPAddr) *Player {
	return &Player{
		ID:       id,
		Addr:     addr,
		Position: 0,
	}
}
