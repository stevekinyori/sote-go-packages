package sNATS

import (
	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sJetStream"

	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
)

func TestNewNC(t *testing.T) {
	if jsmm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if !jsmm.Manager.IsJetStreamEnabled() {
			t.Errorf("Close Failed: IsJetStreamEnabled to return true")
		}
		jsmm.Close()
		if jsmm.Manager.IsJetStreamEnabled() {
			t.Errorf("Close Failed: IsJetStreamEnabled to return false")
		}
	} else {
		t.Errorf("NewNC Failed: Expected error code to be nil")
	}
}
func TestNewNCExpectErrorCode(t *testing.T) {
	if _, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("NewNC Failed: Expected error code to be nil")
	}
}
func TestNewNCExpect200513(t *testing.T) {
	if _, soteErr := NewNC("", sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode != 200513 {
		t.Errorf("NewNC Failed: Expected error code of 200513")
	}
}
func TestNewNCExpect209110(t *testing.T) {
	if _, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, "", "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode != 209110 {
		t.Errorf("NewNC Failed: Expected error code of 209110")
	}
}
func TestNewNCExpect210090(t *testing.T) {
	if _, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", "", 1, 250*time.Millisecond); soteErr.ErrCode != 210090 {
		t.Errorf("NewNC Failed: Expected error code of 210090")
	}
}
func TestNewNCExpectNil(t *testing.T) {
	if _, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 0, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("NewNC Failed: Expected error code to be nil")
	}

	if _, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 200, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("NewNC Failed: Expected error code to be nil")
	}

	if _, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 1*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("NewNC Failed: Expected error code to be nil")
	}

	if _, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 2*time.Minute); soteErr.ErrCode != nil {
		t.Errorf("NewNC Failed: Expected error code to be nil")
	}
}

func TestNatsManager_SPublish(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = ncm.SPublish("greeting", []byte("Hello world")); soteErr.ErrCode != nil {
			t.Errorf("SPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

}

func TestSPublishExpect208310(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = ncm.SPublish("", []byte("Hello world")); soteErr.ErrCode != 208310 {
			t.Errorf("SPublishExpect208310 Failed: Expected error code to be 208310 got %v", soteErr.FmtErrMsg)
		}
	}
}

func TestNatsManager_SPublishMsg(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = ncm.SPublishMsg(&nats.Msg{Data: []byte("SPublishMsg sending hello"), Subject: "greeting"}); soteErr.ErrCode != nil {
			t.Errorf("SPublishMsg Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}
}

func TestSPublishMsgExpect208310(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = ncm.SPublishMsg(&nats.Msg{}); soteErr.ErrCode != 208310 {
			t.Errorf("SPublishExpect208310 Failed: Expected error code to be 208310 got %v", soteErr.FmtErrMsg)
		}
	}
}

func TestNatsManager_SPublishRequest(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = ncm.SPublishRequest("subject", "I have received help", []byte("I need help")); soteErr.ErrCode != nil {
			t.Errorf("SPublishRequest Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}
}

func TestNatsManager_SSubscribe(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = ncm.SSubscribe("greeting"); soteErr.ErrCode != nil {
			t.Errorf("SSubscribe Failed: Expected error code nil but got %v", soteErr.FmtErrMsg)
		}

		if soteErr = ncm.SPublish("greeting", []byte("Good morning")); soteErr.ErrCode != nil {
			t.Errorf("SPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}
}

func TestSSubscribeExpect208310(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = ncm.SSubscribe(""); soteErr.ErrCode != 208310 {
			t.Errorf("SSubscribeExpect208310 Failed: Expected error code nil but got %v", soteErr.FmtErrMsg)
		}

		if soteErr = ncm.SPublish("greeting", []byte("Good morning")); soteErr.ErrCode != nil {
			t.Errorf("SPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}
}

func TestNatsManager_SNextMsg(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		if sub, err := ncm.NC.SubscribeSync("synchronous"); err == nil {

			if soteErr = ncm.SPublish("synchronous", []byte("Good morning!")); soteErr.ErrCode != nil {
				t.Errorf("SPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			}

			if soteErr = ncm.SNextMsg(sub, time.Second); soteErr.ErrCode != nil {
				t.Errorf("SNextMsg Failed: Expected error code nil but got %v", soteErr.FmtErrMsg)
			}
		}

	}
}

func TestNatsManager_SRequest(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = ncm.SRequestReply("help", []byte("Help is on the way")); soteErr.ErrCode != nil {
			t.Errorf("SRequestReply Failed: Expected error code nil but got %v", soteErr.FmtErrMsg)
		}

		if soteErr = ncm.SRequest("help", []byte("I need help"), time.Second); soteErr.ErrCode != nil {
			t.Errorf("SRequest Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}
}


func TestNatsManager_SRequestMsg(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		m := &nats.Msg{Data: []byte("SRequestMsg requesting help"), Subject: "help"}

		if _, soteErr = ncm.SRequestReply("help", []byte("Help is not available")); soteErr.ErrCode != nil {
			t.Errorf("SRequestReply Failed: Expected error code nil but got %v", soteErr.FmtErrMsg)
		}

		if soteErr = ncm.SRequestMsg(m, time.Second); soteErr.ErrCode != nil {
			t.Errorf("SRequest Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}
}

func TestNatsManager_SRequestReply(t *testing.T) {
	if ncm, soteErr := NewNC(sJetStream.TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", sJetStream.TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = ncm.SRequestReply("help", []byte("More help is on the way")); soteErr.ErrCode != nil {
			t.Errorf("SRequestReply Failed: Expected error code nil but got %v", soteErr.FmtErrMsg)
		}

		if soteErr = ncm.SRequest("help", []byte("I need help"), time.Second); soteErr.ErrCode != nil {
			t.Errorf("SRequest Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}
}