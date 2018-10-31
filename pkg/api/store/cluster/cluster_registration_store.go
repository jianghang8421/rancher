package cluster

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/randomtoken"
	"github.com/sirupsen/logrus"
)

type RegistrationTokenStore struct {
	types.Store
}

func (r *RegistrationTokenStore) Create(apiContext *types.APIContext, schema *types.Schema, data map[string]interface{}) (map[string]interface{}, error) {
	if data != nil {
		token, err := randomtoken.Generate()
		if err != nil {
			return nil, err
		}
		data["token"] = token

		logrus.Infof("jianghang %s", apiContext)
		logrus.Infof("jianghang %s", schema)
		logrus.Infof("jianghang %s", data)
	}

	return r.Store.Create(apiContext, schema, data)
}
