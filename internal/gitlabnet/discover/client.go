package discover

import (
	"fmt"
	"net/http"
	"net/url"

	"gitlab.com/gitlab-org/gitlab-shell/client"
	"gitlab.com/gitlab-org/gitlab-shell/internal/command/commandargs"
	"gitlab.com/gitlab-org/gitlab-shell/internal/config"
	"gitlab.com/gitlab-org/gitlab-shell/internal/gitlabnet"
)

type Client struct {
	config *config.Config
	client *client.GitlabNetClient
}

type Response struct {
	UserId   int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func NewClient(config *config.Config) (*Client, error) {
	client, err := gitlabnet.GetClient(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating http client: %v", err)
	}

	return &Client{config: config, client: client}, nil
}

func (c *Client) GetByCommandArgs(args *commandargs.Shell) (*Response, error) {
	params := url.Values{}
	if args.GitlabUsername != "" {
		params.Add("username", args.GitlabUsername)
	} else if args.GitlabKeyId != "" {
		params.Add("key_id", args.GitlabKeyId)
	} else {
		// There was no 'who' information, this  matches the ruby error
		// message.
		return nil, fmt.Errorf("who='' is invalid")
	}

	return c.getResponse(params)
}

func (c *Client) getResponse(params url.Values) (*Response, error) {
	path := "/discover?" + params.Encode()

	response, err := c.client.Get(path)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return parse(response)
}

func parse(hr *http.Response) (*Response, error) {
	response := &Response{}
	if err := gitlabnet.ParseJSON(hr, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *Response) IsAnonymous() bool {
	return r.UserId < 1
}
