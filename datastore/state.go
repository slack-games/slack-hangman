package datastore

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/slack-games/slack-hangman"
)

// TODO: Move words list into some DB, or use some compressed form
var wordList = []string{"receiver", "ritual", "insect", "interrupt", "salmon", "trading", "magic", "superior", "combat", "stem", "surgeon", "acceptable", "physics", "rape", "counsel", "jeans", "hunt", "continuous", "log", "echo", "pill", "excited", "sculpture", "compound", "integrate", "flour", "bitter", "bare", "slope", "rent", "presidency", "serving", "subtle", "greatly", "bishop", "drinking", "acceptance", "pump", "candy", "evil", "pleased", "medal", "beg", "sponsor", "ethical", "secondary", "slam", "export", "experimental", "melt", "midnight", "curve", "integrity", "entitle", "evident", "logic", "essence", "exclude", "harsh", "closet", "suburban", "greet", "interior", "corridor", "retail", "pitcher", "march", "snake", "excuse", "weakness"}

type State struct {
	StateID  string    `db:"state_id"`
	Word     string    `db:"word"`
	Guess    string    `db:"guess"`
	Current  string    `db:"current"`
	Mode     string    `db:"mode"`
	UserID   string    `db:"user_id"`
	ParentID string    `db:"parent_state_id"`
	Created  time.Time `db:"created_at"`
}

func (s *State) isGameOver() bool {
	return s.Mode == fmt.Sprintf("%s", hangman.GameOverState) ||
		s.Mode == fmt.Sprintf("%s", hangman.WinState)
}

func (s State) String() string {
	return fmt.Sprintf("#[%s] - %s %s %s %s %s",
		s.StateID, s.Word, s.Guess, s.Mode, s.UserID, s.Created)
}

func GetNewState(userID string) State {
	newWord := getNewWord()
	currentWord := randomizeWord(newWord)

	return State{
		Word:     newWord,
		Guess:    "",
		Current:  currentWord,
		Mode:     "Turn",
		UserID:   userID,
		ParentID: "00000000-0000-0000-0000-000000000000",
		Created:  time.Now(),
	}
}

func getNewWord() string {
	rand.Seed(int64(time.Now().Nanosecond()))
	index := rand.Intn(len(wordList))
	return wordList[index]
}

func randomizeWord(word string) string {
	rand.Seed(int64(time.Now().Nanosecond()))
	numVisible := rand.Intn(2) + 1
	newWord := ""

	// Create a new string with same size and hidden chars
	for range word {
		newWord += "_"
	}

	for i := 0; i < numVisible; i++ {
		index := rand.Intn(len(word))
		newWord = newWord[:index] + string(word[index]) + newWord[index+1:]
	}
	return newWord
}

func GetState(db *sqlx.DB, id string) (State, error) {
	state := State{}

	err := db.Get(&state, `SELECT * FROM hng.states WHERE state_id=$1 LIMIT 1`, id)
	return state, err
}

func GetUserLastState(db *sqlx.DB, id string) (State, error) {
	state := State{}

	query := `
		SELECT *
		FROM hng.states
		WHERE
			user_id=$1
		ORDER BY created_at DESC LIMIT 1;
	`

	err := db.Get(&state, query, id)
	return state, err
}

func NewState(db *sqlx.DB, state State) (string, error) {
	sql := `
		INSERT INTO hng.states
			(word, guess, current, mode, user_id, parent_state_id)
		VALUES
			(:word, :guess, :current, :mode, :user_id, :parent_state_id)
		RETURNING state_id
	`
	var id string

	rows, err := db.NamedQuery(sql, state)
	if err != nil {
		return id, err
	}

	if rows.Next() {
		rows.Scan(&id)
	}
	return id, err
}
