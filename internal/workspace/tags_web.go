package workspace

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tofutf/tofutf/internal/http/decode"
	"github.com/tofutf/tofutf/internal/http/html"
	"github.com/tofutf/tofutf/internal/http/html/paths"
)

func (h *webHandlers) addTagHandlers(r *mux.Router) {
	r = html.UIRouter(r)

	r.HandleFunc("/workspaces/{workspace_id}/create-tag", h.createTag).Methods("POST")
	r.HandleFunc("/workspaces/{workspace_id}/delete-tag", h.deleteTag).Methods("POST")
}

func (h *webHandlers) createTag(w http.ResponseWriter, r *http.Request) {
	var params struct {
		WorkspaceID *string `schema:"workspace_id,required"`
		TagName     *string `schema:"tag_name,required"`
	}
	if err := decode.All(&params, r); err != nil {
		h.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err := h.client.AddTags(r.Context(), *params.WorkspaceID, []TagSpec{{Name: *params.TagName}})
	if err != nil {
		h.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html.FlashSuccess(w, "created tag: "+*params.TagName)
	http.Redirect(w, r, paths.Workspace(*params.WorkspaceID), http.StatusFound)
}

func (h *webHandlers) deleteTag(w http.ResponseWriter, r *http.Request) {
	var params struct {
		WorkspaceID *string `schema:"workspace_id,required"`
		TagName     *string `schema:"tag_name,required"`
	}
	if err := decode.All(&params, r); err != nil {
		h.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err := h.client.RemoveTags(r.Context(), *params.WorkspaceID, []TagSpec{{Name: *params.TagName}})
	if err != nil {
		h.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html.FlashSuccess(w, "removed tag: "+*params.TagName)
	http.Redirect(w, r, paths.Workspace(*params.WorkspaceID), http.StatusFound)
}
