package neon

func InitNeon() *Service {
	service := NewClient()
	service = InitDB(service)
	return service
}
