package activity

//activity represents an activity state
type Activity struct {
	ID          string   `json:"id"`          //Activity ID*
	Date        uint     `json:"date"`        //Event Day
	Participant []string `json:"participant"` //Activity Participants list
	Place       string   `json:"place"`       //Event Place
	Summary     string   `json:"summary"`     //Short Summary for event
	Kind        string   `json:"kind"`        //Event Type
	Picture     []string `json:"picture"`     //Event Pictures, Path of Picture
}

//New makes an Activity
//Making new Activity needs Authentication with e-mail
func New(id string, date uint, place string, summary string, kind string) (*Activity, error) {
	//TODO
	//e-mail Verification (Auth)
	return &Activity{
		ID:          id,
		Date:        date,
		Participant: []string{},
		Place:       place,
		Summary:     summary,
		Kind:        kind,
		Picture:     []string{},
	}, nil
}
