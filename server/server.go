package server

import (
    "github.com/kataras/iris"
    "github.com/BluePecker/JwtAuth/server/router"
    "fmt"
    "github.com/Sirupsen/logrus"
    //"net/url"
    //"github.com/kataras/iris/core/host"
)

type TLS struct {
    Cert string
    Key  string
}

type Options struct {
    Host string
    Port int
    Tls  *TLS
}

type Server struct {
    app *iris.Application
}

func (s *Server) initHttpApp() {
    if s.app == nil {
        s.app = iris.New()
    }
}

func (s *Server) Accept(options Options) {
    s.initHttpApp()
    
    if options.Tls != nil {
        //是否将80端口的请求转发到443
        //target, _ := url.Parse("https://127.0.0.1:443")
        //go host.NewProxy("127.0.0.1:80", target).ListenAndServe()
        var addr string = fmt.Sprintf("%s:%s", options.Host, options.Port)
        if err := s.app.Run(iris.TLS(addr, options.Tls.Cert, options.Tls.Key)); err != nil {
            logrus.Error(err)
        }
        
    } else {
        var addr string = fmt.Sprintf("%s:%s", options.Host, options.Port)
        if err := s.app.Run(iris.Addr(addr)); err != nil {
            logrus.Error(err)
        }
    }
}

func (s *Server) AddRouter(routers... router.Router) {
    s.initHttpApp()
    
    for _, item := range routers {
        item.Routes(s.app)
    }
}