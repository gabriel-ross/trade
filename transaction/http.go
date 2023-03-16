package transaction

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	arango "github.com/arangodb/go-driver"
	"github.com/gabriel-ross/trade"
	"github.com/go-chi/chi"
)

// request represents a request body containing transaction data.
type request struct {
	Quantities  map[string]float64 `json:"quantities"`
	SenderID    string             `json:"senderID"`
	RecipientID string             `json:"recipientID"`
}

func (s *service) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()
		reqData := trade.Transaction{}

		err = bindRequest(r, &reqData)
		if err != nil {
			s.renderer.RenderError(w, r, err, http.StatusBadRequest, "%s", err.Error())
			return
		}

		resp, err := s.database.Create(ctx, reqData)
		if err != nil {
			s.renderer.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
			return
		}

		s.renderer.RenderJSON(w, r, http.StatusCreated, resp)
	}
}

func (s *service) handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()

		resp, err := s.database.List(ctx)
		if err != nil {
			s.renderer.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
			return
		}

		s.renderer.RenderJSON(w, r, http.StatusOK, resp)
	}
}

func (s *service) handleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()

		resp, err := s.database.Get(ctx, chi.URLParam(r, "id"))
		if err != nil {
			s.renderer.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
			return
		}

		s.renderer.RenderJSON(w, r, http.StatusOK, resp)
	}
}

func (s *service) handlePut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()
		data := trade.Transaction{}

		err = bindRequest(r, &data)
		if err != nil {
			s.renderer.RenderError(w, r, err, http.StatusBadRequest, "%s", err.Error())
			return
		}

		_, err = s.database.Update(ctx, chi.URLParam(r, "id"), data)
		if err != nil {
			if arango.IsNotFoundGeneral(err) {
				s.renderer.RenderError(w, r, err, http.StatusNotFound, "%s", err.Error())
				return
			} else {
				s.renderer.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *service) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()

		err = s.database.Delete(ctx, chi.URLParam(r, "id"))
		if err != nil {
			if arango.IsNotFoundGeneral(err) {
				s.renderer.RenderError(w, r, err, http.StatusNotFound, "%s", err.Error())
				return
			} else {
				s.renderer.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// bindRequest is a helper function for binding data from a request to a
// transaction object.
func bindRequest(r *http.Request, t *trade.Transaction) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var reqBody request
	err = json.Unmarshal(body, &reqBody)

	t.Quantities = reqBody.Quantities
	t.SenderID = reqBody.SenderID
	t.RecipientID = reqBody.RecipientID
	t.Timestamp = time.Now()

	return nil
}
