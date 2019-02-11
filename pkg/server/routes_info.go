package server

import "github.com/sirupsen/logrus"

func PrintRoutesInfo(server *Server, logger *logrus.Entry) {
	for _, route := range server.Routes() {
		logger.Infof("route -> path=%+v, method=%+v", route.Path, route.Method)
	}
}
