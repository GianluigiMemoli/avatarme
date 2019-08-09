package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/GianluigiMemoli/avatarme/generator"
	"image/color"
	"image/png"
	"os"
)

func main(){

	toBeHashed := "Sesso anale"
	h := sha1.New()
	h.Write([]byte(toBeHashed))
	hashed := h.Sum(nil)
	id := generator.New(256, color.RGBA{hashed[0], hashed[1], hashed[2], 0xff})
	id.SetMargins()
	id.Render(hashed)
	id.MirrorHorizontally()

	f, err := os.Create("image.png")
	if err != nil {
		fmt.Printf("Error Occurred: %v", err)
	} else {
		png.Encode(f, id.GetImg())
	}


}


