package models

type APIResourceList struct {
	TotalResults int           `json:"totalResults"`
	Links        []interface{} `json:"links"`
	APIResources []APIResource `json:"apiResources"`
}

type APIResource struct {
	ID                    string     `json:"id"`
	Name                  string     `json:"name"`
	Identifier            string     `json:"identifier"`
	Type                  string     `json:"type"`
	RequiresAuthorization bool       `json:"requiresAuthorization"`
	Properties            []Property `json:"properties"`
	Self                  string     `json:"self"`
}

type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
