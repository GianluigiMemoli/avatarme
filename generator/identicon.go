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
	w           int
	h          	int
	marginLeftRight int
	marginTopDown   int
	cell            int
	backgroundColor color.Color
	asImage			*image.RGBA
}

func setCellSide(w, h int) int {
	halfW := w/2
	halfH := h/2
	div1 := halfW / (_COLS/2)
	div2 := halfH / (_ROWS)
	if div1 < div2{
		return div1
	}
	return div2
}

func getMargins(w, h, cell int) (left_right, top_down int) {
	marginLR := w - (cell * _COLS)
	marginTD := h - (cell * _ROWS)
	log.Printf("Margin RLside: %d\nMargin TD: %d\n", marginLR, marginTD)
	return marginLR, marginTD
}

func getImage(w, h int) *image.RGBA{
	ptTop := image.Point{0, 0}
	ptBtm := image.Point{w, h}
	return image.NewRGBA(image.Rectangle{ptTop, ptBtm})
}

func New(w, h int) *identicon{
	//Defining cell's side length = min(halfW/_COLS, halfH/_ROWS)
	//Half because the image is written on the first half than is mirrored
	cell := setCellSide(w, h)
	//Getting some padding
	marginLR, marginTD := getMargins(w, h, cell)
	//Make an image
	img := getImage(w, h)
	log.Printf("cell: %d\n", cell)
	return &identicon{
		w:           w,
		h:         	 h,
		marginLeftRight: marginLR,
		marginTopDown:   marginTD,
		cell:            cell,
		backgroundColor: color.White,
		asImage: img,
	}
}

func (that *identicon) PrintCell(x ,y int, aColor color.RGBA) (availableX int){
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
	img := that.GetImg()
	halfH := that.h / 2
	halfW := that.w / 2
	for i :=0; i < halfH; i++{
		for j := 0; j < halfW; j++ {
			img.SetRGBA(that.w - i, j, img.RGBAAt(i, j))
		}
	}
}

func (that *identicon) SetMargins(){
	//img := that.GetImg()

}

func (that *identicon) GetImg()  *image.RGBA{
	return that.asImage
}

func (that *identicon) Render(hash []uint8) {
	hashSliced := hash[3:]
	var x, y int
	x = 0
	y = 0
	currByte := 0
	for i:=0; i < _ROWS; i++{
		for j := 0; j < _COLS/2; j++{
			log.Printf("Considering byte: %d", hashSliced[currByte])
			if hashSliced[currByte] % 2 == 0{
				x = that.PrintCell(x,y, color.RGBA{
					R: 255,
					G: 0,
					B: 0,
					A: 0xff,
				})
			} else {
				x = that.PrintCell(x,y, color.RGBA{
					R: 0,
					G: 0,
					B: 255,
					A: 0xff,
				})
			}
			currByte++
		}
		x = 0
		y += that.cell
	}
}