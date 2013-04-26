// Test goamz
// https://wiki.ubuntu.com/goamz
// http://godoc.org/launchpad.net/goamz
// AWS_ACCESS_KEY_ID
// AWS_SECRET_ACCESS_KEY
// http://docs.aws.amazon.com/AWSEC2/latest/APIReference/ApiReference-query-RequestSpotInstances.html

package main

import (
	"fmt"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"os"
)

func main() {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err)
	}
	s := s3.New(auth, aws.USEast)
	pre := "tianl/Glaucoma"
	b := s.Bucket("caglifeng1")
	l, _ := b.List(pre, "", "", 5000)
	fmt.Printf("%+v\n", l)

	cycle := 0
	for l.IsTruncated {
		ll := len(l.Contents)
		fmt.Fprintf(os.Stderr, "Contents[%d].Key=%s\n", ll, l.Contents[ll-1].Key)
		for i, v := range l.Contents {
			fmt.Printf("%d Key:%s\n", i+cycle*1000, v.Key)
		}

		l, _ = b.List(pre, "", l.Contents[ll-1].Key, 5000)
		cycle++
	}

	for i, v := range l.Contents {
		fmt.Printf("%d Key:%s\n", i+cycle*1000, v.Key)
	}
}
