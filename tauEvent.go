package main

import "time"

type TauEvent struct {
	ID        string    `json:"id"`
	Created   time.Time `json:"created"`
	EventData struct {
		ID                   string `json:"id"`
		Bits                 int64  `json:"bits,string"`
		BroadcasterUserID    int64  `json:"broadcaster_user_id,string"`
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
		BroadcasterUserName  string `json:"broadcaster_user_name"`
		CategoryID           int64  `json:"category_id"`
		CategoryName         string `json:"category_name"`
		Data                 struct {
			Message struct {
				ChannelID        int64  `json:"channel_id,string"`
				ChannelName      string `json:"channel_name"`
				Context          string `json:"context"`
				CumulativeMonths int64  `json:"cumulative_months,string"`
				DisplayName      string `json:"display_name"`
				IsGift           bool   `json:"is_gift"`
				StreakMonths     int64  `json:"streak_months,string"`
				SubMessage       struct {
					Emotes  []interface{} `json:"emotes"`
					Message string        `json:"message"`
				} `json:"sub_message"`
				SubPlan     int64     `json:"sub_plan,string"`
				SubPlanName string    `json:"sub_plan_name"`
				Time        time.Time `json:"time"`
				UserID      int64     `json:"user_id,string"`
				UserName    string    `json:"user_name"`
			} `json:"message"`
			Topic string `json:"topic"`
		} `json:"data"`
		Message struct {
			ChannelID        int64  `json:"channel_id,string"`
			ChannelName      string `json:"channel_name"`
			Context          string `json:"context"`
			CumulativeMonths int64  `json:"cumulative_months,string"`
			DisplayName      string `json:"display_name"`
			IsGift           bool   `json:"is_gift"`
			StreakMonths     int64  `json:"streak_months,string"`
			SubMessage       struct {
				Emotes  []interface{} `json:"emotes"`
				Message string        `json:"message"`
			} `json:"sub_message"`
			SubPlan     int64     `json:"sub_plan,string"`
			SubPlanName string    `json:"sub_plan_name"`
			Time        time.Time `json:"time"`
			UserID      int64     `json:"user_id,string"`
			UserName    string    `json:"user_name"`
		} `json:"message"`
		Topic                    string    `json:"topic"`
		FromBroadcasterUserID    int64     `json:"from_broadcaster_user_id,string"`
		FromBroadcasterUserLogin string    `json:"from_broadcaster_user_login"`
		FromBroadcasterUserName  string    `json:"from_broadcaster_user_name"`
		IsAnonymous              bool      `json:"is_anonymous"`
		IsMature                 bool      `json:"is_mature"`
		Languate                 string    `json:"languate"`
		RedeemedAt               time.Time `json:"redeemed_at"`
		Reward                   struct {
			ID     string `json:"id"`
			Cost   int64  `json:"cost"`
			Prompt string `json:"prompt"`
			Title  string `json:"title"`
		} `json:"reward"`
		Status                 string `json:"status"`
		Title                  string `json:"title"`
		ToBroadcasterUserID    int64  `json:"to_broadcaster_user_id,string"`
		ToBroadcasterUserLogin string `json:"to_broadcaster_user_login"`
		ToBroadcasterUserName  string `json:"to_broadcaster_user_name"`
		Type                   string `json:"type"`
		UserID                 int64  `json:"user_id"`
		UserInput              string `json:"user_input"`
		UserLogin              string `json:"user_login"`
		UserName               string `json:"user_name"`
		Viewers                int64  `json:"viewers,string"`
	} `json:"event_data"`
	EventID     string `json:"event_id"`
	EventSource string `json:"event_source"`
	EventType   string `json:"event_type"`
	Origin      string `json:"origin"`
}
