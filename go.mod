module github.com/blukai/auth.svc

go 1.12

require (
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-contrib/sessions v0.0.0-20190226023029-1532893d996f
	github.com/gin-gonic/gin v1.4.0
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/markbates/goth v1.51.0
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
