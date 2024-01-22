package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"logger/internal/parser"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var c map[string]func()

func init() {
	c = make(map[string]func())
	c["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	c["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := c[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func main() {
	zap.L().Info("Welcome to log parser, do you wanna see the log of")
	processor := parser.NewGameParser()

	if err := processor.ReadFile("assets/games.log"); err != nil {
		zap.L().Fatal("ReadFile", zap.Error(err))
	}

	CallClear()

	if len(os.Args) == 1 {
		_, _ = options(processor, "4")
		return
	}

	for {
		printMenu()

		exit, err := options(processor, getUserChoice())
		if exit {
			break
		}
		if err != nil {
			zap.L().Fatal("getUserChoice", zap.Error(err))
		}
	}
}

func options(p *parser.GameParser, choice string) (bool, error) {
	switch {
	case choice == "1":
		CallClear()
		fmt.Println("Exiting the program.")
		return true, nil
	case choice == "2":
		CallClear()
		players, err := json.MarshalIndent(p.GlobalPlayersRanking, "", "  ")
		if err != nil {
			return false, fmt.Errorf("json.MarshalIndent choice 2: %w", err)
		}
		fmt.Println(string(players))
	case choice == "3":
		CallClear()
		fmt.Println("players_ranking: {")

		for _, pr := range p.SortPlayersByKills(p.GlobalPlayersRanking) {
			if _, err := fmt.Printf("  \"%s\": %d\n", pr.Player, pr.Kills); err != nil {
				return false, fmt.Errorf("pretty.Println chpoce 3: %w", err)
			}
		}

		fmt.Println("}")
	case choice == "4":
		CallClear()
		games, err := json.MarshalIndent(p.SortGames(), "", "  ")
		if err != nil {
			fmt.Println("Error:", err)
			return false, fmt.Errorf("json.MarshalIndent choice 5: %w", err)
		}

		fmt.Println(string(games))
	case strings.Contains(choice, "game_"):
		CallClear()
		games, err := json.MarshalIndent(p.GamesInfo[choice], "", "  ")
		if err != nil {
			return false, fmt.Errorf("json.MarshalIndent choice get game: %w", err)
		}

		fmt.Println(string(games))
	default:
		CallClear()
		return false, errors.New("invalid choice. Please enter a valid option")
	}
	return false, nil
}

func printMenu() {
	fmt.Println("1. Exit")
	fmt.Println("2. View Player Ranking")
	fmt.Println("3. View sorted Player Ranking")
	fmt.Println("4. View Games Information")
	fmt.Println("game_*. View a specific game sorted Information * = {MIN: 1, MAX: 21}. Example game_1")
	fmt.Print("Enter your choice: ")
}

func getUserChoice() string {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		return scanner.Text()
	}
}
