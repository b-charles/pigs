package confsources

type SimpleConfigSource struct {
	Priority int
	Env      map[string]string
}

func (self *SimpleConfigSource) GetPriority() int {
	return self.Priority
}

func (self *SimpleConfigSource) LoadEnv() map[string]string {
	return self.Env
}
