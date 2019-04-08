package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	var srcDir, dstDir string
	var rtl bool

	flag.StringVar(&srcDir, "src", "", "src dir of png /j(e)pg images")
	flag.StringVar(&dstDir, "dst", "", "dst dir")
	flag.BoolVar(&rtl, "rtl", false, "to enable RIGHT >> LEFT order")
	flag.Parse()

	if len(srcDir) == 0 || len(dstDir) == 0 {
		fmt.Println("Check your args")
		return
	}

	fmt.Printf("\nsrc: %s, dst: %s, rtl=%v\n", srcDir, dstDir, rtl)
	_, err := os.Stat(dstDir)
	if os.IsNotExist(err) {
		os.MkdirAll(dstDir, 0755)
	}

	file, _ := os.Open(srcDir)
	files, _ := file.Readdir(-1)
	for _, f := range files {
		imgPath := path.Join(srcDir, f.Name())
		workOnImage(imgPath, dstDir, rtl)
	}
}

func workOnImage(imgPath, targetDir string, rtl bool) {
	extName := path.Ext(imgPath)
	baseName := path.Base(imgPath)
	bareName := baseName[0 : len(baseName)-len(extName)]

	var img image.Image
	var err error

	imgFile, err := os.Open(imgPath)
	if err != nil {
		logError(err, "* open img file")
	}
	defer imgFile.Close()

	switch extName {
	case ".png":
		img, err = png.Decode(imgFile)
	case ".jpg":
		fallthrough
	case ".jpeg":
		img, err = jpeg.Decode(imgFile)
	default:
		logError(fmt.Errorf("Unknown image format: %s", extName), "*")
		return
	}

	if err != nil {
		logError(err, "* error when decode png")
		return
	}

	imgRect := img.Bounds()
	w := imgRect.Dx()
	h := imgRect.Dy()
	m := w / 2
	fmt.Printf("%dx%d - %s\n", w, h, imgPath)

	if h >= w {
		fmt.Print("---- Skip\n")
		copyFile(imgPath, targetDir)
		return
	}

	rgba := image.NewRGBA(imgRect)
	for x := imgRect.Min.X; x < imgRect.Max.X; x++ {
		for y := imgRect.Min.Y; y < imgRect.Max.Y; y++ {
			c := img.At(x, y)
			rgba.Set(x, y, c)
		}
	}

	rl := image.Rect(0, 0, m, h)
	imgL := rgba.SubImage(rl)
	imgR := rgba.SubImage(image.Rect(m, 0, w, h))

	tl, tr := getSuffix(rtl)
	pl := path.Join(targetDir, bareName+tl+extName)
	pr := path.Join(targetDir, bareName+tr+extName)
	fmt.Printf("---- %s\n", pl)
	fmt.Printf("---- %s\n", pr)

	fl, err := os.OpenFile(pl, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		logError(err, "error at open left file")
	}

	err = png.Encode(fl, imgL)
	if err != nil {
		logError(err, "error at write left image")
	}
	fl.Close()

	fr, err := os.OpenFile(pr, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		logError(err, "error at open left file")
	}

	err = png.Encode(fr, imgR)
	if err != nil {
		logError(err, "error at write left image")
	}
	fr.Close()
}

// fullpath to file, dir path
func copyFile(fullName, dstPath string) {
	srcFile, err := os.Open(fullName)
	if err != nil {
		logError(err, "* src file ???")
	}
	defer srcFile.Close()

	baseName := path.Base(fullName)
	dstName := path.Join(dstPath, baseName)

	dstFile, err := os.Create(dstName)
	if err != nil {
		logError(err, "* dst file exists?")
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		logError(err, "* copy file")
	}
}

func logError(err error, message string) {
	fmt.Printf("* %s", message)
	log.Println(err)
}

func getSuffix(rtl bool) (string, string) {
	if rtl {
		return "_1", "_0"
	}

	return "_0", "_1"
}
