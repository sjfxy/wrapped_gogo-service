package service

import (
	"encoding/json"
	"fmt"
	"github.com/cloudnativego/gogo-engine"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
)

func createMatchHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newMatchRequest newMatchRequest
		err := json.Unmarshal(payload, &newMatchRequest)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse matcg request")
			return
		}
		if !newMatchRequest.isValid() {
			formatter.Text(w, http.StatusBadRequest, "Invalid new match request")
			return
		}
		// 上面是进行对Reuqt进程处理
		//然后使用daytaprocess
		//返回了对应的处理的match对象
		//然后repo.add
		//然后返回newrospeo
		// matche
		//进行复制即可
		//机会自定生成 对应的初始化的neamcye
		// 即可
		//然后添加即可
		//request->servcielohgcierproceewwairep->newMctehrObjc
		// REPO.add
		// reopspem
		// por pr
		// w.je
		//rof,json
		//
		newMatch := gogo.NewMatch(newMatchRequest.GridSize, newMatchRequest.PlayerBlack, newMatchRequest.PlayerWhite)
		repo.addMatch(newMatch)
		var mr NewMatchResponse
		mr.copyMatch(newMatch)
		w.Header().Add("Location", "/matches/"+newMatch.ID)
		formatter.JSON(w, http.StatusOK, &mr)

	}
}
func getMatchListHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		repoMatches, err := repo.getMatches()
		if err == nil {
			matches := make([]NewMatchResponse, len(repoMatches))
			for idx, match := range repoMatches {
				matches[idx].copyMatch(match)
			}
			formatter.JSON(w, http.StatusOK, matches)
		} else {
			formatter.JSON(w, http.StatusNotFound, err.Error())
		}
	}
}
func getMatchDetailsHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		matchID := vars["id"]
		match, err := repo.getMatch(matchID)
		if err != nil {
			formatter.JSON(w, http.StatusNotFound, err.Error())
		} else {
			var mdr matchDetailsResponse
			mdr.copyMatch(match)
			formatter.JSON(w, http.StatusOK, &mdr)
		}
	}
}
func addMoveHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		matchID := vars["id"]
		match, err := repo.getMatch(matchID)
		if err != nil {
			formatter.JSON(w, http.StatusNotFound, err.Error())
		} else {
			payload, _ := ioutil.ReadAll(request.Body)
			//fmt.Println(payload)
			var moveRequest newMoveRequest
			err := json.Unmarshal(payload, &moveRequest)
			//fmt.Println(moveRequest.Position)
			newBoard, err := match.GameBoard.PerformMove(
				gogo.Move{Player: moveRequest.Player,
					Position: gogo.Coordinate{X: moveRequest.Position.X, Y: moveRequest.Position.Y}})
			if err != nil {
				fmt.Println(err.Error())
				formatter.JSON(w, http.StatusBadRequest, err.Error())
			} else {
				match.GameBoard = newBoard
				err = repo.updateMatch(matchID, match)
				if err != nil {
					formatter.JSON(w, http.StatusInternalServerError, err.Error())
				} else {
					var mdr matchDetailsResponse
					mdr.copyMatch(match)
					formatter.JSON(w, http.StatusCreated, &mdr)
				}
			}
		}

	}
}
