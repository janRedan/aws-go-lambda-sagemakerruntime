package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"
)

var svc *sagemakerruntime.SageMakerRuntime

const (
	//ENDPOINT value of sagemaker endpoint
	ENDPOINT = "ENDPOINT HERE"
)

//ReturnType will be the holder for SageMaker response
type ReturnType []string

//InputEvent is the type for Handler() input
type InputEvent struct {
	UserID string `json:"userid"`
}

//SageMakerInput struct to be passed in invokeendpoint
type SageMakerInput struct {
	Users []string `json:"users"`
}

func init() {
	mySession := session.Must(session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: aws.String("ap-southeast-1"),
			},
		}))
	svc = sagemakerruntime.New(mySession)
}

func main() {
	lambda.Start(Handler)
}

//Handler for the lambda request
func Handler(inputType InputEvent) (ReturnType, error) {

	// if len(inputType.Users) != 1 {
	// 	return "", fmt.Errorf("invalid user data")
	// }

	sageMakerInput := SageMakerInput{
		Users: []string{
			inputType.UserID,
		},
	}

	sageJSONInput, err := json.Marshal(sageMakerInput)
	checkError(err)

	x := &sagemakerruntime.InvokeEndpointInput{
		Body:         sageJSONInput,
		EndpointName: aws.String(ENDPOINT),
		ContentType:  aws.String("application/json"),
	}

	output, err := svc.InvokeEndpoint(x) //how to get output body?
	checkError(err)

	z := BytesToString(output.Body)
	var dat ReturnType
	//var dat1 ReturnType
	err = json.Unmarshal([]byte(z), &dat)
	checkError(err)

	return dat, err
	//depends on header if sagemaker has return

}

//BytesToString converts the given byte to string
func BytesToString(data []byte) string {
	return string(data[:])
}
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
