package main

//Create an Amazon EC2 instance with tags and log into the instance

import (
	"fmt"
	"log"
	"net/url"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	pairname string
	svc      *ec2.EC2
)

const info = `
          Application %s starting.
		  The binary was build by GO: %s`

func main() {

	// Allows you to define the profile if not profile is given the default profile will be used

	//var ip address
	log.Printf(info, "Launch Instance: ", runtime.Version())

	cred := credentials.NewSharedCredentials("credentials", "rai")
	// Create a new AWS Session with Options based on if a profile was given

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(*region)},
		Profile: *profile,
	}))

	/*
		sess.Handlers.Send.PushFront(func(r *request.Request) {
			// Log every request made and its payload
			fmt.Println("Request: %s/%s, Payload: %s",
				r.ClientInfo.ServiceName, r.Operation, r.Params)
		})
	*/

	//Get a handle on EC2 Service
	svc := ec2.New(sess)

	launchResult, err := svc.RunInstances(&ec2.RunInstancesInput{

		ImageId:      aws.String(*imageid),
		InstanceType: aws.String(*instanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		log.Println("Instance could not be created", err)
		return
	}

	log.Println("Created instance", *launchResult.Instances[0].InstanceId)

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{launchResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("MyFirstInstance"),
			},
		},
	})
	if errtag != nil {
		log.Println("Could not create tags for instance", *launchResult.Instances[0].InstanceId, errtag)
		return
	}

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		},
	}

	resp, _ := svc.DescribeInstances(params)

	for idx, _ := range resp.Reservations {

		for _, inst := range resp.Reservations[idx].Instances {

			name := "None"

			for _, keys := range inst.Tags {
				if *keys.Key == "Name" {
					name = url.QueryEscape(*keys.Value)
				}
			}

			important_vals := []*string{
				inst.InstanceId,
				&name,
				inst.PrivateIpAddress,
				inst.InstanceType,
				inst.PublicIpAddress,
			}

			// Convert any nil value to a printable string in case it doesn't
			// doesn't exist, which is the case with certain values
			output_vals := []string{}
			for _, val := range important_vals {
				if val != nil {
					output_vals = append(output_vals, *val)
				} else {
					output_vals = append(output_vals, "None")
				}
			}
			// The values that we care about, in the order we want to print them

			fmt.Println(strings.Join(output_vals, " "))
		}

	}

}
