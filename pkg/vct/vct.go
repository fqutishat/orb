/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vct

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/piprate/json-gold/ld"
	"github.com/trustbloc/vct/pkg/client/vct"
	"github.com/trustbloc/vct/pkg/controller/command"

	orberrors "github.com/trustbloc/orb/pkg/errors"
	"github.com/trustbloc/orb/pkg/vcsigner"
)

const (
	ctxSecurity = "https://w3id.org/security/v1"
)

type signer interface {
	Sign(vc *verifiable.Credential, opts ...vcsigner.Opt) (*verifiable.Credential, error)
	Context() []string
}

type metricsProvider interface {
	WitnessAddProofVctNil(value time.Duration)
	WitnessAddVC(value time.Duration)
	WitnessAddProof(value time.Duration)
	WitnessWebFinger(value time.Duration)
	WitnessVerifyVCTSignature(value time.Duration)
	AddProofParseCredential(value time.Duration)
	AddProofSign(value time.Duration)
}

// HTTPClient represents HTTP client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type endpointRetriever interface {
	GetLogEndpoint() (string, error)
}

// Client represents VCT client.
type Client struct {
	signer            signer
	http              HTTPClient
	documentLoader    ld.DocumentLoader
	endpointRetriever endpointRetriever
	authReadToken     string
	authWriteToken    string
	metrics           metricsProvider
}

// Option is a config client instance option.
type Option func(opts *Client)

// WithHTTPClient allows providing HTTP client.
func WithHTTPClient(client HTTPClient) Option {
	return func(o *Client) {
		o.http = client
	}
}

// WithDocumentLoader allows providing document loader.
func WithDocumentLoader(loader ld.DocumentLoader) Option {
	return func(o *Client) {
		o.documentLoader = loader
	}
}

// WithAuthReadToken add auth token.
func WithAuthReadToken(authToken string) Option {
	return func(o *Client) {
		o.authReadToken = authToken
	}
}

// WithAuthWriteToken add auth token.
func WithAuthWriteToken(authToken string) Option {
	return func(o *Client) {
		o.authWriteToken = authToken
	}
}

// New returns the client.
func New(endpointRetriever endpointRetriever, signer signer, metrics metricsProvider, opts ...Option) *Client {
	client := &Client{
		endpointRetriever: endpointRetriever,
		signer:            signer,
		metrics:           metrics,
		http: &http.Client{
			Timeout: time.Minute,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) addProof(endpoint string, anchorCred []byte, timestamp int64) (*verifiable.Credential, error) {
	parseCredentialStartTime := time.Now()

	vc, err := verifiable.ParseCredential(anchorCred,
		verifiable.WithDisabledProofCheck(),
		verifiable.WithNoCustomSchemaCheck(),
		verifiable.WithJSONLDDocumentLoader(c.documentLoader),
	)
	if err != nil {
		if strings.Contains(err.Error(), "http request unsuccessful") {
			// The server is probably down. Return a transient error so that it may be retried.
			return nil, orberrors.NewTransient(fmt.Errorf("http error during parse credential: %w", err))
		}

		return nil, fmt.Errorf("parse credential: %w", err)
	}

	c.metrics.AddProofParseCredential(time.Since(parseCredentialStartTime))

	opts := []vcsigner.Opt{
		vcsigner.WithCreated(time.Unix(0, timestamp)),
		vcsigner.WithSignatureRepresentation(verifiable.SignatureProofValue),
	}

	if endpoint != "" {
		opts = append(opts, vcsigner.WithDomain(endpoint))
	}

	signStartTime := time.Now()

	for _, signerCtx := range c.signer.Context() {
		exist := false

		for _, vcCtx := range vc.Context {
			if vcCtx == signerCtx {
				exist = true

				break
			}
		}

		if !exist {
			vc.Context = append(vc.Context, signerCtx)
		}
	}

	// adds linked data proof
	vc, err = c.signer.Sign(vc, opts...) // sets created time from the VCT.

	c.metrics.AddProofSign(time.Since(signStartTime))

	if err != nil {
		return nil, fmt.Errorf("add proof to credential: %w", err)
	}

	return vc, nil
}

// HealthCheck check health.
func (c *Client) HealthCheck() error {
	endpoint, err := c.endpointRetriever.GetLogEndpoint()
	if err != nil {
		return fmt.Errorf("failed to get log endpoint: %w", err)
	}

	vctClient := vct.New(endpoint, vct.WithHTTPClient(c.http),
		vct.WithAuthReadToken(c.authReadToken), vct.WithAuthWriteToken(c.authWriteToken))

	return vctClient.HealthCheck(context.Background())
}

// Witness credentials.
func (c *Client) Witness(anchorCred []byte) ([]byte, error) { // nolint: funlen,gocyclo,cyclop
	endpoint, err := c.endpointRetriever.GetLogEndpoint()
	if err != nil {
		return nil, fmt.Errorf("failed to get log endpoint for witness: %w", err)
	}

	if endpoint == "" {
		addProofStartTime := time.Now()

		vc, innnerErr := c.addProof(endpoint, anchorCred, time.Now().UnixNano())
		if innnerErr != nil {
			return nil, fmt.Errorf("add proof: %w", innnerErr)
		}

		ctx := []string{ctxSecurity}

		ctx = append(ctx, c.signer.Context()...)

		c.metrics.WitnessAddProofVctNil(time.Since(addProofStartTime))

		return json.Marshal(Proof{
			Context: ctx,
			Proof:   vc.Proofs[len(vc.Proofs)-1], // gets the latest proof
		})
	}

	addVCStartTime := time.Now()

	vctClient := vct.New(endpoint, vct.WithHTTPClient(c.http),
		vct.WithAuthReadToken(c.authReadToken), vct.WithAuthWriteToken(c.authWriteToken))

	resp, err := vctClient.AddVC(context.Background(), anchorCred)
	if err != nil {
		return nil, orberrors.NewTransientf("failed to add VC: %w", err)
	}

	c.metrics.WitnessAddVC(time.Since(addVCStartTime))

	addProofStartTime := time.Now()

	vc, err := c.addProof(endpoint, anchorCred, int64(resp.Timestamp)*int64(time.Millisecond))
	if err != nil {
		return nil, fmt.Errorf("add proof: %w", err)
	}

	c.metrics.WitnessAddProof(time.Since(addProofStartTime))

	webFingerStartTime := time.Now()

	webResp, err := vctClient.Webfinger(context.Background())
	if err != nil {
		return nil, fmt.Errorf("webfinger: %w", err)
	}

	c.metrics.WitnessWebFinger(time.Since(webFingerStartTime))

	pubKeyRaw, ok := webResp.Properties[command.PublicKeyType]
	if !ok {
		return nil, fmt.Errorf("no public key")
	}

	pubKeyStr, ok := pubKeyRaw.(string)
	if !ok {
		return nil, fmt.Errorf("public key is not a string")
	}

	pubKey, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, fmt.Errorf("decode public key: %w", err)
	}

	// gets the latest proof
	proof := vc.Proofs[len(vc.Proofs)-1]

	createdAt, ok := proof["created"].(string)
	if !ok {
		return nil, errors.New("created time is not a string")
	}

	timestampTime, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse time: %w", err)
	}

	timestamp := uint64(timestampTime.UnixNano()) / uint64(time.Millisecond)

	verifyVCTStartTime := time.Now()

	// verifies the signature by given timestamp from the proof and original credentials.
	err = vct.VerifyVCTimestampSignature(resp.Signature, pubKey, timestamp, vc)
	if err != nil {
		return nil, fmt.Errorf("verify VC timestamp signature: %w", err)
	}

	c.metrics.WitnessVerifyVCTSignature(time.Since(verifyVCTStartTime))

	ctx := []string{ctxSecurity}

	ctx = append(ctx, c.signer.Context()...)

	return json.Marshal(Proof{
		Context: ctx,
		Proof:   proof,
	})
}

// Proof represents response.
type Proof struct {
	Context interface{}      `json:"@context"`
	Proof   verifiable.Proof `json:"proof"`
}
