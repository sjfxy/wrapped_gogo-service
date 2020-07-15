package service

import "github.com/cloudnativego/gogo-engine"

type NewMatchResponse struct {
	ID          string `json:"id"`
	StartedAt   int64  `json:"started_at"`
	GridSize    int    `json:"gridsize"`
	PlayerWhite string `json:"playerWhite"`
	PlayerBlack string `json:"playerBlack"`
	Turn        int    `json:"turn,omitempty"`
}

//定义返回结果 详情返回结果
func (m *NewMatchResponse) copyMatch(match gogo.Match) {
	m.ID = match.ID
	m.StartedAt = match.StartTime.Unix()
	m.GridSize = match.GridSize
	m.PlayerWhite = match.PlayerWhite
	m.PlayerBlack = match.PlayerBlack
	m.Turn = match.TurnCount
}

type matchDetailsResponse struct {
	ID          string   `json:"id"`
	StartedAt   int64    `json:"started_at"`
	GridSize    int      `json:"gridsize"`
	PlayerWhite string   `json:"playerWhite"`
	PlayerBlack string   `json:"playerBlack"`
	Turn        int      `json:"turn,omitempty"`
	GameBoard   [][]byte `json:"gameboard"`
}

func (m *matchDetailsResponse) copyMatch(match gogo.Match) {
	m.ID = match.ID
	m.StartedAt = match.StartTime.Unix()
	m.GridSize = match.GridSize
	m.PlayerWhite = match.PlayerWhite
	m.PlayerBlack = match.PlayerBlack
	m.Turn = match.TurnCount
	m.GameBoard = match.GameBoard.Positions
}

//定义 Match 请求对象
type newMatchRequest struct {
	GridSize    int    `json:"gridsize"`
	PlayerWhite string `json:"playerWhite"`
	PlayerBlack string `json:"playerBlack"`
}
type boardPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

//定义请求对象
type newMoveRequest struct {
	Player   byte          `json:"player"`
	Position boardPosition `json:"position"`
}

//定义接口
type matchRepository interface {
	addMatch(match gogo.Match) (err error)
	getMatches() (matches []gogo.Match, err error)
	getMatch(id string) (match gogo.Match, err error)
	updateMatch(id string, match gogo.Match) (err error)
}

//请求参数过滤
// newMatchReqesutToekn-validate
// newMoevreq
// newDtyro
//
func (request newMatchRequest) isValid() (valid bool) {
	valid = true
	if request.GridSize != 19 && request.GridSize != 13 && request.GridSize != 9 {
		valid = false
	}
	if request.PlayerWhite == "" {
		valid = false
	}
	if request.PlayerBlack == "" {
		valid = false
	}
	return valid
}
