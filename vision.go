package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	//"flag"

	"cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

var allFeatures = []*pb.Feature{
	{Type: pb.Feature_TYPE_UNSPECIFIED, MaxResults: 3},
	{Type: pb.Feature_LANDMARK_DETECTION, MaxResults: 5},
	{Type: pb.Feature_LABEL_DETECTION, MaxResults: 3},
	{Type: pb.Feature_FACE_DETECTION, MaxResults: 3},
	{Type: pb.Feature_LOGO_DETECTION, MaxResults: 3},
	{Type: pb.Feature_TEXT_DETECTION, MaxResults: 3},
	{Type: pb.Feature_DOCUMENT_TEXT_DETECTION, MaxResults: 3},
	{Type: pb.Feature_IMAGE_PROPERTIES, MaxResults: 3},
}

func newClient() (context.Context, *vision.ImageAnnotatorClient, error) {
	ctx := context.Background()

	c, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return ctx, &vision.ImageAnnotatorClient{}, err
	}

	return ctx, c, nil
}

func annotateImage(ctx context.Context, c *vision.ImageAnnotatorClient, features []*pb.Feature, filename string) (*pb.AnnotateImageResponse, error) {
	file, err := os.Open(filename)
	if err != nil {
		return &pb.AnnotateImageResponse{}, err
	}

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return &pb.AnnotateImageResponse{}, err
	}

	res, err := c.AnnotateImage(ctx, &pb.AnnotateImageRequest{
		Image:    image,
		Features: allFeatures,
	})
	if err != nil {
		return &pb.AnnotateImageResponse{}, err
	}

	return res, nil
}

func runImage(filename string) (*pb.AnnotateImageResponse, error) {
	ctx, c, err := newClient()
	if err != nil {
		return &pb.AnnotateImageResponse{}, err
	}

	res, err := annotateImage(ctx, c, allFeatures, filename)
	if err != nil {
		return &pb.AnnotateImageResponse{}, err
	}

	return res, nil
}

// loadResponseFromFile returns a AnnotateImageResponse loaded from a JSON file in order to test without using up API requests
func loadResponseFromFile(filename string) (*pb.AnnotateImageResponse, error) {
	file, err := os.Open(filename)
	if err != nil {
		return &pb.AnnotateImageResponse{}, err
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return &pb.AnnotateImageResponse{}, err
	}

	var res pb.AnnotateImageResponse

	err = json.Unmarshal(fileBytes, &res)
	if err != nil {
		return &pb.AnnotateImageResponse{}, err
	}

	return &res, nil

}

func responseToJSON(res *pb.AnnotateImageResponse) (string, error) {
	resBytes, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(resBytes), nil
}
