package service

import (
	"bytes"
	"encoding/json"
	"github.com/cloudnativego/gogo-engine"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	fakeMatchLocationResult = "/matches/5a003b78-409e-4452-b456-a6f0dcee05bd"
)

var (
	formatter = render.New(render.Options{
		Directory:                   "",
		FileSystem:                  nil,
		Asset:                       nil,
		AssetNames:                  nil,
		Layout:                      "",
		Extensions:                  nil,
		Funcs:                       nil,
		Delims:                      render.Delims{},
		Charset:                     "",
		DisableCharset:              false,
		IndentJSON:                  true,
		IndentXML:                   false,
		PrefixJSON:                  nil,
		PrefixXML:                   nil,
		BinaryContentType:           "",
		HTMLContentType:             "",
		JSONContentType:             "",
		JSONPContentType:            "",
		TextContentType:             "",
		XMLContentType:              "",
		IsDevelopment:               false,
		UnEscapeHTML:                false,
		StreamingJSON:               false,
		RequirePartials:             false,
		RequireBlocks:               false,
		DisableHTTPErrorRendering:   false,
		RenderPartialsWithoutPrefix: false,
	})
)

//func CreateMatchRespondsToBadDaata(t *testing.T)  {
//	client :=&http.Client{}
//	repo :=newInmemoryRepository()
//}
func TestCreateMatchRespondsToBadData(t *testing.T) {
	client := &http.Client{}
	repo := NewInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(createMatchHandler(formatter, repo)))
	defer server.Close()
	body1 := []byte("this is not valid json")
	body2 := []byte("{\"test\":\"this is valid json, but doesn't conform to server expectations.\"}")
	// Send invalid JSON
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body1))
	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to createMatchHandler: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Error("Sending invalid JSON should result in a bad request from server.")
	}

	req2, err2 := http.NewRequest("POST", server.URL, bytes.NewBuffer(body2))
	if err2 != nil {
		t.Errorf("Error in creating second POST request for invalid data on create match: %v", err2)
	}
	req2.Header.Add("Content-Type", "application/json")
	res2, _ := client.Do(req2)
	defer res2.Body.Close()
	if res2.StatusCode != http.StatusBadRequest {
		t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
	}
}

// client := &http.Client
// repo ：= Newinmoery
// repo := mnemogo
// var fakeMatcher
// bvatfocoekl
// newcolr
//server
// router
// handler

func TestCreateMatch(t *testing.T) {
	client := &http.Client{} //客户端 repo server
	repo := NewInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(createMatchHandler(formatter, repo)))
	defer server.Close()
	body := []byte("{\n  \"gridsize\": 19,\n  \"playerWhite\": \"bob\",\n  \"playerBlack\": \"alfred\"\n}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))

	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to createMatchHandler: %v", err)
	}

	defer res.Body.Close()

	payload, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response status 201, received %s", res.Status)
	}

	loc, headerOk := res.Header["Location"]

	if !headerOk {
		t.Error("Location header is not set")
	} else {
		if !strings.Contains(loc[0], "/matches/") {
			t.Errorf("Location header should contain '/matches/'")
		}
		if len(loc[0]) != len(fakeMatchLocationResult) {
			t.Errorf("Location value does not contain guid of new match")
		}

	}

	var matchResponse NewMatchResponse

	err = json.Unmarshal(payload, &matchResponse)

	if err != nil {
		t.Errorf("Could not unmarshal payload into newMatchResponse object")
	}

	if matchResponse.ID == "" || !strings.Contains(loc[0], matchResponse.ID) {
		t.Error("matchResponse.Id does not match Location header")
	}

	// After creating a match match repository shuld have 1 iutem in it
	matches, err := repo.getMatches()

	if err != nil {
		t.Errorf("Unexpected error in getMatches(): %s", err)
	}
	if len(matches) != 1 {
		t.Errorf("Expected a match repo of 1 match, got size %d", len(matches))
	}
	var match gogo.Match

	match = matches[0]

	if match.GridSize != matchResponse.GridSize {
		t.Errorf("Expected repo match and HTTP response gridsize to match. Got %d and %d", match.GridSize, matchResponse.GridSize)
	}
	if match.PlayerWhite != "bob" {
		t.Errorf("Repository match, white player should be bob, got %s", match.PlayerWhite)
	}
	if matchResponse.PlayerWhite != "bob" {
		t.Errorf("Expected white player to be bob, got %s", matchResponse.PlayerWhite)
	}

	if matchResponse.PlayerBlack != "alfred" {
		t.Errorf("Expected black player to be alfred, got %s", matchResponse.PlayerBlack)
	}
}
func TestGetMatchListReturnsEmptyArrayForNoMatches(t *testing.T) {
	client := &http.Client{}
	repo := NewInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(getMatchListHandler(formatter, repo)))
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Errored when sending request to the server", err)
		return
	}
	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("Failed to read response from server", err)
	}
	var matchList []NewMatchResponse
	err = json.Unmarshal(payload, &matchList)
	if err != nil {
		t.Errorf("Could not unmarshal payload into []newMatchResponse slice")
	}
	if len(matchList) != 0 {
		t.Errorf("Expected an empty list of match responses, got %d", len(matchList))
	}
}
func TestGetMatchListReturnWhatsInRepository(t *testing.T) {
	client := &http.Client{}
	repo := NewInMemoryRepository()
	repo.addMatch(gogo.NewMatch(19, "black", "white"))
	repo.addMatch(gogo.NewMatch(13, "bl", "wh"))
	repo.addMatch(gogo.NewMatch(19, "b", "w"))
	server := httptest.NewServer(http.HandlerFunc(getMatchListHandler(formatter, repo)))
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Errored when sending request to the server", err)
		return
	}
	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("Failed to read response from server", err)
	}
	var matchList []NewMatchResponse
	err = json.Unmarshal(payload, &matchList)
	if err != nil {
		t.Errorf("Could not unmarshal payload into []newMatchResponse slice")
	}

	repoMatches, err := repo.getMatches()
	if err != nil {
		t.Errorf("Unexpected error in getMatches(): %s", err)
	}
	if len(matchList) != len(repoMatches) {
		t.Errorf("Match response size should have equaled repo size, sizes were: %d and %d", len(matchList), len(repoMatches))
	}
	for idx := 0; idx < 3; idx++ {
		if matchList[idx].GridSize != repoMatches[idx].GridSize {
			t.Errorf("Gridsize mismatch at index %d. Got %d and %d", idx, matchList[idx].GridSize, repoMatches[idx].GridSize)
		}
		if matchList[idx].PlayerBlack != matchList[idx].PlayerBlack {
			t.Errorf("PlayerBlack mismatch at index %d. Got %s and %s", idx, matchList[idx].PlayerBlack, repoMatches[idx].PlayerBlack)
		}
		if matchList[idx].PlayerWhite != matchList[idx].PlayerWhite {
			t.Errorf("PlayerWhite mismatch at index %d. Got %s and %s", idx, matchList[idx].PlayerWhite, repoMatches[idx].PlayerWhite)
		}
	}
}
func TestGetMatchDetailsReturnsExistingMatch(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)
	repo := NewInMemoryRepository()
	server := MakeTestServer(repo)
	targetMatch := gogo.NewMatch(19, "black", "white")
	repo.addMatch(targetMatch)
	targetMatchID := targetMatch.ID
	recorder = httptest.NewRecorder()
	body := []byte("{\n  \"player\": 2,\n  \"position\": {\n    \"x\": 3,\n    \"y\": 10\n  }\n}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/matches/"+targetMatchID+"/moves", reader)
	server.ServeHTTP(recorder, request)
	//  fmt.Println(targetMatchID)
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new move to return 201, got %d", recorder.Code)
	}

	recorder2 := httptest.NewRecorder()
	request2, _ := http.NewRequest("GET", "/matches/"+targetMatchID, nil)
	server.ServeHTTP(recorder2, request2)
	if recorder2.Code != http.StatusOK {
		t.Errorf("Should've gotten a 200 querying match details, got %d", recorder.Code)
	}
	// server.sERVEHTTO recor euq
	payload := recorder2.Body.Bytes()
	var matchDetails matchDetailsResponse
	err := json.Unmarshal(payload, &matchDetails)
	if err != nil {
		t.Errorf("Could not unmarshal payload into match details response.")
	}
	if len(matchDetails.GameBoard[0]) != 19 {
		t.Errorf("Game board size isn't 19, got %d", len(matchDetails.GameBoard[0]))
	}
	//fmt.Println(matchDetails.GameBoard)
	if matchDetails.GameBoard[3][10] != gogo.PlayerWhite {
		t.Errorf("Game board did not reflect added move to 3,10. Board: %v", matchDetails.GameBoard)
	}
	recorder3 := httptest.NewRecorder()
	body2 := []byte("{\n  \"player\": 1,\n  \"position\": {\n    \"x\": 8,\n    \"y\": 8\n  }\n}")
	reader2 := bytes.NewReader(body2)
	//fmt.Println(targetMatchID)
	request3, _ := http.NewRequest("POST", "/matches/"+targetMatchID+"/moves", reader2)

	server.ServeHTTP(recorder3, request3)
	if recorder3.Code != http.StatusCreated {
		t.Errorf("Expected 201(Created) for 2nd move, got %d", recorder3.Code)
	}

	payload = recorder3.Body.Bytes()
	var matchDetails2 matchDetailsResponse
	err = json.Unmarshal(payload, &matchDetails2)
	if err != nil {
		t.Errorf("Could not unmarshal response for 2nd move add, %s", err.Error())
	}
	if matchDetails2.GameBoard[8][8] != gogo.PlayerBlack {
		t.Errorf("Added move should belong to black at 8,8 - belongs to %d", matchDetails2.GameBoard[8][8])
	}

}
func MakeTestServer(repo matchRepository) *negroni.Negroni {
	server := negroni.New() // 不需要need all the miduer
	mx := mux.NewRouter()
	initRoutes(mx, formatter, repo)
	server.UseHandler(mx)
	return server
}

//僧才 队队员 classcontroller

// request httop.resu
// rcore htp.reoce
// server.Serhttp
// recor
// request
// serverr.nerm
//
func MakeFakseServer(repo matchRepository) *negroni.Negroni {
	server := negroni.New()
	mx := mux.NewRouter()
	initRoutes(mx, formatter, repo)
	server.UseHandler(mx)
	return server
}
