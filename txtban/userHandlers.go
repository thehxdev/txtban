package txtban

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/thehxdev/txtban/models"
	"github.com/thehxdev/txtban/tberr"
)

func (t *Txtban) useraddHandler(w http.ResponseWriter, r *http.Request) {
	var jdata JsonData

	body, err := readRequestBody(r.Body)
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := parseJsonBody(body, &jdata); err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errBadJsonData, http.StatusBadRequest)
		return
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}

	pass := jdata.Password
	minPassLen := viper.GetInt("limits.minPasswordLen")
	if len(pass) < minPassLen {
		err := tberr.New(
			fmt.Sprintf("password length is less than %d characters", minPassLen),
			fmt.Sprintf("choose a strong password with more than %d characters", minPassLen))
		sendError(w, err, http.StatusBadRequest)
		return
	}

	authKey := models.CreateAuthKey(uuid.String(), pass)
	err = t.DB.CreateUser(uuid.String(), pass, authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}

	sendJson(w, map[string]string{
		"uuid":    uuid.String(),
		"authKey": authKey,
	})
}

func (t *Txtban) userdelHandler(w http.ResponseWriter, r *http.Request) {
	authKey := r.Header.Get("Authorization")
	if authKey == "" {
		sendError(w, errEmptyAuthorizationHeader, http.StatusBadRequest)
		return
	}

	user, err := t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	err = t.DB.DeleteUser(user.ID)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errInternalServerError, http.StatusInternalServerError)
	}
}

func (t *Txtban) whoamiHandler(w http.ResponseWriter, r *http.Request) {
	var jdata JsonData

	body, err := readRequestBody(r.Body)
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := parseJsonBody(body, &jdata); err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errBadJsonData, http.StatusBadRequest)
		return
	}

	user, err := t.DB.AuthenticateByPassword(jdata.UserId, jdata.Password)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errUnauthorized, http.StatusInternalServerError)
		return
	}

	sendJson(w, map[string]string{
		"authKey": user.AuthKey,
	})
}

func (t *Txtban) passwdHandler(w http.ResponseWriter, r *http.Request) {
	var jdata JsonData

	body, err := readRequestBody(r.Body)
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := parseJsonBody(body, &jdata); err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errBadJsonData, http.StatusBadRequest)
		return
	}

	user, err := t.DB.AuthenticateByPassword(jdata.UserId, jdata.OldPassword)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	newPass := jdata.NewPassword
	minPassLen := viper.GetInt("limits.minPasswordLen")
	if len(newPass) < minPassLen {
		err := tberr.New(
			fmt.Sprintf("new password length is less than %d characters", minPassLen),
			fmt.Sprintf("choose a strong password with more than %d characters", minPassLen))
		sendError(w, err, http.StatusBadRequest)
		return
	}

	newAuthKey := models.CreateAuthKey(user.UUID, newPass)
	err = t.DB.UpdateUserPassword(user.ID, newPass, newAuthKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}

	sendJson(w, map[string]string{
		"authKey": newAuthKey,
	})
}
