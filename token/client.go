package token

import (
	"errors"
	"net/url"

	. "github.com/stripe/stripe-go"
)

// Client is used to invoke /Tokens APIs.
type Client struct {
	B     Backend
	Token string
}

var c *Client

// Create POSTs a new card or bank account.
// For more details see https://stripe.com/docs/api#create_card_Token and https://stripe.com/docs/api#create_bank_account_token.
func Create(params *TokenParams) (*Token, error) {
	refresh()
	return c.Create(params)
}

func (c *Client) Create(params *TokenParams) (*Token, error) {
	body := &url.Values{}
	token := c.Token

	if len(params.Customer) > 0 {
		if len(params.AccessToken) == 0 {
			err := errors.New("Invalid Token params: an access token is required for customer")
			return nil, err
		}

		body.Add("customer", params.Customer)
		token = params.AccessToken
	}

	if params.Card != nil {
		params.Card.AppendDetails(body, true)
	} else if params.Bank != nil {
		params.Bank.AppendDetails(body)
	} else if len(params.Customer) == 0 {
		err := errors.New("Invalid Token params: either Card or Bank need to be set")
		return nil, err
	}

	params.AppendTo(body)

	tok := &Token{}
	err := c.B.Call("POST", "/tokens", token, body, tok)

	return tok, err
}

// Get returns the details of a Token.
// For more details see https://stripe.com/docs/api#retrieve_Token.
func Get(id string, params *TokenParams) (*Token, error) {
	refresh()
	return c.Get(id, params)
}

func (c *Client) Get(id string, params *TokenParams) (*Token, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}
		params.AppendTo(body)
	}

	token := &Token{}
	err := c.B.Call("GET", "/tokens/"+id, c.Token, body, token)

	return token, err
}

func refresh() {
	if c == nil {
		c = &Client{B: GetBackend()}
	}

	c.Token = Key
}