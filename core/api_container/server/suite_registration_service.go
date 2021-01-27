/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package server

import (
	"context"
	"errors"
	"github.com/kurtosis-tech/kurtosis/api_container/api/bindings"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

// Service that allows for a testsuite container to register itself exactly once
type suiteRegistrationService struct {
	mutex *sync.Mutex

	mainService ApiContainerServerService

	suiteRegistered bool

	// The action the newly-registered suite should take
	suiteAction bindings.SuiteAction

	suiteRegistrationChan chan interface{}
}

func newSuiteRegistrationService(suiteAction bindings.SuiteAction, mainService ApiContainerServerService, suiteRegistrationChan chan interface{}) *suiteRegistrationService {
	return &suiteRegistrationService{
		mutex:                 &sync.Mutex{},
		mainService:           mainService,
		suiteRegistered:       false,
		suiteAction:           suiteAction,
		suiteRegistrationChan: suiteRegistrationChan,
	}
}

func (service *suiteRegistrationService) RegisterSuite(_ context.Context, _ *emptypb.Empty) (*bindings.SuiteRegistrationResponse, error) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	if service.suiteRegistered {
		// Don't use stacktrace so we don't leak internal info
		return nil, errors.New("suite has already been registered")
	}

	service.suiteRegistered = true
	service.suiteRegistrationChan <- "Suite registered"
	response := &bindings.SuiteRegistrationResponse{SuiteAction: service.suiteAction}
	return response, nil
}