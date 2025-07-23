package registry

func Merge(regs ...Registry) (Registry, error) {
	r := InMemory()
	for _, reg := range regs {
		for _, option := range reg.All() {
			if err := r.Define(option); err != nil {
				return nil, err
			}
		}
	}
	return r, nil
}
