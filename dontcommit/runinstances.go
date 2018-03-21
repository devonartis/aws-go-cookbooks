package main

import (
	"fmt"
	"log"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

)

var (
	pairname string
	svc      *ec2.EC2
)



//func writeFile(k string) {
//	fmt.Println("In writeFile %s", k)

//}

func main() {

	pairname = "GOkey34"
	svc = ec2.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))

	// Specify EC2 instance you want to launch

	keyresult, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(pairname),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.Duplicate" {
			exitErrorf("Keypair %q already exists.", pairname)
		}
		exitErrorf("Unable to create key pair: %s, %v.", pairname, err)
	}

	//writeFile(*keyresult.KeyMaterial)

	fmt.Printf("Created key pair %q %s\n%s\n",
		*keyresult.KeyName, *keyresult.KeyFingerprint,
		*keyresult.KeyMaterial)

	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String("ami-9eb4b1e5"),
		InstanceType: aws.String("t2.nano"),
		KeyName:      aws.String(pairname),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		log.Println("Could not create instance", err)
		return
	}

	log.Println("Created Instance", *runResult.Instances[0].InstanceId)

	// Add tags to the created instance

	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("MyFristGoInstance"),
			},
		},
	})

	if errtag != nil {
		log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return
	}

	log.Println("Succesfully tagged instance")

}
