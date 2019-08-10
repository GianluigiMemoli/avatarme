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
	marginX         int
	marginY         int
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
	marginX, marginY := SetPadding(size, cell)
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
		marginX: marginX,
		marginY:marginY,
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
	for x := icon.marginX; x < halfSize; x++{
		for y := icon.marginX; y < icon.size; y++ {
			img.SetRGBA(icon.size - x, y, img.RGBAAt(x, y))
		}
	}
}

func SetPadding(size, cell int) (marginX, marginY int){
	/*usedArea := (icon.cell * icon.cell) * (_COLS * _ROWS)
	unusedArea := (icon.size * icon.size) - usedArea
	//Consider unusedArea as 4 rectangle, one per side
	rectArea := unusedArea / 4
	rectHeight := rectArea / icon.size
	//Adding 3 squares that are lost in the marginX composition
	rectArea += (rectHeight * icon.size)
	rectHeight = rectArea / icon.size*/
	marginX = size/2 - (cell * _COLS)/2
	//Top margin
	usedHeight := cell * _ROWS
	marginY = size/2 - usedHeight/2

	log.Printf("MarginX: %d MarginY: %d", marginX, marginY)

	return
}

func (icon *identicon) GetImg()  *image.RGBA{
	return icon.asImage
}

func (icon *identicon) Render(hash []uint8) {
	hashSliced := hash[3:]
	var x, y int
	x = icon.marginX
	y = icon.marginY
	log.Printf("topmargin: %d", y)
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
		x = icon.marginX
		y += icon.cell
	}
}

func (icon *identicon) RenderBackground(){
	for x := 0; x < icon.size/2; x++ {
		for y := 0; y < icon.size; y++ {

		}
	}
}