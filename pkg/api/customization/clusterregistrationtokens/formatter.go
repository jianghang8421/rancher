package clusterregistrationtokens

import (
	"fmt"
	"net/url"

	"github.com/rancher/norman/api/access"
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/image"
	"github.com/rancher/rancher/pkg/settings"
	"github.com/rancher/rancher/pkg/systemtemplate"
	"github.com/rancher/types/apis/management.cattle.io/v3/schema"
	"github.com/rancher/types/client/management/v3"
)

const (
	commandFormat         = "kubectl apply -f %s"
	insecureCommandFormat = "curl --insecure -sfL %s | kubectl apply -f -"
	nodeCommandFormat     = "sudo docker run -d --privileged --restart=unless-stopped --net=host -v /etc/kubernetes:/etc/kubernetes -v /var/run:/var/run %s --server %s --token %s%s"

	windowsNodeCommandFormat = `PowerShell -Sta -NoLogo -NonInteractive -Command "& {docker run --rm -v C:/:C:/host --isolation hyperv %s -server %s -token %s%s; if($?){& c:/etc/rancher/run.ps1;}}"`
)

func Formatter(request *types.APIContext, resource *types.RawResource) {
	var (
		caNonWindows = ""
		caWindows    = ""
	)
	ca := systemtemplate.CAChecksum()
	if ca != "" {
		caNonWindows = " --ca-checksum " + ca
		caWindows = " -caChecksum " + ca
	}

	token, _ := resource.Values["token"].(string)
	if token != "" {
		var clusterArch string
		var crtList []client.ClusterRegistrationToken
		access.List(request, &schema.Version, client.ClusterRegistrationTokenType, &types.QueryOptions{}, &crtList)
		for _, crt := range crtList {
			if crt.Token == token {
				var cluster client.Cluster
				access.ByID(request, &schema.Version, client.ClusterType, crt.ClusterID, &cluster)
				clusterArch = cluster.Arch
				break
			}
		}
		url := getURL(request, token, clusterArch)
		resource.Values["insecureCommand"] = fmt.Sprintf(insecureCommandFormat, url)
		resource.Values["command"] = fmt.Sprintf(commandFormat, url)
		resource.Values["nodeCommand"] = fmt.Sprintf(nodeCommandFormat,
			image.Resolve(settings.AgentImage.Get()),
			getRootURL(request),
			token,
			caNonWindows)

		resource.Values["armNodeCommand"] = fmt.Sprintf(nodeCommandFormat,
			image.Resolve(settings.ArmAgentImage.Get()),
			getRootURL(request),
			token,
			caNonWindows)

		resource.Values["token"] = token
		resource.Values["manifestUrl"] = url
		resource.Values["windowsNodeCommand"] = fmt.Sprintf(windowsNodeCommandFormat,
			image.Resolve(settings.WindowsAgentImage.Get()),
			getRootURL(request),
			token,
			caWindows)
	}
}

func NodeCommand(token, architecture string) string {
	ca := systemtemplate.CAChecksum()
	if ca != "" {
		ca = " --ca-checksum " + ca
	}

	switch architecture {
	case "amd64":
		return fmt.Sprintf(nodeCommandFormat,
			image.Resolve(settings.AgentImage.Get()),
			getRootURL(nil),
			token,
			ca)
	case "arm64":
		return fmt.Sprintf(nodeCommandFormat,
			image.Resolve(settings.ArmAgentImage.Get()),
			getRootURL(nil),
			token,
			ca)
	default:
		return fmt.Sprintf(nodeCommandFormat,
			image.Resolve(settings.AgentImage.Get()),
			getRootURL(nil),
			token,
			ca)
	}
}

func getRootURL(request *types.APIContext) string {
	serverURL := settings.ServerURL.Get()
	if serverURL == "" {
		if request != nil {
			serverURL = request.URLBuilder.RelativeToRoot("")
		}
	} else {
		u, err := url.Parse(serverURL)
		if err != nil {
			if request != nil {
				serverURL = request.URLBuilder.RelativeToRoot("")
			}
		} else {
			u.Path = ""
			serverURL = u.String()
		}
	}

	return serverURL
}

func getURL(request *types.APIContext, token, clusterArch string) string {
	var path string
	switch clusterArch {
	case "arm64":
		path = "/v3/import/arm64/" + token + ".yaml"
	default:
		path = "/v3/import/" + token + ".yaml"
	}
	serverURL := settings.ServerURL.Get()
	if serverURL == "" {
		serverURL = request.URLBuilder.RelativeToRoot(path)
	} else {
		u, err := url.Parse(serverURL)
		if err != nil {
			serverURL = request.URLBuilder.RelativeToRoot(path)
		} else {
			u.Path = path
			serverURL = u.String()
		}
	}

	return serverURL
}
