// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package cache_test

import (
	"time"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/cache"
	"github.com/juju/juju/testing"
)

func NewStringsWatcherC(c *gc.C, watcher cache.StringsWatcher) StringsWatcherC {
	return StringsWatcherC{
		C:       c,
		Watcher: watcher,
	}
}

type StringsWatcherC struct {
	*gc.C
	Watcher cache.StringsWatcher
}

// AssertOneChange fails if no change is sent before a long time has passed; or
// if, subsequent to that, any further change is sent before a short time has
// passed.
func (c StringsWatcherC) AssertOneChange(expected []string) {
	select {
	case obtained, ok := <-c.Watcher.Changes():
		c.Assert(ok, jc.IsTrue)
		c.Assert(obtained, jc.SameContents, expected)
	case <-time.After(testing.LongWait):
		c.Fatalf("watcher did not send change")
	}
	c.AssertNoChange()
}

// AssertNoChange fails if it manages to read a value from Changes before a
// short time has passed.
func (c StringsWatcherC) AssertNoChange() {
	select {
	case _, ok := <-c.Watcher.Changes():
		if ok {
			c.Fatalf("watcher sent unexpected change")
		}
		c.Fatalf("watcher changes channel closed")
	case <-time.After(testing.ShortWait):
	}
}

// AssertStops Kills the watcher and asserts (1) that Wait completes without
// error before a long time has passed; and (2) that Changes channel is closed.
func (c StringsWatcherC) AssertStops() {
	c.Watcher.Kill()
	wait := make(chan error)
	go func() {
		wait <- c.Watcher.Wait()
	}()
	select {
	case <-time.After(testing.LongWait):
		c.Fatalf("watcher never stopped")
	case err := <-wait:
		c.Assert(err, jc.ErrorIsNil)
	}

	select {
	case _, ok := <-c.Watcher.Changes():
		if ok {
			c.Fatalf("watcher sent unexpected change")
		}
	default:
		c.Fatalf("channel not closed")
	}
}
