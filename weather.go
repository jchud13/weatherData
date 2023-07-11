package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event interface{}) (WeatherResponse, error) {
	fmt.Printf("%v\n", event)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	tableName := "WeatherData"
	movieName := "The Big New Movie"
	movieYear := "2015"

	tableRes, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Year": {
				N: aws.String(movieYear),
			},
			"Title": {
				S: aws.String(movieName),
			},
		},
	})
	if err != nil {
		fmt.Println("Got error calling GetItem: ", err)
	}

	item := WeatherData{}

	err = dynamodbattribute.UnmarshalMap(tableRes.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	var weatherURL string = "https://api.openweathermap.org/data/2.5/weather?lat=39.961178&lon=-82.998795&units=imperial&appid=fc6ef58dc2fa7c67cebf4ea828ccac75"
	resp, err := http.Get(weatherURL)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(err)
	fmt.Println(string(body))
	var result WeatherResponse
	res := json.Unmarshal(body, &result)
	if res != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	fmt.Println(PrettyPrint(result))
	return result, nil
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
