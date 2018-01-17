package zendesk

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"time"
)

const (
	apiTempl = "https://%s.zendesk.com/api/v2/%s"
)

var (
	httpClient *fasthttp.Client
)

type Client struct {
	conf *Config
}

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Domain   string `json:"Domain"`
}

type Organization struct {
	ID                 int64              `json:"id,omitempty"`
	Name               string             `json:"name,omitempty"`
	ExternalId         string             `json:"external_id,omitempty"`
	CreatedAt          string             `json:"created_at,omitempty"`
	UpdatedAt          string             `json:"updated_at,omitempty"`
	GroupID            string             `json:"group_id,omitempty"`
	OrganizationFields OrganizationFields `json:"organization_fields,omitempty"`
}

type OrganizationFields struct {
	EffectiveAt      time.Time `json:"effective_date,omitempty"`
	RenewalAt        time.Time `json:"renewal_date,omitempty"`
	SubscriptionType string    `json:"subscription_type,omitempty"`
	LicenseKey       string    `json:"license_key,omitempty"`
}

type Entitlement struct {
	OrganizationFields
	Product string `json:"product,omitempty"`
}
type User struct {
	ID             int64  `json:"id,omitempty"`
	Email          string `json:"email,omitempty"`
	Name           string `json:"name,omitempty"`
	OrganizationId int64  `json:"ognaization_id,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitEmpty"`
}

type Cluster struct {
	ID int64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	OrganizationID int64 `json:"organization_id,omitempty"`
	Nodes int`json:"nodes,omitempty"`
	Stage string `json:"stage,omitempty"`
	Version string `json:"version,omitempty"`
	Java string `json:"jre, omitempty"`
}

type Adoption struct {
	OrganizationID	int64 `json:"organizaton_id,omitempty"`
	Component string `json:"product,omitempty"`
	stage string `json:"stage,omitempty"`
}

type Error400 struct {
	Error       string                       `json:"error"`
	Description string                       `json:"description"`
	Details     map[string][]ErrorDetails400 `json:"details"`
}

//TODO: actually write an error handling procedure
type ErrorDetails400 struct {
	Description string `json:"description"`
	Error       string `json:"error"`
}

func LoadCredentialsFile(path string) (conf *Config) {
	file, _ := ioutil.ReadFile(path)
	if json.Unmarshal(file, &conf) != nil {
		log.Panic("Failed to read config")
	}
	return
}

func marshaler(payload interface{}) (data []byte) {
	var err error
	if data, err = json.Marshal(payload); err != nil {
		log.Printf("Error serializing resource: ", err)
		return nil
	}
	return data
}

func NewDefaultClient(conf *Config) *Client {
	httpClient = &fasthttp.Client{}
	return &Client{
		conf,
	}
}

func setBasicAuth(request *fasthttp.Request, user string, password string) {
	auth := []byte(user + ":" + password)
	request.Header.Set("Authorization",
		"Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
}

func newRequest(c Client, method string, resource string) (request *fasthttp.Request) {
	request = fasthttp.AcquireRequest()

	request.Header.Set("Content-Type", "application/json")
	setBasicAuth(request, c.conf.User, c.conf.Password)
	request.Header.SetMethod(method)

	request.SetRequestURI(fmt.Sprintf(apiTempl, c.conf.Domain, resource+".json"))

	return
}

func (c Client) Create(resource string, payload []byte, response *fasthttp.Response) {
	request := newRequest(c, "POST", resource+"/create_or_update")
	request.SetBody(payload)

	defer func() {
		fasthttp.ReleaseRequest(request)
	}()

	if err := httpClient.Do(request, response); err != nil {
		log.Printf("Error creating resource:  %s\n", err)
		return
	}

	if response.StatusCode() >= 400 && response.StatusCode() < 500 {
		var Err Error400
		json.Unmarshal(response.Body(), &Err)
		log.Printf("ERROR: Failed to create %s, %+v", resource, Err)
		response.SetBody([]byte(fmt.Sprintf("%s : %+v", Err.Description, Err.Details["name"])))
	}
}
