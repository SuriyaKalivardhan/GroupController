package main

type AllocateRequest struct {
	ControllerId   string `json:"controllerId"`
	DesiredWorkers int    `json:"desiredWorkers"`
	RedisAddress   string `json:"redisAddress"`
	RedisPassword  string `json:"redisPassword"`
	RedisUseSSL    bool   `json:"redisUseSSL"`
}

type WorkerMessage struct {
	WorkerId             string `json:"workerId"`
	ListenerChannel      string `json:"listenerChannel"`
	Controller           string `json:"controller"`
	PlatformFaultDomain  string `json:"platformFaultDomain"`
	PlatformUpdateDomain string `json:"platformUpdateDomain"`
	Zone                 string `json:"zone"`
}

type ControlMessage struct {
	ControllerId  string `json:"controllerId"`
	Method        string `json:"method"`
	RedisAddress  string `json:"redisAddress"`
	RedisPassword string `json:"redisPassword"`
	RedisUseSSL   bool   `json:"redisUseSSL"`
}
