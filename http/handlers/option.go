package handlers

import (
	"github.com/pressly/chi"
	"go.uber.org/zap"
	"net/http"
	"iris.arke.works/forum/http/helper"
	"iris.arke.works/forum/http/resources"
)

func OptionHandler(w http.ResponseWriter, r *http.Request) {
	log, err := helper.GetLog(r)
	if err != nil {
		helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
		return
	}
	if r.Method != http.MethodOptions {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resourceName := chi.URLParam(r, "resource")
	var resourceEP resources.ResourceEndpoint

	if resourceEP, err = resources.GetEndpoint(resourceName); err != nil {
		log.Warn("Resource not found", zap.String("resource-name", resourceName))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Debug("Checking for resource options")
	var methods = http.MethodOptions
	if _, ok := resourceEP.(resources.ResourceEndpointNew); ok {
		log.Debug("Resource has POST Method")
		methods = methods + "," + http.MethodPost
	}
	if _, ok := resourceEP.(resources.ResourceEndpointFind); ok {
		log.Debug("Resource has GET/HEAD Method")
		methods = methods + "," + http.MethodGet + "," + http.MethodHead
	}
	if _, ok := resourceEP.(resources.ResourceEndpointUpdate); ok {
		log.Debug("Resource has PUT/PATCH method")
		methods = methods + "," + http.MethodPut + "," + http.MethodPatch
	}
	if _, ok := resourceEP.(resources.ResourceEndpointDelete); ok {
		log.Debug("Resource has DELETE method")
		methods = methods + "," + http.MethodDelete
	}
	w.Header().Add("Allow", methods)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
