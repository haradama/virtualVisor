package main

import (
	"os"
	"fmt"
	"math"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

var (
	img gocv.Mat
	hexGrid *HexGrid
	deviceID = 0
	xmlFile = "./haarcascade_frontalface_alt2.xml"
)

const (
	WindowPropertyFullscreen gocv.WindowPropertyFlag = 0
    WindowFullscreen gocv.WindowFlag = 1
)

type Hexagon struct {
	centerX int
	centerY int
	radius int
	isFill bool
}

func NewHexagon(centerX int, centerY int, radius int) *Hexagon {
	hexagon := &Hexagon{
		centerX: centerX,
		centerY: centerY,
		radius: radius,
		isFill: false,
	}

	return hexagon
}

type HexGrid struct {
	colNum int
    rowNum int
	radius int
	grid [][]*Hexagon
}

func NewHexGrid(colNum int, rowNum int, radius int) *HexGrid {
	x := int(math.Sqrt(3) * float64(radius))
	y := int(radius)

	grid := [][]*Hexagon{}
	for i := 0; i < colNum; i++ {
		row := []*Hexagon{}
		for j := 0; j < rowNum; j++ {
			hexagon := NewHexagon(y, x, radius)
			row = append(row, hexagon)
			x += int(float64(radius) * math.Sqrt(3))
		}
		grid = append(grid, row)
		y += (radius* 3) / 2;

		if ((i + 1) % 2 == 0) {
			x = int(float64(radius) * math.Sqrt(3))
		} else {
			x = int(float64(radius) * math.Sqrt(3) / 2);
		}
	}

	hexGrid := &HexGrid{
		colNum: colNum,
		rowNum: rowNum,
		radius: radius,
		grid: grid,
	}

	return hexGrid
}

func (hexGrid *HexGrid) draw(rect image.Rectangle) {
	for i := 0; i < hexGrid.rowNum; i++ {
		for j := 0; j < hexGrid.colNum; j++ {
			hexagon := hexGrid.grid[j][i]
			position := image.Point{hexagon.centerX, hexagon.centerY}
			if (isInRectangle(position, rect)) {
				ngon := 6
				points := []image.Point{}
				for k := 0; k < ngon; k++ {
					rate := float64(k) / float64(ngon)
					sx := int(float64(hexagon.radius) * math.Cos(2 * math.Pi * rate)) + hexagon.centerX
					sy := int(float64(hexagon.radius) * math.Sin(2 * math.Pi * rate)) + hexagon.centerY

					points = append(points,
						image.Point{sx, sy},
					)
				}
				poly := [][]image.Point{points}
				gocv.FillPoly(&img, poly, color.RGBA{0, 0, 0, 0})
			}
		}
	}
}

func isInRectangle(point image.Point, rectangle image.Rectangle) bool {
	min := rectangle.Min
	max := rectangle.Max
	if (min.X <= point.X && point.X <= max.X && min.Y <= point.Y && point.Y <= max.Y) {
		return true
	}
	return false
}

func main() {
	radius := 40
	screenWidth:= 800
	screenHeight := 600

	window := gocv.NewWindow("Virtual Visor")
	window.SetWindowProperty(WindowPropertyFullscreen, WindowFullscreen)

	colNum := int(screenWidth / radius)
	rowNum := int(screenHeight / radius)
	fmt.Println(colNum, rowNum)
	hexGrid = NewHexGrid(colNum, rowNum, radius)
	
	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		panic(err)
	}
	defer webcam.Close()

	classifier := gocv.NewCascadeClassifier()
    defer classifier.Close()

    if !classifier.Load(xmlFile) {
        fmt.Printf("Error reading cascade file: %v\n", xmlFile)
        return
    }

	img = gocv.NewMat()
	defer img.Close()

	white := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{800, 600},
	})
	for x := 0; x < screenWidth; x++ {
		for y := 0; y < screenHeight; y++ {
			bg.Set(x, y, color.White)
		}
	}
	bg, _ = gocv.ImageToMatRGBA(white)
	for {
		webcam.Read(&img)
		gocv.Flip(img, &img, 1)
		rects := classifier.DetectMultiScale(img)

		if (len(rects) > 0) {
			rect := rects[0]
			hexGrid.draw(rect)
		}

		window.IMShow(mat)
		key := window.WaitKey(1)
		if (key == 27) {
			window.Close()
			os.Exit(1)
		}
	}
}