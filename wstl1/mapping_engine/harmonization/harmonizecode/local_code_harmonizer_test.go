// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package harmonizecode

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp" /* copybara-comment: cmp */
)

func buildTestLocalHarmonizer(rawMaps []json.RawMessage) (CodeHarmonizer, error) {
	local := NewLocalCodeHarmonizer()
	for _, m := range rawMaps {
		cm, err := unmarshalR3ConceptMap(m)
		if err != nil {
			return nil, fmt.Errorf("unmarshal failed with error: %v", err)
		}
		if err := local.Cache(cm); err != nil {
			return nil, err
		}
	}
	return local, nil
}

func TestHarmonize(t *testing.T) {
	tests := []struct {
		name           string
		rawConceptMap  json.RawMessage
		sourceCode     string
		sourceSystem   string
		sourceName     string
		version        string
		expectedOutput []HarmonizedCode
	}{
		{
			name: "single target code",
			rawConceptMap: json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code":"abc",
								"target":[
									{
										"code":"def",
										"display": "DEF",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"target": "xyz"
					}
				],
				"id": "foo",
				"version": "bar",
   			"resourceType":"ConceptMap"
			}`),
			sourceCode:   "abc",
			sourceSystem: "foo",
			sourceName:   "foo",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "def",
					System:  "xyz",
					Display: "DEF",
					Version: "bar",
				},
			},
		},
		{
			name: "multiple target codes",
			rawConceptMap: json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def1",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"target": "xyz1"
					},
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def2",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"target": "xyz2"
					}
				],
				"id": "foo",
				"version": "bar",
   			"resourceType":"ConceptMap"
			}`),
			sourceCode:   "abc",
			sourceSystem: "foo",
			sourceName:   "foo",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "def1",
					System:  "xyz1",
					Version: "bar",
				},
				HarmonizedCode{
					Code:    "def2",
					System:  "xyz2",
					Version: "bar",
				},
			},
		},
		{
			name: "single target code in second group",
			rawConceptMap: json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code": "blah",
								"target":[
									{
										"code": "def",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"target": "xyz"
					},
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"target": "xyz"
					}
				],
				"id": "foo",
				"version": "bar",
   			"resourceType":"ConceptMap"
			}`),
			sourceCode:   "abc",
			sourceSystem: "foo",
			sourceName:   "foo",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "def",
					System:  "xyz",
					Version: "bar",
				},
			},
		},
		{
			name: "no matches found",
			rawConceptMap: json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"target": "xyz"
					}
				],
				"id": "foo",
				"version": "bar",
				"resourceType": "ConceptMap"
			}`),
			sourceCode:   "unmatched",
			sourceSystem: "foo",
			sourceName:   "foo",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "unmatched",
					System:  "foo-unharmonized",
					Version: "bar",
				},
			},
		},
		{
			name: "no matches found with provided unmapped mode",
			rawConceptMap: json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"unmapped": {
							"mode": "provided"
						},
						"target": "xyz",
						"source": "foo"
					}
				],
				"id": "map-id",
				"version": "bar",
				"resourceType": "ConceptMap"
			}`),
			sourceCode:   "unmatched",
			sourceSystem: "foo",
			sourceName:   "map-id",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "unmatched",
					Display: "unmatched",
					System:  "xyz",
					Version: "bar",
				},
			},
		},
		{
			name: "no matches found with fixed unmapped mode",
			rawConceptMap: json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"unmapped": {
							"mode": "fixed",
							"code": "unknown",
							"display": "Unknown Code"
						},
						"target": "xyz",
						"source": "foo"
					}
				],
				"id": "map-id",
				"version": "bar",
				"resourceType": "ConceptMap"
			}`),
			sourceCode:   "unmatched",
			sourceSystem: "foo",
			sourceName:   "map-id",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "unknown",
					Display: "Unknown Code",
					System:  "xyz",
					Version: "bar",
				},
			},
		},
		{
			name: "mixture of unmatched and matched codes",
			rawConceptMap: json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"unmapped": {
							"mode": "fixed",
							"code": "unknown",
							"display": "Unknown Code"
						},
						"target": "xyz1",
						"source": "foo1"
					},
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"unmapped": {
							"mode": "fixed",
							"code": "unknown",
							"display": "Unknown Code"
						},
						"target": "xyz2",
						"source": "foo2"
					},
					{
						"element":[
							{
								"code": "source-code",
								"target":[
									{
										"code": "def",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"unmapped": {
							"mode": "fixed",
							"code": "unknown",
							"display": "Unknown Code"
						},
						"target": "xyz3",
						"source": "foo3"
					}
				],
				"id": "map-id",
				"version": "bar",
				"resourceType": "ConceptMap"
			}`),
			sourceCode:   "source-code",
			sourceSystem: "foo2",
			sourceName:   "map-id",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "unknown",
					System:  "xyz2",
					Version: "bar",
					Display: "Unknown Code",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			harmonizer, err := buildTestLocalHarmonizer([]json.RawMessage{test.rawConceptMap})
			if err != nil {
				t.Fatalf("buildTestLocalHarmonizer returned unexpected error: %v", err)
			}

			actualOutput, err := harmonizer.Harmonize(test.sourceCode, test.sourceSystem, test.sourceName)
			if err != nil {
				t.Fatalf("Harmonize(%s, %s, %s) returned unexpected error: %v", test.sourceCode, test.sourceSystem, test.sourceName, err)
			}

			if diff := cmp.Diff(test.expectedOutput, actualOutput); diff != "" {
				t.Errorf("Harmonize(%s, %s, %s) => diff -%v +%v\n%s", test.sourceCode, test.sourceSystem, test.sourceName, test.expectedOutput, actualOutput, diff)
			}
		})
	}
}

func TestHarmonize_Errors(t *testing.T) {
	tests := []struct {
		name          string
		rawConceptMap json.RawMessage
	}{
		{
			name: "missing group",
			rawConceptMap: json.RawMessage(`{
				"id": "foo",
				"version": "bar",
   			"resourceType":"ConceptMap"
			}`),
		},
		{
			name: "missing target in element",
			rawConceptMap: json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code": "abc",
							}
						],
						"target": "xyz"
					}
				],
				"id": "foo",
				"version": "bar",
   			"resourceType":"ConceptMap"
			}`),
		},
		{
			name: "missing id",
			rawConceptMap: json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code": "abc",
							}
						],
						"target": "xyz"
					}
				],
				"version": "bar",
   			"resourceType":"ConceptMap"
			}`),
		},
		{
			name: "wrong fhir resource",
			rawConceptMap: json.RawMessage(`{
				"id": "abc",
				"gender": "male",
   			"resourceType":"Patient"
			}`),
		},
		{
			name: "unsupported unmapped mode",
			rawConceptMap: json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"unmapped": {
						  "mode": "other-map"
						},
						"target": "xyz"
					}
				],
				"id": "abc",
				"version": "bar",
   			"resourceType":"ConceptMap"
			}`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := buildTestLocalHarmonizer([]json.RawMessage{test.rawConceptMap})
			if err == nil {
				t.Fatalf("Parsing concept map in test %s expected error but received no errors.", test.name)
			}
		})
	}
}

func TestHarmonizeWithTarget(t *testing.T) {
	testConceptMap1 := json.RawMessage(`{
				"group":[
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def1",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"source": "s1",
						"target": "t1"
					},
					{
						"element":[
							{
								"code": "abc",
								"target":[
									{
										"code": "def2",
										"equivalence": "EQUIVALENT"
									}
								]
							}
						],
						"source": "s2",
						"target": "t2"
					}
				],
				"id": "foo",
				"version": "bar",
   			"resourceType":"ConceptMap"
			}`)
	tests := []struct {
		name           string
		rawConceptMap  json.RawMessage
		sourceCode     string
		sourceSystem   string
		targetSystem   string
		sourceName     string
		version        string
		expectedOutput []HarmonizedCode
	}{
		{
			name:          "match source and target systems 1",
			rawConceptMap: testConceptMap1,
			sourceCode:    "abc",
			sourceSystem:  "s1",
			targetSystem:  "t1",
			sourceName:    "foo",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "def1",
					System:  "t1",
					Version: "bar",
				},
			},
		},
		{
			name:          "match source and mismatch target 2",
			rawConceptMap: testConceptMap1,
			sourceCode:    "abc",
			sourceSystem:  "s2",
			targetSystem:  "t2",
			sourceName:    "foo",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "def2",
					System:  "t2",
					Version: "bar",
				},
			},
		},
		{
			name:          "match source and mismatch target",
			rawConceptMap: testConceptMap1,
			sourceCode:    "abc",
			sourceSystem:  "s1",
			targetSystem:  "t2",
			sourceName:    "foo",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "abc",
					System:  "foo-unharmonized",
					Version: "bar",
				},
			},
		},
		{
			name:          "match source and empty target",
			rawConceptMap: testConceptMap1,
			sourceCode:    "abc",
			sourceSystem:  "s2",
			targetSystem:  "",
			sourceName:    "foo",
			expectedOutput: []HarmonizedCode{
				HarmonizedCode{
					Code:    "def2",
					System:  "t2",
					Version: "bar",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			harmonizer, err := buildTestLocalHarmonizer([]json.RawMessage{test.rawConceptMap})
			if err != nil {
				t.Fatalf("buildTestLocalHarmonizer returned unexpected error: %v", err)
			}

			actualOutput, err := harmonizer.HarmonizeWithTarget(test.sourceCode, test.sourceSystem, test.targetSystem, test.sourceName)
			if err != nil {
				t.Fatalf("HarmonizeWithTarget(%s, %s, %s, %s) returned unexpected error: %v", test.sourceCode, test.sourceSystem, test.targetSystem, test.sourceName, err)
			}

			if diff := cmp.Diff(test.expectedOutput, actualOutput); diff != "" {
				t.Errorf("HarmonizeWithTarget(%s, %s, %s, %s) => diff -%v +%v\n%s", test.sourceCode, test.sourceSystem, test.targetSystem, test.sourceName, test.expectedOutput, actualOutput, diff)
			}
		})
	}
}
