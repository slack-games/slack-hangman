package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

const APIBaseURL = "https://slack.com/api"

// ActionDefault default style
const ActionDefault = "default"

// ActionPrimary primary style
const ActionPrimary = "primary"

// ActionDanger danger style
const ActionDanger = "danger"

// SlackTeam is a slack registered team
type SlackTeam struct {
	TeamID      string `json:"id"`
	Name        string `json:"name"`
	Domain      string `json:"domain"`
	EmailDomain string `json:"email_domain"`
}

// Action buttons for attachments
type Action struct {
	Name    string `json:"name"`
	Text    string `json:"text"`
	Style   string `json:"style,omitempty"`
	Type    string `json:"type,omitempty"`
	Value   string `json:"value,omitempty"`
	Confirm string `json:"confirm,omitempty"`
}

// Attachment is meant for extra text or image in slack response
type Attachment struct {
	Title      string   `json:"title,omitempty"`
	Text       string   `json:"text"`
	Fallback   string   `json:"fallback"`
	ImageURL   string   `json:"image_url,omitempty"`
	Color      string   `json:"color,omitempty"`
	CallbackID string   `json:"callback_id,omitempty"`
	Actions    []Action `json:"actions,omitempty"`
}

// ResponseMessage is slack response for the actions
type ResponseMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type SlackTeamResponse struct {
	Team SlackTeam `json:"team"`
}

func GetTeamInfo(client *http.Client, token *oauth2.Token) (*SlackTeamResponse, error) {
	response, err := client.Get(fmt.Sprintf("%s/team.info?token=%s", APIBaseURL, token.AccessToken))

	if err != nil {
		log.Printf("Could not get the user information based on the token %s\n", err)
		return nil, err
	}
	defer response.Body.Close()

	var teamResponse SlackTeamResponse
	err = json.NewDecoder(response.Body).Decode(&teamResponse)
	if err != nil {
		return nil, err
	}

	return &teamResponse, nil
}

// TextOnly creates new response message with text only
func TextOnly(text string) ResponseMessage {
	return ResponseMessage{
		Text:        text,
		Attachments: []Attachment{},
	}
}
