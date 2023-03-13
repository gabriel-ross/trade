package user

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gabriel-ross/trade"
)

// request represents a request body containing user data.
type request struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

func (s *service) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.TODO()
		reqData := trade.User{}

		err = bindRequest(r, &reqData)
		if err != nil {
			trade.RenderError(w, r, http.StatusBadRequest, err, "%s", err.Error())
			return
		}

		resp, err := s.database.Create(ctx, reqData)
		if err != nil {
			trade.RenderError(w, r, http.StatusInternalServerError, err, "%s", err.Error())
			return
		}

		trade.RenderJSON(w, r, http.StatusCreated, resp)
	}
}

func (s *service) handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *service) handleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *service) handlePut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *service) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// bindRequest is a helper function for binding data from a request to a user
// object.
func bindRequest(r *http.Request, u *trade.User) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var reqBody request
	err = json.Unmarshal(body, &reqBody)

	u.Name = reqBody.Name
	u.Email = reqBody.Email
	u.PhoneNumber = reqBody.PhoneNumber
	return nil
}
