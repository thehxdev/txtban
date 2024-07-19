package txtban

import (
	"github.com/go-chi/chi/v5"
	"github.com/thehxdev/txtban/tberr"
	"net/http"
)

func (t *Txtban) readHandler(w http.ResponseWriter, r *http.Request) {
	txtid := chi.URLParam(r, "txtid")
	content, err := t.DB.GetTxtContentById(txtid)
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, tberr.New("Not Found"), http.StatusNotFound)
		return
	}

	if acceptsGzip(r.Header.Values("Accept-Encoding")) {
		sendCompressed(w, []byte(content))
	} else {
		w.Write([]byte(content))
	}
}

func (t *Txtban) teeHandler(w http.ResponseWriter, r *http.Request) {
	authKey := r.Header.Get("Authorization")
	if authKey == "" {
		sendError(w, errEmptyAuthorizationHeader, http.StatusBadRequest)
		return
	}

	txtName := r.URL.Query().Get("name")
	if len(txtName) == 0 {
		sendError(w,
			tberr.New("txt name is empty", "choose a name for the txt you want to create"),
			http.StatusBadRequest)
		return
	}

	user, err := t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	body, err := readRequestBody(r.Body)
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		sendError(w, tberr.New("txt content is empty"), http.StatusBadRequest)
		return
	}

	txtid, err := t.DB.CreateTxt(user.ID, txtName, string(body))
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}

	sendJson(w, map[string]string{
		"id":   txtid,
		"name": txtName,
	})
}

func (t *Txtban) rmHandler(w http.ResponseWriter, r *http.Request) {
	authKey := r.Header.Get("Authorization")
	if authKey == "" {
		sendError(w, errEmptyAuthorizationHeader, http.StatusBadRequest)
		return
	}

	txtid := r.URL.Query().Get("txtid")
	if txtid == "" {
		sendError(w, errEmptyTxtID, http.StatusBadRequest)
		return
	}

	_, err := t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	err = t.DB.DeleteTxt(txtid)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}
}

func (t *Txtban) lsHandler(w http.ResponseWriter, r *http.Request) {
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

	txts, err := t.DB.GetAllTxts(user.ID)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}

	sendJson(w, txts)
}

func (t *Txtban) chtxtHandler(w http.ResponseWriter, r *http.Request) {
	authKey := r.Header.Get("Authorization")
	if authKey == "" {
		sendError(w, errEmptyAuthorizationHeader, http.StatusBadRequest)
		return
	}

	txtid := r.URL.Query().Get("txtid")
	if txtid == "" {
		sendError(w, errEmptyTxtID, http.StatusBadRequest)
		return
	}

	_, err := t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	body, err := readRequestBody(r.Body)
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = t.DB.ChangeTxtContent(txtid, string(body))
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}
}

func (t *Txtban) mvHandler(w http.ResponseWriter, r *http.Request) {
	authKey := r.Header.Get("Authorization")
	if authKey == "" {
		sendError(w, errEmptyAuthorizationHeader, http.StatusBadRequest)
		return
	}

	txtid := r.URL.Query().Get("txtid")
	if txtid == "" {
		sendError(w, errEmptyTxtID, http.StatusBadRequest)
		return
	}

	_, err := t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err)
		sendError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	newId, err := t.DB.ChangeTxtId(txtid)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}

	sendJson(w, map[string]string{
		"id": newId,
	})
}

func (t *Txtban) renameHandler(w http.ResponseWriter, r *http.Request) {
	authKey := r.Header.Get("Authorization")
	if authKey == "" {
		sendError(w, errEmptyAuthorizationHeader, http.StatusBadRequest)
		return
	}

	txtid := r.URL.Query().Get("txtid")
	if txtid == "" {
		sendError(w, errEmptyTxtID, http.StatusBadRequest)
		return
	}

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

	if len(jdata.Name) == 0 {
		sendError(w, errEmptyTxtName, http.StatusBadRequest)
		return
	}

	_, err = t.DB.AuthenticateByAuthKey(authKey)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	err = t.DB.ChangeTxtName(txtid, jdata.Name)
	if err != nil {
		t.ErrLogger.Println(err.Error())
		sendError(w, errInternalServerError, http.StatusInternalServerError)
		return
	}
}
