package parser

import (
	"bufio"
	"fmt"
	"logger/internal/domain"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const worldKill = "<world>"
const game = "game"

type GameParser struct {
	GameID               int
	MatchName            string
	TotalKills           int
	ListOfPlayers        []string
	KillsByPlayerMap     map[string]int
	KillsByMeansMap      map[string]int
	GlobalPlayersRanking map[string]int
	GamesInfo            map[string]domain.Game
}

func NewGameParser() *GameParser {
	return &GameParser{
		KillsByPlayerMap:     make(map[string]int),
		KillsByMeansMap:      make(map[string]int),
		GlobalPlayersRanking: make(map[string]int),
		GamesInfo:            make(map[string]domain.Game),
	}
}

func (gp *GameParser) AddPlayer(player string) {
	var found bool
	for _, activePlayer := range gp.ListOfPlayers {
		if activePlayer == player {
			found = true
			break
		}
	}
	if !found && player != worldKill {
		gp.ListOfPlayers = append(gp.ListOfPlayers, player)
	}
}

func (gp *GameParser) ReadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrorWhenOpenFile, err)
	}
	defer func() {
		if fileErr := file.Close(); fileErr != nil {
			err = fmt.Errorf("%w: %w", domain.ErrorWhenCloseFile, fileErr)
		}
	}()

	r := bufio.NewReader(file)

	for {
		line, _, readLineErr := r.ReadLine()
		if readLineErr != nil {
			gp.generateGameInfo()
			break
		}

		if len(line) > 0 {
			if err := gp.parseInfo(string(line)); err != nil {
				return err
			}
		}
	}

	gp.generatePlayerRanking()

	return nil
}

func (gp *GameParser) parseInfo(logLine string) error {
	if strings.Contains(logLine, "InitGame") {
		switch gp.MatchName {
		case "":
			gp.GameID = 1
		default:
			gp.generateGameInfo()
			gp.GameID++
			gp.KillsByPlayerMap = map[string]int{}
			gp.KillsByMeansMap = map[string]int{}
			gp.TotalKills = 0
			gp.ListOfPlayers = nil
		}
		gp.MatchName = fmt.Sprintf("%s_%d", game, gp.GameID)
	}

	if strings.Contains(logLine, "ClientUserinfoChanged") {
		gp.validateAndAddPlayer(logLine)
	}

	if strings.Contains(logLine, "Kill") {
		gp.TotalKills++
		if err := gp.populateKillerAndMeanOfDeath(logLine); err != nil {
			return err
		}
	}

	return nil
}

func (gp *GameParser) populateKillerAndMeanOfDeath(logLine string) error {
	killInfo := regexp.MustCompile(`Kill: \d+ \d+ \d+: (.+) killed (.+) by (.+)`).
		FindStringSubmatch(logLine)
	killer := killInfo[1]
	victim := killInfo[2]
	meanOfDeath := killInfo[3]

	if len(killInfo) == 4 {
		switch killer {
		case worldKill:
			kills := gp.KillsByPlayerMap[victim]
			gp.KillsByPlayerMap[victim] = kills - 1
		default:
			kills := gp.KillsByPlayerMap[killer]
			gp.KillsByPlayerMap[killer] = kills + 1
		}
	}

	mean, err := domain.Parse(meanOfDeath)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrorMeanOfDeathParse, err)
	}

	gp.KillsByMeansMap[mean.String()] += 1

	return nil
}

func (gp *GameParser) SortMeansByKills() []domain.KillsByMeans {
	var killsSlice []domain.KillsByMeans
	for m, k := range gp.KillsByMeansMap {
		killsSlice = append(killsSlice, domain.KillsByMeans{MeansOfDeath: m, Kills: k})
	}

	if len(killsSlice) > 0 {
		sort.Slice(killsSlice, func(i, j int) bool {
			return killsSlice[i].Kills > killsSlice[j].Kills
		})
	}

	return killsSlice
}

func (gp *GameParser) SortPlayersByKills(k map[string]int) []domain.KillsByPlayer {
	var killsSlice []domain.KillsByPlayer
	for p, k := range k {
		killsSlice = append(killsSlice, domain.KillsByPlayer{Player: p, Kills: k})
	}

	if len(killsSlice) > 0 {
		sort.Slice(killsSlice, func(i, j int) bool {
			return killsSlice[i].Kills > killsSlice[j].Kills
		})
	}

	return killsSlice
}

func (gp *GameParser) SortGames() []domain.Game {
	var gamesSlice []domain.Game
	for _, g := range gp.GamesInfo {
		gamesSlice = append(gamesSlice, g)
	}

	if len(gamesSlice) > 0 {
		sort.Slice(gamesSlice, func(i, j int) bool {
			f, _ := strconv.Atoi(strings.Split(gamesSlice[i].Name, "_")[1])
			s, _ := strconv.Atoi(strings.Split(gamesSlice[j].Name, "_")[1])
			return f < s
		})
	}

	return gamesSlice
}

func (gp *GameParser) generateGameInfo() {
	gp.addMissingPlayersInKillList()

	gp.GamesInfo[gp.MatchName] = domain.Game{
		Name:         gp.MatchName,
		TotalKills:   gp.TotalKills,
		Players:      gp.ListOfPlayers,
		Kills:        gp.KillsByPlayerMap,
		KillsByMeans: gp.KillsByMeansMap,
	}
}

func (gp *GameParser) addMissingPlayersInKillList() {
	for _, activePlayer := range gp.ListOfPlayers {
		var found bool
		for player := range gp.KillsByPlayerMap {
			if activePlayer == player {
				found = true
				break
			}
		}

		if !found {
			gp.KillsByPlayerMap[activePlayer] = 0
		}
	}

}

func (gp *GameParser) generatePlayerRanking() {
	for _, g := range gp.GamesInfo {
		for player, kills := range g.Kills.(map[string]int) {
			totalKills := gp.GlobalPlayersRanking[player]
			gp.GlobalPlayersRanking[player] = totalKills + kills
		}
	}
}

func (gp *GameParser) validateAndAddPlayer(logLine string) {
	player := regexp.MustCompile(`n\\([^\\]+)`).
		FindStringSubmatch(logLine)[1]
	gp.AddPlayer(player)
}
