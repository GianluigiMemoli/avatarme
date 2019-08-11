package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/GianluigiMemoli/avatarme/generator"
	"os"
	"strconv"
)

func makeHash(usrInputString string) []byte {
	h := sha1.New()
	h.Write([]byte(usrInputString))
	hashed := h.Sum(nil)
	return hashed
}


func main(){

	args := os.Args
	if len(args[1:]) != 3 {
		fmt.Println("BAD INPUT\n" +
			"********************************************" +
			"\navatarme produces a squared image that is the graphical representation of an hash, so \n" +
			"run: avatarme <word to represent> <size of square side> <filename output>\n" +
			"\nside MUST be at least 128!\n\n"+
			"********************************************")
		return
	}
	wordToHash := args[1]
	sideStr := args[2]
	fileName := args[3]
	side, err := strconv.Atoi(sideStr)
	if err != nil || side < 128{
		fmt.Println(
			"BAD INPUT\n" +
			"****************************************\n"+
			"Provided a side value < 128 or non numeric!")
		return
	}


	icon := generator.New(side)
	errCreating := icon.Create(makeHash(wordToHash), fileName)
	if errCreating != nil {
		fmt.Printf("Error occurred: %s", err)
	}
	fmt.Printf("SUCCESS! \n Graphical Hash of %s is in file %s in dim of %d", wordToHash, fileName, side)
}


