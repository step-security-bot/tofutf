package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leg100/go-tfe"
	"github.com/leg100/jsonapi"
	"github.com/leg100/ots"
)

func (s *Server) CreateRun(w http.ResponseWriter, r *http.Request) {
	opts := tfe.RunCreateOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, &opts); err != nil {
		WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	obj, err := s.RunService.CreateRun(&opts)
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	WriteResponse(w, r, obj, WithCode(http.StatusCreated))
}

func (s *Server) ApplyRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	opts := &tfe.RunApplyOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, opts); err != nil {
		WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := s.RunService.ApplyRun(vars["id"], opts); err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) GetRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	obj, err := s.RunService.GetRun(vars["id"])
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	WriteResponse(w, r, obj)
}

func (s *Server) ListRuns(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var opts tfe.RunListOptions
	if err := DecodeQuery(&opts, r.URL.Query()); err != nil {
		WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	SanitizeListOptions(&opts.ListOptions)

	obj, err := s.RunService.ListRuns(vars["workspace_id"], opts)
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	WriteResponse(w, r, obj)
}

func (s *Server) DiscardRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	opts := &tfe.RunDiscardOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, opts); err != nil {
		WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	err := s.RunService.DiscardRun(vars["id"], opts)
	if err == ots.ErrRunDiscardNotAllowed {
		WriteError(w, http.StatusConflict, err)
		return
	} else if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) CancelRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	opts := &tfe.RunCancelOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, opts); err != nil {
		WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	err := s.RunService.CancelRun(vars["id"], opts)
	if err == ots.ErrRunCancelNotAllowed {
		WriteError(w, http.StatusConflict, err)
		return
	} else if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) ForceCancelRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	opts := &tfe.RunForceCancelOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, opts); err != nil {
		WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	err := s.RunService.ForceCancelRun(vars["id"], opts)
	if err == ots.ErrRunForceCancelNotAllowed {
		WriteError(w, http.StatusConflict, err)
		return
	} else if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}