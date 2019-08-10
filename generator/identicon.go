package generator

import (
	"image"
	"image/color"
	"log"
)

const (
	_ROWS = 5
	_COLS = 6
)

type identicon struct {
	size            int
	cell            int
	backgroundColor color.RGBA
	foregroundColor color.RGBA
	asImage         *image.RGBA
	padding         int
}

func setCellSide(size int) int {
	/*
		The cell size is chosen on the half size of image because it will be
		horizontally mirrored
	*/
	halfSize := size/2
	div1 := halfSize / (_COLS/2)
	div2 := halfSize / (_ROWS)
	if div1 < div2{
		return div1
	}
	return div2
}

func initImg(size int) *image.RGBA{
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
	img := initImg(size)
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

func (icon *identicon) renderCell(x, y int, aColor color.RGBA) (availableX int){
	//Render a cell given x,y and returns last x available
	img := icon.asImage
	var i,j int
	log.Printf("Printing a cell x:%d y:%d", x, y)
	for i = x; i - x <= icon.cell; i++ {
		for j = y; j - y <= icon.cell; j++ {
			img.SetRGBA(i, j, aColor)
		}
	}
	return i
}

func (icon *identicon) MirrorHorizontally(){
	img := icon.asImage
	halfSize := icon.size/2
	for x := icon.padding; x < halfSize; x++{
		for y := icon.padding; y < icon.size; y++ {
			img.SetRGBA(icon.size - x, y, img.RGBAAt(x, y))
		}
	}
}

func (icon *identicon) SetPadding(){
	usedArea := (icon.cell * icon.cell) * (_COLS * _ROWS)
	unusedArea := (icon.size * icon.size) - usedArea
	//Consider unusedArea as 4 rectangle, one per side
	rectArea := unusedArea / 4
	rectHeight := rectArea / icon.size
	//Seems work, im adding 3 squared icon are lost in the padding composition
	rectArea += (rectHeight * rectHeight)
	rectHeight = rectArea / icon.size
	log.Printf("UsedArea:%d\nUnusedArea:%d\nRectSide: %d", usedArea, unusedArea, rectHeight)
	icon.padding = rectHeight
}

func (icon *identicon) GetImg()  *image.RGBA{
	return icon.asImage
}

func (icon *identicon) Render(hash []uint8) {
	hashSliced := hash[3:]
	var x, y int
	x = icon.padding
	y = icon.padding
	currByte := 0
	for i := 0; i < _ROWS; i++ {
		for j := 0; j < _COLS/2; j++ {
			log.Printf("Considering byte: %d", hashSliced[currByte])
			if hashSliced[currByte]%2 == 0 {
				x = icon.renderCell(x, y, icon.foregroundColor)
			} else {
				x = icon.renderCell(x, y, icon.backgroundColor)
			}
			currByte++

		}
		x = icon.padding
		y += icon.cell
	}
}