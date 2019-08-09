package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/GianluigiMemoli/avatarme/generator"
	"image/png"
	"os"
)

func main(){
	id := generator.New(256, 256)
	toBeHashed := "Diventerò un hash"
	h := sha1.New()
	h.Write([]byte(toBeHashed))
	hashed := h.Sum(nil)
	id.Render(hashed)
	id.MirrorHorizontally()
	f, err := os.Create("image.png")
	if err != nil {
		fmt.Printf("Error Occurred: %v", err)
	} else {
		png.Encode(f, id.GetImg())
	}


}


