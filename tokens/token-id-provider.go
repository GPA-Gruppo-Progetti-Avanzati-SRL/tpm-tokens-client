package tokens

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"strings"
)

type DefaultTokenIdProvider struct {
}

func (p *DefaultTokenIdProvider) NewId(ctxId string, unique bool, act *Action) (string, error) {
	oid := util.NewObjectId().String()
	return strings.Join([]string{ctxId, oid}, ":"), nil
}

func (p *DefaultTokenIdProvider) Close() error {
	return nil
}
