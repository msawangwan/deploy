package dock

// import (
// 	"testing"
// )

// type mockAPICommandURL struct {
// 	MockEndPoint string
// 	MockParams   []string
// }

// func (m mockAPICommandURL) Build() []byte {
// 	return []byte(
// 		`{{- with . -}}
// 			mock endpoint: {{- .MockEndPoint -}} 
// 			mock params: {{- range .MockParams -}}
// 				{{- . -}}&
// 			{{- end -}}
// 		{{- end -}}`,
// 	)
// }

// func TestBuildAPIURLStrings(t *testing.T) {
// 	m := &mockAPICommandURL{
// 		"SOME_ENDPOINT",
// 		[]string{"ONE_STRING", "TWO_STRING", "THREE_STRING"},
// 	}

// 	result, err := BuildAPIURLString(m)

// 	if err != nil {
// 		t.Errorf("%s", err)
// 	} else {
// 		t.Logf("%s", result)
// 	}
// }

// func TestBuildAPIURLs(t *testing.T) {
// 	apiurls := []APIStringBuilder{
// 		&ContainerCommand{
// 			URLComponents{
// 				Command: "containers",
// 				Option:  "json",
// 			},
// 		},
// 		&ContainerCommand{
// 			URLComponents{
// 				Command: "containers",
// 				Option:  "create",
// 				Parameters: map[string]string{
// 					"name": "container_name",
// 				},
// 			},
// 		},
// 		&ContainerCommandByID{
// 			URLComponents{
// 				Command: "containers",
// 				Option:  "json",
// 			},
// 			"1234598765abcdefg",
// 		},
// 		&ContainerCommandByID{
// 			URLComponents{
// 				Command: "containers",
// 				Option:  "start",
// 			},
// 			"1234598765abcdefg",
// 		},
// 		&ContainerCommandByID{
// 			URLComponents{
// 				Command: "containers",
// 				Option:  "stop",
// 				Parameters: map[string]string{
// 					"some_param":    "some_value",
// 					"another_param": "another_value",
// 				},
// 			},
// 			"1234598765abcdefg",
// 		},
// 		&BuildCommand{
// 			QueryStrings: map[string]string{
// 				"somekey":    "someval",
// 				"anotherkey": "anotherval",
// 				"finalkey":   "finalval",
// 				"somebool":   "true",
// 			},
// 		},
// 		NewBuildDockerfileCommand(map[string]string{"a": "b", "c": "d"}),
// 		NewContainerCommandByID("DELETE", "containers", "", "1234567899k"),
// 	}

// 	for _, apiurl := range apiurls {
// 		res, err := BuildAPIURLString(apiurl)
// 		if err != nil {
// 			t.Errorf("%s", err)
// 		} else {
// 			t.Logf("%s", res)
// 		}
// 	}
// }
