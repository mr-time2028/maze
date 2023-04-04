package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func GetMazeCoordinates() (mazeCol, mazeRow int) {
	col, row := 0, 0
	fmt.Print("Enter col and row of the maze: ")
	_, err := fmt.Scan(&col, &row)
	if err != nil {
		log.Fatal(err)
	}
	return col, row
}

// GetInputSlice get the maze elements from user
func GetInputSlice() []string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.Fields(scanner.Text())
}

// GetMazeElements stores the maze elements taken by GetInputSlice in a slice
func GetMazeElements(mazeRow int) []string {
	var mazeElements []string
	for i := 0; i < mazeRow; i++ {
		mazeElements = append(mazeElements, GetInputSlice()...)
	}
	return mazeElements
}

func main() {
	col, row := GetMazeCoordinates()
	fmt.Println(col, row)
	mazeElements := GetMazeElements(row)
	fmt.Println(mazeElements)
}
