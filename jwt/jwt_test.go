package jwt

import (
	"testing"
	"time"

	"google.golang.org/appengine/aetest"
)

func TestJWT(t *testing.T) {
	c, done, err := aetest.NewContext()
	defer done()

	now := time.Now()
	cls1 := claims{Sub: "user-id", Exp: now.Add(1 * time.Hour).Unix()}

	jwt, err := toJWT(&cls1, c)
	if err != nil {
		t.Fatal("could not make JWT from user")
	}

	cls2, err := fromJWT(jwt, c)
	if err != nil {
		t.Fatalf("could not get user from JWT (%v)", err)
	}
	if cls2 == nil {
		t.Fatalf("should have decoded claims")
	}
	if cls2.Sub != cls1.Sub {
		t.Errorf("should have the same sub")
	}
	if cls2.Exp != cls1.Exp {
		t.Errorf("should have the same exp")
	}
}
