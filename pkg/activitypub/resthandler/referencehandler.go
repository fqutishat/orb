/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package resthandler

import (
	"net/http"

	"github.com/trustbloc/orb/pkg/activitypub/store/spi"
	"github.com/trustbloc/orb/pkg/activitypub/store/storeutil"
	"github.com/trustbloc/orb/pkg/activitypub/vocab"
)

// Followers implements the 'followers' REST handler that retrieves a service's list of followers.
type Followers struct {
	*reference
}

// NewFollowers returns a new 'followers' REST handler.
func NewFollowers(cfg *Config, activityStore spi.Store) *Followers {
	return &Followers{
		reference: newReference(FollowersPath, spi.Follower, cfg, activityStore),
	}
}

// Following implements the 'following' REST handler that retrieves a service's list of following.
type Following struct {
	*reference
}

// NewFollowing returns a new 'following' REST handler.
func NewFollowing(cfg *Config, activityStore spi.Store) *Followers {
	return &Followers{
		reference: newReference(FollowingPath, spi.Following, cfg, activityStore),
	}
}

// Witnesses implements the 'witnesses' REST handler that retrieves a service's list of witnesses.
type Witnesses struct {
	*reference
}

// NewWitnesses returns a new 'witnesses' REST handler.
func NewWitnesses(cfg *Config, activityStore spi.Store) *Followers {
	return &Followers{
		reference: newReference(WitnessesPath, spi.Witness, cfg, activityStore),
	}
}

// Witnessing implements the 'witnessing' REST handler that retrieves collection of the services that the
// local service is witnessing.
type Witnessing struct {
	*reference
}

// NewWitnessing returns a new 'witnessing' REST handler.
func NewWitnessing(cfg *Config, activityStore spi.Store) *Followers {
	return &Followers{
		reference: newReference(WitnessingPath, spi.Witnessing, cfg, activityStore),
	}
}

type reference struct {
	*handler

	refType spi.ReferenceType
}

func newReference(path string, refType spi.ReferenceType, cfg *Config, activityStore spi.Store) *reference {
	h := &reference{
		refType: refType,
	}

	h.handler = newHandler(path, cfg, activityStore, h.handle, pageParam, pageNumParam)

	return h
}

func (h *reference) handle(w http.ResponseWriter, req *http.Request) {
	if h.isPaging(req) {
		h.handleReferencePage(w, req)
	} else {
		h.handleReference(w, req)
	}
}

//nolint:dupl
func (h *reference) handleReference(rw http.ResponseWriter, _ *http.Request) {
	coll, err := h.getReference()
	if err != nil {
		logger.Errorf("[%s] Error retrieving %s for service IRI [%s]: %s",
			h.endpoint, h.refType, h.ServiceIRI, err)

		h.writeResponse(rw, http.StatusInternalServerError, nil)

		return
	}

	collBytes, err := h.marshal(coll)
	if err != nil {
		logger.Errorf("[%s] Unable to marshal %s collection for service IRI [%s]: %s",
			h.endpoint, h.refType, h.ServiceIRI, err)

		h.writeResponse(rw, http.StatusInternalServerError, nil)

		return
	}

	h.writeResponse(rw, http.StatusOK, collBytes)
}

func (h *reference) handleReferencePage(rw http.ResponseWriter, req *http.Request) {
	var page *vocab.CollectionPageType

	var err error

	pageNum, ok := h.getPageNum(req)
	if ok {
		page, err = h.getPage(spi.WithPageSize(h.PageSize), spi.WithPageNum(pageNum))
	} else {
		page, err = h.getPage(spi.WithPageSize(h.PageSize))
	}

	if err != nil {
		logger.Errorf("[%s] Error retrieving page for service IRI [%s]: %s",
			h.endpoint, h.ServiceIRI, err)

		h.writeResponse(rw, http.StatusInternalServerError, nil)

		return
	}

	pageBytes, err := h.marshal(page)
	if err != nil {
		logger.Errorf("[%s] Unable to marshal page for service IRI [%s]: %s",
			h.endpoint, h.ServiceIRI, err)

		h.writeResponse(rw, http.StatusInternalServerError, nil)

		return
	}

	h.writeResponse(rw, http.StatusOK, pageBytes)
}

func (h *reference) getReference() (*vocab.CollectionType, error) {
	it, err := h.activityStore.QueryReferences(h.refType,
		spi.NewCriteria(
			spi.WithActorIRI(h.ServiceIRI),
		),
	)
	if err != nil {
		return nil, err
	}

	defer it.Close()

	firstURL, err := h.getPageURL(-1)
	if err != nil {
		return nil, err
	}

	lastURL, err := h.getPageURL(getLastPageNum(it.TotalItems(), h.PageSize, spi.SortAscending))
	if err != nil {
		return nil, err
	}

	return vocab.NewCollection(nil,
		vocab.WithContext(vocab.ContextActivityStreams),
		vocab.WithID(h.id),
		vocab.WithFirst(firstURL),
		vocab.WithLast(lastURL),
		vocab.WithTotalItems(it.TotalItems()),
	), nil
}

//nolint:dupl
func (h *reference) getPage(opts ...spi.QueryOpt) (*vocab.CollectionPageType, error) {
	it, err := h.activityStore.QueryReferences(
		h.refType,
		spi.NewCriteria(spi.WithActorIRI(h.ServiceIRI)),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	defer it.Close()

	options := storeutil.GetQueryOptions(opts...)

	refs, err := storeutil.ReadReferences(it, options.PageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*vocab.ObjectProperty, len(refs))

	for i, ref := range refs {
		items[i] = vocab.NewObjectProperty(vocab.WithIRI(ref))
	}

	id, prev, next, err := h.getIDPrevNextURL(it.TotalItems(), options)
	if err != nil {
		return nil, err
	}

	return vocab.NewCollectionPage(items,
		vocab.WithContext(vocab.ContextActivityStreams),
		vocab.WithID(id),
		vocab.WithPrev(prev),
		vocab.WithNext(next),
		vocab.WithTotalItems(it.TotalItems()),
	), nil
}