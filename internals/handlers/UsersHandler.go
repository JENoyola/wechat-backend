package handlers

import (
	"net/http"
	"strconv"
	"time"
	"wechat-back/internals/database"
	"wechat-back/internals/generators"
	"wechat-back/internals/logger"
	"wechat-back/internals/models"
	"wechat-back/internals/server"
	"wechat-back/internals/tools"
)

/*
NewUserAccountEP
creates a new user account and inserts it on the database
*/
func NewUserAccountEP(w http.ResponseWriter, r *http.Request, db *database.DB) {

	log := logger.StartLogger()

	w.Header().Set("Content-Type", "application/json")

	var user models.User

	err := tools.ReadJSON(w, r, &user)
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.BAD_FIELD, err))
		return
	}

	// get the user if exist
	_, exist, err := db.FindUserDB(user.Email)
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusConflict, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}
	if exist {
		log.WarningLogger("not allowed")
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatCustomErrResponse("not allowed", server.NOT_ALLOWED))
		return
	}

	// proceed to format user model
	user = *models.FormatUserModel(&user)

	// save it to database
	_, err = db.InsertUserDB(user)
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}
	// send code verification via email **TO BE IMPLEMENTED

	// return email entered and send redirection
	tools.WriteJSON(w, http.StatusOK, tools.FormatSuccessResponse(user.Email, server.REDIRECTION, "ok"))
}

/*
RequestNewCode
request a new code to the user to login
*/
func RequestNewCodeEP(w http.ResponseWriter, r *http.Request, db *database.DB) {

	log := logger.StartLogger()

	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("e")
	if len(email) < 1 {
		log.ErrorLog("no target")
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatCustomErrResponse("no target", server.BAD_REQUEST))
		return
	}

	user, exist, err := db.FindUserDB(email)
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}
	if !exist {
		log.WarningLogger("not allowed")
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatCustomErrResponse("not allowed", server.NO_DOCUMENTS))
		return
	}

	update := make(map[string]any)
	update["temp_code"], err = generators.Generate6DigitCode()
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusInternalServerError, tools.FormatErrResponse(server.SERVER_ERROR, err))
		return
	}
	update["valid_code"] = models.CODE_VALID
	update["code_timestamp"] = time.Now()

	err = db.UpdateUserAccountDB(update, user.ID.Hex())
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	// send new code to the user

	tools.WriteJSON(w, http.StatusOK, tools.FormatSuccessResponse(user.Email, server.OK, "ok"))

}

/*
UserCodeVerificationEP
Verifies the user account and signs the user
*/
func UserCodeVerificationEP(w http.ResponseWriter, r *http.Request, db *database.DB) {

	log := logger.StartLogger()

	w.Header().Set("Content-Type", "application/json")

	var user models.User

	err := tools.ReadJSON(w, r, &user)
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	DBuser, exist, err := db.FindUserDB(user.Email)
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	if !exist {
		log.WarningLogger("not allowed")
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatCustomErrResponse("not allowed", server.NOT_ALLOWED))
		return
	}

	// check weather the code is valid
	if DBuser.ValidCode == models.CODE_NOT_VALID {
		log.WarningLogger("not allowed")
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatCustomErrResponse("not allowed", server.NOT_ALLOWED))
		return
	}

	if user.TempCode != DBuser.TempCode {
		log.WarningLogger("bad credentials")
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatCustomErrResponse("bad credentials", server.BAD_CREDENTIALS))
		return
	}

	if DBuser.CodeTimestamp.Add(5 * time.Minute).Before(time.Now()) {
		log.WarningLogger("not allowed")
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatCustomErrResponse("not allowed", server.NOT_ALLOWED))
		return
	}

	// change variables on database
	changes := make(map[string]any)
	changes["valid_code"] = models.CODE_NOT_VALID
	changes["temp_code"] = 000000
	err = db.UpdateUserAccountDB(changes, DBuser.ID.Hex())
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	var res struct {
		Email string `json:"email" bson:"email"`
		UI    string `json:"ui" bson:"ui"`
		PI    string `json:"pi" bson:"pi"`
	}

	res.Email = DBuser.Email
	res.UI = DBuser.ID.Hex()
	res.PI = DBuser.Credentials.PushToken

	tools.WriteJSON(w, http.StatusOK, tools.FormatSuccessResponse(res, server.REDIRECTION, "ok"))

}

/*
LoginEP
logs the user and sends a verification code to the user email
*/
func LoginEP(w http.ResponseWriter, r *http.Request, db *database.DB) {

	log := logger.StartLogger()

	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("e")

	user, exist, err := db.FindUserDB(email)
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.BAD_FIELD, err))
		return
	}

	if !exist {
		log.WarningLogger("not allowed")
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatCustomErrResponse("not allowed", server.NOT_ALLOWED))
		return
	}

	newCode, err := generators.Generate6DigitCode()
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusInternalServerError, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	updated := make(map[string]any)
	updated["temp_code"] = newCode
	updated["valid_code"] = models.CODE_VALID
	updated["code_timestamp"] = time.Now()

	err = db.UpdateUserAccountDB(updated, user.ID.Hex())
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	// send new code to user via smtp

	tools.WriteJSON(w, http.StatusOK, tools.FormatSuccessResponse(user.Email, server.REDIRECTION, "proceed"))
}

/*
SearchUsersEP
gets a list of users that match the desired query
*/
func SearchUsersEP(w http.ResponseWriter, r *http.Request, db *database.DB) {

	log := logger.StartLogger()

	w.Header().Set("Content-Type", "application/json")

	p := r.URL.Query().Get("pg")
	if p == "" {
		p = "1"
	}

	pg, err := strconv.Atoi(p)
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.BAD_FIELD, err))
		return
	}

	query := r.URL.Query().Get("q")

	users, err := db.GetUsers(pg, query)
	if err != nil {
		log.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	tools.WriteJSON(w, http.StatusOK, tools.FormatSuccessResponse(users, server.OK, "ok"))

}
