package generator

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
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
	image           *image.RGBA
	marginX         int
	marginY         int
}

func New(size int) *identicon{
	//Defining cell's side length = min(halfW/_COLS, halfH/_ROWS)
	//Half because the image is written on the first half than is mirrored
	cell := setCellSide(size)

	//Make an image
	img := initImg(size)
	marginX, marginY := getMargins(size, cell)
	return &identicon{
		size:         	 size,
		cell:            cell,
		backgroundColor: color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 0xff,
		},
		image:   img,
		marginX: marginX,
		marginY: marginY,
	}
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


func getMargins(size, cell int) (marginX, marginY int){
	// left/right margin
	marginX = size/2 - (cell * _COLS)/2
	// top/bottom margin
	marginY = size/2 - (cell * _ROWS)/2
	return
}

func (icon *identicon) renderCell(x, y int, aColor color.RGBA) (availableX int){
	//render a cell given x,y and returns last x available
	img := icon.image
	var i,j int
	for i = x; i - x <= icon.cell; i++ {
		for j = y; j - y <= icon.cell; j++ {
			img.SetRGBA(i, j, aColor)
		}
	}
	return i
}

func (icon *identicon) mirrorHorizontally(){
	img := icon.image
	halfSize := icon.size/2
	for x := icon.marginX; x < halfSize; x++{
		for y := icon.marginY; y < icon.size; y++ {
			img.SetRGBA(icon.size - x, y, img.RGBAAt(x, y))
		}
	}
}


func (icon *identicon) render(hash []uint8, waiter *sync.WaitGroup) {
	hashSliced := hash[3:]
	var x, y int
	x = icon.marginX
	y = icon.marginY
	currByte := 0
	for i := 0; i < _ROWS; i++ {
		for j := 0; j < _COLS/2; j++ {
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
	if waiter != nil{
		waiter.Done()
	}
}

func (icon *identicon) renderBackground(waiter *sync.WaitGroup) {
	img := icon.image
	//drawing top & bottom background
	for x := 0; x < icon.size; x++ {
		for y := 0; y < icon.marginY; y++ {
			img.SetRGBA(x, y, icon.backgroundColor)
			img.SetRGBA(x, icon.size-y, icon.backgroundColor)
		}
	}
	//drawing left/right background
	for y := icon.marginY; y <= icon.size-icon.marginY; y++ {
		for x:= 0; x <= icon.marginX; x++ {
			img.SetRGBA(x, y, icon.backgroundColor)
			img.SetRGBA(icon.size - x, y, icon.backgroundColor)
		}
	}

	if waiter != nil{
		waiter.Done()
	}

}



func (icon *identicon) Create(hash []uint8, fileName string) interface{} {
	//Using first 3 hash bytes for color
	icon.foregroundColor = color.RGBA{hash[0], hash[1], hash[2], 255}
	var waiter  sync.WaitGroup
	waiter.Add(2)
	go icon.render(hash, &waiter)
	go icon.renderBackground(&waiter)
	waiter.Wait()
	icon.mirrorHorizontally()
	f, err := os.Create(fileName+".png")
	if err != nil {
		return err
	} else {
		err = png.Encode(f, icon.image)
	}
	return err
}

func (icon *identicon) NoCreate(hash []uint8, fileName string) interface{} {
	//Using first 3 hash bytes for color
	icon.foregroundColor = color.RGBA{hash[0], hash[1], hash[2], 255}

	icon.render(hash, nil)
	icon.renderBackground(nil)

	icon.mirrorHorizontally()
	f, err := os.Create(fileName+".png")
	if err != nil {
		return err
	} else {
		err = png.Encode(f, icon.image)
	}
	return err
}

