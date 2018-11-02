package clusterregistrationtokens

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/urlbuilder"
	"github.com/rancher/rancher/pkg/image"
	"github.com/rancher/rancher/pkg/settings"
	"github.com/rancher/rancher/pkg/systemtemplate"
	"github.com/rancher/types/apis/management.cattle.io/v3/schema"
)

func ClusterImportHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain")
	token := mux.Vars(req)["token"]
	arch := req.URL.Query().Get("arch")

	logrus.Infof("jianghang urlquery arch %s", arch)

	urlBuilder, err := urlbuilder.New(req, schema.Version, types.NewSchemas())
	if err != nil {
		resp.WriteHeader(500)
		resp.Write([]byte(err.Error()))
		return
	}

	url := urlBuilder.RelativeToRoot("")
	var agentimage string
	switch arch {
	case "arm64":
		agentimage = settings.ArmAgentImage.Get()
	default:
		agentimage = settings.AgentImage.Get()
	}

	if err := systemtemplate.SystemTemplate(resp, image.Resolve(agentimage), token, url); err != nil {
		resp.WriteHeader(500)
		resp.Write([]byte(err.Error()))
	}
}
