package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

// Request の構造体定義
type PostHelloRequest struct {
	Name string `json:"name"`
}

// Response の構造体定義
type HelloMessageResponse struct {
	Message string `json:"message"`
}

// バリデーション設定
var ValidateHelloMessageSettings = []*ValidatorSetting{
	{ArgName: "name", ValidateTags:"required"},
}

func PostHello(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// バリデーション
	validErr := ValidateBody(request.Body, ValidateHelloMessageSettings)
	if validErr != nil {
		return Response400(*validErr)
	}
	// HTTP ボディ部の JSON を受け取る
	body := request.Body

	// JSON から構造体に変換する
	var req PostHelloRequest
	err := json.Unmarshal([]byte(body), &req)
	if err != nil {
		return Response500(err)
	}

	// レスポンスのメッセージを作成
	msg := fmt.Sprintf("Hello!%s", req.Name)
	res := &HelloMessageResponse{Message:msg}

	return Response200(res)
}


