package board

type Board struct {
	BState map[uint8]string // Board state as a simple map
}

func NewBoard() *Board {
	return &Board{
		BState: map[uint8]string{
			// White pieces
			1:  "wR", 2: "wN", 3: "wB", 4: "wQ", 5: "wK", 6: "wB", 7: "wN", 8: "wR",
			9:  "wP", 10: "wP", 11: "wP", 12: "wP", 13: "wP", 14: "wP", 15: "wP", 16: "wP",
			// Black pieces
			57: "bR", 58: "bN", 59: "bB", 60: "bQ", 61: "bK", 62: "bB", 63: "bN", 64: "bR",
			49: "bP", 50: "bP", 51: "bP", 52: "bP", 53: "bP", 54: "bP", 55: "bP", 56: "bP",
		},
	}
}

