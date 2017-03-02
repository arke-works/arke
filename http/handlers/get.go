package handlers

import (
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"go.uber.org/zap"
	"iris.arke.works/forum/http/ctxkeys"
	"iris.arke.works/forum/http/helper"
	"iris.arke.works/forum/http/resources"
	"iris.arke.works/forum/snowflakes"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	log, err := helper.GetLog(r)
	if err != nil {
		helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var (
		resourceEP   resources.ResourceEndpoint
		resourceID   int64
		resourceName string
	)

	resourceName = chi.URLParam(r, "resource")
	{
		resourceIDString := chi.URLParam(r, "snowflake")
		resourceID = -1
		if resourceIDString != "" {
			resourceID, err = snowflakes.EncodedToID(resourceIDString)
			if err != nil {
				helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
				return
			}
		}
	}

	if resourceEP, err = resources.GetEndpoint(resourceName); err != nil {
		log.Warn("Resource not found", zap.String("resource-name", resourceName))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if find, ok := resourceEP.(resources.ResourceEndpointFind); ok {

		if r.Method == http.MethodHead {
			// If the client did a HEAD request we give them only the headers and no body
			// This saves some bandwidth on the network (but not any CPU on the server)
			w.WriteHeader(http.StatusOK)
			return
		}

		if resourceID > 0 {
			resource, err := find.Find(resourceID)
			if err != nil {
				helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
				return
			}
			render.JSON(w, r, resource)
			return
		} else if resourceID == -1 {
			pivot, size := r.Context().Value(ctxkeys.CtxPivotIDKey).(int64), r.Context().Value(ctxkeys.CtxSizeKey).(int64)
			resources, err := find.FindAll(pivot, size)
			if err != nil {
				helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
				return
			}
			render.JSON(w, r, resources)
			return
		}
		helper.ErrorStringWriter(w, r, http.StatusInternalServerError, "The requested resource ID is not valid")
		return
	} else {
		helper.ErrorStringWriter(w, r, http.StatusBadRequest, "The requested resource is not readable")
		return
	}

}
