// Package pocketic_test provides a client for the "hello" canister.
// Do NOT edit this file. It was automatically generated by https://github.com/aviate-labs/agent-go.
package pocketic_test

import (
    "github.com/aviate-labs/agent-go"
    
    "github.com/aviate-labs/agent-go/principal"
)

// Agent is a client for the "hello" canister.
type Agent struct {
    *agent.Agent
    CanisterId principal.Principal
}

// NewAgent creates a new agent for the "hello" canister.
func NewAgent(canisterId principal.Principal, config agent.Config) (*Agent, error) {
    a, err := agent.New(config)
    if err != nil {
        return nil, err
    }
    return &Agent{
        Agent:      a,
        CanisterId: canisterId,
    }, nil
}

// HelloQuery calls the "helloQuery" method on the "hello" canister.
func (a Agent) HelloQuery(arg0 string) (*string, error) {
    var r0 string
    if err := a.Agent.Query(
        a.CanisterId,
        "helloQuery",
        []any{arg0},
        []any{&r0},
    ); err != nil {
        return nil, err
    }
    return &r0, nil
}

// HelloQueryQuery creates an indirect representation of the "helloQuery" method on the "hello" canister.
func (a Agent) HelloQueryQuery(arg0 string) (*agent.CandidAPIRequest, error) {
    return a.Agent.CreateCandidAPIRequest(
        agent.RequestTypeQuery,
        a.CanisterId,
        "helloQuery",
        arg0,
    )
}

// HelloUpdate calls the "helloUpdate" method on the "hello" canister.
func (a Agent) HelloUpdate(arg0 string) (*string, error) {
    var r0 string
    if err := a.Agent.Call(
        a.CanisterId,
        "helloUpdate",
        []any{arg0},
        []any{&r0},
    ); err != nil {
        return nil, err
    }
    return &r0, nil
}

// HelloUpdateCall creates an indirect representation of the "helloUpdate" method on the "hello" canister.
func (a Agent) HelloUpdateCall(arg0 string) (*agent.CandidAPIRequest, error) {
    return a.Agent.CreateCandidAPIRequest(
        agent.RequestTypeCall,
        a.CanisterId,
        "helloUpdate",
        arg0,
    )
}