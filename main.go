package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// GetMazeCoordinates get col and row of the maze from user
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
func GetMazeElements(
	mazeRow int,
) []string {
	var mazeElements []string
	for i := 0; i < mazeRow; i++ {
		mazeElements = append(GetInputSlice(), mazeElements...)
	}
	return mazeElements
}

// RemoveSpecialCharacter removes all barriers, in this case "%" sign
func RemoveSpecialCharacter(
	maze []string,
	mazeCol int,
	mazeRow int,
	element string,
) ([]string, []int) {
	var mazeElementsIndex []int
	var mazeElementsRemovedChar []string
	for j, v := range maze {
		if v != element {
			mazeElementsRemovedChar = append(mazeElementsRemovedChar, v)
			mazeElementsIndex = append(mazeElementsIndex, j)
		}
		if ((j+1)%mazeCol) == 0 && j+1 >= mazeCol && j+1 < mazeRow*mazeCol {
			mazeElementsRemovedChar = append(mazeElementsRemovedChar, "up")
			mazeElementsIndex = append(mazeElementsIndex, -1)
		}
	}

	return mazeElementsRemovedChar, mazeElementsIndex
}

// FindMazeStartIndex finds index of first point that maze pointer start moving from that, in this case "S"
func FindMazeStartIndex(
	maze []string,
	startChr string,
) (int, error) {
	for i, v := range maze {
		if v == startChr {
			return i, nil
		}
	}
	return -1, errors.New("there is no start point")
}

// FindMazeEndIndex finds index of end point that maze pointer stop moving when catch that, in this case "G"
func FindMazeEndIndex(
	maze []string,
	endChar string,
) (int, error) {
	for i, v := range maze {
		if v == endChar {
			return i, nil
		}
	}
	return -1, errors.New("there is no end point")
}

// IsExistsUpWay checks if there is up move for the maze pointer
func IsExistsUpWay(
	mazeIndex []int,
	mazeCol int,
	mazeSolution []string,
	mazeStartIndex int,
) (bool, int) {
	isExistsUpIndex := false
	mazeUpIndex := 0
	for i := mazeStartIndex; i < len(mazeIndex); i++ {
		if mazeIndex[i] == mazeIndex[mazeStartIndex]+mazeCol {
			isExistsUpIndex = true
			mazeUpIndex = i
			break
		}
	}
	return isExistsUpIndex &&
			(len(mazeSolution) == 0 || len(mazeSolution) > 0 && mazeSolution[len(mazeSolution)-1] != "D"),
		mazeUpIndex
}

// IsExistsDownWay checks if there is down move for the maze pointer
func IsExistsDownWay(
	mazeIndex []int,
	mazeCol int,
	mazeStartIndex int,
) (bool, int) {
	isExistsDownIndex := false
	mazeDownIndex := 0
	for i := mazeStartIndex; i >= 0; i-- {
		if mazeIndex[i] == mazeIndex[mazeStartIndex]-mazeCol {
			isExistsDownIndex = true
			mazeDownIndex = i
			break
		}
	}

	return isExistsDownIndex, mazeDownIndex
}

// IsExistsRightWay checks if there is right move for the maze pointer
func IsExistsRightWay(
	mazeIndex []int,
	mazeSolution []string,
	mazeStartIndex int,
	mazeEndIndex int,
) bool {
	return mazeStartIndex < mazeEndIndex &&
		mazeIndex[mazeStartIndex+1]-mazeIndex[mazeStartIndex] == 1 &&
		(len(mazeSolution) == 0 || len(mazeSolution) > 0 && mazeSolution[len(mazeSolution)-1] != "L")
}

// IsExistsLeftWay checks if there is left move for the maze pointer
func IsExistsLeftWay(
	mazeIndex []int,
	mazeStartIndex int,
) bool {
	return mazeIndex[mazeStartIndex]-mazeIndex[mazeStartIndex-1] == 1
}

// MazeSolution finds the maze solution with the help of other defined functions
func MazeSolution(
	mazeIndex []int,
	mazeCol int,
	mazeStartIndex int,
	mazeEndIndex int,
) []string {
	var mazeSolution []string
	for mazeStartIndex != mazeEndIndex {
		isExistsUpWay, mazeUpIndex := IsExistsUpWay(mazeIndex, mazeCol, mazeSolution, mazeStartIndex)
		if isExistsUpWay {
			mazeStartIndex = mazeUpIndex
			mazeSolution = append(mazeSolution, "U")
		} else if IsExistsRightWay(mazeIndex, mazeSolution, mazeStartIndex, mazeEndIndex) {
			mazeStartIndex += 1
			mazeSolution = append(mazeSolution, "R")
		} else if IsExistsLeftWay(mazeIndex, mazeStartIndex) {
			mazeStartIndex -= 1
			mazeSolution = append(mazeSolution, "L")
		} else {
			isExistsDownWay, mazeDownIndex := IsExistsDownWay(mazeIndex, mazeCol, mazeStartIndex)
			if isExistsDownWay {
				mazeStartIndex = mazeDownIndex
				mazeSolution = append(mazeSolution, "D")
			}
		}
	}
	return mazeSolution
}

func main() {
	col, row := GetMazeCoordinates()
	mazeSlice := GetMazeElements(row)
	mazeRemovedPercent, mazeRemovedPercentIndex := RemoveSpecialCharacter(mazeSlice, col, row, "%")

	startIndex, startIndexErr := FindMazeStartIndex(mazeRemovedPercent, "S")
	endIndex, endIndexErr := FindMazeEndIndex(mazeRemovedPercent, "G")
	if startIndexErr != nil {
		log.Println(startIndexErr)
	}
	if endIndexErr != nil {
		log.Println(endIndexErr)
	}

	start := time.Now()
	result := MazeSolution(mazeRemovedPercentIndex, col, startIndex, endIndex)

	fmt.Printf("the maze solution: %v\n", result)
	fmt.Printf("program processing time: %v\n", time.Since(start))
}
