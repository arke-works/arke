package handlers

import (
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"io/ioutil"
	"iris.arke.works/forum/http/ctxkeys"
	"iris.arke.works/forum/http/helper"
	"iris.arke.works/forum/http/resources"
	"iris.arke.works/forum/snowflakes"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var (
		resourceEP      resources.ResourceEndpoint
		resourceName    string
		resourceFactory resources.ResourceFactory
		fountain        snowflakes.Fountain
		ok              bool
		err             error
	)

	fountain, ok = r.Context().Value(ctxkeys.CtxFountainKey).(snowflakes.Fountain)
	if !ok || fountain == nil {
		helper.ErrorStringWriter(w, r, http.StatusInternalServerError, "No ID Fountain")
		return
	}

	resourceName = chi.URLParam(r, "resource")
	if resourceEP, err = resources.GetEndpoint(resourceName); err != nil {
		helper.ErrorWriter(w, r, http.StatusNotFound, err)
		return
	}

	if resourceFactory, err = resources.GetResourceFactory(resourceName); err != nil {
		helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
		return
	}

	if newRes, ok := resourceEP.(resources.ResourceEndpointNew); ok {
		res, err := resourceFactory(nil)
		if err != nil {
			helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
			return
		}
		reqData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helper.ErrorWriter(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()
		err = res.UnmarshalJSON(reqData)
		if err != nil {
			helper.ErrorWriter(w, r, http.StatusBadRequest, err)
			return
		}
		err = res.StripReadOnly()
		if err != nil {
			helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
			return
		}
		storeResource, err := resourceFactory(fountain)
		if err != nil {
			helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
			return
		}
		storeResource.Merge(res)
		err = newRes.New(storeResource)
		if err != nil {
			helper.ErrorWriter(w, r, http.StatusInternalServerError, err)
			return
		}
		render.JSON(w, r, newRes)
		return
	}
	helper.ErrorStringWriter(w, r, http.StatusBadRequest, "The requested resource is not creatable")
}
