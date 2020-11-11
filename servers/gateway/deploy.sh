# calls the build.sh script
bash build.sh
ssh -i ~/.ssh/MyKeyPair.pem ec2-user@ec2-3-12-128-39.us-east-2.compute.amazonaws.com < toPipe.sh