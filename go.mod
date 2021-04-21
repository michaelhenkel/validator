module github.com/michaelhenkel/validator

go 1.15

replace ssd-git.juniper.net/contrail/cn2/third_party/apiserver-builder-alpha => ../../../ssd-git.juniper.net/contrail/cn2/third_party/apiserver-builder-alpha

replace github.com/kubernetes-incubator/reference-docs => github.com/kubernetes-incubator/reference-docs v0.0.0-20170929004150-fcf65347b256

replace github.com/markbates/inflect => github.com/markbates/inflect v1.0.4

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1

replace k8s.io/apimachinery => k8s.io/apimachinery v0.20.1

replace k8s.io/client-go => k8s.io/client-go v0.20.2

replace k8s.io/api => k8s.io/api v0.20.2

replace k8s.io/kubectl => k8s.io/kubectl v0.20.2

require (
	github.com/spf13/cobra v1.1.3
	k8s.io/apimachinery v0.21.0
	k8s.io/apiserver v0.20.1
	k8s.io/client-go v0.21.0
	k8s.io/kubectl v0.21.0
	ssd-git.juniper.net/contrail/cn2/contrail v0.0.0-20210420065613-ff1f72447bc8
)
