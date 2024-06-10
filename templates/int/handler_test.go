package int

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handlerEkyc_CreateEkycProcesses(t *testing.T) {
	type fields struct {
		DB   *sqlx.DB
		TxID string
	}
	tests := []struct {
		name           string
		fields         fields
		mockDBPrep     func(*sqlx.DB)
		bodyReq        string
		expectedResult string
	}{
		{
			name:           "ErrorParsingRequestBody",
			bodyReq:        `{`,
			expectedResult: `{"Error":true,"Code":1,"Type":"error","Msg":"El cuerpo de la petición realizada no es valida"}`,
		},
		//TODO: Other test cases including Success, Failures etc
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handlerEkyc{
				DB:   tt.fields.DB,
				TxID: tt.fields.TxID,
			}

			// Mocking DB Prep
			if tt.mockDBPrep != nil {
				tt.mockDBPrep(h.DB)
			}

			app := fiber.New()
			app.Post("/ekyc/processes", h.CreateEkycProcesses)

			req := httptest.NewRequest("POST", "/ekyc/processes", strings.NewReader(tt.bodyReq))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusAccepted, resp.StatusCode)
			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResult, string(body))

			ecatch.ResetMock()
		})
	}
}
