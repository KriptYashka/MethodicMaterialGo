package main

import "fmt"

func addVote(results map[string]int, name string) {
	results[name]++
}

func winner(results map[string]int) string {
	var maxName string
	maxVotes := -1
	for name, votes := range results {
		if votes > maxVotes {
			maxVotes = votes
			maxName = name
		}
	}
	return maxName
}

func totalVotes(results map[string]int) int {
	total := 0
	for _, votes := range results {
		total += votes
	}
	return total
}

func main() {
	var n int
	fmt.Scan(&n)

	names := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&names[i])
	}

	results := make(map[string]int)
	for _, name := range names {
		addVote(results, name)
	}

	w := winner(results)
	fmt.Printf("Лучший сотрудник: %s, голосов: %d из %d\n", w, results[w], totalVotes(results))
}
