package config

type Service struct {
	Name string
	Port string
}

var Services = struct {
	Gateway     Service
	Order       Service
	Inventory   Service
	Payment     Service
	Notification Service
	Saga        Service
}{
	Gateway: Service{
		Name: "gateway",
		Port: ":8080",
	},
	Order: Service{
		Name: "order",
		Port: ":8081",
	},
	Inventory: Service{
		Name: "inventory",
		Port: ":8082",
	},
	Payment: Service{
		Name: "payment",
		Port: ":8083",
	},
	Notification: Service{
		Name: "notification",
		Port: ":8084",
	},
	Saga : Service{
		Name : "saga",
		Port: ":8086",
	},
}