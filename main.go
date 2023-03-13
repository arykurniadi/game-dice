package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Dice struct {
	TopSideValue int
}

func (d *Dice) Roll() *Dice {
	min := 1
	max := 7

	rand.Seed(time.Now().UnixNano())
	d.TopSideValue = rand.Intn(max-min) + min

	return d
}

func (d *Dice) GetTopSideValue() int {
	return d.TopSideValue
}

type Player struct {
	Name     string
	Position int
	Point    int
	Dices    []Dice
}

func NewPlayer(numberOfDice int, position int, name string) *Player {
	p := &Player{
		Name:     name,
		Position: position,
		Point:    0,
	}

	for i := 0; i < numberOfDice; i++ {
		p.Dices = append(p.Dices, Dice{})
	}

	return p
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) GetPosition() int {
	return p.Position
}

func (p *Player) AddPoint(point int) {
	p.Point += point
}

func (p *Player) GetPoint() int {
	return p.Point
}

func (p *Player) InsertDice(dice Dice) {
	p.Dices = append(p.Dices, dice)
}

func (p *Player) RemoveDice(key int) {
	dices := make([]Dice, 0)

	for _, dice := range p.Dices {
		if dice.GetTopSideValue() != key {
			dices = append(dices, dice)
		}
	}

	p.Dices = dices
}

func (p *Player) Play() {
	for i := range p.Dices {
		p.Dices[i].Roll()
	}
}

//=====================================
type Game struct {
	Round                 int
	NumberOfPlayer        int
	NumberOfDicePerPlayer int
	Players               []*Player
}

func NewGame(numberOfPlayer int, numberOfDicePerPlayer int) *Game {
	game := &Game{
		Round:                 0,
		NumberOfPlayer:        numberOfPlayer,
		NumberOfDicePerPlayer: numberOfDicePerPlayer,
	}

	for i := 0; i < game.NumberOfPlayer; i++ {
		game.Players = append(game.Players, NewPlayer(game.NumberOfDicePerPlayer, i, string(65+i)))
	}

	return game
}

func (g *Game) DisplayRound() *Game {
	fmt.Printf("===== Giliran %v =====\n", g.Round)
	return g
}

func (g *Game) DisplayTopSideDice(title string) {
	fmt.Printf("%v \n", title)

	for _, player := range g.Players {
		fmt.Printf("Pemain #%v: ", player.GetName())

		diceTopSide := ""
		for _, dice := range player.Dices {
			diceTopSide += fmt.Sprintf("%v, ", dice.GetTopSideValue())
		}

		if len(diceTopSide) > 0 {
			fmt.Printf("%v\n", diceTopSide[:len(diceTopSide)-2])
		} else {
			fmt.Printf("\n")
		}
	}
}

func (g *Game) GetWinner() *Player {
	var (
		winner    *Player
		highscore int = 0
	)

	for _, player := range g.Players {
		if player.GetPoint() >= highscore {
			highscore = player.GetPoint()
			winner = player
		}
	}

	return winner
}

func (g *Game) DisplayWinner(player *Player) {
	fmt.Printf("\nPemenang\n")
	fmt.Printf("Pemain %v \n", player.GetName())
}

func (g *Game) Start() {
	fmt.Printf("Pemain = %v, Dadu = %v\n\r\n", g.NumberOfPlayer, g.NumberOfDicePerPlayer)

	for {
		g.Round++

		for _, player := range g.Players {
			player.Play()
		}

		g.DisplayRound().DisplayTopSideDice("Lempar Dadu")

		diceCarryForward := make(map[int][]Dice, 0)

		for playerIndex, player := range g.Players {
			tempDiceArr := []Dice{}
			for _, dice := range player.Dices {
				if dice.GetTopSideValue() == 6 {
					player.AddPoint(1)
					player.RemoveDice(dice.GetTopSideValue())
				}

				if dice.GetTopSideValue() == 1 {
					if player.GetPosition() == (g.NumberOfPlayer - 1) {
						g.Players[0].InsertDice(dice)
						player.RemoveDice(1)
					} else {
						tempDiceArr = append(tempDiceArr, dice)
						player.RemoveDice(1)
					}
				}
			}

			diceCarryForward[playerIndex+1] = tempDiceArr

			if _, ok := diceCarryForward[playerIndex]; ok && len(diceCarryForward[playerIndex]) > 0 {
				for _, dice := range diceCarryForward[playerIndex] {
					player.InsertDice(dice)
				}

				diceCarryForward = make(map[int][]Dice)
			}
		}

		fmt.Printf("\n")
		g.DisplayTopSideDice("Setelah Evaluasi")

		playerHasDice := g.NumberOfPlayer

		for _, player := range g.Players {
			if len(player.Dices) <= 0 {
				playerHasDice--
			}
		}

		if playerHasDice == 1 {
			g.DisplayWinner(g.GetWinner())
			break
		}
	}
}

func main() {
	numberPlayer := 3
	numberOfDice := 4

	game := NewGame(numberPlayer, numberOfDice)
	game.Start()
}
