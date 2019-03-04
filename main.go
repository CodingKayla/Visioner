package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"

	"github.com/fogleman/gg"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func main() {

	imageFileName := flag.String("imageFileName", "", "imageFileName file to get annotations for")
	storedJSON := flag.String("json", "", "JSON file containing a imageFileName annotation response")
	flag.Parse()

	if (*imageFileName != "" && *storedJSON != "") || (*imageFileName == "" && *storedJSON == "") {
		fmt.Println("Use either -imageFileName or -json, not both or neither")
		return
	}
	var res *pb.AnnotateImageResponse
	var err error

	if *imageFileName != "" {
		res, err = runImage(*imageFileName)
		if err != nil {
			panic(err)
		}
	}

	if *storedJSON != "" {
		res, err = loadResponseFromFile(*storedJSON)
		if err != nil {
			panic(err)
		}
	}

	_ = res

	sourceFile, err := os.Open("images/basketball.jpg")
	if err != nil {
		panic(err)
	}

	imageData, _, err := image.Decode(sourceFile)

	dc := gg.NewContextForImage(imageData)
	dc.DrawCircle(500, 500, 400)
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	dc.SavePNG("out.png")

	faces := res.GetFaceAnnotations()

	for _, face := range faces {
		poly := face.GetBoundingPoly()
		verts := poly.GetVertices()
		for _, vert := range verts {
			fmt.Printf("%v %v\n", vert.GetX(), vert.GetY())

		}
	}

	/*jsonString, err := responseToJSON(res)
	if err != nil {
		panic(err)
	}

	fmt.Println(jsonString)*/
}
