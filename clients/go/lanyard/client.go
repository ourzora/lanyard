// An API client for [lanyard.org].
//
// [lanyard.org]: https://lanyard.org
package lanyard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/contextwtf/lanyard/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/xerrors"
)

var ErrNotFound error = xerrors.New("resource not found")

type Client struct {
	httpClient *http.Client
	url        string
}

type ClientOpt func(*Client)

func WithURL(url string) ClientOpt {
	return func(c *Client) {
		c.url = url
	}
}

func WithClient(hc *http.Client) ClientOpt {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// Uses https://lanyard.org/api/v1 for a default url
// and http.Client with a 30s timeout unless specified
// using [WithURL] or [WithClient]
func New(opts ...ClientOpt) *Client {
	const url = "https://lanyard.org/api/v1"
	c := &Client{
		url: url,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) sendRequest(
	ctx context.Context,
	method, path string,
	body, destination any,
) error {
	var (
		req *http.Request
		err error
	)

	url := c.url + path

	if body == nil {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)

		if err != nil {
			return xerrors.Errorf("error creating request: %w", err)
		}
	} else {
		b, err := json.Marshal(body)
		if err != nil {
			return xerrors.Errorf("failed to marshal body: %w", err)
		}
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewReader(b))
		if err != nil {
			return xerrors.Errorf("error creating request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("User-Agent", "lanyard-go+v1.0.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return xerrors.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode >= 400 {
		// special case 404s to make consuming client API easier
		if resp.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}

		return xerrors.Errorf("error making http request: %s", resp.Status)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&destination); err != nil {
		return xerrors.Errorf("failed to decode response: %w", err)
	}

	return nil
}

type CreateTreeRequest struct {
	// UnhashedLeaves is a slice of addresses or ABI encoded types
	UnhashedLeaves []hexutil.Bytes `json:"unhashedLeaves"`

	// LeafTypeDescriptor describes the abi-encoded types of the leaves, and
	// is required if leaves are not address types
	LeafTypeDescriptor []string `json:"leafTypeDescriptor,omitempty"`

	// PackedEncoding is true by default
	PackedEncoding types.JsonNullBool `json:"packedEncoding,omitempty"` // what's sent over the wire
}

type CreateResponse struct {
	// MerkleRoot is the root of the created merkle tree
	MerkleRoot hexutil.Bytes `json:"merkleRoot"`
}

// If you have a list of addresses for an allowlist, you can
// create a Merkle tree using CreateTree. Any Merkle tree
// published on Lanyard will be publicly available to any
// user of Lanyard’s API, including minting interfaces such
// as Zora or mint.fun.
func (c *Client) CreateTree(
	ctx context.Context,
	req CreateTreeRequest,
) (*CreateResponse, error) {
	resp := &CreateResponse{}

	err := c.sendRequest(ctx, http.MethodPost, "/tree", req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type TreeResponse struct {
	// UnhashedLeaves is a slice of addresses or ABI encoded types
	UnhashedLeaves []hexutil.Bytes `json:"unhashedLeaves"`

	// LeafTypeDescriptor describes the abi-encoded types of the leaves, and
	// is required if leaves are not address types
	LeafTypeDescriptor []string `json:"leafTypeDescriptor,omitempty"`

	// PackedEncoding is true by default
	PackedEncoding types.JsonNullBool `json:"packedEncoding,omitempty"`

	LeafCount int `json:"leafCount"`
}

// If a Merkle tree has been published to Lanyard, GetTreeFromRoot
// will return the entire tree based on the root.
// This endpoint will return ErrNotFound if the tree
// associated with the root has not been published.
func (c *Client) GetTreeFromRoot(
	ctx context.Context,
	root hexutil.Bytes,
) (*TreeResponse, error) {
	resp := &TreeResponse{}

	err := c.sendRequest(
		ctx, http.MethodGet,
		fmt.Sprintf("/tree?root=%s", root.String()),
		nil, resp,
	)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

type ProofResponse struct {
	UnhashedLeaf hexutil.Bytes   `json:"unhashedLeaf"`
	Proof        []hexutil.Bytes `json:"proof"`
}

// If the tree has been published to Lanyard, GetProof will
// return the proof associated with an unHashedLeaf.
// This endpoint will return ErrNotFound if the tree
// associated with the root has not been published.
func (c *Client) GetProofFromLeaf(
	ctx context.Context,
	root hexutil.Bytes,
	unhashedLeaf hexutil.Bytes,
) (*ProofResponse, error) {
	resp := &ProofResponse{}

	err := c.sendRequest(
		ctx, http.MethodGet,
		fmt.Sprintf("/proof?root=%s&unhashedLeaf=%s",
			root.String(), unhashedLeaf.String(),
		),
		nil, resp,
	)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

type RootResponse struct {
	Root hexutil.Bytes `json:"root"`
}

// If a Merkle tree has been published to Lanyard,
// GetRootFromLeaf will return the root of the tree
// based on a proof of a leaf. This endpoint will return
// ErrNotFound if the tree associated with the
// leaf has not been published.
func (c *Client) GetRootFromProof(
	ctx context.Context,
	proof []hexutil.Bytes,
) (*RootResponse, error) {
	resp := &RootResponse{}

	if len(proof) == 0 {
		return nil, xerrors.New("proof must not be empty")
	}

	var pq []string
	for _, p := range proof {
		pq = append(pq, p.String())
	}

	err := c.sendRequest(
		ctx, http.MethodGet,
		fmt.Sprintf("/root?proof=%s",
			strings.Join(pq, ","),
		),
		nil, resp,
	)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
