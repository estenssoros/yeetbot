set -x

GOOS=linux go build -o app
if [ $? -ne 0 ]; then
    echo ERROR: failed to build go binary
    exit 1
fi

rm lambda.zip
zip lambda.zip ./app
aws lambda update-function-code --function-name yeetbot --zip-file fileb://lambda.zip

rm lambda.zip   
rm app