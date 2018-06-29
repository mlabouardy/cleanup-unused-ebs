package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func listAvailableVolumes(cfg aws.Config) ([]string, error) {
	volumes := make([]string, 0)
	svc := ec2.New(cfg)
	req := svc.DescribeVolumesRequest(&ec2.DescribeVolumesInput{
		Filters: []ec2.Filter{
			ec2.Filter{
				Name:   aws.String("status"),
				Values: []string{"available"},
			},
		},
	})
	res, err := req.Send()
	if err != nil {
		return volumes, err
	}

	for _, volume := range res.Volumes {
		volumes = append(volumes, *volume.VolumeId)
	}

	return volumes, nil
}

func deleteMovie(cfg aws.Config, volumeId string) error {
	svc := ec2.New(cfg)
	req := svc.DeleteVolumeRequest(&ec2.DeleteVolumeInput{
		VolumeId: &volumeId,
	})
	_, err := req.Send()
	return err
}

func handler() error {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}

	volumes, err := listAvailableVolumes(cfg)
	if err != nil {
		return err
	}

	for _, volume := range volumes {
		log.Println("Removing volume: ", volume)
		if err != deleteMovie(cfg, volume) {
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
