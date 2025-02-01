package nps

// ParksService handles communication with
// the /parks related methods of the NPS API.
type ParksService struct {
	client *Client
}

type Parks struct {
	Total string  `json:"total"`
	Limit string  `json:"limit"`
	Start string  `json:"start"`
	Data  []*Park `json:"data"`
}
type Activities struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Topics struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type PhoneNumbers struct {
	PhoneNumber string `json:"phoneNumber"`
	Description string `json:"description"`
	Extension   string `json:"extension"`
	Type        string `json:"type"`
}
type EmailAddresses struct {
	Description  string `json:"description"`
	EmailAddress string `json:"emailAddress"`
}
type Contacts struct {
	PhoneNumbers   []*PhoneNumbers   `json:"phoneNumbers"`
	EmailAddresses []*EmailAddresses `json:"emailAddresses"`
}
type Exceptions struct {
	ExceptionHours *Hours `json:"exceptionHours"`
	StartDate      string `json:"startDate"`
	Name           string `json:"name"`
	EndDate        string `json:"endDate"`
}
type Hours struct {
	Wednesday string `json:"wednesday"`
	Monday    string `json:"monday"`
	Thursday  string `json:"thursday"`
	Sunday    string `json:"sunday"`
	Tuesday   string `json:"tuesday"`
	Friday    string `json:"friday"`
	Saturday  string `json:"saturday"`
}
type OperatingHours struct {
	Exceptions    []*Exceptions `json:"exceptions"`
	Description   string        `json:"description"`
	StandardHours *Hours        `json:"standardHours"`
	Name          string        `json:"name"`
}
type Address struct {
	PostalCode            string `json:"postalCode"`
	City                  string `json:"city"`
	StateCode             string `json:"stateCode"`
	CountryCode           string `json:"countryCode"`
	ProvinceTerritoryCode string `json:"provinceTerritoryCode"`
	Line1                 string `json:"line1"`
	Line2                 string `json:"line2"`
	Line3                 string `json:"line3"`
	Type                  string `json:"type"`
}
type Image struct {
	Credit  string `json:"credit"`
	Title   string `json:"title"`
	AltText string `json:"altText"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
}

type EntranceFee struct {
	Cost        string `json:"cost"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

type EntrancePass struct {
	Cost        string `json:"cost"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

type Multimedia struct {
	Title string `json:"title"`
	ID    string `json:"id"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

type Park struct {
	Activities     []*Activities     `json:"activities"`
	Addresses      []*Address        `json:"addresses"`
	Contacts       *Contacts         `json:"contacts"`
	Description    string            `json:"description"`
	Designation    string            `json:"designation"`
	DirectionsInfo string            `json:"directionsInfo"`
	DirectionsURL  string            `json:"directionsUrl"`
	EntranceFees   []*EntranceFee    `json:"entranceFees"`
	EntrancePasses []*EntrancePass   `json:"entrancePasses"`
	FullName       string            `json:"fullName"`
	ID             string            `json:"id"`
	Images         []*Image          `json:"images"`
	LatLong        string            `json:"latLong"`
	Latitude       string            `json:"latitude"`
	Longitude      string            `json:"longitude"`
	Multimedia     []Multimedia      `json:"multimedia"`
	Name           string            `json:"name"`
	OperatingHours []*OperatingHours `json:"operatingHours"`
	ParkCode       string            `json:"parkCode"`
	RelevanceScore int               `json:"relevanceScore"`
	States         string            `json:"states"`
	Topics         []*Topics         `json:"topics"`
	URL            string            `json:"url"`
	WeatherInfo    string            `json:"weatherInfo"`
}
