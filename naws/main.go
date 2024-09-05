package main

import (
	"net/http"

	"naws/lambda"
	"naws/utils"
)

func main() {
	http.HandleFunc("/", utils.LambdaProxy(lambda.User))
	http.HandleFunc("/posts", utils.LambdaProxy(lambda.Posts))
	http.ListenAndServe(":8080", nil)
}
