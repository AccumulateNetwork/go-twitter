package twitter

import "fmt"

// OAUTH 2.0 requirements
//dm.write
//dm.read
//tweet.read
//users.read

// https://api.twitter.com/2/dm_conversations/:dm_conversation_id/messages

type Attachment struct {
	MediaId string `json:"media_id"`
}
type Message struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

// GroupDM is a stuct used to create a group DM conversation
// https://api.twitter.com/2/dm_conversations/
type GroupDM struct {
	ConversationType string   `json:"conversation_type"`
	ParticipantIds   []string `json:"participant_ids"`
	Message          Message  `json:"message"`
}

type GroupDMResponse struct {
	ConversationId int `json:"conversation_id"`
	EventId        int `json:"event_id"`
}

// CreateGroupDMConversation will create a group conversation from a list of users IDs and send
// initial message to the group, if there is no message then it must include media attachments.
func CreateGroupDMConversation(participantIds []string, message string, mediaIds []string) *GroupDM {
	groupDM := &GroupDM{
		ConversationType: "Group",
	}
	groupDM.ParticipantIds = append(groupDM.ParticipantIds, participantIds...)
	groupDM.Message.Text = message
	if len(mediaIds) > 0 {
		for _, v := range mediaIds {
			groupDM.Message.Attachments = append(groupDM.Message.Attachments, Attachment{MediaId: v})
		}
	}

	return groupDM
}

func CreateDMMessage(message string) *Message {
	return &Message{Text: message}
}

func (group *GroupDM) Validate() error {
	if group.ConversationType != "Group" {
		return fmt.Errorf("conversation type must be set to 'Group'")
	}

	if len(group.ParticipantIds) == 0 {
		return fmt.Errorf("no participants provided")
	}

	if len(group.ParticipantIds) > 50 {
		return fmt.Errorf("only up to 50 participants are permitted")
	}

	return group.Message.Validate()
}

func (m *Message) Validate() error {

	if len(m.Text) == 0 {
		if len(m.Attachments) == 0 {
			return fmt.Errorf("no message test or attachments provided")
		}

		if len(m.Attachments) != 1 {
			return fmt.Errorf("only 1 media attachment supported at this time")
		}

		if len(m.Attachments[0].MediaId) == 0 {
			return fmt.Errorf("media id has not been set inside attachment")
		}
	} else if len(m.Text) > 10000 {
		return fmt.Errorf("maximum length of message text permitted is 10000 characters")
	}

	return nil
}
