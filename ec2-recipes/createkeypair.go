package main

import (
	"fmt"
	"io/ioutil"
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

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func savePem(f string, k string) error {

	return ioutil.WriteFile(f, []byte(k), 0666)

}



func main() {
	//pairname := flag.String("Keyname", "~/.ssh/", "Enter the name of the")
	pairname = "testWB"
	svc = ec2.New(session.New(&aws.Config{Region: aws.String("us-west-1")}))

	keyresult, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(pairname),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.Duplicate" {
			exitErrorf("Keypair %q already exists.", pairname)
		}
		exitErrorf("Unable to create key pair: %s, %v.", pairname, err)
	}

	savePem(pairname+".pem", *keyresult.KeyMaterial)

	fmt.Printf("Created key pair %q %s\n%s\n",
		*keyresult.KeyName, *keyresult.KeyFingerprint,
		*keyresult.KeyMaterial)

}
