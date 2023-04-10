package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// CreateFileScanner open the file and creates a file pointer
func CreateFileScanner(
	filePath string,
) *bufio.Scanner {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return bufio.NewScanner(file)
}

// GetMazeCoordinatesFromFile returns the maze col and the maze row
func GetMazeCoordinatesFromFile(
	fileScanner *bufio.Scanner,
) (int, int) {
	fileScanner.Scan()
	res := strings.Split(fileScanner.Text(), ", ")
	mazeRow, err := strconv.Atoi(res[0])
	if err != nil {
		log.Fatal(err)
	}

	mazeCol, err := strconv.Atoi(res[1])
	if err != nil {
		log.Fatal(err)
	}

	return mazeRow, mazeCol
}

// GetMazeElementsFromFile returns the maze elements
func GetMazeElementsFromFile(
	mazeElements []string,
	mazeRow int,
	fileScanner *bufio.Scanner,
) []string {

	for i := 0; i < mazeRow; i++ {
		fileScanner.Scan()
		mazeElements = append(strings.Fields(fileScanner.Text()), mazeElements...)
	}

	return mazeElements
}

// GetMazeFromFile returns the maze col and row and also the maze elements
func GetMazeFromFile(
	filePath string,
) [][]string {
	var mazes [][]string
	var mazeCol int
	var mazeRow int
	var mazeElements []string
	fileScanner := CreateFileScanner(filePath)

	for fileScanner.Scan() {
		if fileScanner.Text() == "MAZE" {
			mazeElements = []string{}
			mazeRow, mazeCol = GetMazeCoordinatesFromFile(fileScanner)
			mazeElements = append(mazeElements, strconv.Itoa(mazeRow), strconv.Itoa(mazeCol))
		}
		mazeElements = GetMazeElementsFromFile(mazeElements, mazeRow, fileScanner)
		mazes = append(mazes, mazeElements)
	}

	return mazes
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
	for i, v := range maze {
		if v != element {
			mazeElementsRemovedChar = append(mazeElementsRemovedChar, v)
			mazeElementsIndex = append(mazeElementsIndex, i)
		}
		if ((i+1)%mazeCol) == 0 && i+1 >= mazeCol && i+1 < mazeRow*mazeCol {
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
	lastBadMove string,
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
			lastBadMove != "U" &&
			(len(mazeSolution) == 0 || len(mazeSolution) > 0 && mazeSolution[len(mazeSolution)-1] != "D"),
		mazeUpIndex
}

// IsExistsDownWay checks if there is down move for the maze pointer
func IsExistsDownWay(
	mazeIndex []int,
	mazeCol int,
	mazeSolution []string,
	mazeStartIndex int,
	lastBadMove string,
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

	return isExistsDownIndex &&
			lastBadMove != "D" &&
			(len(mazeSolution) == 0 || len(mazeSolution) > 0 && mazeSolution[len(mazeSolution)-1] != "U"),
		mazeDownIndex
}

// IsExistsRightWay checks if there is right move for the maze pointer
func IsExistsRightWay(
	mazeIndex []int,
	mazeSolution []string,
	mazeStartIndex int,
	mazeEndIndex int,
	lastBadMove string,
) bool {
	return mazeStartIndex < mazeEndIndex &&
		lastBadMove != "R" &&
		mazeIndex[mazeStartIndex+1]-mazeIndex[mazeStartIndex] == 1 &&
		(len(mazeSolution) == 0 || len(mazeSolution) > 0 && mazeSolution[len(mazeSolution)-1] != "L")
}

// IsExistsLeftWay checks if there is left move for the maze pointer
func IsExistsLeftWay(
	mazeIndex []int,
	mazeSolution []string,
	mazeStartIndex int,
	lastBadMove string,
) bool {
	return mazeIndex[mazeStartIndex]-mazeIndex[mazeStartIndex-1] == 1 &&
		lastBadMove != "L" &&
		(len(mazeSolution) == 0 || len(mazeSolution) > 0 && mazeSolution[len(mazeSolution)-1] != "R")
}

// RemoveLastElement remove last element of a string slice
func RemoveLastElement(
	mazeSolution []string,
) (string, []string) {
	lastElement := mazeSolution[len(mazeSolution)-1]
	removedLastElementSlice := mazeSolution[0 : len(mazeSolution)-1]
	return lastElement, removedLastElementSlice
}

// MazeSolution finds the maze solution with the help of other defined functions
func MazeSolution(
	mazeIndex []int,
	mazeCol int,
	mazeStartIndex int,
	mazeEndIndex int,
) []string {
	var lastMove string
	var lastBadMove string
	var mazeSolution []string

	for mazeStartIndex != mazeEndIndex {
		isExistsUpWay, mazeUpIndex := IsExistsUpWay(mazeIndex, mazeCol, mazeSolution, mazeStartIndex, lastBadMove)
		isExistsDownWay, mazeDownIndex := IsExistsDownWay(mazeIndex, mazeCol, mazeSolution, mazeStartIndex, lastBadMove)

		if isExistsUpWay && lastBadMove != "U" {
			mazeStartIndex = mazeUpIndex
			mazeSolution = append(mazeSolution, "U")
		} else if IsExistsRightWay(mazeIndex, mazeSolution, mazeStartIndex, mazeEndIndex, lastBadMove) {
			mazeStartIndex += 1
			mazeSolution = append(mazeSolution, "R")
		} else if IsExistsLeftWay(mazeIndex, mazeSolution, mazeStartIndex, lastBadMove) && lastBadMove != "L" {
			mazeStartIndex -= 1
			mazeSolution = append(mazeSolution, "L")
		} else if isExistsDownWay && lastBadMove != "D" {
			mazeStartIndex = mazeDownIndex
			mazeSolution = append(mazeSolution, "D")
		} else {
			lastMove, mazeSolution = RemoveLastElement(mazeSolution)
			switch lastMove {
			case "U":
				_, mazeDownIndex = IsExistsDownWay(mazeIndex, mazeCol, mazeSolution, mazeStartIndex, lastBadMove)
				mazeStartIndex = mazeDownIndex
			case "R":
				mazeStartIndex -= 1
			case "L":
				mazeStartIndex += 1
			case "D":
				_, mazeUpIndex = IsExistsUpWay(mazeIndex, mazeCol, mazeSolution, mazeStartIndex, lastBadMove)
				mazeStartIndex = mazeUpIndex
			}
		}

		if len(lastMove) == 0 {
			lastBadMove = ""
		} else {
			lastBadMove = lastMove
			lastMove = ""
		}
	}
	return mazeSolution
}

func main() {
	var mazesSlice [][]string
	var mazeCol int
	var mazeRow int

	mazesSlice = GetMazeFromFile("env.txt")

	for i := 0; i < len(mazesSlice); i++ {
		mazeCol, _ = strconv.Atoi(mazesSlice[i][len(mazesSlice[i])-1])
		mazeRow, _ = strconv.Atoi(mazesSlice[i][len(mazesSlice[i])-2])
		_, mazesSlice[i] = RemoveLastElement(mazesSlice[i])
		_, mazesSlice[i] = RemoveLastElement(mazesSlice[i])

		mazeRemovedPercent, mazeRemovedPercentIndex := RemoveSpecialCharacter(mazesSlice[i], mazeCol, mazeRow, "%")
		startIndex, startIndexErr := FindMazeStartIndex(mazeRemovedPercent, "S")
		endIndex, endIndexErr := FindMazeEndIndex(mazeRemovedPercent, "G")
		if startIndexErr != nil {
			log.Println(startIndexErr)
		}
		if endIndexErr != nil {
			log.Println(endIndexErr)
		}

		start := time.Now()
		result := MazeSolution(mazeRemovedPercentIndex, mazeCol, startIndex, endIndex)

		fmt.Printf("----------------------------- maze %d (%dx%d) --------------------------------\n", i+1, mazeRow, mazeCol)
		fmt.Printf("the maze solution: %v\n", result)
		fmt.Printf("program processing time: %v\n", time.Since(start))
		if i == len(mazesSlice)-1 {
			fmt.Println("-----------------------------------------------------------------------------")
		}
	}
}
