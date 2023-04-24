# S3 Unzipper 
## Lambda function
This lambda function unzipped a zip file that was uploaded to the AWS S3 bucket. The zip file content file(s) content-type is preserved.

###### Instructions
- Create S3 bucket 
- Create `Go` lambda function 
- Build application locally into a zip file using 
`GOOS=linux GOARCH=amd64 go build -o main && zip main.zip main`
- Upload the main.zip file to the lambda function code
- Add S3 trigger to the lambda function with Prefix: `/folderName` and Suffix: `.zip`
- Add `AmazonS3FullAccess` permissions policy to the lambda function role name

