module github.com/michaelhenkel/validator

go 1.15

replace ssd-git.juniper.net/contrail/cn2/third_party/apiserver-builder-alpha => ../../../ssd-git.juniper.net/contrail/cn2/third_party/apiserver-builder-alpha

replace ssd-git.juniper.net/contrail/cn2/deployer => ../../../ssd-git.juniper.net/contrail/cn2/deployer

replace ssd-git.juniper.net/contrail/cn2/contrail => ../../../ssd-git.juniper.net/contrail/cn2/contrail

replace github.com/kubernetes-incubator/reference-docs => github.com/kubernetes-incubator/reference-docs v0.0.0-20170929004150-fcf65347b256

replace github.com/markbates/inflect => github.com/markbates/inflect v1.0.4

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1

replace k8s.io/apimachinery => k8s.io/apimachinery v0.20.1

replace k8s.io/client-go => k8s.io/client-go v0.20.2

replace k8s.io/api => k8s.io/api v0.20.2

replace k8s.io/kubectl => k8s.io/kubectl v0.20.2

require (
	github.com/spf13/cobra v1.1.3
	gonum.org/v1/gonum v0.0.0-20190331200053-3d26580ed485
	gopkg.in/ini.v1 v1.51.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
	ssd-git.juniper.net/contrail/cn2/contrail v1.0.0
	ssd-git.juniper.net/contrail/cn2/deployer v0.0.0-00010101000000-000000000000
)
