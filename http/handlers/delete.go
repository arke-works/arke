package handlers

import (
	"github.com/pressly/chi"
	"go.uber.org/zap"
	"iris.arke.works/forum/http/helper"
	"iris.arke.works/forum/http/resources"
	"iris.arke.works/forum/snowflakes"
	"net/http"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	log, err := helper.GetLog(r)
	if err != nil {
		helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
		return
	}
	if r.Method != http.MethodDelete {
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

	if del, ok := resourceEP.(resources.ResourceEndpointDelete); ok {
		if resourceID >= 0 {
			del.SoftDelete(resourceID)
		} else {
			del.HardDelete(resourceID * -1)
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
