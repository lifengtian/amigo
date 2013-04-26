// Test goamz
// https://wiki.ubuntu.com/goamz
// http://godoc.org/launchpad.net/goamz
// AWS_ACCESS_KEY_ID
// AWS_SECRET_ACCESS_KEY
// http://docs.aws.amazon.com/AWSEC2/latest/APIReference/ApiReference-query-RequestSpotInstances.html

package main

import (
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
)

func main() {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err)
	}
	e := ec2.New(auth, aws.USEast)

	options := ec2.RunInstances{
		ImageId:      "ami-ccf405a5", // Ubuntu Maverick, i386, EBS store
		InstanceType: "t1.micro",
	}
	resp, err := e.RunInstances(&options)
	if err != nil {
		panic(err)
	}

	for _, instance := range resp.Instances {
		println("Now running", instance.InstanceId)
	}
	println("Make sure you terminate instances to stop the cash flow.")
}
