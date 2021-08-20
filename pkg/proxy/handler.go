package proxy

import (
	"errors"
	"net/http"

	"github.com/justcompile/tnl/pkg/socketserver"
	"github.com/justcompile/tnl/pkg/types"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Hub *socketserver.Hub
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := req.Host

	client := h.Hub.GetClientForDomain(host)
	if client == nil {
		serverError(w, errors.New("proxy not found"))
		return
	}

	msg, err := types.SerializeRequest(req)
	if err != nil {
		log.Error(err.Error())
		serverError(w, err)
		return
	}

	if err := client.Send(msg); err != nil {
		log.Error(err.Error())
		serverError(w, err)
		return
	}

	resp := <-client.Messages

	for key, value := range resp.Headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(resp.Status)

	_, _ = w.Write(resp.Body)
}

func serverError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadGateway)
	_, _ = w.Write([]byte(err.Error()))
}
