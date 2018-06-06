package main

var (
	managerUri = "http://127.0.0.1:3303/"
	tokenHeaderName = "WMQX_API_TOKEN" // wmqx.toml [api]
	token = "guest" // wmqx.toml [api]
)


// example

func main()  {

	message := map[string]string{
		"name": "ada",
		"comment": "test",
		"durable": "1",
		"is_need_token": "1",
		"mode": "topic",
		"token": "test_token",
	}

	// add a message
	addMessage(message)

	consumer := map[string]string{
		"url": "http://127.0.0.1:80/test.php",
		"route_key": "test222",
		"timeout": "2000",
		"code": "200",
		"check_code": "1",
		"comment": "test consumer",
	}

	// add a consumer to message
	addConsumer(consumer)
}