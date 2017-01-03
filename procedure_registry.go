// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sunrpc

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

/*
From RFC 5531:
    The RPC call message has three unsigned-integer fields -- remote
    program number, remote program version number, and remote procedure
    number -- that uniquely identify the procedure to be called.
*/

// ProcedureID uniquely identifies a remote procedure
type ProcedureID struct {
	ProgramNumber   uint32
	ProgramVersion  uint32
	ProcedureNumber uint32
}

// pMap is looked up in ServerCodec to map ProcedureID to method name.
// rMap is looked up in ClientCodec to map method name to ProcedureID.
var procedureRegistry = struct {
	sync.RWMutex
	pMap map[ProcedureID]string
	rMap map[string]ProcedureID
}{
	pMap: make(map[ProcedureID]string),
	rMap: make(map[string]ProcedureID),
}

func isExported(name string) bool {
	firstRune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(firstRune)
}

func isValidProcedureName(procedureName string) bool {
	// procedureName must be of the format 'T.MethodName' to satisfy
	// criteria set by 'net/rpc' package for remote functions.

	procedureTypeName := strings.Split(procedureName, ".")
	if len(procedureTypeName) != 2 {
		return false
	}

	for _, name := range procedureTypeName {
		if !isExported(name) {
			return false
		}
	}

	return true
}

// RegisterProcedure will register the procedure name which will be uniquely
// indentified by (ProgramNumber, ProgramVersion, ProcedureNumber) pair.
func RegisterProcedure(procedureID ProcedureID, procedureName string) error {

	if !isValidProcedureName(procedureName) {
		return errors.New("Invalid procedure name")
	}

	procedureRegistry.Lock()
	defer procedureRegistry.Unlock()

	procedureRegistry.pMap[procedureID] = procedureName
	// Create reverse mapping too
	procedureRegistry.rMap[procedureName] = procedureID
	return nil
}

// GetProcedureName will return a string containing procedure name and a bool
// value which is set to true only if the procedure is found in registry.
func GetProcedureName(procedureID ProcedureID) (string, bool) {
	procedureRegistry.RLock()
	defer procedureRegistry.RUnlock()

	procedureName, ok := procedureRegistry.pMap[procedureID]
	return procedureName, ok
}

// GetProcedureID will return a struct containing (ProgramNumber, ProgramVersion, ProcedureNumber)
// pair, given the method name. It also returns a bool which is set to true only if the procedure
// is found in registry.
func GetProcedureID(procedureName string) (ProcedureID, bool) {
	procedureRegistry.RLock()
	defer procedureRegistry.RUnlock()

	procedureID, ok := procedureRegistry.rMap[procedureName]
	return procedureID, ok
}

// DumpProcedureRegistry will print the entire procedure map.
// Use this for logging/debugging.
func DumpProcedureRegistry() {
	procedureRegistry.RLock()
	defer procedureRegistry.RUnlock()

	for key, value := range procedureRegistry.rMap {
		fmt.Printf("%s : %+v\n", key, value)
	}
}