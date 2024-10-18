package ip_white

import (
	"sync"
)

type option struct {
	WhiteList []string
	sync.Mutex
}

type Option func(*option)

func WithIpWhite(ips []string) Option {
	return func(o *option) {
		o.WhiteList = ips
	}
}

//type option struct {
//	WhiteList []string
//	*sync.Mutex
//}
//
//// Option specifies instrumentation configuration options.
//type Option interface {
//	apply(*option)
//}
//
//type optionFunc func(*option)
//
//func (o optionFunc) apply(c *option) {
//	o(c)
//}
//
//// WithIpWhite  ip white
//func WithIpWhite(ips []string) Option {
//	return optionFunc(func(cfg *option) {
//		cfg.WhiteList = ips
//	})
//}
