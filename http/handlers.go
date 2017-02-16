package http

import (
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func optionHandler(w http.ResponseWriter, r *http.Request) {
	log, err := getLog(r)
	if err != nil {
		errorWriter(w, r, http.StatusInternalServerError, err)
		return
	}
	if r.Method != http.MethodOptions {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resourceName := chi.URLParam(r, "resource")
	var (
		resourceEP ResourceEndpoint
		ok         bool
	)
	if resourceEP, ok = resourceEndpoints[resourceName]; !ok {
		log.Warn("Resource not found", zap.String("resource-name", resourceName))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var methods = http.MethodOptions
	if _, ok := resourceEP.(ResourceEndpointNew); ok {
		methods = methods + "," + http.MethodPost
	}
	if _, ok := resourceEP.(ResourceEndpointFind); ok {
		methods = methods + "," + http.MethodGet + "," + http.MethodHead
	}
	if _, ok := resourceEP.(ResourceEndpointUpdate); ok {
		methods = methods + "," + http.MethodPut + "," + http.MethodPatch
	}
	if _, ok := resourceEP.(ResourceEndpointDelete); ok {
		methods = methods + "," + http.MethodDelete
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Allow", methods)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	log, err := getLog(r)
	if err != nil {
		errorWriter(w, r, http.StatusInternalServerError, err)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var (
		resourceEP   ResourceEndpoint
		resourceID   int64
		resourceName string
		ok           bool
	)

	resourceName = chi.URLParam(r, "resource")
	{
		resourceIDString := chi.URLParam(r, "snowflake")
		resourceID = -1
		if resourceIDString != "" {
			resourceID, err = strconv.ParseInt(resourceIDString, 10, 63)
			if err != nil {
				errorWriter(w, r, http.StatusInternalServerError, err)
				return
			}
		}
	}

	if resourceEP, ok = resourceEndpoints[resourceName]; !ok {
		log.Warn("Resource not found", zap.String("resource-name", resourceName))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if find, ok := resourceEP.(ResourceEndpointFind); ok {

		if r.Method == http.MethodHead {
			// If the client did a HEAD request we give them only the headers and no body
			// This saves some bandwidth on the network (but not any CPU on the server)
			w.WriteHeader(http.StatusOK)
			return
		}

		if resourceID > 0 {
			resource, err := find.Find(resourceID)
			if err != nil {
				errorWriter(w, r, http.StatusInternalServerError, err)
				return
			}
			render.JSON(w, r, resource)
		} else if resourceID == -1 {
			page, size := r.Context().Value(ctxPageKey).(int64), r.Context().Value(ctxSizeKey).(int64)
			resources, err := find.FindAll(page, size)
			if err != nil {
				errorWriter(w, r, http.StatusInternalServerError, err)
				return
			}
			render.JSON(w, r, resources)
		}
		errorStringWriter(w, r, http.StatusInternalServerError, "The requested resource ID is not valid")
	} else {
		errorStringWriter(w, r, http.StatusBadRequest, "The requested resource is not readable")
	}

}
