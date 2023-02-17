package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	database "github.com/KMRLAppPro/backend/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func UploadFile() gin.HandlerFunc {
    return func(c *gin.Context) {
		filename := c.Param("image_name")
		formFile, err := c.FormFile("file")
		if err != nil {
        	log.Fatal(err)
    	}
		openedFile, err := formFile.Open()
		if err != nil {
        	log.Fatal(err)
    	}
		data, err := ioutil.ReadAll(openedFile)
    	if err != nil {
        	log.Fatal(err)
    	}
    	bucket, err := gridfs.NewBucket(
        	database.Client.Database("kmrl"),
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
    	log.Printf("Write file to DB was successful. %s File size: %d M\n", filename,fileSize)
		c.JSON(http.StatusOK, uploadStream)
	}
}


func DownloadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileName := c.Param("image_name")
		// For CRUD operations, here is an example
    	db := database.Client.Database("kmrl")
    	fsFiles := db.Collection("fs.files")
    	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    	var results bson.M
		defer cancel()
    	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
    	if err != nil {
        	log.Fatal(err)
			//log.Fatal("findone")
    	}
    	// you can print out the results
    	//fmt.Println(results)

    	bucket, _ := gridfs.NewBucket(
        	db,
    	)
    	var buf bytes.Buffer
    	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
    	if err != nil {
        	log.Fatal(err)
			//log.Fatal("download")
    	}
    	fmt.Printf("File size to download: %v\n", dStream)
    	ioutil.WriteFile(fileName, buf.Bytes(), 0600)
		c.JSON(http.StatusOK, buf.Bytes())

		// imgbyt := buf.Bytes()
		
		// img, _, err := image.Decode(bytes.NewReader(imgbyt))
    	// if err != nil {
		// 	log.Println("image")
        // 	log.Fatalln(err)

    	// }
		// c.JSON(http.StatusOK, img)
	}
}


func DFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Find data
		db := database.Client.Database("kmrl")
		coll := db.Collection("fs.files")
		uid:=c.Value("uid")
		ucoll:=db.Collection("user")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		cursor, err := coll.Find(ctx, bson.D{},
  			options.Find().SetProjection(bson.D{
    			{"_id", 1},
    			{"filename", 1},
			}),
  			options.Find().SetSort(bson.D{{"uploadDate", -1}}))
		if err != nil {
		  log.Fatal(err)
		}
	
		for cursor.Next(context.TODO()) {
			// A new result variable should be declared for each document.
			var res bson.M
			if err := cursor.Decode(&res); err != nil {
				log.Fatal(err)
			}
			//fmt.Println(res)
			oid:=res["_id"]
			otid := oid.(primitive.ObjectID).Hex()
			//fmt.Println(i,"=",otid)  //i = _ in for 
		

			filter := bson.M{
    			"user_id": uid,
    			"viewed": bson.M{"$ne": otid},
			}
			update := bson.M{
    			"$push": bson.M{"viewed" : otid},
			}

			

			// filter := bson.M{
    		// 	"user_id": uid,
			// }
			// update := bson.M{
    		// 	"$addToSet": bson.M{
        	// 		"viewed": otid,
			// 	},
    		// }
			

			result, err := ucoll.UpdateOne(ctx, filter, update)
			if err != nil {
    			log.Fatal(err)
			}
			if result.MatchedCount == 0 {
    			continue
			}
			if result.ModifiedCount == 1 {
				//fmt.Println("not viewed")
    			fn :=res["filename"]
				fileName := fmt.Sprint(fn)
				bucket, _ := gridfs.NewBucket(
        			db,
    			)
    			var buf bytes.Buffer
    			dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
    			if err != nil {
        			log.Fatal(err)
					//log.Fatal("download")
    			}
    			fmt.Printf("File size to download: %v\n", dStream)
    			ioutil.WriteFile(fileName, buf.Bytes(), 0600)
				c.JSON(http.StatusOK, buf.Bytes())

				break
			}
		}
		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
		defer cancel()
    	
    	// you can print out the results
    	//fmt.Println(results)

    	

		// imgbyt := buf.Bytes()
		
		// img, _, err := image.Decode(bytes.NewReader(imgbyt))
    	// if err != nil {
		// 	log.Println("image")
        // 	log.Fatalln(err)

    	// }
		// c.JSON(http.StatusOK, img)
	}
}

//err := fsFiles.find({}).skip(1).limit(50).sort({"date":-1}) 

// // Requires the MongoDB Go Driver
// // https://go.mongodb.org/mongo-driver
// ctx := context.TODO()

// // Set client options
// clientOptions := options.Client().ApplyURI("mongodb+srv://samanuai:backendkmrl@kmrl.phjuman.mongodb.net/test")

// // Connect to MongoDB
// client, err := mongo.Connect(ctx, clientOptions)
// if err != nil {
//   log.Fatal(err)
// }
// defer func() {
//   if err := client.Disconnect(ctx); err != nil {
//     log.Fatal(err)
//   }
// }()

// // Find data
// coll := client.Database("kmrl").Collection("fs.files")
// _, err = coll.Find(ctx, bson.D{},
//   options.Find().SetProjection(bson.D{{"_id", 1}}),
//   options.Find().SetSort(bson.D{{"uploadDate", -1}}))
// if err != nil {
//   log.Fatal(err)
// }



			// //fmt.Println(res)
			// oid:=res["_id"]
			// otid := oid.(primitive.ObjectID).Hex()
			// //fmt.Println(i,"=",otid)  //i = _ in for 
		

			// filter := bson.M{
    		// 	"user_id": uid,
    		// 	"viewed": bson.M{"$ne": otid},
			// }
			// update := bson.M{
    		// 	"$push": bson.M{"viewed" : otid},
			// }

			

			// // filter := bson.M{
    		// // 	"user_id": uid,
			// // }
			// // update := bson.M{
    		// // 	"$addToSet": bson.M{
        	// // 		"viewed": otid,
			// // 	},
    		// // }
			

			// result, err := ucoll.UpdateOne(ctx, filter, update)
			// if err != nil {
    		// 	log.Fatal(err)
			// }
			// if result.MatchedCount == 0 {
    		// 	continue
			// }
			// if result.ModifiedCount == 1 {
			// 	//fmt.Println("not viewed")
    		// 	fn :=res["filename"]
			// 	fileName := fmt.Sprint(fn)
			// 	bucket, _ := gridfs.NewBucket(
        	// 		db,
    		// 	)
    		// 	var buf bytes.Buffer
    		// 	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
    		// 	if err != nil {
        	// 		log.Fatal(err)
			// 		//log.Fatal("download")
    		// 	}
    		// 	fmt.Printf("File size to download: %v\n", dStream)
    		// 	ioutil.WriteFile(fileName, buf.Bytes(), 0600)
			// 	c.JSON(http.StatusOK, buf.Bytes())

			// 	break
			// }