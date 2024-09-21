package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"wechat-back/internals/database"
	"wechat-back/internals/logger"
	"wechat-back/internals/models"
	"wechat-back/internals/server"
	"wechat-back/internals/tools"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
CreateNewGroupEP
Creates a new group and inserts it to the database
*/
func CreateNewGroupEP(w http.ResponseWriter, r *http.Request, db database.DBHUB) {

	alog := logger.StartLogger()
	w.Header().Set("Content-Type", "multipart/form-data")

	// get group data through formdata
	val := r.FormValue("data")
	var group models.Group
	err := json.Unmarshal([]byte(val), &group)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusNotAcceptable, tools.FormatErrResponse(server.BAD_FIELD, err))
		return
	}

	// send image to services be uploaded

	group = *models.FormatGroup(&group)

	_, err = db.InsertGroupDB(group)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	// add chat to the user current chats

	// return  group info
	group.ID = primitive.NilObjectID
	tools.WriteJSON(w, http.StatusCreated, tools.FormatSuccessResponse(group, server.COMPLETED, "ok"))

}

/*
UpdateGroupInfoEP
updates the group info of the specific group
*/
func UpdateGroupInfoEP(w http.ResponseWriter, r *http.Request, db database.DBHUB) {

	alog := logger.StartLogger()

	w.Header().Set("Content-Type", "application/json")

	var group models.Group

	err := tools.ReadJSON(w, r, &group)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusNotAcceptable, tools.FormatErrResponse(server.BAD_FIELD, err))
		return
	}

	DBgroup, err := db.GetGroupDB(group.GroupID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			alog.ErrorLog(err.Error())
			tools.WriteJSON(w, http.StatusNotFound, tools.FormatErrResponse(server.BAD_REQUEST, err))
			return
		}
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	update := make(map[string]any)
	update["name"] = group.Name
	update["description"] = group.Description

	err = db.UpdateGroupDB(update, DBgroup.ID)
	if err != nil {

		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	DBgroup.ID = primitive.NilObjectID
	DBgroup.Name = group.Name
	DBgroup.Description = group.Description

	tools.WriteJSON(w, http.StatusOK, tools.FormatSuccessResponse(DBgroup, server.OK, "ok"))

}

/*
UpdateGroupAdminsEP
Updates the current user admins of the group, if the user exist on the array it will get deleted, if the user not exist on the array it will get added
*/
func UpdateGroupAdminsEP(w http.ResponseWriter, r *http.Request, db database.DBHUB) {
	alog := logger.StartLogger()

	w.Header().Set("Content-Type", "application/json")

	admin := r.URL.Query().Get("ai") // name
	operationType := r.URL.Query().Get("ot")
	gID := r.URL.Query().Get("gi")
	targets := strings.Split(r.URL.Query().Get("ads"), ",")
	if len(targets) == 0 || r.URL.Query().Get("ads") == "" {
		alog.ErrorLog("no targets")
		tools.WriteJSON(w, http.StatusNotAcceptable, tools.FormatCustomErrResponse("no targets", server.BAD_FIELD))
		return
	}

	DBgroup, err := db.GetGroupDB(gID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			alog.ErrorLog(err.Error())
			tools.WriteJSON(w, http.StatusNotFound, tools.FormatErrResponse(server.BAD_REQUEST, err))
			return
		}
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	// convert targets to primitive
	tars := []primitive.ObjectID{}
	for _, t := range targets {

		tar, err := primitive.ObjectIDFromHex(t)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.BAD_CREDENTIALS, err))
			return
		}
		tars = append(tars, tar)
	}

	newAdmins := []primitive.ObjectID{}
	update := make(map[string]any)
	if operationType == OPERATION_ADD {
		newAdmins = append(newAdmins, tools.AddSliceValues(DBgroup.Admins, tars...)...)
		update["admins"] = newAdmins
	} else {
		newAdmins = append(newAdmins, tools.FilterSliceValues(DBgroup.Admins, tars)...)
		update["admins"] = newAdmins
	}

	err = db.UpdateGroupDB(update, DBgroup.ID)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	// notify all users that they have been added or remove as admins
	alog.InfoLogger(fmt.Sprintf("admin %s has {operation} as admin", admin))

	DBgroup.ID = primitive.NilObjectID
	DBgroup.Admins = newAdmins

	tools.WriteJSON(w, http.StatusOK, tools.FormatSuccessResponse(DBgroup, server.OK, "ok"))

}

/*
UpdateGroupParticipantsEP
Updates the current user users of the group, if the user exist on the array it will get deleted, if the user not exist on the array it will get added
*/
func UpdateGroupParticipantsEP(w http.ResponseWriter, r *http.Request, db database.DBHUB) {

	alog := logger.StartLogger()

	w.Header().Set("Content-Type", "application/json")

	admin := r.URL.Query().Get("ai") // name
	operationType := r.URL.Query().Get("ot")
	gID := r.URL.Query().Get("gi")
	targets := strings.Split(r.URL.Query().Get("usrs"), ",")
	if len(targets) == 0 || r.URL.Query().Get("usrs") == "" {
		alog.ErrorLog("no targets")
		tools.WriteJSON(w, http.StatusNotAcceptable, tools.FormatCustomErrResponse("no targets", server.BAD_FIELD))
		return
	}

	DBgroup, err := db.GetGroupDB(gID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			alog.ErrorLog(err.Error())
			tools.WriteJSON(w, http.StatusNotFound, tools.FormatErrResponse(server.BAD_REQUEST, err))
			return
		}
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	// convert targets to primitive
	tars := []primitive.ObjectID{}
	for _, t := range targets {

		tar, err := primitive.ObjectIDFromHex(t)
		if err != nil {
			alog.ErrorLog(err.Error())
			tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.BAD_CREDENTIALS, err))
			return
		}
		tars = append(tars, tar)
	}

	newParticipants := []primitive.ObjectID{}
	update := make(map[string]any)
	if operationType == OPERATION_ADD {
		newParticipants = append(newParticipants, tools.AddSliceValues(DBgroup.Admins, tars...)...)
		update["participants"] = newParticipants
	} else {
		newParticipants = append(newParticipants, tools.FilterSliceValues(DBgroup.Admins, tars)...)
		update["participants"] = newParticipants
	}

	err = db.UpdateGroupDB(update, DBgroup.ID)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	// notify all users that they have been added or remove as admins
	alog.InfoLogger(fmt.Sprintf("admin %s has added {operation} as participants", admin))

	DBgroup.ID = primitive.NilObjectID
	DBgroup.Participants = newParticipants

	tools.WriteJSON(w, http.StatusOK, tools.FormatSuccessResponse(DBgroup, server.OK, "ok"))

}

/* FUNCTION TO BE IMPLEMENTED
func UpdateGroupImage() {

}
*/

/*
DeleteGroupEP
deletes all the chats logs and deletes the group
*/
func DeleteGroupEP(w http.ResponseWriter, r *http.Request, db database.DBHUB) {

	alog := logger.StartLogger()

	w.Header().Set("Content-Type", "application/json")
	// get group
	groupID := r.URL.Query().Get("gi")
	if groupID == "" {
		alog.ErrorLog("no target")
		tools.WriteJSON(w, http.StatusNotAcceptable, tools.FormatCustomErrResponse("no targets", server.BAD_FIELD))
		return
	}

	DBgroup, err := db.GetGroupDB(groupID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			alog.ErrorLog(err.Error())
			tools.WriteJSON(w, http.StatusNotFound, tools.FormatErrResponse(server.BAD_REQUEST, err))
			return
		}
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	// delete chatlogs TO BE IMPLEMENTED
	// CHATLOG CODE HERE
	err = db.DeleteGroupDB(DBgroup.ID.Hex())
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}
	// return redirection

	tools.WriteJSON(w, http.StatusContinue, tools.FormatSuccessResponse(DBgroup.GroupID, server.COMPLETED, "done"))
}

/*
SearchGroupsEP
Search groups based on the name of the group
*/
func SearchGroupsEP(w http.ResponseWriter, r *http.Request, db database.DBHUB) {

	alog := logger.StartLogger()

	w.Header().Set("Content-Type", "application/json")

	p := r.URL.Query().Get("pg")
	if p == "" {
		p = "1"
	}

	pg, err := strconv.Atoi(p)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.BAD_FIELD, err))
		return
	}

	query := r.URL.Query().Get("q")

	groups, err := db.SearchGroups(pg, query)
	if err != nil {
		alog.ErrorLog(err.Error())
		tools.WriteJSON(w, http.StatusBadRequest, tools.FormatErrResponse(server.DB_ERROR, err))
		return
	}

	tools.WriteJSON(w, http.StatusOK, tools.FormatSuccessResponse(groups, server.OK, "ok"))

}
