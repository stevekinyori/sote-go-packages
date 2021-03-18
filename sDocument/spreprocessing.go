package sDocument


import (
	"fmt"
	"gocv.io/x/gocv"
	leptonica "gopkg.in/GeertJohan/go.leptonica.v1"
	tesseract "gopkg.in/GeertJohan/go.tesseract.v1"
	"log"
	"os"
	"path/filepath"
)

func main() {
	img := gocv.IMRead("/Users/fionamurie/Desktop/sote/practice/Testing Materials/Container Guarantee Form back.jpg", 0)
	window := gocv.NewWindow("Hello")
	window.IMShow(img)
	window.WaitKey(0)
	// get the image to try
	//flag.Parse()
	//image := flag.Arg(0)
	image := "/Users/fionamurie/Desktop/sote/practice/Testing Materials/Container Guarantee Form back.jpg"

	fmt.Println(tesseract.Version())

	// create new tess instance and point it to the tessdata location. Set language to english.
	tessdata_prefix := os.Getenv("TESSDATA_PREFIX")
	if tessdata_prefix == "" {
		tessdata_prefix = "/opt/homebrew/Cellar/tesseract/4.1.1/share"
	}
	t, err := tesseract.NewTess(filepath.Join(tessdata_prefix, "tessdata"), "eng")
	if err != nil {
		log.Fatalf("Error while initializing Tess: %s\n", err)
	}
	defer t.Close()

	// open a new Pix from file with leptonica
	pix, err := leptonica.NewPixFromFile(image)
	if err != nil {
		log.Fatalf("Error while getting pix from file: %s\n", err)
	}
	defer pix.Close() // remember to cleanup

	// set the page seg mode to autodetect
	t.SetPageSegMode(tesseract.PSM_AUTO_OSD)

	// setup a whitelist of all basic ascii
	err = t.SetVariable("tessedit_char_whitelist", ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_abcdefghijklmnopqrstuvwxyz{|}~`+"`")
	if err != nil {
		log.Fatalf("Failed to SetVariable: %s\n", err)
	}

	// set the image to the tesseract instance
	t.SetImagePix(pix)
	//t.SetInputName(image)

	// retrieve text from the tesseract instance
	fmt.Println(t.Text())

	// // retrieve text from the tesseract instance
	// fmt.Println(t.HOCRText(0))

	// retrieve text from the tesseract instance
	//fmt.Println(t.BoxText(0))
	//
	//// now select just the first two columns (if using FelixScan.jpg)
	//t.SetRectangle(30, 275, 1120, 1380)
	//fmt.Println(t.Text())
	//fmt.Println(t.BoxText(0))

	// // retrieve text from the tesseract instance
	// fmt.Println(t.UNLVText())

	// dump variables for info
	// t.DumpVariables()

	//spew.Dump(t.AvailableLanguages())

}

