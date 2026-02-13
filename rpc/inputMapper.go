package rpc

import (
	"time"
)

// Activity holds the data for discord rich presence
type Activity struct {
	// The activity's name
	Name string
	// What the player is currently doing
	Details string
	// URL to open when clicking on the details text
	DetailsUrl string
	// The user's current party status
	State string
	// URL to open when clicking on the state text
	StateUrl string
	// The id for a large asset of the activity, usually a snowflake
	LargeImage string
	// Text displayed when hovering over the large image of the activity
	LargeText string
	// URL to open when clicking on the large image
	LargeUrl string
	// The id for a small asset of the activity, usually a snowflake
	SmallImage string
	// Text displayed when hovering over the small image of the activity
	SmallText string
	// URL to open when clicking on the small image
	SmallUrl string
	// Information for the current party of the player
	Party *Party
	// Unix timestamps for start and/or end of the game
	Timestamps *Timestamps
	// Secrets for Rich Presence joining and spectating
	Secrets *Secrets
	// Clickable buttons that open a URL in the browser
	Buttons []*Button
	// The type of the activity, defaults to 0 (Playing) if not set
	Type ActivityType
	// Which field is displayed in the status text in the member list
	StatusDisplayType StatusDisplayType
}

// Button holds a label and the corresponding URL that is opened on press
type Button struct {
	// The label of the button
	Label string
	// The URL of the button
	Url string
}

// Party holds information for the current party of the player
type Party struct {
	// The ID of the party
	ID string
	// Used to show the party's current size
	Players int
	// Used to show the party's maximum size
	MaxPlayers int
}

// Timestamps holds unix timestamps for start and/or end of the game
type Timestamps struct {
	// unix time (in milliseconds) of when the activity started
	Start *time.Time
	// unix time (in milliseconds) of when the activity ends
	End *time.Time
}

// Secrets holds secrets for Rich Presence joining and spectating
type Secrets struct {
	// The secret for a specific instanced match
	Match string
	// The secret for joining a party
	Join string
	// The secret for spectating a game
	Spectate string
}

type ActivityType int

const (
	ActivityTypePlaying   ActivityType = 0
	ActivityTypeListening ActivityType = 2
	ActivityTypeWatching  ActivityType = 3
	ActivityTypeCompeting ActivityType = 5
)

type StatusDisplayType int

const (
	// "Listening to Spotify"
	StatusDisplayTypeName StatusDisplayType = 0
	// "Listening to Rick Astley"
	StatusDisplayTypeState StatusDisplayType = 1
	// "Listening to Never Gonna Give You Up"
	StatusDisplayTypeDetails StatusDisplayType = 2
)

func mapActivity(activity *Activity) *PayloadActivity {
	final := &PayloadActivity{
		Name:       activity.Name,
		Details:    activity.Details,
		DetailsUrl: activity.DetailsUrl,
		State:      activity.State,
		StateUrl:   activity.StateUrl,
		Assets: PayloadAssets{
			LargeImage: activity.LargeImage,
			LargeText:  activity.LargeText,
			LargeUrl:   activity.LargeUrl,
			SmallImage: activity.SmallImage,
			SmallText:  activity.SmallText,
			SmallUrl:   activity.SmallUrl,
		},
		Type:              activity.Type,
		StatusDisplayType: activity.StatusDisplayType,
	}

	if activity.Timestamps != nil && activity.Timestamps.Start != nil {
		start := uint64(activity.Timestamps.Start.UnixNano() / 1e6)
		final.Timestamps = &PayloadTimestamps{
			Start: &start,
		}
		if activity.Timestamps.End != nil {
			end := uint64(activity.Timestamps.End.UnixNano() / 1e6)
			final.Timestamps.End = &end
		}
	}

	if activity.Party != nil {
		final.Party = &PayloadParty{
			ID:   activity.Party.ID,
			Size: [2]int{activity.Party.Players, activity.Party.MaxPlayers},
		}
	}

	if activity.Secrets != nil {
		final.Secrets = &PayloadSecrets{
			Join:     activity.Secrets.Join,
			Match:    activity.Secrets.Match,
			Spectate: activity.Secrets.Spectate,
		}
	}

	if len(activity.Buttons) > 0 {
		for _, btn := range activity.Buttons {
			final.Buttons = append(final.Buttons, &PayloadButton{
				Label: btn.Label,
				Url:   btn.Url,
			})
		}
	}

	return final
}
