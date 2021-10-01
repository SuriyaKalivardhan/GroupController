package main

type AllocateRequest struct {
	Id            string `json:"id"`
	Target        int    `json:"target"`
	RedisHost     string `json:"redisHost"`
	RedisPassword string `json:"redisPassword"`
}

type BootChannelMessage struct {
	Id              string `json:"id"`
	ListenerChannel string `json:"listenerChannel"`
	BindedClient    string `json:"client"`
}
