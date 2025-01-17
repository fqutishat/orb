/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package didorbgenerator

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/trustbloc/edge-core/pkg/log"

	"github.com/trustbloc/orb/pkg/activitypub/vocab"
	"github.com/trustbloc/orb/pkg/anchor/subject"
	"github.com/trustbloc/orb/pkg/document/util"
	"github.com/trustbloc/orb/pkg/hashlink"
	"github.com/trustbloc/orb/pkg/linkset"
)

var logger = log.New("anchorevent")

const (
	// ID specifies the ID of the generator.
	ID = "https://w3id.org/orb#v0"

	// Namespace specifies the namespace of the generator.
	Namespace = "did:orb"

	// Version specifies the version of the generator.
	Version = uint64(0)

	multihashPrefix  = "did:orb"
	unpublishedLabel = "uAAA"

	separator     = ":"
	hashlinkParts = 3
)

// Generator generates a content object for did:orb anchor events.
type Generator struct {
	*options
}

// Opt defines an option for the generator.
type Opt func(opts *options)

type options struct {
	id        *url.URL
	namespace string
	version   uint64
}

// WithNamespace sets the namespace of the generator.
func WithNamespace(ns string) Opt {
	return func(opts *options) {
		opts.namespace = ns
	}
}

// WithVersion sets the version of the generator.
func WithVersion(version uint64) Opt {
	return func(opts *options) {
		opts.version = version
	}
}

// WithID sets the ID of the generator.
func WithID(id *url.URL) Opt {
	return func(opts *options) {
		opts.id = id
	}
}

// New returns a new generator.
func New(opts ...Opt) *Generator {
	optns := &options{
		id:        vocab.MustParseURL(ID),
		namespace: Namespace,
		version:   Version,
	}

	for _, opt := range opts {
		opt(optns)
	}

	return &Generator{
		options: optns,
	}
}

// ID returns the ID of the generator.
func (g *Generator) ID() *url.URL {
	return g.id
}

// Namespace returns the Namespace for the DID method.
func (g *Generator) Namespace() string {
	return g.namespace
}

// Version returns the Version of this generator.
func (g *Generator) Version() uint64 {
	return g.version
}

// CreateContentObject creates a content object from the given payload.
func (g *Generator) CreateContentObject(payload *subject.Payload) (vocab.Document, error) {
	if payload.CoreIndex == "" {
		return nil, fmt.Errorf("payload is missing core index")
	}

	if len(payload.PreviousAnchors) == 0 {
		return nil, fmt.Errorf("payload is missing previous anchors")
	}

	anchorInfo, err := hashlink.New().ParseHashLink(payload.CoreIndex)
	if err != nil {
		return nil, fmt.Errorf("parse core index hashlink [%s]: %w", payload.CoreIndex, err)
	}

	// Don't include the metadata of the hashlink since it is mutable. The metadata will be included in the 'via'
	// field of the 'related' link in the anchor Linkset.
	anchorURI, err := url.Parse(hashlink.HLPrefix + anchorInfo.ResourceHash)
	if err != nil {
		return nil, fmt.Errorf("parse core index URI [%s]: %w", payload.CoreIndex, err)
	}

	authorURI, err := url.Parse(payload.AnchorOrigin)
	if err != nil {
		return nil, fmt.Errorf("parse anchor origin URI [%s]: %w", payload.AnchorOrigin, err)
	}

	var items []*linkset.Item

	for _, value := range payload.PreviousAnchors {
		item, e := newItem(value)
		if e != nil {
			return nil, e
		}

		items = append(items, item)
	}

	link := linkset.NewAnchorLink(anchorURI, authorURI, g.id, items)

	anchorLinksetDoc, err := vocab.MarshalToDoc(linkset.New(link))
	if err != nil {
		return nil, fmt.Errorf("marshal anchor linkset to document: %w", err)
	}

	return anchorLinksetDoc, nil
}

func newItem(value *subject.SuffixAnchor) (*linkset.Item, error) {
	logger.Debugf("Resource - Key [%s] Value [%s]", value.Suffix, value.Anchor)

	if value.Anchor == "" {
		hrefURI, e := url.Parse(fmt.Sprintf("%s:%s:%s", multihashPrefix, unpublishedLabel, value.Suffix))
		if e != nil {
			return nil, fmt.Errorf("parse item HRef URI: %w", e)
		}

		return linkset.NewItem(hrefURI, nil), nil
	}

	parts := strings.Split(value.Anchor, separator)

	if len(parts) != hashlinkParts {
		return nil, fmt.Errorf("invalid number of parts for previous anchor hashlink[%s] for suffix[%s]: expected 3, got %d", value, value.Suffix, len(parts)) //nolint:lll
	}

	pos := strings.LastIndex(value.Anchor, ":")
	if pos == -1 {
		return nil, fmt.Errorf("invalid previous anchor hashlink[%s] - must contain separator ':'", value)
	}

	prevAnchor := parts[0] + separator + parts[1]

	hrefURI, e := url.Parse(fmt.Sprintf("%s:%s:%s", multihashPrefix, parts[1], value.Suffix))
	if e != nil {
		return nil, fmt.Errorf("parse item HRef URI: %w", e)
	}

	prevURI, e := url.Parse(prevAnchor)
	if e != nil {
		return nil, fmt.Errorf("parse item previous URI: %w", e)
	}

	return linkset.NewItem(hrefURI, prevURI), nil
}

// CreatePayload creates a payload from the given document.
func (g *Generator) CreatePayload(doc vocab.Document, coreIndexURI *url.URL,
	anchors []*url.URL) (*subject.Payload, error) {
	anchorLinkset := &linkset.Linkset{}

	err := vocab.UnmarshalFromDoc(doc, anchorLinkset)
	if err != nil {
		return nil, fmt.Errorf("unmarshal anchor Linkset: %w", err)
	}

	anchorLink := anchorLinkset.Link()
	if anchorLink == nil {
		return nil, fmt.Errorf("empty anchor Linkset")
	}

	items := anchorLink.Items()

	operationCount := uint64(len(items))

	prevAnchors, err := getPreviousAnchors(items, anchors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse previous anchors: %w", err)
	}

	if coreIndexURI == nil || !strings.HasPrefix(coreIndexURI.String(), anchorLink.Anchor().String()) {
		return nil, fmt.Errorf("URI [%s] is not related to core index URI [%s]",
			coreIndexURI, anchorLink.Anchor())
	}

	return &subject.Payload{
		Namespace:       g.namespace,
		Version:         g.version,
		CoreIndex:       coreIndexURI.String(),
		OperationCount:  operationCount,
		PreviousAnchors: prevAnchors,
		AnchorOrigin:    anchorLink.Author().String(),
	}, nil
}

func getPreviousAnchors(resources []*linkset.Item, previous []*url.URL) ([]*subject.SuffixAnchor, error) {
	var previousAnchors []*subject.SuffixAnchor

	for _, res := range resources {
		suffix, err := util.GetSuffix(res.HRef().String())
		if err != nil {
			return nil, err
		}

		prevAnchor := &subject.SuffixAnchor{Suffix: suffix}

		if res.Previous() != nil {
			prevAnchor, err = getPreviousAnchorForResource(suffix, res.Previous().String(), previous)
			if err != nil {
				return nil, fmt.Errorf("get previous anchor for resource: %w", err)
			}
		}

		previousAnchors = append(previousAnchors, prevAnchor)
	}

	return previousAnchors, nil
}

func getPreviousAnchorForResource(suffix, res string, previous []*url.URL) (*subject.SuffixAnchor, error) {
	for _, prev := range previous {
		if !strings.HasPrefix(prev.String(), res) {
			continue
		}

		logger.Debugf("Found previous anchor [%s] for suffix [%s]", prev, suffix)

		return &subject.SuffixAnchor{Suffix: suffix, Anchor: prev.String()}, nil
	}

	return nil, fmt.Errorf("resource[%s] not found in previous anchor list", res)
}
