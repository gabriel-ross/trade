package account

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

// request represents a request body containing user data.
type request struct {
	Owner string `json:"owner"`
}

func (s *service) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()
		reqData := trade.Account{}

		err = bindRequest(r, &reqData)
		if err != nil {
			trade.RenderError(w, r, err, http.StatusBadRequest, "%s", err.Error())
			return
		}

		resp, err := s.database.Create(ctx, reqData)
		if err != nil {
			trade.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
			return
		}

		trade.RenderJSON(w, r, http.StatusCreated, resp)
	}
}

func (s *service) handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()

		resp, err := s.database.List(ctx)
		if err != nil {
			trade.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
			return
		}

		trade.RenderJSON(w, r, http.StatusOK, resp)
	}
}

func (s *service) handleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()

		resp, err := s.database.Get(ctx, chi.URLParam(r, "id"))
		if err != nil {
			trade.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
			return
		}

		trade.RenderJSON(w, r, http.StatusOK, resp)
	}
}

func (s *service) handlePut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()
		data := trade.Account{}

		err = bindRequest(r, &data)
		if err != nil {
			trade.RenderError(w, r, err, http.StatusBadRequest, "%s", err.Error())
			return
		}

		_, err = s.database.Update(ctx, chi.URLParam(r, "id"), data)
		if err != nil {
			if arango.IsNotFoundGeneral(err) {
				trade.RenderError(w, r, err, http.StatusNotFound, "%s", err.Error())
				return
			} else {
				trade.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
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
				trade.RenderError(w, r, err, http.StatusNotFound, "%s", err.Error())
				return
			} else {
				trade.RenderError(w, r, err, http.StatusInternalServerError, "%s", err.Error())
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// bindRequest is a helper function for binding data from a request to a user
// object.
func bindRequest(r *http.Request, u *trade.Account) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var reqBody request
	err = json.Unmarshal(body, &reqBody)

	u.Owner = reqBody.Owner
	u.Balances = map[string]float64{}
	u.Reputation = 100
	u.CreationTimestamp = time.Now()

	return nil
}