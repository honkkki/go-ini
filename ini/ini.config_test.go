package ini

import (
	"io/ioutil"
	"testing"
)


func TestIniConfig(t *testing.T) {
	data, err := ioutil.ReadFile("./config.ini")
	if err != nil {
		t.Error("failed to read file")
	}

	conf := &Config{}
	err = UnMarshal(data, conf)
	if err != nil {
		t.Error(err)
	}

	t.Logf("success: %#v\n", conf)
}
