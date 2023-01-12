package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	database "github.com/KMRLAppPro/backend/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func UploadFile() gin.HandlerFunc {
    return func(c *gin.Context) {
		filename := c.Param("image_name")
		formFile, _ := c.FormFile("file")
		openedFile, _ := formFile.Open()
		data, err := ioutil.ReadAll(openedFile)
    	if err != nil {
        	log.Fatal(err)
    	}
    	bucket, err := gridfs.NewBucket(
        	database.Client.Database("myfiles"),
    	)
    	if err != nil {
        	log.Fatal(err)
        	return
    	}
    	uploadStream, err := bucket.OpenUploadStream(
        	filename,
    	)
    	if err != nil {
        	return
    	}
    	defer uploadStream.Close()

    	fileSize, err := uploadStream.Write(data)
    	if err != nil {
        	log.Fatal(err)
       		return
    	}
    	log.Printf("Write file to DB was successful. File size: %d M\n", fileSize)

	}
}


func DownloadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileName := c.Param("image_name")
		// For CRUD operations, here is an example
    	db := database.Client.Database("myfiles")
    	fsFiles := db.Collection("fs.files")
    	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    	var results bson.M
		defer cancel()
    	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
    	if err != nil {
        	log.Fatal(err)
    	}
    	// you can print out the results
    	fmt.Println(results)

    	bucket, _ := gridfs.NewBucket(
        	db,
    	)
    	var buf bytes.Buffer
    	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
    	if err != nil {
        	log.Fatal(err)
    	}
    	fmt.Printf("File size to download: %v\n", dStream)
    	ioutil.WriteFile(fileName, buf.Bytes(), 0600)

	}
}