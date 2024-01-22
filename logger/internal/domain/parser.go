package domain

type (
	Game struct {
		Name         string   `json:"name"`
		TotalKills   int      `json:"total_kills"`
		Players      []string `json:"players"`
		Kills        any      `json:"kills"`
		KillsByMeans any      `json:"kills_by_means"`
	}

	KillsByPlayer struct {
		Player string
		Kills  int
	}

	KillsByMeans struct {
		MeansOfDeath string
		Kills        int
	}
)
