/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package proof

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/piprate/json-gold/ld"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/sidetree-core-go/pkg/canonicalizer"

	"github.com/trustbloc/orb/pkg/activitypub/vocab"
	"github.com/trustbloc/orb/pkg/anchor/util"
	"github.com/trustbloc/orb/pkg/anchor/vcpubsub"
	proofapi "github.com/trustbloc/orb/pkg/anchor/witness/proof"
	"github.com/trustbloc/orb/pkg/datauri"
	"github.com/trustbloc/orb/pkg/linkset"
	"github.com/trustbloc/orb/pkg/vct"
)

var logger = log.New("proof-handler")

type pubSub interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
}

type anchorLinkPublisher interface {
	Publish(anchorLinkset *linkset.Linkset) error
}

type metricsProvider interface {
	WitnessAnchorCredentialTime(duration time.Duration)
}

// New creates new proof handler.
func New(providers *Providers, pubSub pubSub, dataURIMediaType datauri.MediaType) *WitnessProofHandler {
	return &WitnessProofHandler{
		Providers:        providers,
		publisher:        vcpubsub.NewPublisher(pubSub),
		dataURIMediaType: dataURIMediaType,
	}
}

// Providers contains all of the providers required by the handler.
type Providers struct {
	AnchorLinkStore anchorEventStore
	StatusStore     statusStore
	WitnessStore    witnessStore
	WitnessPolicy   witnessPolicy
	MonitoringSvc   monitoringSvc
	DocLoader       ld.DocumentLoader
	Metrics         metricsProvider
}

// WitnessProofHandler handles an anchor credential witness proof.
type WitnessProofHandler struct {
	*Providers
	publisher        anchorLinkPublisher
	dataURIMediaType vocab.MediaType
}

type witnessStore interface {
	AddProof(anchorID string, witness *url.URL, p []byte) error
	Get(anchorID string) ([]*proofapi.WitnessProof, error)
}

type anchorEventStore interface {
	Get(id string) (*linkset.Link, error)
}

type statusStore interface {
	AddStatus(anchorEventID string, status proofapi.AnchorIndexStatus) error
	GetStatus(anchorEventID string) (proofapi.AnchorIndexStatus, error)
}

type monitoringSvc interface {
	Watch(vc *verifiable.Credential, endTime time.Time, domain string, created time.Time) error
}

type witnessPolicy interface {
	Evaluate(witnesses []*proofapi.WitnessProof) (bool, error)
}

// HandleProof handles proof.
func (h *WitnessProofHandler) HandleProof(witness *url.URL, anchor string, endTime time.Time, proof []byte) error { //nolint:lll
	logger.Debugf("received proof for anchor [%s] from witness[%s], proof: %s",
		anchor, witness.String(), string(proof))

	serverTime := time.Now().Unix()

	if endTime.Unix() < serverTime {
		// proof came after expiry time so nothing to do here
		// clean up process for witness store and Sidetree batch files will have to be initiated differently
		// since we can have scenario that proof never shows up
		return nil
	}

	status, err := h.StatusStore.GetStatus(anchor)
	if err != nil {
		return fmt.Errorf("failed to get status for anchor [%s]: %w", anchor, err)
	}

	if status == proofapi.AnchorIndexStatusCompleted {
		logger.Infof("Received proof from [%s] but witness policy has already been satisfied for anchor[%s]",
			witness, anchor, string(proof))

		// witness policy has been satisfied and witness proofs added to verifiable credential - nothing to do
		return nil
	}

	var witnessProof vct.Proof

	err = json.Unmarshal(proof, &witnessProof)
	if err != nil {
		return fmt.Errorf("failed to unmarshal incoming witness proof for anchor [%s]: %w", anchor, err)
	}

	anchorLink, err := h.AnchorLinkStore.Get(anchor)
	if err != nil {
		return fmt.Errorf("failed to retrieve anchor link [%s]: %w", anchor, err)
	}

	err = h.WitnessStore.AddProof(anchor, witness, proof)
	if err != nil {
		return fmt.Errorf("failed to add witness[%s] proof for anchor [%s]: %w",
			witness.String(), anchor, err)
	}

	vc, err := util.VerifiableCredentialFromAnchorLink(anchorLink,
		verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDDocumentLoader(h.DocLoader),
	)
	if err != nil {
		return fmt.Errorf("failed get verifiable credential from anchor: %w", err)
	}

	err = h.setupMonitoring(witnessProof, vc, endTime)
	if err != nil {
		return fmt.Errorf("failed to setup monitoring for anchor [%s]: %w", anchor, err)
	}

	return h.handleWitnessPolicy(anchorLink, vc)
}

func (h *WitnessProofHandler) setupMonitoring(wp vct.Proof, vc *verifiable.Credential, endTime time.Time) error {
	var created string
	if createdVal, ok := wp.Proof["created"].(string); ok {
		created = createdVal
	}

	createdTime, err := time.Parse(time.RFC3339, created)
	if err != nil {
		return fmt.Errorf("parse created: %w", err)
	}

	var domain string
	if domainVal, ok := wp.Proof["domain"].(string); ok {
		domain = domainVal
	}

	return h.MonitoringSvc.Watch(vc, endTime, domain, createdTime)
}

func (h *WitnessProofHandler) handleWitnessPolicy(anchorLink *linkset.Link, vc *verifiable.Credential) error { //nolint:funlen,gocyclo,cyclop,lll
	anchorID := anchorLink.Anchor().String()

	logger.Debugf("Handling witness policy for anchor link [%s]", anchorID)

	witnessProofs, err := h.WitnessStore.Get(anchorID)
	if err != nil {
		return fmt.Errorf("failed to get witness proofs for anchor [%s]: %w", anchorID, err)
	}

	ok, err := h.WitnessPolicy.Evaluate(witnessProofs)
	if err != nil {
		return fmt.Errorf("failed to evaluate witness policy for anchor [%s]: %w", anchorID, err)
	}

	if !ok {
		// Witness policy has not been satisfied - wait for other witness proofs to arrive ...
		logger.Infof("Witness policy has not been satisfied for anchor [%s]. Waiting for other proofs.", anchorID)

		return nil
	}

	// Witness policy has been satisfied so add witness proofs to anchor, set 'complete' status for anchor
	// publish witnessed anchor to batch writer channel for further processing
	logger.Infof("Witness policy has been satisfied for anchor [%s]", anchorID)

	vc, err = addProofs(vc, witnessProofs)
	if err != nil {
		return fmt.Errorf("failed to add witness proofs: %w", err)
	}

	status, err := h.StatusStore.GetStatus(anchorID)
	if err != nil {
		return fmt.Errorf("failed to get status for anchor [%s]: %w", anchorID, err)
	}

	logger.Debugf("Current status for VC [%s] is [%s]", anchorID, status)

	if status == proofapi.AnchorIndexStatusCompleted {
		logger.Infof("VC status has already been marked as completed for [%s]", anchorID)

		return nil
	}

	// Publish the VC before setting the status to completed since, if the publisher returns a transient error,
	// then this handler would be invoked on another server instance. So, we want the status to remain in-process,
	// otherwise the handler on the other instance would not publish the VC because it would think that is has
	// already been processed.
	logger.Debugf("Publishing anchor [%s]", anchorID)

	vcBytes, err := canonicalizer.MarshalCanonical(vc)
	if err != nil {
		return fmt.Errorf("create new object with document: %w", err)
	}

	vcDataURI, err := datauri.New(vcBytes, h.dataURIMediaType)
	if err != nil {
		return fmt.Errorf("create data URI from VC: %w", err)
	}

	// Create a new anchor with the updated verifiable credential.
	anchorLink = linkset.NewLink(
		anchorLink.Anchor(), anchorLink.Author(), anchorLink.Profile(),
		anchorLink.Original(), anchorLink.Related(),
		linkset.NewReference(vcDataURI, linkset.TypeJSONLD),
	)

	err = h.publisher.Publish(linkset.New(anchorLink))
	if err != nil {
		return fmt.Errorf("publish credential[%s]: %w", anchorID, err)
	}

	logger.Debugf("Setting status to [%s] for [%s]", proofapi.AnchorIndexStatusCompleted, anchorID)

	err = h.StatusStore.AddStatus(anchorID, proofapi.AnchorIndexStatusCompleted)
	if err != nil {
		return fmt.Errorf("failed to change status to 'completed' for anchor [%s]: %w", anchorID, err)
	}

	if vc.Issued != nil {
		h.Metrics.WitnessAnchorCredentialTime(time.Since(vc.Issued.Time))
	}

	return nil
}

func addProofs(vc *verifiable.Credential, proofs []*proofapi.WitnessProof) (*verifiable.Credential, error) {
	for _, p := range proofs {
		if p.Proof != nil {
			var witnessProof vct.Proof

			err := json.Unmarshal(p.Proof, &witnessProof)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal stored witness proof for anchor credential[%s]: %w", vc.ID, err)
			}

			if !proofExists(vc.Proofs, witnessProof.Proof) {
				logger.Debugf("Adding witness proof: %s", witnessProof.Proof)

				vc.Proofs = append(vc.Proofs, witnessProof.Proof)
			} else {
				logger.Debugf("Not adding witness proof since it already exists: %s", witnessProof.Proof)
			}
		}
	}

	return vc, nil
}

func proofExists(proofs []verifiable.Proof, proof verifiable.Proof) bool {
	for _, p := range proofs {
		if reflect.DeepEqual(p, proof) {
			return true
		}
	}

	return false
}
