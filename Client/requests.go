package main

type AllocateRequest struct {
	Id              string `json:"id"`
	Target          int    `json:"target"`
	RedisHost       string `json:"redisHost"`
	RedisPort       int    `json:"redisPort"`
	RedisPassword   string `json:"redisPassword"`
	RegisterChannel string `json:"registerChannel"`
}

type BootChannelMessage struct {
	Id              string `json:"id"`
	ListenerChannel string `json:"listenerChannel"`
	Method          string `json:"method"`
}
