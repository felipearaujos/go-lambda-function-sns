package main

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

type MyEvent struct {
    Bucket string `json:"bucket"`
    Key    string `json:"key"`
}

func handler(ctx context.Context, snsEvent events.SNSEvent) error {
    for _, record := range snsEvent.Records {
        snsRecord := record.SNS
        var myEvent MyEvent
        err := json.Unmarshal([]byte(snsRecord.Message), &myEvent)
        if err != nil {
            return fmt.Errorf("Error unmarshalling SNS message: %s", err.Error())
        }
        fmt.Printf("Received event: %+v\n", myEvent)

        // Initialize the S3 client
        sess := session.Must(session.NewSession())
        svc := s3.New(sess)

        // Copy the object from the source bucket to the destination bucket
        _, err = svc.CopyObject(&s3.CopyObjectInput{
            Bucket:     aws.String("<destination-bucket-name>"),
            CopySource: aws.String(fmt.Sprintf("%s/%s", myEvent.Bucket, myEvent.Key)),
            Key:        aws.String(myEvent.Key),
        })
        if err != nil {
            return fmt.Errorf("Error copying object: %s", err.Error())
        }
        fmt.Printf("Successfully copied object from %s/%s to <destination-bucket-name>/%s\n", myEvent.Bucket, myEvent.Key, myEvent.Key)
    }
    return nil
}

func main() {
    lambda.Start(handler)
}
