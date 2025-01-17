// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"net/url"
	"sync"

	"github.com/trustbloc/orb/pkg/activitypub/store/spi"
	"github.com/trustbloc/orb/pkg/activitypub/vocab"
)

type ActivityStore struct {
	AddActivityStub        func(*vocab.ActivityType) error
	addActivityMutex       sync.RWMutex
	addActivityArgsForCall []struct {
		arg1 *vocab.ActivityType
	}
	addActivityReturns struct {
		result1 error
	}
	addActivityReturnsOnCall map[int]struct {
		result1 error
	}
	AddReferenceStub        func(spi.ReferenceType, *url.URL, *url.URL, ...spi.RefMetadataOpt) error
	addReferenceMutex       sync.RWMutex
	addReferenceArgsForCall []struct {
		arg1 spi.ReferenceType
		arg2 *url.URL
		arg3 *url.URL
		arg4 []spi.RefMetadataOpt
	}
	addReferenceReturns struct {
		result1 error
	}
	addReferenceReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteReferenceStub        func(spi.ReferenceType, *url.URL, *url.URL) error
	deleteReferenceMutex       sync.RWMutex
	deleteReferenceArgsForCall []struct {
		arg1 spi.ReferenceType
		arg2 *url.URL
		arg3 *url.URL
	}
	deleteReferenceReturns struct {
		result1 error
	}
	deleteReferenceReturnsOnCall map[int]struct {
		result1 error
	}
	GetActivityStub        func(*url.URL) (*vocab.ActivityType, error)
	getActivityMutex       sync.RWMutex
	getActivityArgsForCall []struct {
		arg1 *url.URL
	}
	getActivityReturns struct {
		result1 *vocab.ActivityType
		result2 error
	}
	getActivityReturnsOnCall map[int]struct {
		result1 *vocab.ActivityType
		result2 error
	}
	GetActorStub        func(*url.URL) (*vocab.ActorType, error)
	getActorMutex       sync.RWMutex
	getActorArgsForCall []struct {
		arg1 *url.URL
	}
	getActorReturns struct {
		result1 *vocab.ActorType
		result2 error
	}
	getActorReturnsOnCall map[int]struct {
		result1 *vocab.ActorType
		result2 error
	}
	PutActorStub        func(*vocab.ActorType) error
	putActorMutex       sync.RWMutex
	putActorArgsForCall []struct {
		arg1 *vocab.ActorType
	}
	putActorReturns struct {
		result1 error
	}
	putActorReturnsOnCall map[int]struct {
		result1 error
	}
	QueryActivitiesStub        func(*spi.Criteria, ...spi.QueryOpt) (spi.ActivityIterator, error)
	queryActivitiesMutex       sync.RWMutex
	queryActivitiesArgsForCall []struct {
		arg1 *spi.Criteria
		arg2 []spi.QueryOpt
	}
	queryActivitiesReturns struct {
		result1 spi.ActivityIterator
		result2 error
	}
	queryActivitiesReturnsOnCall map[int]struct {
		result1 spi.ActivityIterator
		result2 error
	}
	QueryReferencesStub        func(spi.ReferenceType, *spi.Criteria, ...spi.QueryOpt) (spi.ReferenceIterator, error)
	queryReferencesMutex       sync.RWMutex
	queryReferencesArgsForCall []struct {
		arg1 spi.ReferenceType
		arg2 *spi.Criteria
		arg3 []spi.QueryOpt
	}
	queryReferencesReturns struct {
		result1 spi.ReferenceIterator
		result2 error
	}
	queryReferencesReturnsOnCall map[int]struct {
		result1 spi.ReferenceIterator
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ActivityStore) AddActivity(arg1 *vocab.ActivityType) error {
	fake.addActivityMutex.Lock()
	ret, specificReturn := fake.addActivityReturnsOnCall[len(fake.addActivityArgsForCall)]
	fake.addActivityArgsForCall = append(fake.addActivityArgsForCall, struct {
		arg1 *vocab.ActivityType
	}{arg1})
	stub := fake.AddActivityStub
	fakeReturns := fake.addActivityReturns
	fake.recordInvocation("AddActivity", []interface{}{arg1})
	fake.addActivityMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ActivityStore) AddActivityCallCount() int {
	fake.addActivityMutex.RLock()
	defer fake.addActivityMutex.RUnlock()
	return len(fake.addActivityArgsForCall)
}

func (fake *ActivityStore) AddActivityCalls(stub func(*vocab.ActivityType) error) {
	fake.addActivityMutex.Lock()
	defer fake.addActivityMutex.Unlock()
	fake.AddActivityStub = stub
}

func (fake *ActivityStore) AddActivityArgsForCall(i int) *vocab.ActivityType {
	fake.addActivityMutex.RLock()
	defer fake.addActivityMutex.RUnlock()
	argsForCall := fake.addActivityArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ActivityStore) AddActivityReturns(result1 error) {
	fake.addActivityMutex.Lock()
	defer fake.addActivityMutex.Unlock()
	fake.AddActivityStub = nil
	fake.addActivityReturns = struct {
		result1 error
	}{result1}
}

func (fake *ActivityStore) AddActivityReturnsOnCall(i int, result1 error) {
	fake.addActivityMutex.Lock()
	defer fake.addActivityMutex.Unlock()
	fake.AddActivityStub = nil
	if fake.addActivityReturnsOnCall == nil {
		fake.addActivityReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addActivityReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ActivityStore) AddReference(arg1 spi.ReferenceType, arg2 *url.URL, arg3 *url.URL, arg4 ...spi.RefMetadataOpt) error {
	fake.addReferenceMutex.Lock()
	ret, specificReturn := fake.addReferenceReturnsOnCall[len(fake.addReferenceArgsForCall)]
	fake.addReferenceArgsForCall = append(fake.addReferenceArgsForCall, struct {
		arg1 spi.ReferenceType
		arg2 *url.URL
		arg3 *url.URL
		arg4 []spi.RefMetadataOpt
	}{arg1, arg2, arg3, arg4})
	stub := fake.AddReferenceStub
	fakeReturns := fake.addReferenceReturns
	fake.recordInvocation("AddReference", []interface{}{arg1, arg2, arg3, arg4})
	fake.addReferenceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ActivityStore) AddReferenceCallCount() int {
	fake.addReferenceMutex.RLock()
	defer fake.addReferenceMutex.RUnlock()
	return len(fake.addReferenceArgsForCall)
}

func (fake *ActivityStore) AddReferenceCalls(stub func(spi.ReferenceType, *url.URL, *url.URL, ...spi.RefMetadataOpt) error) {
	fake.addReferenceMutex.Lock()
	defer fake.addReferenceMutex.Unlock()
	fake.AddReferenceStub = stub
}

func (fake *ActivityStore) AddReferenceArgsForCall(i int) (spi.ReferenceType, *url.URL, *url.URL, []spi.RefMetadataOpt) {
	fake.addReferenceMutex.RLock()
	defer fake.addReferenceMutex.RUnlock()
	argsForCall := fake.addReferenceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *ActivityStore) AddReferenceReturns(result1 error) {
	fake.addReferenceMutex.Lock()
	defer fake.addReferenceMutex.Unlock()
	fake.AddReferenceStub = nil
	fake.addReferenceReturns = struct {
		result1 error
	}{result1}
}

func (fake *ActivityStore) AddReferenceReturnsOnCall(i int, result1 error) {
	fake.addReferenceMutex.Lock()
	defer fake.addReferenceMutex.Unlock()
	fake.AddReferenceStub = nil
	if fake.addReferenceReturnsOnCall == nil {
		fake.addReferenceReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addReferenceReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ActivityStore) DeleteReference(arg1 spi.ReferenceType, arg2 *url.URL, arg3 *url.URL) error {
	fake.deleteReferenceMutex.Lock()
	ret, specificReturn := fake.deleteReferenceReturnsOnCall[len(fake.deleteReferenceArgsForCall)]
	fake.deleteReferenceArgsForCall = append(fake.deleteReferenceArgsForCall, struct {
		arg1 spi.ReferenceType
		arg2 *url.URL
		arg3 *url.URL
	}{arg1, arg2, arg3})
	stub := fake.DeleteReferenceStub
	fakeReturns := fake.deleteReferenceReturns
	fake.recordInvocation("DeleteReference", []interface{}{arg1, arg2, arg3})
	fake.deleteReferenceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ActivityStore) DeleteReferenceCallCount() int {
	fake.deleteReferenceMutex.RLock()
	defer fake.deleteReferenceMutex.RUnlock()
	return len(fake.deleteReferenceArgsForCall)
}

func (fake *ActivityStore) DeleteReferenceCalls(stub func(spi.ReferenceType, *url.URL, *url.URL) error) {
	fake.deleteReferenceMutex.Lock()
	defer fake.deleteReferenceMutex.Unlock()
	fake.DeleteReferenceStub = stub
}

func (fake *ActivityStore) DeleteReferenceArgsForCall(i int) (spi.ReferenceType, *url.URL, *url.URL) {
	fake.deleteReferenceMutex.RLock()
	defer fake.deleteReferenceMutex.RUnlock()
	argsForCall := fake.deleteReferenceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ActivityStore) DeleteReferenceReturns(result1 error) {
	fake.deleteReferenceMutex.Lock()
	defer fake.deleteReferenceMutex.Unlock()
	fake.DeleteReferenceStub = nil
	fake.deleteReferenceReturns = struct {
		result1 error
	}{result1}
}

func (fake *ActivityStore) DeleteReferenceReturnsOnCall(i int, result1 error) {
	fake.deleteReferenceMutex.Lock()
	defer fake.deleteReferenceMutex.Unlock()
	fake.DeleteReferenceStub = nil
	if fake.deleteReferenceReturnsOnCall == nil {
		fake.deleteReferenceReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReferenceReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ActivityStore) GetActivity(arg1 *url.URL) (*vocab.ActivityType, error) {
	fake.getActivityMutex.Lock()
	ret, specificReturn := fake.getActivityReturnsOnCall[len(fake.getActivityArgsForCall)]
	fake.getActivityArgsForCall = append(fake.getActivityArgsForCall, struct {
		arg1 *url.URL
	}{arg1})
	stub := fake.GetActivityStub
	fakeReturns := fake.getActivityReturns
	fake.recordInvocation("GetActivity", []interface{}{arg1})
	fake.getActivityMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ActivityStore) GetActivityCallCount() int {
	fake.getActivityMutex.RLock()
	defer fake.getActivityMutex.RUnlock()
	return len(fake.getActivityArgsForCall)
}

func (fake *ActivityStore) GetActivityCalls(stub func(*url.URL) (*vocab.ActivityType, error)) {
	fake.getActivityMutex.Lock()
	defer fake.getActivityMutex.Unlock()
	fake.GetActivityStub = stub
}

func (fake *ActivityStore) GetActivityArgsForCall(i int) *url.URL {
	fake.getActivityMutex.RLock()
	defer fake.getActivityMutex.RUnlock()
	argsForCall := fake.getActivityArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ActivityStore) GetActivityReturns(result1 *vocab.ActivityType, result2 error) {
	fake.getActivityMutex.Lock()
	defer fake.getActivityMutex.Unlock()
	fake.GetActivityStub = nil
	fake.getActivityReturns = struct {
		result1 *vocab.ActivityType
		result2 error
	}{result1, result2}
}

func (fake *ActivityStore) GetActivityReturnsOnCall(i int, result1 *vocab.ActivityType, result2 error) {
	fake.getActivityMutex.Lock()
	defer fake.getActivityMutex.Unlock()
	fake.GetActivityStub = nil
	if fake.getActivityReturnsOnCall == nil {
		fake.getActivityReturnsOnCall = make(map[int]struct {
			result1 *vocab.ActivityType
			result2 error
		})
	}
	fake.getActivityReturnsOnCall[i] = struct {
		result1 *vocab.ActivityType
		result2 error
	}{result1, result2}
}

func (fake *ActivityStore) GetActor(arg1 *url.URL) (*vocab.ActorType, error) {
	fake.getActorMutex.Lock()
	ret, specificReturn := fake.getActorReturnsOnCall[len(fake.getActorArgsForCall)]
	fake.getActorArgsForCall = append(fake.getActorArgsForCall, struct {
		arg1 *url.URL
	}{arg1})
	stub := fake.GetActorStub
	fakeReturns := fake.getActorReturns
	fake.recordInvocation("GetActor", []interface{}{arg1})
	fake.getActorMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ActivityStore) GetActorCallCount() int {
	fake.getActorMutex.RLock()
	defer fake.getActorMutex.RUnlock()
	return len(fake.getActorArgsForCall)
}

func (fake *ActivityStore) GetActorCalls(stub func(*url.URL) (*vocab.ActorType, error)) {
	fake.getActorMutex.Lock()
	defer fake.getActorMutex.Unlock()
	fake.GetActorStub = stub
}

func (fake *ActivityStore) GetActorArgsForCall(i int) *url.URL {
	fake.getActorMutex.RLock()
	defer fake.getActorMutex.RUnlock()
	argsForCall := fake.getActorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ActivityStore) GetActorReturns(result1 *vocab.ActorType, result2 error) {
	fake.getActorMutex.Lock()
	defer fake.getActorMutex.Unlock()
	fake.GetActorStub = nil
	fake.getActorReturns = struct {
		result1 *vocab.ActorType
		result2 error
	}{result1, result2}
}

func (fake *ActivityStore) GetActorReturnsOnCall(i int, result1 *vocab.ActorType, result2 error) {
	fake.getActorMutex.Lock()
	defer fake.getActorMutex.Unlock()
	fake.GetActorStub = nil
	if fake.getActorReturnsOnCall == nil {
		fake.getActorReturnsOnCall = make(map[int]struct {
			result1 *vocab.ActorType
			result2 error
		})
	}
	fake.getActorReturnsOnCall[i] = struct {
		result1 *vocab.ActorType
		result2 error
	}{result1, result2}
}

func (fake *ActivityStore) PutActor(arg1 *vocab.ActorType) error {
	fake.putActorMutex.Lock()
	ret, specificReturn := fake.putActorReturnsOnCall[len(fake.putActorArgsForCall)]
	fake.putActorArgsForCall = append(fake.putActorArgsForCall, struct {
		arg1 *vocab.ActorType
	}{arg1})
	stub := fake.PutActorStub
	fakeReturns := fake.putActorReturns
	fake.recordInvocation("PutActor", []interface{}{arg1})
	fake.putActorMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ActivityStore) PutActorCallCount() int {
	fake.putActorMutex.RLock()
	defer fake.putActorMutex.RUnlock()
	return len(fake.putActorArgsForCall)
}

func (fake *ActivityStore) PutActorCalls(stub func(*vocab.ActorType) error) {
	fake.putActorMutex.Lock()
	defer fake.putActorMutex.Unlock()
	fake.PutActorStub = stub
}

func (fake *ActivityStore) PutActorArgsForCall(i int) *vocab.ActorType {
	fake.putActorMutex.RLock()
	defer fake.putActorMutex.RUnlock()
	argsForCall := fake.putActorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ActivityStore) PutActorReturns(result1 error) {
	fake.putActorMutex.Lock()
	defer fake.putActorMutex.Unlock()
	fake.PutActorStub = nil
	fake.putActorReturns = struct {
		result1 error
	}{result1}
}

func (fake *ActivityStore) PutActorReturnsOnCall(i int, result1 error) {
	fake.putActorMutex.Lock()
	defer fake.putActorMutex.Unlock()
	fake.PutActorStub = nil
	if fake.putActorReturnsOnCall == nil {
		fake.putActorReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.putActorReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ActivityStore) QueryActivities(arg1 *spi.Criteria, arg2 ...spi.QueryOpt) (spi.ActivityIterator, error) {
	fake.queryActivitiesMutex.Lock()
	ret, specificReturn := fake.queryActivitiesReturnsOnCall[len(fake.queryActivitiesArgsForCall)]
	fake.queryActivitiesArgsForCall = append(fake.queryActivitiesArgsForCall, struct {
		arg1 *spi.Criteria
		arg2 []spi.QueryOpt
	}{arg1, arg2})
	stub := fake.QueryActivitiesStub
	fakeReturns := fake.queryActivitiesReturns
	fake.recordInvocation("QueryActivities", []interface{}{arg1, arg2})
	fake.queryActivitiesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ActivityStore) QueryActivitiesCallCount() int {
	fake.queryActivitiesMutex.RLock()
	defer fake.queryActivitiesMutex.RUnlock()
	return len(fake.queryActivitiesArgsForCall)
}

func (fake *ActivityStore) QueryActivitiesCalls(stub func(*spi.Criteria, ...spi.QueryOpt) (spi.ActivityIterator, error)) {
	fake.queryActivitiesMutex.Lock()
	defer fake.queryActivitiesMutex.Unlock()
	fake.QueryActivitiesStub = stub
}

func (fake *ActivityStore) QueryActivitiesArgsForCall(i int) (*spi.Criteria, []spi.QueryOpt) {
	fake.queryActivitiesMutex.RLock()
	defer fake.queryActivitiesMutex.RUnlock()
	argsForCall := fake.queryActivitiesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ActivityStore) QueryActivitiesReturns(result1 spi.ActivityIterator, result2 error) {
	fake.queryActivitiesMutex.Lock()
	defer fake.queryActivitiesMutex.Unlock()
	fake.QueryActivitiesStub = nil
	fake.queryActivitiesReturns = struct {
		result1 spi.ActivityIterator
		result2 error
	}{result1, result2}
}

func (fake *ActivityStore) QueryActivitiesReturnsOnCall(i int, result1 spi.ActivityIterator, result2 error) {
	fake.queryActivitiesMutex.Lock()
	defer fake.queryActivitiesMutex.Unlock()
	fake.QueryActivitiesStub = nil
	if fake.queryActivitiesReturnsOnCall == nil {
		fake.queryActivitiesReturnsOnCall = make(map[int]struct {
			result1 spi.ActivityIterator
			result2 error
		})
	}
	fake.queryActivitiesReturnsOnCall[i] = struct {
		result1 spi.ActivityIterator
		result2 error
	}{result1, result2}
}

func (fake *ActivityStore) QueryReferences(arg1 spi.ReferenceType, arg2 *spi.Criteria, arg3 ...spi.QueryOpt) (spi.ReferenceIterator, error) {
	fake.queryReferencesMutex.Lock()
	ret, specificReturn := fake.queryReferencesReturnsOnCall[len(fake.queryReferencesArgsForCall)]
	fake.queryReferencesArgsForCall = append(fake.queryReferencesArgsForCall, struct {
		arg1 spi.ReferenceType
		arg2 *spi.Criteria
		arg3 []spi.QueryOpt
	}{arg1, arg2, arg3})
	stub := fake.QueryReferencesStub
	fakeReturns := fake.queryReferencesReturns
	fake.recordInvocation("QueryReferences", []interface{}{arg1, arg2, arg3})
	fake.queryReferencesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ActivityStore) QueryReferencesCallCount() int {
	fake.queryReferencesMutex.RLock()
	defer fake.queryReferencesMutex.RUnlock()
	return len(fake.queryReferencesArgsForCall)
}

func (fake *ActivityStore) QueryReferencesCalls(stub func(spi.ReferenceType, *spi.Criteria, ...spi.QueryOpt) (spi.ReferenceIterator, error)) {
	fake.queryReferencesMutex.Lock()
	defer fake.queryReferencesMutex.Unlock()
	fake.QueryReferencesStub = stub
}

func (fake *ActivityStore) QueryReferencesArgsForCall(i int) (spi.ReferenceType, *spi.Criteria, []spi.QueryOpt) {
	fake.queryReferencesMutex.RLock()
	defer fake.queryReferencesMutex.RUnlock()
	argsForCall := fake.queryReferencesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ActivityStore) QueryReferencesReturns(result1 spi.ReferenceIterator, result2 error) {
	fake.queryReferencesMutex.Lock()
	defer fake.queryReferencesMutex.Unlock()
	fake.QueryReferencesStub = nil
	fake.queryReferencesReturns = struct {
		result1 spi.ReferenceIterator
		result2 error
	}{result1, result2}
}

func (fake *ActivityStore) QueryReferencesReturnsOnCall(i int, result1 spi.ReferenceIterator, result2 error) {
	fake.queryReferencesMutex.Lock()
	defer fake.queryReferencesMutex.Unlock()
	fake.QueryReferencesStub = nil
	if fake.queryReferencesReturnsOnCall == nil {
		fake.queryReferencesReturnsOnCall = make(map[int]struct {
			result1 spi.ReferenceIterator
			result2 error
		})
	}
	fake.queryReferencesReturnsOnCall[i] = struct {
		result1 spi.ReferenceIterator
		result2 error
	}{result1, result2}
}

func (fake *ActivityStore) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addActivityMutex.RLock()
	defer fake.addActivityMutex.RUnlock()
	fake.addReferenceMutex.RLock()
	defer fake.addReferenceMutex.RUnlock()
	fake.deleteReferenceMutex.RLock()
	defer fake.deleteReferenceMutex.RUnlock()
	fake.getActivityMutex.RLock()
	defer fake.getActivityMutex.RUnlock()
	fake.getActorMutex.RLock()
	defer fake.getActorMutex.RUnlock()
	fake.putActorMutex.RLock()
	defer fake.putActorMutex.RUnlock()
	fake.queryActivitiesMutex.RLock()
	defer fake.queryActivitiesMutex.RUnlock()
	fake.queryReferencesMutex.RLock()
	defer fake.queryReferencesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ActivityStore) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ spi.Store = new(ActivityStore)
