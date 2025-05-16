package sciencedirect

import "github.com/go-kratos/kratos/v2/log"

type ScienceDirect struct {
	log *log.Helper
}

func NewScienceDirect() *ScienceDirect {
	return &ScienceDirect{
		log: log.NewHelper(log.DefaultLogger),
	}
}

func (s *ScienceDirect) List() error {
	return nil
}

func (s *ScienceDirect) Detail() error {
	return nil
}
