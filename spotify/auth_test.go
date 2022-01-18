package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	type Args struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	var goodArg Args
	f, _ := os.Open("config.json")
	j, _ := ioutil.ReadAll(f)
	err := json.Unmarshal(j, &goodArg)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(goodArg)

	badArg := Args{
		ClientID:     "sdfsdfsdfsfsd",
		ClientSecret: "dsfasdfsdfsdf",
	}

	tests := []struct {
		name    string
		args    Args
		wantErr bool
	}{
		{name: "goodKey", args: goodArg, wantErr: false},
		{name: "badKey", args: badArg, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Auth(tt.args.ClientID, tt.args.ClientSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("resp: %#v", got)
			t.Logf("err: %#v", err)
		})
	}
}
