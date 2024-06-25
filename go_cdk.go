package main

import (
	"go_cdk/stack"
	"go_cdk/config"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	stack.NewGoCdkStack(app, "GoCdkStack", &stack.GoCdkStackProps{
		StackProps: awscdk.StackProps{
			Env: config.Env(),
		},
	})

	app.Synth(nil)
}