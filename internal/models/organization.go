package models

type Organization struct {
	ID           string                  `json:"id,omitempty"`
	Name         string                  `json:"name,omitempty"`
	Description  string                  `json:"description,omitempty"`
	Status       string                  `json:"status,omitempty"`
	Type         string                  `json:"type,omitempty"`
	Created      string                  `json:"created,omitempty"`
	LastModified string                  `json:"lastModified,omitempty"`
	ParentID     string                  `json:"parentId,omitempty"`
	Attributes   []OrganizationAttribute `json:"attributes,omitempty"`
}

type OrganizationAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OrganizationList struct {
	TotalResults  int            `json:"totalResults"`
	StartIndex    int            `json:"startIndex"`
	Count         int            `json:"count"`
	Organizations []Organization `json:"organizations"`
	Links         []Link         `json:"links"`
}

type OrganizationPatch struct {
	Operation string      `json:"operation"`
	Path      string      `json:"path"`
	Value     interface{} `json:"value"`
}
