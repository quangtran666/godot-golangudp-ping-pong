package protocol

import "encoding/json"

const (
	MsgTypeJoinGame    = "join_game"
	MsgTypeGameStart   = "game_start"
	MsgTypeGameEnd     = "game_end"
	MsgTypePlayerInput = "player_input"
	MsgTypeBallUpdate  = "ball_update"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // json.RawMessage is a byte slice that holds a JSON encoded value
}

// JoinGameMessage when player wants to join a game
type JoinGameMessage struct {
	PlayerID string `json:"player_id"`
}

// GameStartMessage sent when enough players have joined
type GameStartMessage struct {
	Player1ID string `json:"player1_id"`
	Player2ID string `json:"player2_id"`
}

type PlayerInputMessage struct {
	PlayerID string  `json:"player_id"`
	Position float64 `json:"position"`
}

type BallUpdateMessage struct {
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	DirectionX float64 `json:"direction_x"`
	DirectionY float64 `json:"direction_y"`
	Speed      float64 `json:"speed"`
}

type GameEndMessage struct {
	WinnerID string `json:"winner_id"`
}
