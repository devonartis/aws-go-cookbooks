package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func savePem(f string, k string) error {

	return ioutil.WriteFile(f, []byte(k), 0666)

}

func main() {

	log.Printf(info, "Creating Key Pair: ", runtime.Version())

	/*
		TO DO: pairname should not be default of empty string
		add additional value checking for the pairname
		think about adding os.args with flags
	*/

	//flags for CLI commands

	pairname := flag.String("keyname", "", "Enter keypair name to create")
	profile := flag.String("profile", "", "Default profile will be used")
	region := flag.String("region", "us-east-1", "Region defaults to us-east-2")

	flag.Parse()

	// Create a new AWS Session with Options based on if a profile was given

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(*region)},
		Profile: *profile,
	}))

	//Get a handle on EC2 Service

	svc := ec2.New(sess)

	if *pairname == "" {
		exitErrorf("Keyname can not be empty")
	}

	keyresult, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(*pairname),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.Duplicate" {
			exitErrorf("Keypair %q already exists.", pairname)
		}
		exitErrorf("Unable to create key pair: %s, %v.", pairname, err)
	}

	savePem(*pairname+".pem", *keyresult.KeyMaterial)

	fmt.Printf("Created key pair %q %s\n%s\n",
		*keyresult.KeyName, *keyresult.KeyFingerprint,
		*keyresult.KeyMaterial)

}
