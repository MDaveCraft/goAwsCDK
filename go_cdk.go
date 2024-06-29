package main

import (
	"goCdk/stack"
	"goCdk/config"
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()
	app := cdk.NewApp(nil)

	stack.NewGoCdkStack(app, "GoCdkStack", &stack.GoCdkStackProps{
		StackProps: cdk.StackProps{ 
			Env: config.Env(),
		},
	})

	app.Synth(nil)
}