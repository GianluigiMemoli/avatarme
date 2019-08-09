package generator

import (
	//	"fmt"
	//	"image"
	"image"
	"image/color"
	"log"

	//	"image/png"
	//	"os"
)

const (
	_ROWS = 5
	_COLS = 6
)

type identicon struct {
	size           int
	cell            int
	backgroundColor color.RGBA
	foregroundColor color.RGBA
	asImage			*image.RGBA
	margin 			int
}

func setCellSide(size int) int {
	halfSize := size/2
	div1 := halfSize / (_COLS/2)
	div2 := halfSize / (_ROWS)
	if div1 < div2{
		return div1
	}
	return div2
}

func getImage(size int) *image.RGBA{
	ptTop := image.Point{0, 0}
	ptBtm := image.Point{size, size}
	return image.NewRGBA(image.Rectangle{ptTop, ptBtm})
}

func New(size int, fg color.RGBA) *identicon{
	//Defining cell's side length = min(halfW/_COLS, halfH/_ROWS)
	//Half because the image is written on the first half than is mirrored
	log.Printf("Color fg: %v", fg)
	cell := setCellSide(size)
	//Make an image
	img := getImage(size)
	log.Printf("cell: %d\n", cell)
	return &identicon{
		size:         	 size,
		cell:            cell,
		backgroundColor: color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 0xff,
		},
		foregroundColor:	fg,
		asImage: img,
	}
}

func (that *identicon) PrintCell(x, y int, aColor color.RGBA) (availableX int){
	//Render a cell given x,y and returns last x available
	img := that.asImage
	var i,j int
	log.Printf("Printing a cell x:%d y:%d", x, y)
	for i = x; i - x <= that.cell; i++ {
		for j = y; j - y <= that.cell; j++ {
			img.SetRGBA(i, j, aColor)
		}
	}
	return i
}

func (that *identicon) MirrorHorizontally(){
	img := that.asImage
	halfSize := that.size/2
	for x := that.margin; x < halfSize; x++{
		for y := that.margin; y < that.size; y++ {
			img.SetRGBA(that.size - x, y, img.RGBAAt(x, y))
		}
	}
}

func (that *identicon) SetMargins(){
	usedArea := (that.cell * that.cell) * (_COLS * _ROWS)
	unusedArea := (that.size * that.size) - usedArea
	//Consider unusedArea as 4 rectangle, one per side
	rectArea := unusedArea / 4
	rectSide := rectArea / that.size
	//Seems work, im adding 3 squared that are lost in the margin composition
	rectArea += (rectSide * rectSide)*3
	rectSide  = rectArea / that.size
	//draw margins
	img := that.asImage
	log.Printf("UsedArea:%d\nUnusedArea:%d\nRectSide: %d", usedArea, unusedArea, rectSide)
	//top and bottom margins
	for x := 0; x < that.size; x++ {
		for y := 0; y < rectSide; y++ {
			//top
			img.SetRGBA(x, y, color.RGBA{255, 255,255, 0xff})
			//bottom
			img.SetRGBA(x, that.size-y, color.RGBA{255, 255,255, 0xff})

		}
	}
	//left and right margins
	for y := rectSide; y < that.size; y++ {
		for x := 0; x < rectSide; x++{
			//left
			img.SetRGBA(x, y, color.RGBA{255, 255,255, 0xff})
			//right
			img.SetRGBA(that.size - x, y, color.RGBA{255, 255,255, 0xff})
		}

	}
	that.margin  = rectSide
}

func (that *identicon) GetImg()  *image.RGBA{
	return that.asImage
}

func (that *identicon) Render(hash []uint8) {
	hashSliced := hash[3:]
	var x, y int
	x = that.margin
	y = that.margin
	currByte := 0
	for i := 0; i < _ROWS; i++ {
		for j := 0; j < _COLS/2; j++ {
			log.Printf("Considering byte: %d", hashSliced[currByte])
			if hashSliced[currByte]%2 == 0 {
				x = that.PrintCell(x, y, that.foregroundColor)
			} else {
				x = that.PrintCell(x, y, that.backgroundColor)
			}
			currByte++

		}
		x = that.margin
		y += that.cell
	}
}