package boundedkv

import (
	"testing"
)

func TestCM4(t *testing.T) {
	var counter CM4

	counts := [...]struct {
		site   string
		visits int
	}{
		{"Google.com", 130},
		{"Youtube.com", 150},
		{"Facebook.com", 10},
		{"Twitter.com", 50},
		{"Wikipedia.org", 6},
		{"Instagram.com", 254},
		{"Baidu.com", 61},
		{"Yahoo.com", 152},
	}

	for _, c := range counts {
		for i := 0; i < c.visits; i++ {
			counter.Add(c.site)
		}
	}

	for _, c := range counts {
		if est := int(counter.Estimate(c.site)); est > 16 {
			t.Error("violation of the implementation")
		}
	}

}
