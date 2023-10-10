package main

import (
	"fmt"
	"math"
	"sort"
)

type Team struct {
	Name   string
	Rating int
}

func isPowerOfTwo(x int) bool {
	return x != 0 && (x&(x-1)) == 0
}

func findOpponent(index int, used map[int]bool, teams []Team) int {
	for i := index + 1; i < len(teams); i++ {
		if !used[i] {
			return i
		}
	}
	return -1
}

func createFixture(teams []Team) error {
	// Urutkan tim berdasarkan rating dari tertinggi ke terendah
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Rating > teams[j].Rating
	})

	// Hitung jumlah tim dan putaran
	totalTeams := len(teams)
	rounds := int(math.Log2(float64(totalTeams)))

	// menentukan jenis format
	isNormal := isPowerOfTwo(totalTeams)
	if isNormal {
		fmt.Println("masuk dalam gugur normal")

		fixtures := make([]string, len(teams))
		used := make(map[int]bool)

		for i := 0; i < len(teams); i++ {
			if used[i] {
				continue
			}

			opponentIndex := findOpponent(i, used, teams)
			if opponentIndex == -1 {
				// Tim mendapat bye
				fixtures[i] = fmt.Sprintf("%s bye", teams[i].Name)
			} else {
				// Pasangkan tim dengan lawan
				fixtures[i] = fmt.Sprintf("%s x %s", teams[i].Name, teams[opponentIndex].Name)
				used[i] = true
				used[opponentIndex] = true
			}
		}

		// Tampilkan hasil fixture
		fmt.Println("Hasil Fixture:")
		for _, fixture := range fixtures {
			fmt.Println(fixture)
		}
		return nil
	}

	// Tentukan fixture berdasarkan jumlah putaran
	for round := 1; round <= rounds; round++ {
		matches := totalTeams / int(math.Pow(2, float64(round)))

		fmt.Printf("Putaran %d:\n", round)
		for i := 0; i < matches; i++ {
			team1 := teams[i*2].Name
			team2 := teams[i*2+1].Name

			fmt.Printf("%s vs %s\n", team1, team2)
		}

		fmt.Println()
	}

	// Cek apakah ada kontestan yang mendapatkan bye
	if totalTeams > int(math.Pow(2, float64(rounds))) {
		fmt.Println("Kontestan mendapatkan bye:")
		for i := totalTeams / 2; i < totalTeams; i++ {
			fmt.Println(teams[i].Name)
		}
	}

	return nil
}

func main() {
	teams := []Team{
		{"Tim A", 10},
		{"Tim B", 5},
		{"Tim C", 4},
		{"Tim D", 1},
		{"Tim E", 11},
		{"Tim F", 8},
		{"Tim G", 6},
		{"Tim H", 9},
	}

	createFixture(teams)

}
