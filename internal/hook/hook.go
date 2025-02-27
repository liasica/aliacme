// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-08, by liasica

package hook

import (
	"sync"

	"go.uber.org/zap"

	"github.com/liasica/autoacme/internal/g"
)

type Hook struct {
	wg          *sync.WaitGroup
	do          *g.Domain
	privateKey  []byte
	certificate []byte
}

// NewHook 创建 Hook
func NewHook(cfg *g.Domain, priv, cert []byte) *Hook {
	return &Hook{
		wg:          &sync.WaitGroup{},
		do:          cfg,
		privateKey:  priv,
		certificate: cert,
	}
}

// Run 运行 Hook
func (h *Hook) Run() {
	for i := 0; i < len(h.do.Hooks); i++ {
		h.wg.Add(1)
		hook := h.do.Hooks[i]
		switch hook.Name {
		case g.DomainHookNameCDN:
			if hook.CDNHook == nil {
				zap.S().Error("CDN hook is not configured")
				continue
			}
			go h.AliyunCDN(hook.CDNHook)

		case g.DomainHookNameQiniuSSL:
			if hook.QiniuSSLHook == nil {
				zap.S().Error("Qiniu SSL hook is not configured")
				continue
			}
			go h.QiniuSSL(hook.QiniuSSLHook)

		default:
			zap.S().Error("unknown hook", zap.String("hook", string(hook.Name)))
		}
	}
	h.wg.Wait()
}
