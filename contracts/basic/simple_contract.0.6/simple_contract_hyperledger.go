/*******************************************************************************
Copyright (c) 2016 IBM Corporation and other Contributors.


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.


Contributors:

Sumabala Nair - Initial Contribution
Kim Letkeman - Initial Contribution
Sumabala Nair - Updated for hyperledger May 2016
Sumabala Nair - Partial updates added May 2016
******************************************************************************/
//SN: March 2016

// IoT Blockchain Simple Smart Contract v 1.0

// This is a simple contract that creates a CRUD interface to
// create, read, update and delete an asset

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("SimpleContractChainCode")

var samples = `
{
    "event": {
        "assetID": "The ID of a managed asset. The resource focal point for a smart contract.",
        "carrier": "transport entity currently in possession of asset",
        "location": {
            "latitude": 123.456,
            "longitude": 123.456
        },
        "temperature": 123.456
    },
    "initEvent": {
        "nickname": "SIMPLE",
        "version": "The ID of a managed asset. The resource focal point for a smart contract."
    },
    "state": {
        "assetID": "The ID of a managed asset. The resource focal point for a smart contract.",
        "carrier": "transport entity currently in possession of asset",
        "location": {
            "latitude": 123.456,
            "longitude": 123.456
        },
        "temperature": 123.456
    }
}`

var schemas = `
{
    "API": {
        "createAsset": {
            "description": "Create an asset. One argument, a JSON encoded event. AssetID is required with zero or more writable properties. Establishes an initial asset state.",
            "properties": {
                "args": {
                    "description": "args are JSON encoded strings",
                    "items": {
                        "description": "A set of fields that constitute the writable fields in an asset's state. AssetID is mandatory along with at least one writable field. In this contract pattern, a partial state is used as an event.",
                        "properties": {
                            "assetID": {
                                "description": "The ID of a managed asset. The resource focal point for a smart contract.",
                                "type": "string"
                            },
                            "carrier": {
                                "description": "transport entity currently in possession of asset",
                                "type": "string"
                            },
                            "location": {
                                "description": "A geographical coordinate",
                                "properties": {
                                    "latitude": {
                                        "type": "number"
                                    },
                                    "longitude": {
                                        "type": "number"
                                    }
                                },
                                "type": "object"
                            },
                            "temperature": {
                                "description": "Temperature of the asset in CELSIUS.",
                                "type": "number"
                            }
                        },
                        "required": [
                            "assetID"
                        ],
                        "type": "object"
                    },
                    "maxItems": 1,
                    "minItems": 1,
                    "type": "array"
                },
                "function": {
                    "description": "createAsset function",
                    "enum": [
                        "createAsset"
                    ],
                    "type": "string"
                },
                "method": "invoke"
            },
            "type": "object"
        },
        "deleteAsset": {
            "description": "Delete an asset. Argument is a JSON encoded string containing only an assetID.",
            "properties": {
                "args": {
                    "description": "args are JSON encoded strings",
                    "items": {
                        "description": "An object containing only an assetID for use as an argument to read or delete.",
                        "properties": {
                            "assetID": {
                                "description": "The ID of a managed asset. The resource focal point for a smart contract.",
                                "type": "string"
                            }
                        },
                        "type": "object"
                    },
                    "maxItems": 1,
                    "minItems": 1,
                    "type": "array"
                },
                "function": {
                    "description": "deleteAsset function",
                    "enum": [
                        "deleteAsset"
                    ],
                    "type": "string"
                },
                "method": "invoke"
            },
            "type": "object"
        },
        "init": {
            "description": "Initializes the contract when started, either by deployment or by peer restart.",
            "properties": {
                "args": {
                    "description": "args are JSON encoded strings",
                    "items": {
                        "description": "event sent to init on deployment",
                        "properties": {
                            "nickname": {
                                "default": "SIMPLE",
                                "description": "The nickname of the current contract",
                                "type": "string"
                            },
                            "version": {
                                "description": "The ID of a managed asset. The resource focal point for a smart contract.",
                                "type": "string"
                            }
                        },
                        "required": [
                            "version"
                        ],
                        "type": "object"
                    },
                    "maxItems": 1,
                    "minItems": 1,
                    "type": "array"
                },
                "function": {
                    "description": "init function",
                    "enum": [
                        "init"
                    ],
                    "type": "string"
                },
                "method": "deploy"
            },
            "type": "object"
        },
        "readAsset": {
            "description": "Returns the state an asset. Argument is a JSON encoded string. AssetID is the only accepted property.",
            "properties": {
                "args": {
                    "description": "args are JSON encoded strings",
                    "items": {
                        "description": "An object containing only an assetID for use as an argument to read or delete.",
                        "properties": {
                            "assetID": {
                                "description": "The ID of a managed asset. The resource focal point for a smart contract.",
                                "type": "string"
                            }
                        },
                        "type": "object"
                    },
                    "maxItems": 1,
                    "minItems": 1,
                    "type": "array"
                },
                "function": {
                    "description": "readAsset function",
                    "enum": [
                        "readAsset"
                    ],
                    "type": "string"
                },
                "method": "query",
                "result": {
                    "description": "A set of fields that constitute the complete asset state.",
                    "properties": {
                        "assetID": {
                            "description": "The ID of a managed asset. The resource focal point for a smart contract.",
                            "type": "string"
                        },
                        "carrier": {
                            "description": "transport entity currently in possession of asset",
                            "type": "string"
                        },
                        "location": {
                            "description": "A geographical coordinate",
                            "properties": {
                                "latitude": {
                                    "type": "number"
                                },
                                "longitude": {
                                    "type": "number"
                                }
                            },
                            "type": "object"
                        },
                        "temperature": {
                            "description": "Temperature of the asset in CELSIUS.",
                            "type": "number"
                        }
                    },
                    "type": "object"
                }
            },
            "type": "object"
        },
        "readAssetSamples": {
            "description": "Returns a string generated from the schema containing sample Objects as specified in generate.json in the scripts folder.",
            "properties": {
                "args": {
                    "description": "accepts no arguments",
                    "items": {},
                    "maxItems": 0,
                    "minItems": 0,
                    "type": "array"
                },
                "function": {
                    "description": "readAssetSamples function",
                    "enum": [
                        "readAssetSamples"
                    ],
                    "type": "string"
                },
                "method": "query",
                "result": {
                    "description": "JSON encoded object containing selected sample data",
                    "type": "string"
                }
            },
            "type": "object"
        },
        "readAssetSchemas": {
            "description": "Returns a string generated from the schema containing APIs and Objects as specified in generate.json in the scripts folder.",
            "properties": {
                "args": {
                    "description": "accepts no arguments",
                    "items": {},
                    "maxItems": 0,
                    "minItems": 0,
                    "type": "array"
                },
                "function": {
                    "description": "readAssetSchemas function",
                    "enum": [
                        "readAssetSchemas"
                    ],
                    "type": "string"
                },
                "method": "query",
                "result": {
                    "description": "JSON encoded object containing selected schemas",
                    "type": "string"
                }
            },
            "type": "object"
        },
        "updateAsset": {
            "description": "Update the state of an asset. The one argument is a JSON encoded event. AssetID is required along with one or more writable properties. Establishes the next asset state. ",
            "properties": {
                "args": {
                    "description": "args are JSON encoded strings",
                    "items": {
                        "description": "A set of fields that constitute the writable fields in an asset's state. AssetID is mandatory along with at least one writable field. In this contract pattern, a partial state is used as an event.",
                        "properties": {
                            "assetID": {
                                "description": "The ID of a managed asset. The resource focal point for a smart contract.",
                                "type": "string"
                            },
                            "carrier": {
                                "description": "transport entity currently in possession of asset",
                                "type": "string"
                            },
                            "location": {
                                "description": "A geographical coordinate",
                                "properties": {
                                    "latitude": {
                                        "type": "number"
                                    },
                                    "longitude": {
                                        "type": "number"
                                    }
                                },
                                "type": "object"
                            },
                            "temperature": {
                                "description": "Temperature of the asset in CELSIUS.",
                                "type": "number"
                            }
                        },
                        "required": [
                            "assetID"
                        ],
                        "type": "object"
                    },
                    "maxItems": 1,
                    "minItems": 1,
                    "type": "array"
                },
                "function": {
                    "description": "updateAsset function",
                    "enum": [
                        "updateAsset"
                    ],
                    "type": "string"
                },
                "method": "invoke"
            },
            "type": "object"
        }
    },
    "objectModelSchemas": {
        "assetIDKey": {
            "description": "An object containing only an assetID for use as an argument to read or delete.",
            "properties": {
                "assetID": {
                    "description": "The ID of a managed asset. The resource focal point for a smart contract.",
                    "type": "string"
                }
            },
            "type": "object"
        },
        "event": {
            "description": "A set of fields that constitute the writable fields in an asset's state. AssetID is mandatory along with at least one writable field. In this contract pattern, a partial state is used as an event.",
            "properties": {
                "assetID": {
                    "description": "The ID of a managed asset. The resource focal point for a smart contract.",
                    "type": "string"
                },
                "carrier": {
                    "description": "transport entity currently in possession of asset",
                    "type": "string"
                },
                "location": {
                    "description": "A geographical coordinate",
                    "properties": {
                        "latitude": {
                            "type": "number"
                        },
                        "longitude": {
                            "type": "number"
                        }
                    },
                    "type": "object"
                },
                "temperature": {
                    "description": "Temperature of the asset in CELSIUS.",
                    "type": "number"
                }
            },
            "required": [
                "assetID"
            ],
            "type": "object"
        },
        "initEvent": {
            "description": "event sent to init on deployment",
            "properties": {
                "nickname": {
                    "default": "SIMPLE",
                    "description": "The nickname of the current contract",
                    "type": "string"
                },
                "version": {
                    "description": "The ID of a managed asset. The resource focal point for a smart contract.",
                    "type": "string"
                }
            },
            "required": [
                "version"
            ],
            "type": "object"
        },
        "state": {
            "description": "A set of fields that constitute the complete asset state.",
            "properties": {
                "assetID": {
                    "description": "The ID of a managed asset. The resource focal point for a smart contract.",
                    "type": "string"
                },
                "carrier": {
                    "description": "transport entity currently in possession of asset",
                    "type": "string"
                },
                "location": {
                    "description": "A geographical coordinate",
                    "properties": {
                        "latitude": {
                            "type": "number"
                        },
                        "longitude": {
                            "type": "number"
                        }
                    },
                    "type": "object"
                },
                "temperature": {
                    "description": "Temperature of the asset in CELSIUS.",
                    "type": "number"
                }
            },
            "type": "object"
        }
    }
}`

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// CONTRACTSTATEKEY is used to store contract state into world state
const CONTRACTSTATEKEY string = "ContractStateKey"

// MYVERSION must use this to deploy contract
const MYVERSION string = "1.0"

// MRUKEY is used to store the last 10 assetIDs and timestamps
const MRUKEY string = "_MruListKey"

// ************************************
// asset and contract state
// ************************************

// ContractState holds the contract version
type ContractState struct {
	Version string `json:"version"`
}

// Geolocation stores lat and long
type Geolocation struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// AssetState stores current state for any assset
type AssetState struct {
	AssetID     *string      `json:"assetID,omitempty"`     // all assets must have an ID, primary key of contract
	Location    *Geolocation `json:"location,omitempty"`    // current asset location
	Temperature *float64     `json:"temperature,omitempty"` // asset temp
	Carrier     *string      `json:"carrier,omitempty"`     // the name of the carrier
	UpdatedAt   *time.Time   `json:"updatedAt,omitempty"`
}

type AssetUpdatedAt struct {
	AssetID   *string   `json:"assetID,omitempty"`   // all assets must have an ID, primary key of contract
	UpdatedAt time.Time `json:"updatedAt,omitempty"` // when was AssetID last updated
}

type AssetMruList struct {
	List []AssetUpdatedAt `json:"mruList,omitempty"` // the list
}

var contractState = ContractState{MYVERSION}

// ************************************
// deploy callback mode
// ************************************

// Init is called during deploy
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("fmt: Init called")
	logger.Error("logger: Init called")
	var stateArg ContractState
	var err error
	if len(args) != 1 {
		return nil, errors.New("init expects one argument, a JSON string with tagged version string")
	}
	err = json.Unmarshal([]byte(args[0]), &stateArg)
	if err != nil {
		return nil, errors.New("Version argument unmarshal failed: " + fmt.Sprint(err))
	}
	if stateArg.Version != MYVERSION {
		return nil, errors.New("Contract version " + MYVERSION + " must match version argument: " + stateArg.Version)
	}
	contractStateJSON, err := json.Marshal(stateArg)
	if err != nil {
		return nil, errors.New("Marshal failed for contract state" + fmt.Sprint(err))
	}
	err = stub.PutState(CONTRACTSTATEKEY, contractStateJSON)
	if err != nil {
		return nil, errors.New("Contract state failed PUT to ledger: " + fmt.Sprint(err))
	}
	return nil, nil
}

// ************************************
// deploy and invoke callback mode
// ************************************

// Invoke is called when an invoke message is received
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Handle different functions
	if function == "createAsset" {
		// create assetID
		return t.createAsset(stub, args)
	} else if function == "updateAsset" {
		// create assetID
		return t.updateAsset(stub, args)
	} else if function == "deleteAsset" {
		// Deletes an asset by ID from the ledger
		return t.deleteAsset(stub, args)
	}
	return nil, errors.New("Received unknown invocation: " + function)
}

// ************************************
// query callback mode
// ************************************

// Query is called when a query message is received
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Handle different functions
	if function == "readAsset" {
		// gets the state for an assetID as a JSON struct
		return t.readAsset(stub, args)
	} else if function == "readAssetObjectModel" {
		return t.readAssetObjectModel(stub, args)
	} else if function == "readAssetSamples" {
		// returns selected sample objects
		return t.readAssetSamples(stub, args)
	} else if function == "readAssetSchemas" {
		// returns selected sample objects
		return t.readAssetSchemas(stub, args)
	} else if function == "readMruList" {
		// returns list of assets updated
		return t.readMruList(stub, args)
	}
	return nil, errors.New("Received unknown invocation: " + function)
}

/**********main implementation *************/

func main() {
	logger.SetLevel(shim.LogDebug)
	// logLevel, _ := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
	logLevel, _ := shim.LogLevel("DEBUG")
	shim.SetLoggingLevel(logLevel)
	logger.Error("Foo test logger")

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple Chaincode: %s", err)
	}
}

/*****************ASSET CRUD INTERFACE starts here************/

/****************** 'deploy' methods *****************/

/******************** createAsset ********************/

func (t *SimpleChaincode) createAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	_, erval := t.createOrUpdateAsset(stub, args)
	return nil, erval
}

//******************** updateAsset ********************/

func (t *SimpleChaincode) updateAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	_, erval := t.createOrUpdateAsset(stub, args)
	return nil, erval
}

//******************** deleteAsset ********************/

func (t *SimpleChaincode) deleteAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var assetID string // asset ID
	var err error
	var stateIn AssetState

	// validate input data for number of args, Unmarshaling to asset state and obtain asset id
	stateIn, err = t.validateInput(args)
	if err != nil {
		return nil, err
	}
	assetID = *stateIn.AssetID
	// Delete the key / asset from the ledger
	err = stub.DelState(assetID)
	if err != nil {
		err = errors.New("DELSTATE failed! : " + fmt.Sprint(err))
		return nil, err
	}
	return nil, nil
}

/******************* Query Methods ***************/

//********************readAsset********************/

func (t *SimpleChaincode) readAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var assetID string // asset ID
	var err error
	var state AssetState

	// validate input data for number of args, Unmarshaling to asset state and obtain asset id
	stateIn, err := t.validateInput(args)
	if err != nil {
		return nil, errors.New("Asset does not exist!")
	}
	assetID = *stateIn.AssetID
	// Get the state from the ledger
	assetBytes, err := stub.GetState(assetID)
	if err != nil || len(assetBytes) == 0 {
		err = errors.New("Unable to get asset state from ledger")
		return nil, err
	}
	err = json.Unmarshal(assetBytes, &state)
	if err != nil {
		err = errors.New("Unable to unmarshal state data obtained from ledger")
		return nil, err
	}
	return assetBytes, nil
}

//********************readAsset********************/

func (t *SimpleChaincode) readMruList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var state AssetMruList
	var err error
	assetMruListBytes, err := stub.GetState(MRUKEY)
	if err != nil || len(assetMruListBytes) == 0 {
		err = errors.New("Unable to get MruList from ledger")
		return nil, err
	}

	err = json.Unmarshal(assetMruListBytes, &state)
	if err != nil {
		err = errors.New("Unable to unmarshal MruList state data obtained from ledger")
		return nil, err
	}
	return assetMruListBytes, nil
}

//*************readAssetObjectModel*****************/

func (t *SimpleChaincode) readAssetObjectModel(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var state AssetState = AssetState{}

	// Marshal and return
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}
	return stateJSON, nil
}

//*************readAssetSamples*******************/

func (t *SimpleChaincode) readAssetSamples(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return []byte(samples), nil
}

//*************readAssetSchemas*******************/

func (t *SimpleChaincode) readAssetSchemas(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return []byte(schemas), nil
}

// ************************************
// validate input data : common method called by the CRUD functions
// ************************************
func (t *SimpleChaincode) validateInput(args []string) (stateIn AssetState, err error) {
	var assetID string       // asset ID
	var state = AssetState{} // The calling function is expecting an object of type AssetState

	if len(args) != 1 {
		err = errors.New("Incorrect number of arguments. Expecting a JSON strings with mandatory assetID")
		return state, err
	}
	jsonData := args[0]
	assetID = ""
	stateJSON := []byte(jsonData)
	err = json.Unmarshal(stateJSON, &stateIn)
	if err != nil {
		err = errors.New("Unable to unmarshal input JSON data")
		return state, err
		// state is an empty instance of asset state
	}
	// was assetID present?
	// The nil check is required because the asset id is a pointer.
	// If no value comes in from the json input string, the values are set to nil

	if stateIn.AssetID != nil {
		assetID = strings.TrimSpace(*stateIn.AssetID)
		if assetID == "" {
			err = errors.New("AssetID not passed")
			return state, err
		}
	} else {
		err = errors.New("Asset id is mandatory in the input JSON data")
		return state, err
	}

	stateIn.AssetID = &assetID
	return stateIn, nil
}

//****** minInt: helper function for createOrUpdateAsset ******/
func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

//******************** createOrUpdateAsset ********************/

func (t *SimpleChaincode) createOrUpdateAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var assetID string // asset ID                    // used when looking in map
	var err error
	var stateIn AssetState
	var stateStub AssetState

	// validate input data for number of args, Unmarshaling to asset state and obtain asset id

	stateIn, err = t.validateInput(args)
	if err != nil {
		return nil, err
	}
	assetID = *stateIn.AssetID
	// Partial updates introduced here
	// Check if asset record existed in stub
	assetBytes, err := stub.GetState(assetID)
	if err != nil || len(assetBytes) == 0 {
		// This implies that this is a 'create' scenario
		logger.Error(fmt.Sprintf("Error GetState for assetID (%s): %s", assetID, err))
		// fmt.Printf("Error GetState for assetID (%s): %s\n", assetID, err)
		stateStub = stateIn // The record that goes into the stub is the one that cme in
	} else {
		// This is an update scenario
		err = json.Unmarshal(assetBytes, &stateStub)
		if err != nil {
			fmt.Printf("Error Unmarshaling assetID (%s): %s\n", assetID, err)
			err = errors.New("Unable to unmarshal JSON data from stub")
			return nil, err
			// state is an empty instance of asset state
		}
		// Merge partial state updates
		stateStub, err = t.mergePartialState(stateStub, stateIn)
		if err != nil {
			fmt.Printf("Error Unable to merge state assetID (%s): %s\n", assetID, err)
			err = errors.New("Unable to merge state")
			return nil, err
		}
	}
	now := time.Now().UTC()
	stateStub.UpdatedAt = &now
	stateJSON, err := json.Marshal(stateStub)
	if err != nil {
		fmt.Printf("Error Marshaling assetID (%s): %s", assetID, err)
		return nil, errors.New("Marshal failed for contract state" + fmt.Sprint(err))
	}
	// Get existing state from the stub

	// Write the new state to the ledger
	err = stub.PutState(assetID, stateJSON)
	if err != nil {
		fmt.Printf("Error PutState assetID (%s): %s\n", assetID, err)
		err = errors.New("PUT ledger state failed: " + fmt.Sprint(err))
		return nil, err
	}
	// deal with MRUList structure
	var mruList AssetMruList
	mruBytes, err := stub.GetState(MRUKEY)
	logger.Error(fmt.Sprintf("Current mruBytes length: %v", len(mruBytes)))
	// fmt.Printf("Current mruBytes length: %v\n", len(mruBytes))
	if err != nil || len(mruBytes) == 0 { // No MRU List yet
		fmt.Printf("Error GetState MRUKEY (%s): %s\n", MRUKEY, err)
		mruList = AssetMruList{} //make a new one
	} else {
		err = json.Unmarshal(mruBytes, &mruList)
		if err != nil {
			logger.Error(fmt.Sprintf("Error Unmarshaling MRUKEY (%s): %s", MRUKEY, err))
			// fmt.Printf("Error Unmarshaling MRUKEY (%s): %s", MRUKEY, err)
			err = errors.New("Unable to unmarshal JSON data for mruList")
			return nil, err
		}
	}
	var mruEntry AssetUpdatedAt
	mruEntryJSON := fmt.Sprintf("{\"assetID\":\"%v\"}", assetID) // TODO: combine with next line
	json.Unmarshal([]byte(mruEntryJSON), &mruEntry)
	logger.Error(fmt.Sprintf("Current mruList length: %v", len(mruList.List)))
	// fmt.Printf("Current mruList length: %v\n", len(mruList.List))
	mruEntry.UpdatedAt = time.Now().UTC()
	mruList.List = append([]AssetUpdatedAt{mruEntry}, mruList.List...)[:minInt(len(mruList.List)+1, 10)]
	logger.Error(fmt.Sprintf("New mruList length: %v", len(mruList.List)))
	// fmt.Printf("New mruList length: %v\n", len(mruList.List))
	mruListJSON, err := json.Marshal(mruList)
	if err != nil {
		fmt.Printf("Error Marshaling mruList: %s", err)
		err = errors.New("Unable to marshal mruList")
		return nil, err
	}
	logger.Error(fmt.Sprintf("mruListJSON: %s", mruListJSON))
	// fmt.Printf("mruListJSON: %s\n", mruListJSON)
	err = stub.PutState(MRUKEY, mruListJSON)
	if err != nil {
		logger.Error(fmt.Sprintf("Error PutState MRUKEY (%s): %s", MRUKEY, err))
		err = errors.New("PUT ledger state failed: " + fmt.Sprint(err))
		return nil, err
	}

	return nil, nil
}

/*********************************  internal: mergePartialState ****************************/
func (t *SimpleChaincode) mergePartialState(oldState AssetState, newState AssetState) (AssetState, error) {

	old := reflect.ValueOf(&oldState).Elem()
	new := reflect.ValueOf(&newState).Elem()
	for i := 0; i < old.NumField(); i++ {
		oldOne := old.Field(i)
		newOne := new.Field(i)
		if !reflect.ValueOf(newOne.Interface()).IsNil() {
			oldOne.Set(reflect.Value(newOne))
		}
	}
	return oldState, nil
}
