package gaefire

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func newTestSwaggerJson() SwaggerJsonModel {
	if buf, err := NewAssetManager().LoadFile("www/swagger/swagger.json"); err != nil {
		panic(err)
	} else {
		result := SwaggerJsonModel{}
		if err := json.Unmarshal(buf, &result); err != nil {
			panic(err)
		}
		return result
	}
}

func TestSwaggerJsonParse(t *testing.T) {
	swagger := newTestSwaggerJson()
	assert.NotEqual(t, swagger.Host, "")
	assert.Equal(t, swagger.SecurityDefinitions.GoogleIdToken.Issuer, "https://accounts.google.com")
	assert.NotEqual(t, len(swagger.SecurityDefinitions.GoogleIdToken.Audiences), 0)
	assert.Equal(t, swagger.SecurityDefinitions.ApiKey.Name, "key")
	assert.Equal(t, swagger.SecurityDefinitions.ApiKey.In, "query")
}

func TestServiceCheckModel(t *testing.T) {
	JSON := `
	{
	  "operation": {
	    "operationId": "123e4567-e89b-12d3-a456-426655440000",
	    "operationName": "check:123e4567-e89b-12d3-a456-426655440000",
	    "consumerId": "api_key:replace.your.api.key",
	    "startTime": "2017-03-08T00:00:00Z"
	  }
	}
	`

	model := ServiceCheckModel{}
	assert.Nil(t, json.Unmarshal([]byte(JSON), &model))

	assert.Equal(t, model.Operation.OperationId, "123e4567-e89b-12d3-a456-426655440000")
	assert.Equal(t, model.Operation.OperationName, "check:123e4567-e89b-12d3-a456-426655440000")
	assert.Equal(t, model.Operation.ConsumerId, "api_key:replace.your.api.key")
	assert.Equal(t, model.Operation.StartTime, "2017-03-08T00:00:00Z")
}

func TestServiceCheckResultModel(t *testing.T) {
	JSON := `{
		  "operationId": "123e4567-e89b-12d3-a456-426655440000",
		  "checkErrors": [
		    {
		      "code": "API_KEY_INVALID",
		      "detail": "Invalid API key: api.key"
		    }
		  ],
		  "serviceConfigId": "2017-03-08r0",
		  "checkInfo": {
		    "unusedArguments": [
		      "servicecontrol.googleapis.com/caller_ip",
		      "servicecontrol.googleapis.com/user"
		    ]
		  }
	}
	`

	model := ServiceCheckResultModel{}
	assert.Nil(t, json.Unmarshal([]byte(JSON), &model))

	assert.Equal(t, len(*model.CheckErrors), 1)
}

func TestServiceCheckResultModel_NotError(t *testing.T) {
	JSON := `{
		  "operationId": "123e4567-e89b-12d3-a456-426655440000",
		  "serviceConfigId": "2017-03-08r0",
		  "checkInfo": {
		    "unusedArguments": [
		      "servicecontrol.googleapis.com/caller_ip",
		      "servicecontrol.googleapis.com/user"
		    ]
		  }
		}
	`

	model := ServiceCheckResultModel{}
	assert.Nil(t, json.Unmarshal([]byte(JSON), &model))

	assert.True(t, model.CheckErrors == nil)
}