package registry

import regv1 "github.com/naivary/codemark/api/registry/v1"

func Merge(regs ...regv1.Registry) (regv1.Registry, error) {
	reg := InMemory()
	for _, re := range regs {
		for _, opt := range re.All() {
			if err := reg.Define(opt); err != nil {
				return nil, err
			}
		}
	}
	return reg, nil
}
