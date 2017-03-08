package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	flag "github.com/ogier/pflag"
)

var profileName = flag.StringP("profile", "p", "", "AWS profile to use")
var sessionName = flag.StringP("session-name", "n", "", "Name for the session")
var roleArn = flag.StringP("role-arn", "r", "", "ARN for the role to assume")
var duration = flag.Int64P("duration", "d", 3600, "Validity period in seconds")

// Determine the path of AWS shared credentials file to save to.
func credsFilePath() (string, error) {
	if path := os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); path != "" {
		return path, nil
	}

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return "", errors.New("Unable to find home directory.")
	}

	return filepath.Join(homeDir, ".aws", "credentials"), nil
}

// Save the new STS credentials to the shared credentials file.
func saveCreds(role *sts.AssumeRoleOutput, credsPath string) error {
	cfg, err := ini.Load(credsPath)
	if err != nil {
		return err
	}

	profile := cfg.Section(*profileName)
	profile.NewKey("aws_access_key_id", *role.Credentials.AccessKeyId)
	profile.NewKey("aws_secret_access_key", *role.Credentials.SecretAccessKey)
	profile.NewKey("aws_session_token", *role.Credentials.SessionToken)
	profile.NewKey("aws_security_token", *role.Credentials.SessionToken)

	return cfg.SaveTo(credsPath)
}

func main() {
	flag.Parse()
	if *sessionName == "" || *roleArn == "" {
		fmt.Println("Both session-name and role-arn parameters are requried.")
		os.Exit(1)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           *profileName,
	}))

	svc := sts.New(sess)

	roleInput := &sts.AssumeRoleInput{
		DurationSeconds: duration,
		RoleArn:         roleArn,
		RoleSessionName: sessionName,
	}

	role, err := svc.AssumeRole(roleInput)
	if err != nil {
		fmt.Println(err)
		os.Exit(128)
	}

	credsPath, err := credsFilePath()
	if err != nil {
		fmt.Println(err)
		os.Exit(256)
	}

	err = saveCreds(role, credsPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(256)
	}

	fmt.Printf("Credentials saved to %q profile. Expiration: %v\n",
		*profileName, role.Credentials.Expiration)
}
