package main

import (
	"flag"
	"fmt"
	//"image"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func main() {

	image := flag.String("image", "", "Image file to get annotations for")
	storedJSON := flag.String("json", "", "JSON file containing a image annotation response")
	flag.Parse()

	if (*image != "" && *storedJSON != "") || (*image == "" && *storedJSON == "") {
		fmt.Println("Use either -image or -json, not both or neither")
		return
	}

	var res *pb.AnnotateImageResponse
	var err error

	if *image != "" {
		res, err = runImage(*image)
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

	/*faces := res.GetFaceAnnotations()

	for _, face := range faces {
		poly := face.GetBoundingPoly()
		verts := poly.GetVertices()
		for _, vert := range verts {
			fmt.Printf("%v %v\n", vert.GetX(), vert.GetY())
		}
	}*/

	jsonString, err := responseToJSON(res)
	if err != nil {
		panic(err)
	}

	fmt.Println(jsonString)
}
